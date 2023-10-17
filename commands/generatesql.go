package main

import (
	"olx-clone/functions/logger"
	"os"
	"regexp"
	"strings"
	"time"
)

var log = logger.Log

// GenerateSQLFile Generates a file in migrations/scripts/ directory in required migration format for a given tableName
func GenerateSQLFile(tableName string) {
	var sb strings.Builder
	timeString := time.Now().Format("20060102150405.003059_")
	regexString := regexp.MustCompile(`^(.*?)\.(.*)$`)
	replaceString := "${1}$2"
	sb.WriteString("migrations/scripts/")
	sb.WriteString(regexString.ReplaceAllString(timeString, replaceString))
	sb.WriteString(tableName)
	sb.WriteString(".sql")
	fileName := sb.String()
	emptyFile, err := os.Create(fileName)
	if err != nil {
		log.Errorln(err)
	}
	log.Println("Created SQL File:", fileName)
	emptyFile.Close()
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Please specify a table name as first argument")
	}
	fileName := os.Args[1]
	GenerateSQLFile(fileName)
}
