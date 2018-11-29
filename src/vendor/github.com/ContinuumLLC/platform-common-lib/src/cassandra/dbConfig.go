package cassandra

import (
	"time"
)

//DbConfig is a set of configuration which is required to connect to Cassandra db.
type DbConfig struct {
	//Hosts is a mandatory field, port can be passed along with every hosts if it is different than 9042.
	//9042 is the default port Cassandra runs on
	Hosts []string
	//Keyspace is mandatory
	Keyspace string

	TimeoutMillisecond time.Duration
}
