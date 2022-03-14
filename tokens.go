package main

func checkToken(token string) bool {
	if _, ok := allowedTokens[token]; ok {
		if allowedTokens[token] {
			return true
		}
	}
	return false
}

func parseTokens(tokens []string) {
	for _, token := range tokens {
		allowedTokens[token] = true
	}
}
