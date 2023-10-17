/*
To do any of the things below, first import the package

import "olx-clone/infra/db"

--------------------------------------------------------
Run a query and map results to array of struct instances
--------------------------------------------------------
Make sure struct columns are same as column names of table.
You can also skip the parameters following the query in the query below, and add multiple.

database := db.GetDB()
students := []Student{}
err := database.Select(&students, "SELECT id, firstname, lastname FROM student WHERE id > $1", 50)
if err != nil {
	log.Println(err)
}
for _, student := range students{
	log.Println(student)
}

-----------------------------------------------------------
Run a query and map single row result to a struct instance
-----------------------------------------------------------
database := db.GetDB()
student := Student{}
err := database.Get(&student, "SELECT id, firstname, lastname FROM student WHERE id = $1", 50)
if err != nil {
	log.Println(err)
}
log.Println(student)

--------------------------------------------
Execute an insert/update/create table query
--------------------------------------------
database := db.GetDB()
_, err := database.Exec("INSERT INTO student (firstname, lastname) VALUES ($1, $2)", "John", "Doe")
if err != nil {
	log.Println(err)
}

-----------------------------------------------------------------------
Run a query with IN operator and map rows to array of struct instances
-----------------------------------------------------------------------
database := db.GetDB()
var ids = []int{2, 3}
query, args, err := sqlx.In("SELECT column_name1, column_name2, column_name3 FROM student WHERE id IN (?);", ids)
if err != nil {
	log.Println(err)
}
query = database.Rebind(query) // this will convert ? to respective database bindvars, example: $1, $2, etc supported by Postgresql
students := []Student{}
err = database.Select(&students, query, args...)
if err != nil {
	log.Println(err)
}
for _, student := range students{
	log.Println(student)
}
*/

// Package db provides database connection
package db

import (
	"olx-clone/conf"
	"olx-clone/functions/logger"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

// manager structs hold the DB conn object.
type manager struct {
	DBConn *sqlx.DB
}

var Mgr manager
var log = logger.Log

// init initiates a datadog instance.
func init() {
	var err error
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName(conf.ClientENV))
	database, err := sqlxtrace.Open(DBConfig["driver"], DBConfig["conn_string"])
	if err != nil {
		log.Errorln(err)
		panic(err)
	}
	maxOpenConn := 50
	if ENV == ENV_PROD {
		maxOpenConn = 800
	}
	database.SetMaxOpenConns(maxOpenConn)
	database.SetMaxIdleConns(50)
	database.SetConnMaxIdleTime(2 * time.Minute)
	database.SetConnMaxLifetime(5 * time.Minute)
	Mgr = manager{
		DBConn: database,
	}
}
