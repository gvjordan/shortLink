package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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

type uiConfigData struct {
	Host string
	Port string
}

func uipage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("www/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	data := uiConfigData{
		Host: c.Domain,
		Port: c.Port,
	}
	tmpl.Execute(w, data)
}

func handleRequests() {
	fmt.Println("test")
	http.Handle("/", http.HandlerFunc(indexHandler))
	http.Handle("/api", http.HandlerFunc(apiPage))
	http.Handle("/api/add", rateLimitMiddleware(http.HandlerFunc(apiAddHandler)))
	http.Handle("/tokens", http.HandlerFunc(tokenPage))
	http.Handle("/stats", http.HandlerFunc(statsPage))
	http.Handle("/ui", http.HandlerFunc(uipage))
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

var (
	c      *Conf
	l      *logger
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
