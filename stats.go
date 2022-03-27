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

func handleStats(thisStat string) {
	if c.Stats {
		switch thisStat {
		case "Requests":
			sT.Requests++
		case "FailedRequests":
			sT.FailedRequests++
		case "SuccessRequests":
			sT.SuccessRequests++
		case "LinksShortened":
			sT.LinksShortened++
		case "LinksRedirected":
			sT.LinksRedirected++
		case "DatabaseError":
			sT.DatabaseError++
		case "ResolveError":
			sT.ResolveError++
		case "InvalidJSON":
			sT.InvalidJSON++
		case "InvalidToken":
			sT.InvalidToken++
		case "AccessDenied":
			sT.AccessDenied++
		}
	}
}
