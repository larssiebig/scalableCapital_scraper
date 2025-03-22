package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
    // Initialize a Colly collector
    c := colly.NewCollector(
        colly.AllowedDomains("de.scalable.capital"),
    )

    // Handle authentication if necessary
    // Use c.OnRequest and mimic login if the site has a login form or cookies.

    // Extract portfolio table and details
    c.OnHTML("div[role='table']", func(e *colly.HTMLElement) {
        // Extract Name
        name := e.ChildText(".css-1b1ehxh-text")
        fmt.Printf("Name: %s\n", name)

        // Extract Total Value
        totalValue := e.ChildText("div[aria-label='Total value'] span")
        fmt.Printf("Total Value: %s\n", totalValue)

        // Extract Savings Plan Indicator
        savingsPlan := e.ChildText("div[aria-label='Savings plan'] span")
        fmt.Printf("Savings Plan: %s\n", savingsPlan)

        // Extract Price
        price := e.ChildText("div[aria-label*='Mid price for'] div.MuiGrid-root")
        fmt.Printf("Price: %s\n", price)
    })

    // Error handling
    c.OnError(func(r *colly.Response, err error) {
        log.Printf("Request failed with status %d: %s", r.StatusCode, err)
    })

    // Visit the target webpage
    err := c.Visit("https://de.scalable.capital/en/secure-login")
    if err != nil {
        log.Fatalf("Failed to visit webpage: %v", err)
    }
}
