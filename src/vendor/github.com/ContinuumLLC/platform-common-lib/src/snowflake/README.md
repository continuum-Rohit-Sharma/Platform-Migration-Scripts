Snowflake Connection Wrapper

LoadSnowflake(DBConfig) Loads the configuration for connection.

NewConnection() Return connections on successful creation.

ExecDBQuery(query string,values...)([]map[string]interface{}, error)    Executes query and returns map of the result set.

Close() Closes a database connection.
