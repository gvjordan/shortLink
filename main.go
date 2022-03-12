package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

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

type sqlLink struct {
	ID        int
	Name      string
	URL       string
	CreatedAt int
	CreatedBy string
}

func getShortLink(id string) (string, error) {
	db, err := sql.Open("mysql", dbLink)
	if err != nil {
		return "", err
	}

	defer db.Close()
	res, err := db.Query("SELECT URL FROM links WHERE Name = ?", id)
	if err != nil {
		return "", err
	}

	defer res.Close()
	var sL sqlLink
	for res.Next() {
		err := res.Scan(&sL.URL)
		if err != nil {
			return "", err
		}
	}

	if sL.URL == "" {
		return "http://" + c.Domain + ":" + c.Port + "/", nil
	} else {
		return sL.URL, nil
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index")
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}

func statsPage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(sT)
}

type apiSingleLinkResponse struct {
	ShortLink string
	LongLink  string
	Success   bool
}

func returnSingleLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	url, err := getShortLink(key)
	if err != nil {
		errorObject := jsonError{
			Error: "Shortlink error",
		}
		json.NewEncoder(w).Encode(errorObject)
		return
	}

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

type apiAddLinkBody struct {
	URL   string `json:"URL"`
	Token string `json:"Token"`
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func generateRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

type jsonError struct {
	Error string
}

func addNewLink(w http.ResponseWriter, r *http.Request) {
	var sL apiAddLinkBody

	err := json.NewDecoder(r.Body).Decode(&sL)
	if err != nil {
		errorObject := jsonError{
			Error: "Invalid JSON",
		}
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	db, err := sql.Open("mysql", dbLink)
	if err != nil {
		errorObject := jsonError{
			Error: "Database error",
		}
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	defer db.Close()
	// TODO: check if link already exists
	shortKey := generateRandomString(5)

	stmt, err := db.Prepare("INSERT INTO links (Name, URL, CreatedAt, CreatedBy) VALUES (?, ?, ?, ?)")
	if err != nil {
		errorObject := jsonError{
			Error: "Database error",
		}
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	_, err = stmt.Exec(shortKey, sL.URL, int(time.Now().Unix()), r.RemoteAddr)
	if err != nil {
		errorObject := jsonError{
			Error: "Database error",
		}
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	response := apiSingleLinkResponse{
		ShortLink: shortKey,
		LongLink:  sL.URL,
		Success:   true,
	}

	json.NewEncoder(w).Encode(response)

}

func redirectFromShortLink(w http.ResponseWriter, r *http.Request, key string) {
	url, err := getShortLink(key)
	if err != nil {
		errorObject := jsonError{
			Error: "Shortlink error",
		}
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func apiAddHandler(w http.ResponseWriter, r *http.Request) {
	router := strings.Split(r.URL.Path, "/")
	switch router[2] {
	case "add":
		addNewLink(w, r)
	default:
		fmt.Fprintf(w, "API")
	}

}

func tokenPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Token")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	router := strings.Split(r.URL.Path, "/")
	if router[1] == "" {
		indexPage(w, r)
	} else if len(router[1]) == 5 {
		redirectFromShortLink(w, r, router[1])

	}
}

func handleRequests() {

	http.Handle("/", http.HandlerFunc(indexHandler))
	http.Handle("/api", http.HandlerFunc(apiPage))
	http.Handle("/api/add", rateLimitMiddleware(http.HandlerFunc(apiAddHandler)))
	http.Handle("/tokens", http.HandlerFunc(tokenPage))
	http.Handle("/stats", http.HandlerFunc(statsPage))
	log.Fatal(http.ListenAndServe(":"+c.Port, nil))
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rL.inc(ip)
		go rL.expire(ip, 5)

		sT.Requests++

		next.ServeHTTP(w, r)
	})
}

type stats struct {
	Requests        int
	LinksShortened  int
	LinksRedirected int
}

var (
	c      *Conf
	dbLink string
)

var rL = newRateLimit()

var sT = stats{
	Requests:        0,
	LinksShortened:  0,
	LinksRedirected: 0,
}

func main() {
	c = getConf()
	fmt.Println(c)
	dbLink = c.DbUser + ":" + c.DbPassword + "@tcp(" + c.DbHost + ":" + c.DbPort + ")/" + c.DbName
	fmt.Println("Starting server...")
	handleRequests()
}
