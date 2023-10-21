/*
All migration SQL files reside in scripts folder with following format:

YYYYMMDDHHMMSSmillseconds_tablename.sql
*/
package migrations

import (
	"io/ioutil"
	"olx-clone/functions/logger"
	"olx-clone/infra/db"

	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var log = logger.Log

var fileNameRegex *regexp.Regexp = regexp.MustCompile(`scripts/\d{20}_\S+.sql`)

// checks if sql file name is in correct format
func checkFileNameFormat(fileName string) bool {
	matched := fileNameRegex.MatchString(fileName)
	return matched
}

// reads the metadata from migrations_metadata table and return last executed script name
func readMetadata() string {
	database := db.Mgr.DBConn

	// create migrations table if not exists
	_, err := database.Exec("CREATE TABLE IF NOT EXISTS migrations_metadata (migrated_at TIMESTAMP, script_name TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	// read last executed script name
	var lastScriptName string
	err = database.Get(&lastScriptName, "SELECT script_name from migrations_metadata ORDER BY migrated_at DESC LIMIT 1")
	if err != nil {
		return ""
	}
	return lastScriptName
}

// writeMetadata writes a new row in the migrations_metadata table to record our action
func writeMetadata(scriptName string) bool {
	sql := "INSERT INTO migrations_metadata values (NOW(), $1);"
	database := db.Mgr.DBConn
	_, err := database.Exec(sql, scriptName)
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

// MigrateDB finds the last run migration, and run all those after it in order
func MigrateDB() {

	database := db.Mgr.DBConn
	
	// Create a migration table if it doesn't exist.
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id serial PRIMARY KEY,
			version INT,
			description TEXT
		);
	`)

	if err != nil {
		log.Errorln("Error while pinging the DB", err)
		panic(err)
	}

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id serial PRIMARY KEY,
			username VARCHAR (50),
			email VARCHAR (100)
		);
	`)
	if err != nil {
		logger.Log.Panicln("Error while migration: ", err)
	}

	// Update the migrations table to track the applied migration.
	// _, err = database.Exec(`
	// 	INSERT INTO migrations (version, description) VALUES (1, 'Creat Users Table');
	// `)
	// if err != nil {
	// 	logger.Log.Panicln("Error while migration: ", err)
	// }

}

// MigrateDB finds the last run migration, and run all those after it in order
func MigrateDB2() {
	// Get a list of migration files
	files, err := filepath.Glob("scripts/*.sql")
	if err != nil {
		log.Printf("Error running restore %s", err)
		return
	}

	// Sort the list alphabetically
	sort.Strings(files)
	log.Println("files:", files)

	// get last run migration
	log.Println("Reading from Metadata table...")
	lastScriptName := readMetadata()
	log.Println("Last migrated script:", lastScriptName)

	// get database object
	database := db.Mgr.DBConn

	var lastCompleted string
	completedCount := 0

	for _, file := range files {

		// check file name format
		if !checkFileNameFormat(file) {
			log.Println("Invalid file name format for file:", file)
			break
		}

		// if no migrations were made or the migration file is newer than last migrated file
		if lastScriptName == "" || strings.Compare(file, lastScriptName) > 0 {
			log.Println("Running migration:", file)

			// reading contents of SQL file
			content, _ := ioutil.ReadFile(file)
			// Convert []byte to string
			sqlQueries := string(content)

			// Execute queries in a transaction If at any point we fail, rollback it and break
			tx, _ := database.Begin()
			_, err = tx.Exec(sqlQueries)
			if err != nil {
				log.Println(sqlQueries)
				log.Errorln(err)
				tx.Rollback()
				break
			}
			tx.Commit()

			lastCompleted = file
			completedCount += 1

			log.Println("Completed migration:", file)
			writeMetadata(file)
		}
	}

	if completedCount > 0 {
		log.Println(completedCount, "Migrations completed. Last completed:", lastCompleted)
	} else {
		log.Println("No migrations performed")
	}
}
