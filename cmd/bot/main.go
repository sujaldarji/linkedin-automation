package main

import (
	"log"
	"strings"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/search"

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

	runSearchDemo(b)

	// Keep browser open for observation / demo
	time.Sleep(20 * time.Second)
}

func runSearchDemo(b *browser.Browser) {
	input := &search.Input{
		Keywords:  "software engineer",
		PageLimit: 2,
	}
	_ = input.Validate()

	url := search.BuildPeopleSearchURL(input)

	nav := search.NewNavigator(b)

	// ✅ capture the page returned by OpenSearch
	page := nav.OpenSearch(url)

	// Phase 3 actions
	search.LightScroll(page)

	profiles, err := search.ParseProfiles(page)
	
	if err != nil {
		log.Fatalf("failed to parse profiles: %v", err)
	}

	log.Printf("Found %d profiles", len(profiles))
	for _, p := range profiles {
		log.Println(p.URL)
	}
}

