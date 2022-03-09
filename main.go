package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type conf struct {
	domain string `yaml:"domain"`
	port   string `yaml:"port"`
	debug  bool   `yaml:"debug"`
	dbUser string `yaml:"db.User"`
	dbPass string `yaml:"db.Pass"`
	dbName string `yaml:"db.Name"`
	dbHost string `yaml:"db.Host"`
	dbPort string `yaml:"db.Port"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

var dbLink string = ""

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

	if sL.URL == "" {
		return "http://localhost:8888/"
	} else {
		return sL.URL
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index")
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}

type apiSingleLinkResponse struct {
	ShortLink string
	LongLink  string
}

func returnSingleLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	url := getShortLink(key)

	response := apiSingleLinkResponse{
		ShortLink: key,
		LongLink:  url,
	}
	json.NewEncoder(w).Encode(response)
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

var c conf

func main() {
	c.getConf()
	fmt.Println("Starting server...")
	fmt.Println("Domain: " + c.domain)
	fmt.Println("Port: " + c.port)
	fmt.Println("Debug: " + fmt.Sprintf("%t", c.debug))
	fmt.Println("DB User: " + c.dbUser)
	fmt.Println("DB Pass: " + c.dbPass)
	fmt.Println("DB Name: " + c.dbName)
	fmt.Println("DB Host: " + c.dbHost)
	fmt.Println("DB Port: " + c.dbPort)
	dbLink = c.dbUser + ":" + c.dbPass + "@tcp(" + c.dbHost + ":" + c.dbPort + ")/" + c.dbName
	handleRequests()
}
