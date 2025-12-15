package main

import (
	"log"
	"strings"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 3)")

	_ = godotenv.Load() // env is optional

	b := browser.New()
	defer b.Close()

	page := b.NewPage("https://www.linkedin.com")

	// Simple check: navigate to feed and see if we're logged in
	time.Sleep(2 * time.Second)
	
	if strings.Contains(page.MustInfo().URL, "/feed") {
		log.Println("✅ Already logged in via persistent browser profile")
	} else {
		log.Println("No active session, performing login")

		if !auth.Login(page) {
			log.Fatal("Login failed")
		}

		log.Println("Login successful (session persisted by Chrome)")
	}

	// Keep browser open for observation / demo
	time.Sleep(20 * time.Second)
}