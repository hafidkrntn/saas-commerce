package main

import "os"

// MockProviderBaseURL returns the redirect base URL used by the mock payment
// provider (the storefront URL in development).
func MockProviderBaseURL() string {
	url := os.Getenv("STORE_FRONTEND_URL")
	if url == "" {
		url = "http://localhost:3000"
	}
	return url
}
