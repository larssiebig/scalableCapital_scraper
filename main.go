package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

const (
    port = 9515
)

func main() {
    // Set up WebDriver capabilities
    caps := selenium.Capabilities{"browserName": "chrome"}
    chromeOpts := map[string]interface{}{
        "args": []string{
            "--headless", // Run in headless mode
            "--disable-dev-shm-usage",
            "--no-sandbox",
        },
    }
    caps["goog:chromeOptions"] = chromeOpts

    // Start WebDriver
    driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
    if err != nil {
        log.Fatalf("failed to start WebDriver: %v", err)
    }
    defer driver.Quit()

    // Open target webpage
    err = driver.Get("https://de.scalable.capital/en/secure-login")
    if err != nil {
        log.Fatalf("failed to load webpage: %v", err)
    }

    // Add delay for page load
    time.Sleep(60 * time.Second)

    // Locate the parent element of the portfolio details
    parentElement, err := driver.FindElement(selenium.ByXPATH, "//div[@role='table']")
    if err != nil {
        log.Fatalf("failed to find portfolio table: %v", err)
    }

    // Extract portfolio details
    extractPortfolioDetails(parentElement)
}

func extractPortfolioDetails(parent selenium.WebElement) {
    // Extract Name
    nameElement, err := parent.FindElement(selenium.ByXPATH, ".//div[@class='css-1b1ehxh-text']")
    if err != nil {
        log.Printf("name not found: %v", err)
    } else {
        name, _ := nameElement.Text()
        fmt.Printf("Name: %s\n", name)
    }

    // Extract Total Value
    valueElement, err := parent.FindElement(selenium.ByXPATH, ".//div[@aria-label='Total value']//span")
    if err != nil {
        log.Printf("total value not found: %v", err)
    } else {
        totalValue, _ := valueElement.Text()
        fmt.Printf("Total Value: %s\n", totalValue)
    }

    // Extract Savings Plan Indicator
    savingsPlanElement, err := parent.FindElement(selenium.ByXPATH, ".//div[@aria-label='Savings plan']//span")
    if err != nil {
        log.Printf("savings plan not found: %v", err)
    } else {
        savingsPlan, _ := savingsPlanElement.Text()
        fmt.Printf("Savings Plan: %s\n", savingsPlan)
    }

    // Extract Price
    priceElement, err := parent.FindElement(selenium.ByXPATH, ".//div[@aria-label='Mid price for iShares Core S&P 500 (Acc) (ticking)']//div[@class='MuiGrid-root']")
    if err != nil {
        log.Printf("price not found: %v", err)
    } else {
        price, _ := priceElement.Text()
        fmt.Printf("Price: %s\n", price)
    }
}
