package main

import (
	_ "github.com/bmizerany/pq"
	"database/sql"
	"log"
	"bytes"
	"fmt"
	"strconv"
	"go/types"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		return
	}
}

type MyDB struct {
	dbName   string
	username string
	password string
	ip		 string
	port     int
}

func newMyDB(db string, user string, pwd string, ip string, port int) *MyDB {
	return &MyDB{
		dbName: db,
		username: user,
		password: pwd,
		ip: ip,
		port: port,
	}
}

func (d *MyDB) newConn() *sql.DB {
	buf := bytes.Buffer{}
	buf.WriteString("user=")
	buf.WriteString(d.username)
	buf.WriteString(" password=")
	buf.WriteString(d.password)
	buf.WriteString(" dbname=")
	buf.WriteString(d.dbName)
	buf.WriteString(" host=")
	buf.WriteString(d.ip)
	buf.WriteString(" port=")
	buf.WriteString(strconv.Itoa(d.port))
	buf.WriteString(" sslmode=disable")

	s := buf.String()
	fmt.Println("value when open data connection:", s)
	c, err := sql.Open("postgres", s)
	checkError(err)
	return c
}

// need to fit changeable filed
func (d *MyDB) insertData(schema string, table string, data map[string]string, c *sql.DB) sql.Result {
	// combine the sql string
	// exp: INSERT INTO schema.table(field1, field2, ...) VALUES($1, $2, ...)
	var buf bytes.Buffer
	buf.WriteString("INSERT INTO ")
	buf.WriteString(schema)
	buf.WriteString(".")
	buf.WriteString(table)
	buf.WriteString("(")
	i := 0
	// #####################################################h
	for k,_ := range data {
		fmt.Println("i when range:", i)
		// valueList[i] = v
		buf.WriteString(k)
		if i != len(data) -1 {
			buf.WriteString(",")
		}
		i++
	}
	//fmt.Println("value list:", valueList)
	buf.WriteString(") ")
	buf.WriteString("VALUES(")
	for i=1;i<=len(data);i++ {
		buf.WriteString("$")
		buf.WriteString(strconv.Itoa(i))
		if i != len(data) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(") ")
	s := buf.String()
	fmt.Println("sql when insert:", s)

	stmt, err := c.Prepare(s)
	checkError(err)

	tmpBuff := bytes.Buffer{}
	j := 0
	for _, v := range data {
		tmpBuff.WriteString("\"")
		tmpBuff.WriteString(v)
		tmpBuff.WriteString("\"")
		if j != len(data)-1 {
			tmpBuff.WriteString(",")
		}
		j++
	}
	valueString := tmpBuff.String()
	fmt.Println("value when insert to database:", valueString)
	fmt.Printf("value when insert to database: %s %s %s", data["username"], data["departname"], data["created"])
	res, err := stmt.Exec(data["username"], data["departname"])
	checkError(err)
	return res
}

func (d *MyDB) updateData(schema string, table string, v map[string]string, c *sql.DB) sql.Result {
	stmt, err := c.Prepare("UPDATE" + schema + "." + table + "SET username=$1 WHERE uid=$2")
	checkError(err)
	res, err := stmt.Exec("jiamiao", 1)
	checkError(err)
	return res
}

func (d *MyDB) deleteData(scheme string, table string, v map[string]string, c *sql.DB) sql.Result {
	stmt, err := c.Prepare("DELETE FROM " + scheme +"."+table+" WHERE uid=$1")
	checkError(err)
	res, err := stmt.Exec(1)
	checkError(err)
	return res
}

func (d *MyDB) createTable(schema string, table string, c *sql.DB) sql.Result {
	stmt, err := c.Prepare("CREATE TABLE IF NOT EXISTS " + schema + "." + table+"(" +
		"uid SERIAL NOT NULL," +
		"username CHARACTER VARYING(100) NOT NULL," +
		"departname CHARACTER VARYING(500) NOT NULL," +
		"Created DATE," +
		"CONSTRAINT "+ table +"_pkey PRIMARY KEY (uid) )")
	checkError(err)
	res, err := stmt.Exec()
	checkError(err)
	return res
}

func (d *MyDB) createSchema(schema string, c *sql.DB) sql.Result {
	stmt, err := c.Prepare("CREATE SCHEMA IF NOT EXISTS " + schema)
	checkError(err)
	res, err := stmt.Exec()
	checkError(err)
	fmt.Println("create schema:", res)
	return res
}

func (d *MyDB) queryData(schema string, table string, c *sql.DB) *sql.Rows {
	rows, err := c.Query("SELECT * FROM " + schema + "." + table)
	checkError(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created *types.Nil
		err = rows.Scan(&uid, &username, &department, &created)
		checkError(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}
	return rows
}

