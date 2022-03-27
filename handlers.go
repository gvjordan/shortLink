package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func apiAddHandler(w http.ResponseWriter, r *http.Request) {
	router := strings.Split(r.URL.Path, "/")
	switch router[2] {
	case "add":
		addNewLink(w, r)
	default:
		fmt.Fprintf(w, "API")
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	router := strings.Split(r.URL.Path, "/")
	if router[1] == "" {
		indexPage(w, r)
	} else if len(router[1]) == 5 {
		redirectFromShortLink(w, r, router[1])
		handleStats("SuccessRequests")
	} else if len(router[1]) == 6 && router[1][0] == '+' {
		returnSingleLink(w, r)
	}
}

func sendErrorJSON(w http.ResponseWriter, r *http.Request, message string) {
	errorObject := jsonError{
		Error: message,
	}
	json.NewEncoder(w).Encode(errorObject)
}
