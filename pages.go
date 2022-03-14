package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index")
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}

func statsPage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(sT)
}

func tokenPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Token")
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

type uiConfigData struct {
	Host string
	Port string
}
