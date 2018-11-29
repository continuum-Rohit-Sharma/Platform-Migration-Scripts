package snowflake

//DBConfig is config for snowflake connection
type DBConfig struct {
	User      string
	Password  string
	Account   string
	Database  string
	Host      string
	Role      string
	Schema    string
	Warehouse string
}
