package cassandra

import (
	"time"

	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"

	"github.com/gocql/gocql"
)

//Row is a callback function to process database Row
type Row func(map[string]interface{})

//connection is responsible of connecting with Cassandra db
type connection struct {
	session *gocql.Session
	cluster *gocql.ClusterConfig
}

//batchConnection is responsible of connecting with Cassandra db and creating batch
type batchConnection struct {
	Connection connection
	batch      *gocql.Batch
}

//Cassandra Protocol version is set to be 4 if you are using cassandra >= 3.0
const protoVersion = 4

//newConnection is a constructor of connection which will intialize struct and will return an open connection object(if no error) of connection
func newConnection(conf *DbConfig) (*connection, error) {
	db := &connection{}

	if len(conf.Hosts) == 0 || conf.Keyspace == "" {
		return db, exc.New(ErrDbHostsAndKeyspaceRequired, nil)
	}

	db.cluster = gocql.NewCluster()
	db.cluster.Hosts = conf.Hosts
	db.cluster.ProtoVersion = protoVersion
	db.cluster.Keyspace = conf.Keyspace
	db.cluster.Consistency = gocql.Quorum
	db.cluster.Compressor = gocql.SnappyCompressor{}
	db.cluster.NumConns = 20

	if conf.TimeoutMillisecond != 0 {
		db.cluster.Timeout = time.Duration(conf.TimeoutMillisecond) * time.Millisecond
	} else {
		db.cluster.Timeout = 1 * time.Second //default timeout is 1 second
	}

	var err error
	//If there are some error in connecting to the cluster, below method also does a log.Printf() and logs the error
	db.session, err = db.cluster.CreateSession()
	if err != nil {
		err = exc.New(ErrDbUnableToConnect, err)
	}
	return db, err
}

//newBatchDbConnection is a constructor of batchDbConnection which will intialize struct and will return an open connection object(if no error) of batchDbConnection
func newBatchConnection(conf *DbConfig) (*batchConnection, error) {
	batchdb := &batchConnection{}
	db, err := newConnection(conf)
	if err != nil {
		return nil, err
	}
	batchdb.Connection = *db
	batchdb.batch = db.session.NewBatch(gocql.LoggedBatch)
	return batchdb, err

}
func (d batchConnection) BatchExecution(query string, values [][]interface{}) (err error) {
	length := len(values)
	for i := 0; i < length; i++ {
		d.batch.Query(query, values[i]...)
	}
	return d.executeBatch(d.batch)
}

func (d connection) Insert(query string, value ...interface{}) error {
	return d.executeDmlQuery(query, value...)
}

func (d connection) Update(query string, value ...interface{}) error {
	return d.executeDmlQuery(query, value...)
}

func (d connection) Delete(query string, value ...interface{}) error {
	return d.executeDmlQuery(query, value...)
}

func (d batchConnection) executeBatch(b *gocql.Batch) error {

	defer d.Connection.session.Close()
	if d.Connection.session == nil {

		return exc.New(ErrDbNoOpenConnection, nil)
	}
	err := d.Connection.session.ExecuteBatch(b)
	if err != nil {
		return err
	}
	return nil
}

func (d connection) InsertWithScanCas(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	return d.executeScanCasQuery(dest, query, value...)
}

func (d connection) UpdateWithScanCas(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	return d.executeScanCasQuery(dest, query, value...)
}

func (d connection) DeleteWithScanCas(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	return d.executeScanCasQuery(dest, query, value...)
}

//executeScanCasQuery executes a transaction (query with an IF statement).
//If the transaction is successfully executed, it returns true. If the transaction fails
//beacuse the IF condition was not satisfied, it returns false and populates dest(only for
//insert AND update query) with the existing values in cassandra.
//https://godoc.org/github.com/gocql/gocql#Query.MapScanCAS
func (d connection) executeScanCasQuery(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	resultMap := make(map[string]interface{})
	if d.session == nil {
		return false, exc.New(ErrDbNoOpenConnection, nil)
	}
	q := d.session.Query(query, value...)
	isApplied, err := q.MapScanCAS(resultMap)
	defer q.Release()
	if err != nil {
		return false, exc.New(ErrDbDMLFailed, err)
	}
	//This has been purposefully done here so that after resultMap is populated,
	//its data can be copied to dest
	if dest != nil {
		*dest = resultMap
	}
	return isApplied, nil
}

func (d connection) executeDmlQuery(query string, value ...interface{}) error {
	//TODO - need to decide if validation needs to be added for the function parameter
	if d.session == nil {
		return exc.New(ErrDbNoOpenConnection, nil)
	}
	q := d.session.Query(query, value...)
	err := q.Exec()
	defer q.Release()
	if err != nil {
		return exc.New(ErrDbDMLFailed, err)
	}
	return nil
}

func (d connection) Select(query string, value ...interface{}) ([]map[string]interface{}, error) {
	q := d.session.Query(query, value...) //.Consistency(gocql.One)
	iter := q.Iter()
	data, err := iter.SliceMap()
	defer q.Release()
	if err != nil {
		return nil, exc.New(ErrDbUnableToFetchRecord, err)
	}
	err = iter.Close()
	if err != nil {
		return data, exc.New(ErrDbUnableToFetchRecord, err)
	}
	return data, nil
}

func (d connection) SelectWithPaging(page int, callback ProcessRow, query string, value ...interface{}) error {
	q := d.session.Query(query, value...) //.Consistency(gocql.One)
	q.PageSize(page)
	defer q.Release()
	iter := q.Iter()
	m := make(map[string]interface{})
	for {
		if !iter.MapScan(m) {
			break
		}
		callback(m)
		m = make(map[string]interface{})
	}
	return iter.Close()
}

// GetRandomUUID() returns the random generated UUID
func (connection) GetRandomUUID() (string, error) {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return "", exc.New(ErrUUID, err)
	}
	return uuid.String(), nil
}

//Close function closes the connection and does not return error
func (d connection) Close() {
	if d.session != nil {
		d.session.Close()
	}
}

//Close function closes the connection and does not return error
func (d batchConnection) Close() {
	if d.Connection.session != nil {
		d.Connection.session.Close()
	}
}

//Closed function to check is session is closed or not
func (d connection) Closed() bool {
	if d.session != nil {
		return d.session.Closed()
	}
	return true
}

//Closed function to check is session is closed or not
func (d batchConnection) Closed() bool {
	if d.Connection.session != nil {
		return d.Connection.session.Closed()
	}
	return true
}
