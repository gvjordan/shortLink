package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var dbLink string = "nuntius:nuntius@tcp(localhost:3306)/nuntius"

type sLink struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Success bool   `json:"success"`
}
type sqlLink struct {
	ID        int
	Name      string
	URL       string
	CreatedAt int
	CreatedBy string
}

func errHandler(err error, errorType string) {
	if err != nil {
		switch errorType {
		case "db":
			panic(err.Error())
		case "api":
			panic(err.Error())
		default:
			panic(err.Error())
		}
	}
}

func exHandleDb() {
	db, err := sql.Open("mysql", dbLink)
	errHandler(err, "db")
	defer db.Close()

	res, err := db.Query("Select * FROM links")
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()

	for res.Next() {
		var sL sqlLink
		err := res.Scan(&sL.ID, &sL.Name, &sL.URL, &sL.CreatedAt, &sL.CreatedBy)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(sL)
	}

}

func getShortLink(id string) string {
	db, err := sql.Open("mysql", dbLink)
	errHandler(err, "db")
	defer db.Close()
	res, err := db.Query("SELECT URL FROM links WHERE Name = ?", id)
	errHandler(err, "db")
	defer res.Close()
	var sL sqlLink
	for res.Next() {
		err := res.Scan(&sL.URL)
		errHandler(err, "db")
	}
	return sL.URL
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index")
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}

func returnSingleLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Fprintf(w, "Key: "+key)
	fmt.Fprintf(w, "URL: "+getShortLink(key))
}

func redirectFromShortLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println(key)
	http.Redirect(w, r, getShortLink(key), http.StatusMovedPermanently)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexPage)
	router.HandleFunc("/{id}", redirectFromShortLink)
	router.HandleFunc("/api", apiPage)
	router.HandleFunc("/api/get/{id}", returnSingleLink)

	log.Fatal(http.ListenAndServe(":8888", router))
}

func main() {
	handleRequests()
}
