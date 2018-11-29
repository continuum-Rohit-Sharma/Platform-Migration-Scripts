package cassandra

import (
	"fmt"
)

//BatchQueryExecutor is based on classic builder pattern which leverage user to add,
//as many query with ease with differnt level of code execution unit and
//at the end will Execute instruction will instruct batch to execute al queries
type BatchQueryExecutor interface {
	//AddQuery is add query with its args to the batch and return the same instanceso that user
	//can keep on adding or appending queries with its args and once they are done can
	//instruct execute command to execute all queries at once
	AddQuery(query string, args ...interface{}) BatchQueryExecutor

	//Execute to instruct batch to execute all queries at once
	Execute() error
}

//GetBatchQueryExecutor is factory method to get BatchQueryExecutor
func GetBatchQueryExecutor(conf *DbConfig) (BatchQueryExecutor, error) {
	if conf == nil {
		return nil, fmt.Errorf("config : %+v is mandatory for getting executor of batch", conf)
	}
	batchConnection, err := newBatchConnection(conf)
	if err != nil {
		return nil, err
	}
	return &batchQueryExecutorImpl{batchConnection}, nil
}

type batchQueryExecutorImpl struct {
	batchConnection *batchConnection
}

func (b *batchQueryExecutorImpl) AddQuery(query string, args ...interface{}) BatchQueryExecutor {
	b.batchConnection.batch.Query(query, args...)
	return b
}

func (b *batchQueryExecutorImpl) Execute() error {
	defer b.batchConnection.Close()
	return b.batchConnection.executeBatch(b.batchConnection.batch)
}
