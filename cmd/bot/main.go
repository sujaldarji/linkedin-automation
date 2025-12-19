package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/stealth"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 4)")

	_ = godotenv.Load() // env is optional

	b := browser.New()
	defer b.Close()

	// Open LinkedIn home
	page := b.NewPage("https://www.linkedin.com")
	time.Sleep(2 * time.Second)

	// Check login status
	if strings.Contains(page.MustInfo().URL, "/feed") {
		log.Println("✅ Already logged in via persistent browser profile")
	} else {
		log.Println("No active session, performing login")

		if !auth.Login(page) {
			log.Fatal("Login failed")
		}

		log.Println("Login successful (session persisted by Chrome)")
	}

	// Stage 4 – Phase 4 demo (pagination + human behavior)
	runSearchDemo(b)

	// Keep browser open briefly for observation/demo
	time.Sleep(30 * time.Second)
}

func runSearchDemo(b *browser.Browser) {
	log.Println("\n=== Stage 4 | Phase 4: Pagination + Human Behavior Demo ===")

	baseURL := "https://www.linkedin.com/search/results/people/?keywords=software%20engineer"
	maxPages := 3

	allProfiles := make(map[string]search.Profile)

	for pageNum := 1; pageNum <= maxPages; pageNum++ {
		log.Printf("\n[Phase 4] Navigating to page %d", pageNum)

		searchURL := baseURL
		if pageNum > 1 {
			searchURL += "&page=" + strconv.Itoa(pageNum)
		}

		// Open search page
		page := b.NewPage(searchURL)

		// Human-like pause after navigation (think time)
		log.Println("[Behavior] Page loaded, waiting briefly before interaction")
		stealth.RandomDelay(2000, 4000)

		// Light human scroll to trigger lazy loading
		log.Println("[Behavior] Performing light human-like scroll")
		search.LightScroll(page)

		// Human reading pause before parsing
		log.Println("[Behavior] Pausing before reading results")
		stealth.RandomDelay(1000, 1500)

		// Parse visible results (read-only)
		profiles := search.ParseVisibleResults(page)

		// Aggregate with global deduplication
		for _, p := range profiles {
			if _, exists := allProfiles[p.URL]; exists {
				continue
			}
			allProfiles[p.URL] = p
		}

		log.Printf(
			"[Phase 4] Page %d complete — total unique profiles collected so far: %d",
			pageNum,
			len(allProfiles),
		)
	}

	// Final summary
	log.Println("\n=== FINAL RESULTS ===")
	log.Printf("Total unique profiles collected: %d\n", len(allProfiles))

	i := 1
	for _, p := range allProfiles {
		log.Printf("%d. %s", i, p.URL)
		if p.Name != "" {
			log.Printf("   Name: %s", p.Name)
		}
		i++
	}
}
