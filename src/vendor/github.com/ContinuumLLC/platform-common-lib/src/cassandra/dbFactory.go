package cassandra

import (
	"sync"

	"github.com/ContinuumLLC/platform-common-lib/src/web/rest"
)

//Factory ...
type Factory interface {
	GetDbConnector(cfg *DbConfig) (DbConnector, error)
	GetNewDbConnector(cfg *DbConfig) (DbConnector, error)
	Health(cfg *DbConfig) rest.Statuser
}

//FactoryImpl ...
type FactoryImpl struct{}

//NewDbConnection is a factory method which returns the struct implementation of DbConnector
func NewDbConnection(conf *DbConfig) (DbConnector, error) {
	return newConnection(conf)
}

//NewBatchDbConnection is a factory method which returns the struct implementation of BatchDbConnector
func NewBatchDbConnection(conf *DbConfig) (BatchDbConnector, error) {
	return newBatchConnection(conf)
}

//GetDbConnector is a factory for Cassandra single Session Creation
func (FactoryImpl) GetDbConnector(cfg *DbConfig) (DbConnector, error) {
	if !sessionInitialized || session.Closed() {
		mu.Lock()
		defer mu.Unlock()
		if !sessionInitialized || session.Closed() {
			if sessionInitialized {
				session.closeSuper()
			}

			db, err := NewDbConnection(cfg) //GetDbConnector(cfg)
			if err != nil {
				return nil, err
			}
			session = cassandraSession{db}
			sessionInitialized = true
		}
	}
	return session, nil
}

//GetBatchDbConnector is a factory for Cassandra single Session Creation
func (FactoryImpl) GetBatchDbConnector(cfg *DbConfig) (BatchDbConnector, error) {
	if !batchSessionInitialized || batchSession.Closed() {
		mu.Lock()
		defer mu.Unlock()
		if !batchSessionInitialized || batchSession.Closed() {
			if batchSessionInitialized {
				batchSession.closeSuper()
			}

			batchdb, err := NewBatchDbConnection(cfg) //GetBatchDbConnector(cfg)
			if err != nil {
				return nil, err
			}
			batchSession = batchCassandraSession{batchdb}
			batchSessionInitialized = true
		}
	}
	return batchSession, nil
}

//GetNewDbConnector is a factory for Cassandra single Session Creation
func (FactoryImpl) GetNewDbConnector(cfg *DbConfig) (DbConnector, error) {
	if !new1SessionInitialized || new1Session.Closed() {
		mu.Lock()
		defer mu.Unlock()
		if !new1SessionInitialized || new1Session.Closed() {
			if new1SessionInitialized {
				new1Session.closeSuper()
			}

			db, err := NewDbConnection(cfg) //GetDbConnector(cfg)
			if err != nil {
				return nil, err
			}
			new1Session = cassandraSession{db}
			new1SessionInitialized = true
		}
	}
	return new1Session, nil
}

//Health is a function for Cassandra Health
func (FactoryImpl) Health(cfg *DbConfig) rest.Statuser {
	return status{
		conf: cfg,
	}
}

type cassandraSession struct {
	DbConnector
}

type batchCassandraSession struct {
	BatchDbConnector
}

func (d cassandraSession) Close() {
}

func (d cassandraSession) closeSuper() {
	d.DbConnector.Close()
}

func (d batchCassandraSession) Close() {
}

func (d batchCassandraSession) closeSuper() {
	d.BatchDbConnector.Close()
}

var mu sync.Mutex
var session cassandraSession
var batchSession batchCassandraSession
var sessionInitialized = false
var batchSessionInitialized = false
var new1Session cassandraSession
var new1SessionInitialized = false

//Status used for getting status of Casandra connection
type status struct {
	conf *DbConfig
}

//Status used for getting status of Casandra connection
func (c status) Status(conn rest.OutboundConnectionStatus) *rest.OutboundConnectionStatus {
	conn.ConnectionType = "Cassandra"
	conn.ConnectionURLs = c.conf.Hosts

	session, err := NewDbConnection(c.conf)

	if err != nil {
		conn.ConnectionStatus = rest.ConnectionStatusUnavailable
		return &conn
	}

	defer session.Close()

	conn.ConnectionStatus = rest.ConnectionStatusActive

	return &conn
}
