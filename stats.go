package main

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

func handleStats(stat int) {
	if c.Stats {
		stat++
	}
}
