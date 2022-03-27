package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
)

func handleImportSql() {
	fmt.Println("Validating database config options")
	if c.DbHost == "" {
		fmt.Println("Error: DbHost not set")
		os.Exit(1)
	}
	if c.DbPort == "" {
		fmt.Println("Error: DbPort not set")
		os.Exit(1)
	}
	if c.DbUser == "" {
		fmt.Println("Error: DbUser not set")
		os.Exit(1)
	}
	if c.DbPassword == "" {
		fmt.Println("Error: DbPassword not set")
		os.Exit(1)
	}
	if c.DbName == "" {
		fmt.Println("Error: DbName not set")
		os.Exit(1)
	}
	if dbLink == "" {
		dbLink = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
	}

	sqlData, err := ioutil.ReadFile("sql/tables.sql")
	if err != nil {
		fmt.Println("Error: Unable to read sql/tables.sql")
		os.Exit(1)
	}

	sqlStatement := string(sqlData)
	db, err := sql.Open("mysql", dbLink)
	if err != nil {
		fmt.Println("Error: Unable to connect to database")
		os.Exit(1)
	}
	_, err = db.Exec(sqlStatement)
	if err != nil {
		fmt.Println("Error: Unable to execute sql/tables.sql")
		os.Exit(1)
	}

	fmt.Println("Success: sql/tables.sql executed")
}

func addNewToken(token string) {

	if token == "" || token == "0" {
		token = generateUUID()
	}

	sqlStatement := "INSERT INTO tokens (token) VALUES (?)"
	db, err := sql.Open("mysql", dbLink)
	if err != nil {
		fmt.Println("Error: Unable to connect to database")
		os.Exit(1)
	}
	defer db.Close()
	_, err = db.Exec(sqlStatement, token)
	if err != nil {
		fmt.Println("Error: Unable to add token")
		os.Exit(1)
	}
	fmt.Println("Success: Token added")
}
