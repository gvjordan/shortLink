package main

import (
	"log"
	"net/http"
)

func handleRequests() {
	http.Handle("/", http.HandlerFunc(indexHandler))
	http.Handle("/api", http.HandlerFunc(apiPage))
	http.Handle("/api/add", rateLimitMiddleware(http.HandlerFunc(apiAddHandler)))
	http.Handle("/tokens", http.HandlerFunc(tokenPage))

	if c.Stats {
		http.Handle("/stats", http.HandlerFunc(statsPage))
	}

	if c.EnableFrontend {
		http.Handle("/ui", http.HandlerFunc(uipage))
	}

	log.Fatal(http.ListenAndServe(":"+c.Port, nil))
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		sT.Requests++

		if rL.current[ip] == nil {
			rL.inc(ip)
		} else if rL.current[ip].Count >= limitByIP {
			sT.FailedRequests++
			denyAccess(w, r)
			return
		} else {
			rL.inc(ip)
		}
		sT.SuccessRequests++
		go rL.expire(ip, limitByIPExpire)

		next.ServeHTTP(w, r)
	})
}
