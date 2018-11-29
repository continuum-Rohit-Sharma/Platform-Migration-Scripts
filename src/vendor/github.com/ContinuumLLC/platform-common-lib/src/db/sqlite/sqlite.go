package sqlite

//Config is a struct to define sqlite configuration
type Config struct {
	DBName string
}

//Service is an interface holds all the functions related to sqlite
type Service interface {
	Init() error
	Close() error
	CreateTable(Table interface{}) error
	Add(record interface{}) error
	AddAll(records []interface{}) error
	//It will update all fields, even it is not changed
	Update(record interface{}) error
	Delete(record interface{}) error
	DeleteWhere(whereQuery, whereArgs, out interface{}) error
	Get(limit int, out interface{}) error
	GetWhere(limit int, whereQuery, whereArgs, out interface{}) error
	GetWhereObject(where, out interface{}) error
	// Update multiple attributes with `struct`, will only update those changed & non blank fields
	Set(out interface{}) error
	// Similar to Set function but you can specify where clause to update rows
	SetWhere(whereQuery, whereArgs, out interface{}) error
	Execute(out interface{}, sql string, values ...interface{}) error
	FirstOrCreate(where, out interface{}) error
}
