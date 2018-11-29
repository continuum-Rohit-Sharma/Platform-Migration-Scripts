package main

import (
	"fmt"

	"github.com/ContinuumLLC/platform-common-lib/src/cassandra"
)

func main() {

	conf := &cassandra.DbConfig{}
	conf.Hosts = []string{"192.168.137.122:9042"}
	conf.Keyspace = "broker_keyspace"

	//Test-1 : Check with valid hosts/port and keyspace
	fmt.Println("Test-1 : Check with valid hosts/port and keyspace")
	testConnection(conf)

	//Test-2 : Check with blank hosts and/or keyspace
	fmt.Println("Test-2 : Check with blank hosts and/or keyspace")
	conf.Hosts = []string{}
	testConnection(conf)

	//Test-3 : Check with invalid host,port or keyspace
	fmt.Println("Test-3 : Check with invalid host,port or keyspace")
	conf.Hosts = []string{"192.168.137.122:9041"} //9041 is not the port where Cassandra is running
	testConnection(conf)

	//Test-4 : Check insert
	fmt.Println("Test-4 : Check insert")
	conf.Hosts = []string{"192.168.137.122:9042"}
	conf.Keyspace = "broker_keyspace"
	testInsert(conf, "INSERT INTO agent_heartbeat (regid,counter,updated_time_utc) VALUES(?,?,dateof(now()))", 111, 2)

	//Test-5 : Check update
	fmt.Println("Test-5 : Check update")
	testUpdate(conf, "UPDATE agent_heartbeat SET counter=? WHERE regid=?", 50, 111)

	//Test-6 : Check delete
	fmt.Println("Test-6 : Check delete")
	testDelete(conf, "DELETE FROM agent_heartbeat WHERE regid=?", 111)

	//Test-7 : Select check
	fmt.Println("Test-7 : Select check")
	testSelect(conf, " select * from p2_performance_db.perfdatacounter where regid=? and objectname=? and countername=? and dayvalue=?", 2, "Processor", "% Interrupt Time", 20160103)
}

func testConnection(conf *cassandra.DbConfig) {
	db, err := cassandra.New(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connection Success")
	}
	db.Close() //It closes the connection if there was a connection made
	fmt.Println("Connection Closed")
}

func testInsert(conf *cassandra.DbConfig, query string, value ...interface{}) {
	db, err := cassandra.New(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		err = db.Insert(query, value...)
		if err != nil {
			fmt.Printf("Insert Failed : %v", err)
		} else {
			fmt.Println("Insert Success")
		}
		db.Close()
	}
}

func testUpdate(conf *cassandra.DbConfig, query string, value ...interface{}) {
	db, err := cassandra.New(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		err = db.Update(query, value...)
		if err != nil {
			fmt.Printf("Update Failed : %v", err)
		} else {
			fmt.Println("Update Success")
		}
		db.Close()
	}
}

func testDelete(conf *cassandra.DbConfig, query string, value ...interface{}) {
	db, err := cassandra.New(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		err = db.Delete(query, value...)
		if err != nil {
			fmt.Printf("Delete Failed : %v", err)
		} else {
			fmt.Println("Delete Success")
		}
		db.Close()
	}
}

func testSelect(conf *cassandra.DbConfig, query string, value ...interface{}) {
	db, err := cassandra.New(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := db.Select(query, value...)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range data {
		for key, val := range value {
			fmt.Printf("[%v] : [%v]", key, val)
			fmt.Println()
		}
		fmt.Println("------------------")
	}
}
