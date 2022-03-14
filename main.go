package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	c               *Conf
	l               *logger
	dbLink          string
	limitByIP       int
	limitByIPExpire int
	sT              = stats{
		Requests:        0,
		FailedRequests:  0,
		SuccessRequests: 0,
		LinksShortened:  0,
		LinksRedirected: 0,
		DatabaseError:   0,
		ResolveError:    0,
		InvalidJSON:     0,
		InvalidToken:    0,
		AccessDenied:    0,
	}
	allowedTokens = map[string]bool{
		"52fdfc07-2182-654f-163f-5f0f9a621d72": true,
	}
	rL = newRateLimit()
)

func handleGenerateUUID() {
	fmt.Println(generateUUID())
}

func main() {
	configData, err := newConf()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %s\n", err)
	}
	c = configData
	loadConf()
	handleFlags()
	fmt.Println(c)
	fmt.Println("Starting server...")
	handleRequests()
}
