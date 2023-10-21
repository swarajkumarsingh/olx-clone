package db

import (
	"fmt"
	"olx-clone/constants"
	"os"
)

const (
	ENV_PROD  = constants.ENV_PROD
	ENV_UAT   = constants.ENV_UAT
	ENV_DEV   = constants.ENV_DEV
	ENV_LOCAL = constants.ENV_LOCAL
)

/*
DB Configurations
*/
type DBConfigStruct struct {
	ReadOnlyHost string
	Host         string
	User         string
	Password     string
	Port         string
	DBName       string
}

var ENV string = os.Getenv("STAGE")

const (
	ReadOnly = "READ_ONLY"
)

func getDBConnectionString(_ string) string {
	dbConfig := &DBConfigStruct{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	sslConfig := ""
	if ENV == ENV_PROD {
		sslConfig = "sslmode=require"
	} else {
		sslConfig = "sslmode=disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, sslConfig)
}

var DBConfig = map[string]string{
	"driver":                "postgres",
	"conn_string":           getDBConnectionString(""),
	"read_only_conn_string": getDBConnectionString(ReadOnly),
}
