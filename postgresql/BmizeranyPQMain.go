package main

import (
)

func main() {
	db := newMyDB("dbname", "username", "password", "host", 5432)
	c := db.newConn()
	defer c.Close()
	db.createSchema("go_test", c )
	db.createTable("go_test", "userinfo", c)
	data_map := map[string]string {
		"username": "yuexiong",
		"departname": "haoli",
		// "created": "2018-03-25",
		// time.Now().Format("2006-01-21")
	}
	db.insertData("go_test","userinfo", data_map, c)
	db.queryData("go_test", "userinfo", c)
}
