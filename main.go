package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Conf struct {
	Domain        string
	Port          string
	Debug         bool
	DbUser        string
	DbPassword    string
	DbName        string
	DbHost        string
	DbPort        string
	AllowedTokens []string
	Stats         bool
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
		sT.ResolveError++
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
		sT.InvalidJSON++
		json.NewEncoder(w).Encode(errorObject)
		return
	}
	fmt.Println(sL.Token)
	if sL.Token == "" || !checkToken(sL.Token) {
		errorObject := jsonError{
			Error: "Invalid token",
		}
		sT.InvalidToken++
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	db, err := sql.Open("mysql", dbLink)
	if err != nil {
		errorObject := jsonError{
			Error: "Database error",
		}
		sT.DatabaseError++
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
		sT.DatabaseError++
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	_, err = stmt.Exec(shortKey, sL.URL, int(time.Now().Unix()), r.RemoteAddr)
	if err != nil {
		errorObject := jsonError{
			Error: "Database error",
		}
		sT.DatabaseError++
		json.NewEncoder(w).Encode(errorObject)
		return
	}

	response := apiSingleLinkResponse{
		ShortLink: shortKey,
		LongLink:  sL.URL,
		Success:   true,
	}

	sT.LinksShortened++

	json.NewEncoder(w).Encode(response)

}

func redirectFromShortLink(w http.ResponseWriter, r *http.Request, key string) {
	url, err := getShortLink(key)
	if err != nil {
		errorObject := jsonError{
			Error: "Shortlink error",
		}
		sT.ResolveError++
		json.NewEncoder(w).Encode(errorObject)
		return
	}
	sT.LinksRedirected++
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

	} else if len(router[1]) == 6 && router[1][0] == '+' {
		returnSingleLink(w, r)
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

func denyAccess(w http.ResponseWriter, r *http.Request) {
	errorObject := jsonError{
		Error: "Access denied",
	}
	sT.AccessDenied++
	json.NewEncoder(w).Encode(errorObject)
}

var maxRequests int = 5

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		sT.Requests++

		if rL.current[ip] == nil {
			rL.inc(ip)
		} else if rL.current[ip].Count >= maxRequests {
			sT.FailedRequests++
			denyAccess(w, r)
			return
		} else {
			rL.inc(ip)
		}
		sT.SuccessRequests++
		go rL.expire(ip, 5)

		next.ServeHTTP(w, r)
	})
}

type stats struct {
	Requests        int
	FailedRequests  int
	SuccessRequests int
	LinksShortened  int
	LinksRedirected int
	DatabaseError   int
	ResolveError    int
	InvalidJSON     int
	InvalidToken    int
	AccessDenied    int
}

var (
	c      *Conf
	dbLink string
)

var rL = newRateLimit()

var allowedTokens = map[string]bool{
	"52fdfc07-2182-654f-163f-5f0f9a621d72": true,
}

func checkToken(token string) bool {
	if _, ok := allowedTokens[token]; ok {
		if allowedTokens[token] {
			return true
		}
	}
	return false
}

var sT = stats{
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

func parseTokens(tokens []string) {
	for _, token := range tokens {
		allowedTokens[token] = true
	}
}

func handleStats(stat int) {
	if c.Stats {
		stat++
	}
}

func handleGenerateUUID() {
	fmt.Println(generateUUID())
}

func loadConf() {
	c = getConf()
	parseTokens(c.AllowedTokens)
	dbLink = c.DbUser + ":" + c.DbPassword + "@tcp(" + c.DbHost + ":" + c.DbPort + ")/" + c.DbName
}

func handleFlags() {
	debugFlag := flag.Bool("debug", false, "Enable stats")
	tokenFlag := flag.Bool("token", false, "Generate token")

	flag.Parse()
	if *debugFlag {
		c.Debug = true
		return
	} else if *tokenFlag {
		handleGenerateUUID()
		os.Exit(0)
	} else {
		c.Debug = false
		return
	}
}

func main() {
	loadConf()
	handleFlags()
	fmt.Println("Starting server...")
	handleRequests()
}
