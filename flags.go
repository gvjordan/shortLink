package main

import (
	"flag"
	"os"
)

func handleFlags() {
	debugFlag := flag.Bool("debug", false, "Enable debug")
	statsFlag := flag.Bool("stats", false, "Enable stats")
	tokenFlag := flag.Bool("token", false, "Generate token")
	importSqlFlag := flag.Bool("import-sql", false, "Import sql from sql/tables.sql using database info from config file")
	flag.Parse()

	if *tokenFlag {
		handleGenerateUUID()
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
