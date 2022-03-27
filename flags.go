package main

import (
	"flag"
	"os"
)

func handleFlags() {
	debugFlag := flag.Bool("debug", false, "Enable debug")
	statsFlag := flag.Bool("stats", false, "Enable stats")
	tokenFlag := flag.Bool("token", false, "Generate a new token to be manually added into the database")
	tokenAddFlag := flag.Bool("token-add", false, "Generate a new token and add it to the database")
	importSqlFlag := flag.Bool("import-sql", false, "Import sql from sql/tables.sql using database info from config file")
	flag.Parse()

	if *tokenFlag {
		handleGenerateUUID()
		os.Exit(0)
	}

	if *tokenAddFlag {
		addNewToken("")
		os.Exit(0)
	}

	if *importSqlFlag {
		handleImportSql()
		os.Exit(0)
	}

	if *debugFlag {
		c.Debug = true
	}

	if *statsFlag {
		c.Stats = true
	}

}
