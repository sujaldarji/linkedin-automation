package main

import (
	"log"
	"time"
	"linkedin-automation/internal/browser"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 1)")

	b := browser.New()
	defer b.Close()

	page := b.Instance.MustPage("https://www.linkedin.com")

	log.Println("Page title:", page.MustInfo().Title)

	// Keep browser open for observation
	time.Sleep(5 * time.Second)
}
