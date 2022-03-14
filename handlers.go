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

	} else if len(router[1]) == 6 && router[1][0] == '+' {
		returnSingleLink(w, r)
	}
}

func denyAccess(w http.ResponseWriter, r *http.Request) {
	errorObject := jsonError{
		Error: "Rate Limit Exceeded",
	}
	sT.AccessDenied++
	json.NewEncoder(w).Encode(errorObject)
}
