package main

import (
	"linkedin-automation/internal/browser"
	"log"
	"time"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 2)")

	b := browser.New()
	defer b.Close()

	page := b.Instance.MustPage("https://www.linkedin.com")

	log.Println("Page title:", page.MustInfo().Title)

	// Keep browser open for observation
	time.Sleep(10 * time.Second)
}
