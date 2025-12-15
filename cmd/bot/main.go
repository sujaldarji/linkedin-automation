package main

import (
	"log"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 3)")

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}

	b := browser.New()
	defer b.Close()

	page := b.NewPage("https://www.linkedin.com")

	if auth.RestoreSession(page) {
		log.Println("Session restored using cookies")
	} else {
		log.Println("No cookies found, logging in via env vars")

		if !auth.Login(page) {
			log.Fatal("Login failed")
		}

		if err := auth.SaveCookies(page, "cookies.json"); err != nil {
			log.Fatal("Failed to save cookies:", err)
		}

		log.Println("Login successful, cookies saved")
	}

	time.Sleep(20 * time.Second)
}
