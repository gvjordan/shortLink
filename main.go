package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Conf struct {
	Domain     string
	Port       string
	Debug      bool
	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	DbPort     string
}

func getConf() *Conf {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	var config = &Conf{}
	err = viper.Unmarshal(config)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return config
}

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
		return "http://" + c.Domain + ":" + c.Port + "/"
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
	Success   bool
}

func returnSingleLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	url := getShortLink(key)

	response := apiSingleLinkResponse{
		ShortLink: key,
		LongLink:  url,
		Success:   true,
	}

	if response.LongLink == "http://"+c.Domain+":"+c.Port+"/" {
		response.Success = false
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
	router.HandleFunc("/+{id}", returnSingleLink)
	router.HandleFunc("/{id}", redirectFromShortLink)
	router.HandleFunc("/api", apiPage)
	router.HandleFunc("/api/get/{id}", returnSingleLink)

	log.Fatal(http.ListenAndServe(":"+c.Port, router))
}

var (
	c      *Conf
	dbLink string
)

func init() {
	c = getConf()
	fmt.Println(c)
	dbLink = c.DbUser + ":" + c.DbPassword + "@tcp(" + c.DbHost + ":" + c.DbPort + ")/" + c.DbName
	fmt.Println(dbLink)
}

func main() {
	fmt.Println("Starting server...")
	fmt.Println("Domain: " + c.Domain)
	fmt.Println("Port: " + c.Port)
	fmt.Println("Debug: " + fmt.Sprintf("%t", c.Debug))
	fmt.Println("DB User: " + c.DbUser)
	fmt.Println("DB Pass: " + c.DbPassword)
	fmt.Println("DB Name: " + c.DbName)
	fmt.Println("DB Host: " + c.DbHost)
	fmt.Println("DB Port: " + c.DbPort)
	handleRequests()
}
