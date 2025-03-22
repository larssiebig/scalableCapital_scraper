package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create cookie jar for session persistence
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	// Initialize Colly collector with cookie support
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.SetCookieJar(cookieJar)

	// Mimic browser behavior
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		r.Headers.Set("Referer", "https://secure.scalable.capital/login")
	})

	// Perform login to fetch session cookies
	loginURL := "https://secure.scalable.capital/u/login"
	err = c.Post(loginURL, map[string]string{
		"username": os.Getenv("SCALABLE_LOGIN"),
		"password": os.Getenv("SCALABLE_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Check if login was successful by looking for session cookies
	sessionCookies := cookieJar.Cookies(&url.URL{Scheme: "https", Host: "secure.scalable.capital"})
	if len(sessionCookies) == 0 {
		log.Fatal("Login failed: No session cookies found")
	} else {
		log.Println("Login successful, session cookies stored!")
	}

	// New API endpoint for the request
	portfolioURL := "https://de.scalable.capital/broker/api/data"

	// Create the GET request with the new API URL
	req, err := http.NewRequest("GET", portfolioURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Attach cookies to the request to authenticate
	for _, cookie := range sessionCookies {
		req.AddCookie(cookie)
	}

	// Add additional headers if needed
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Referer", "https://secure.scalable.capital")

	// Perform the GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the Content-Type header
	contentType := resp.Header.Get("Content-Type")
	log.Printf("Response Content-Type: %s", contentType)

	// Read the response body to log it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Log the response body for inspection
	log.Printf("Response Body: %s", string(body))

	// Try to decode the JSON response
	if contentType == "application/json" {
		var responseData map[string]any
		if err := json.Unmarshal(body, &responseData); err != nil {
			log.Fatalf("Failed to decode JSON response: %v", err)
		}

		// Print the decoded response
		fmt.Println("Response Data:")
		fmt.Printf("%+v\n", responseData)
	} else {
		log.Println("Response is not in JSON format, check the body for errors.")
	}
}
