package snowflake

import (
	"database/sql"

	sf "github.com/snowflakedb/gosnowflake"
)

const (
	snowflakeDriver = "snowflake"
)

var snowFlakeConfig *sf.Config

type dbConnection struct {
	session *sql.DB
}

//Load config loads config to ssnowflake config
func LoadConfig(conf DBConfig) {
	snowFlakeConfig = &sf.Config{
		User:      conf.User,
		Password:  conf.Password,
		Account:   conf.Account,
		Database:  conf.Database,
		Host:      conf.Host,
		Role:      conf.Role,
		Schema:    conf.Schema,
		Warehouse: conf.Warehouse,
	}
}

//NewConnection creates connection with snowflake
func NewConnection() (*dbConnection, error) {
	db := &dbConnection{}
	dsn, err := sf.DSN(snowFlakeConfig)
	if err != nil {
		return nil, err
	}
	db.session, err = sql.Open(snowflakeDriver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//convertToMap convert rows to map
func convertToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	var response []map[string]interface{}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		response = append(response, m)
	}
	return response, nil
}

//ExecDBQuery helps execute query on snowflake database
func (db dbConnection) ExecDBQuery(query string, values ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.session.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res, err := convertToMap(rows)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//Close closes a connection
func (db dbConnection) Close() {
	db.session.Close()
}
