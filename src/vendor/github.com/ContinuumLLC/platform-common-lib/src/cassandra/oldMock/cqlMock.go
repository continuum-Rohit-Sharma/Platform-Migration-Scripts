package oldMock

import (
	"github.com/ContinuumLLC/platform-common-lib/src/cassandra"
	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/gocql/gocql"
)

type CqlMock struct {
	Data *CqlMockData
}

type CqlMockData struct {
	GetDbConnErr        error
	InsertErr           error
	InsertCalled        *int
	UpdateErr           error
	DeleteErr           error
	SelectResult        []map[string]interface{}
	SelectErr           error
	SelectWithPagingErr error
}

func (mock CqlMock) GetDbConnector(cfg *cassandra.DbConfig) (cassandra.DbConnector, error) {
	return &mock, (*mock.Data).GetDbConnErr
}

func (mock CqlMock) Insert(query string, value ...interface{}) error {
	*(*mock.Data).InsertCalled = *(*mock.Data).InsertCalled + 1
	return (*mock.Data).InsertErr
}
func (mock CqlMock) Update(query string, value ...interface{}) error {
	return (*mock.Data).UpdateErr
}
func (mock CqlMock) Delete(query string, value ...interface{}) error {
	return (*mock.Data).DeleteErr
}
func (mock CqlMock) Close() {}

func (mock CqlMock) Select(query string, value ...interface{}) ([]map[string]interface{}, error) {
	return (*mock.Data).SelectResult, (*mock.Data).SelectErr
}

func (mock CqlMock) SelectWithPaging(page int, callback cassandra.ProcessRow, query string, value ...interface{}) error {
	return (*mock.Data).SelectWithPagingErr
}

func (mock CqlMock) GetRandomUUID() (string, error) {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return "", exc.New("ErrUUID", err)
	}
	return uuid.String(), nil
}

func (mock CqlMock) Closed() bool { return true }
