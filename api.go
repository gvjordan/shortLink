package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		handleStats("ResolveError")
		sendErrorJSON(w, r, "Invalid shortLink")
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
		handleStats("InvalidJSON")
		json.NewEncoder(w).Encode(errorObject)
		return
	}
	fmt.Println(sL.Token)
	if sL.Token == "" || !checkToken(sL.Token) {
		handleStats("InvalidToken")
		sendErrorJSON(w, r, "Invalid Token")
		return
	}

	db, err := sql.Open("mysql", dbLink)
	if err != nil {
		handleStats("DatabaseError")
		sendErrorJSON(w, r, "Database error")
		return
	}

	defer db.Close()
	// TODO: check if link already exists
	shortKey := generateRandomString(5)

	stmt, err := db.Prepare("INSERT INTO links (Name, URL, CreatedAt, CreatedBy) VALUES (?, ?, ?, ?)")
	if err != nil {
		handleStats("DatabaseError")
		sendErrorJSON(w, r, "Database error")
		return
	}

	_, err = stmt.Exec(shortKey, sL.URL, int(time.Now().Unix()), r.RemoteAddr)
	if err != nil {
		handleStats("DatabaseError")
		sendErrorJSON(w, r, "Database error")
		return
	}

	response := apiSingleLinkResponse{
		ShortLink: shortKey,
		LongLink:  sL.URL,
		Success:   true,
	}

	handleStats("LinksShortened")

	json.NewEncoder(w).Encode(response)

}

func redirectFromShortLink(w http.ResponseWriter, r *http.Request, key string) {
	url, err := getShortLink(key)
	if err != nil {
		handleStats("ResolveError")
		sendErrorJSON(w, r, "Invalid shortLink")
		return
	}
	handleStats("LinksRedirected")
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
