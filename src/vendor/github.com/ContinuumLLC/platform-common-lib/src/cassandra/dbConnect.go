package cassandra

import (
	"time"

	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"

	"github.com/gocql/gocql"
)

//dbConnection is responsible of connecting with Cassandra db
type dbConnection struct {
	session *gocql.Session
	cluster *gocql.ClusterConfig
}

//newDbConnection is a constructor of dbConnection which will intialize struct and will return an open connection object(if no error) of dbConnection
func newDbConnection(conf *DbConfig) (*dbConnection, error) {
	db := &dbConnection{}

	if len(conf.Hosts) == 0 || conf.Keyspace == "" {
		return db, exc.New(ErrDbHostsAndKeyspaceRequired, nil)
	}

	db.cluster = gocql.NewCluster()
	db.cluster.Hosts = conf.Hosts
	db.cluster.ProtoVersion = protoVersion
	db.cluster.Keyspace = conf.Keyspace

	if conf.TimeoutMillisecond != 0 {
		db.cluster.Timeout = conf.TimeoutMillisecond * time.Millisecond
	}

	var err error
	//If there are some error in connecting to the cluster, below method also does a log.Printf() and logs the error
	db.session, err = db.cluster.CreateSession()
	if err != nil {
		err = exc.New(ErrDbUnableToConnect, err)
	}
	return db, err
}

func (d dbConnection) Insert(query string, value ...interface{}) error {
	return d.executeDmlQuery(query, value...)
}

func (d dbConnection) Update(query string, value ...interface{}) error {
	return d.executeDmlQuery(query, value...)
}

func (d dbConnection) Delete(query string, value ...interface{}) error {
	return d.executeDmlQuery(query, value...)
}

func (d dbConnection) executeDmlQuery(query string, value ...interface{}) error {
	//TODO - need to decide if validation needs to be added for the function parameter
	if d.session == nil {
		return exc.New(ErrDbNoOpenConnection, nil)
	}
	err := d.session.Query(query, value...).Exec()
	if err != nil {
		return exc.New(ErrDbDMLFailed, err)
	}
	return nil
}

func (d dbConnection) InsertWithScanCas(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	return d.executeScanCasQuery(dest, query, value...)
}

func (d dbConnection) UpdateWithScanCas(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	return d.executeScanCasQuery(dest, query, value...)
}

func (d dbConnection) DeleteWithScanCas(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	return d.executeScanCasQuery(dest, query, value...)
}

//executeScanCasQuery executes a transaction (query with an IF statement).
//If the transaction is successfully executed, it returns true. If the transaction fails
//beacuse the IF condition was not satisfied, it returns false and populates dest(only for
//insert AND update query) with the existing values in cassandra
//https://godoc.org/github.com/gocql/gocql#Query.MapScanCAS
func (d dbConnection) executeScanCasQuery(dest *map[string]interface{}, query string, value ...interface{}) (bool, error) {
	resultMap := make(map[string]interface{})
	if d.session == nil {
		return false, exc.New(ErrDbNoOpenConnection, nil)
	}
	isApplied, err := d.session.Query(query, value...).MapScanCAS(resultMap)
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

func (d dbConnection) Select(query string, value ...interface{}) ([]map[string]interface{}, error) {
	data, err := d.session.Query(query, value...).Iter().SliceMap()
	if err != nil {
		return nil, exc.New(ErrDbUnableToFetchRecord, err)
	}
	return data, nil
}

func (d dbConnection) SelectWithPaging(page int, callback ProcessRow, query string, value ...interface{}) error {
	q := d.session.Query(query, value...)
	q.PageSize(page)
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
func (dbConnection) GetRandomUUID() (string, error) {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return "", exc.New(ErrUUID, err)
	}
	return uuid.String(), nil
}

//Close function closes the connection and does not return error
func (d dbConnection) Close() {
	if d.session != nil {
		d.session.Close()
	}
}

//Closed function to check is session is closed or not
func (d dbConnection) Closed() bool {
	if d.session != nil {
		return d.session.Closed()
	}
	return true
}
