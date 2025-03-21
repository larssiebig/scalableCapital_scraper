package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
    port = 9515
)

func main() {
    // Chrome capabilities with updated settings
    caps := selenium.Capabilities{"browserName": "chrome"}
    chromeOpts := map[string]interface{}{
        "args": []string{
            "--disable-blink-features=AutomationControlled",
            "--disable-search-engine-choice-screen",
            "--headless",
            "--disable-dev-shm-usage",
            "--no-sandbox",
            "--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
        },
    }
    caps["goog:chromeOptions"] = chromeOpts

    // Create a new remote session (Ensure ChromeDriver is running first)
    driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
    if err != nil {
        log.Fatalf("Failed to start the browser: %v", err)
    }
    defer driver.Quit()

    // Open Scalable Capital Login Page
    err = driver.Get("https://de.scalable.capital/en/secure-login")
    if err != nil {
        log.Fatalf("Failed to load page: %v", err)
    }

    // Handle Cookies in Shadow DOM
    handleCookies(driver)

    // Login
    login := os.Getenv("SCALABLE_LOGIN")
    password := os.Getenv("SCALABLE_PASSWORD")
    err = performLogin(driver, login, password)
    if err != nil {
        log.Printf("Login failed: %v", err)
    } else {
        log.Println("Login successful!")
    }

    // Extract Portfolio Data
    portfolio, err := extractPortfolioData(driver)
    if err != nil {
        log.Fatalf("Error extracting portfolio data: %v", err)
    }

    // Print results
    fmt.Println("Extracted Portfolio Data:")
    for _, asset := range portfolio {
        fmt.Printf("Asset: %s | ISIN: %s | Shares: %s | Value: %s\n", asset["name"], asset["isin"], asset["shares"], asset["value"])
    }
}

// Handle Shadow DOM for Cookie Consent
func handleCookies(driver selenium.WebDriver) {
    driver.ExecuteScript(`
        try {
            var shadow = document.querySelector("#usercentrics-root").shadowRoot;
            var button = shadow.querySelector("button[data-testid='uc-deny-all-button']");
            if (button) { button.click(); }
        } catch (e) {}
    `, nil)
    time.Sleep(2 * time.Second)
}

// Perform Login
func performLogin(driver selenium.WebDriver, login, password string) error {
    usernameField, err := waitForElement(driver, selenium.ByID, "username", 10*time.Second)
    if err != nil {
        return err
    }
    usernameField.SendKeys(login)

    passwordField, err := waitForElement(driver, selenium.ByID, "password", 10*time.Second)
    if err != nil {
        return err
    }
    passwordField.SendKeys(password)

    submitBtn, err := waitForElement(driver, selenium.ByXPATH, ".//*[@type='submit']", 10*time.Second)
    if err != nil {
        return err
    }
    submitBtn.Click()

    time.Sleep(5 * time.Second) // Allow page to load
    return nil
}

// Extract Portfolio Data
func extractPortfolioData(driver selenium.WebDriver) ([]map[string]string, error) {
    var portfolio []map[string]string

    // Scroll multiple times to ensure all content is loaded
    for i := 0; i < 10; i++ {
        driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
        time.Sleep(3 * time.Second)
    }

    // Locate portfolio section based on header or section
    portfolioSection, err := waitForElement(driver, selenium.ByXPATH, "//header[@class='jss122']/following-sibling::section[@aria-label='Security list']", 60*time.Second)
    if err != nil {
        return nil, fmt.Errorf("Portfolio section not found: %w", err)
    }

    // Extract asset names
    assetElements, err := portfolioSection.FindElements(selenium.ByXPATH, ".//div[@role='row']//div[@role='table']")
    if err != nil {
        return nil, fmt.Errorf("Asset names not found: %w", err)
    }

    // Extract ISIN codes
    isinElements, err := portfolioSection.FindElements(selenium.ByXPATH, ".//div[@role='row']//a")
    if err != nil {
        return nil, fmt.Errorf("ISIN codes not found: %w", err)
    }

    if len(assetElements) != len(isinElements) {
        return nil, fmt.Errorf("Mismatch in number of assets and ISIN elements")
    }

    // Process assets
    for i := 0; i < len(assetElements); i++ {
        assetName, _ := assetElements[i].Text()
        isinLink, _ := isinElements[i].GetAttribute("href")

        // Extract ISIN code from URL
        isin := strings.Split(strings.Split(isinLink, "isin=")[1], "&")[0]

        portfolio = append(portfolio, map[string]string{
            "name":  assetName,
            "isin":  isin,
            "value": "0.00", // Placeholder
            "shares": "0",   // Placeholder
        })
    }

    return portfolio, nil
}



// Helper Function: Wait for Element
func waitForElement(driver selenium.WebDriver, by string, value string, timeout time.Duration) (selenium.WebElement, error) {
    for i := 0; i < int(timeout.Seconds()); i++ {
        elem, err := driver.FindElement(by, value)
        if err == nil {
            return elem, nil
        }
        time.Sleep(1 * time.Second)
    }
    return nil, fmt.Errorf("element not found: %s", value)
}
