package main

import (
	"log"
	"strings"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/mousemovement"
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/state"
	"linkedin-automation/internal/stealth"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 4)")

	_ = godotenv.Load()

	// Initialize state store
	store, err := state.NewStore("state.db")
	if err != nil {
		log.Fatalf("failed to initialize state store: %v", err)
	}
	defer store.Close()

	b := browser.New()
	defer b.Close()

	// Open LinkedIn
	page := b.NewPage("https://www.linkedin.com")
	time.Sleep(2 * time.Second)

	// Login check
	if strings.Contains(page.MustInfo().URL, "/feed") {
		log.Println("✅ Already logged in via persistent browser profile")
	} else {
		log.Println("No active session, performing login")

		if !auth.Login(page) {
			log.Fatal("Login failed")
		}
		log.Println("Login successful (session persisted)")
	}

	// Stage 4 execution
	runSearchDemo(b, store)

	time.Sleep(20 * time.Second)
}

func runSearchDemo(b *browser.Browser, store *state.Store) {
	log.Println("\n=== Stage 4 | Search, Pagination & State Persistence Demo ===")

	input := &search.Input{
		Keywords:  "software engineer",
		PageLimit: 3,
	}

	if err := input.Validate(); err != nil {
		log.Fatalf("invalid search input: %v", err)
	}

	for pageNum := 1; pageNum <= input.PageLimit; pageNum++ {
		log.Printf("\n[Phase 4] Navigating to page %d", pageNum)

		searchURL := search.BuildPeopleSearchURL(input, pageNum)
		page := b.NewPage(searchURL)

		// Create mouse controller for this page
		mouse := mousemovement.New(page)

		// Human-like delay
		log.Println("[Behavior] Page loaded, waiting briefly")
		stealth.RandomDelay(1200, 2000)

		// Optional: hover over results container
		if container, err := page.Element("div.search-results-container"); err == nil {
			log.Println("[Behavior] Hovering over results area")
			_ = mouse.Hover(container)
			stealth.RandomDelay(400, 700)
		}

		log.Println("[Behavior] Performing light scroll")
		search.LightScroll(page)
		
		// Mouse wait simulation
		mouse.Wait()

		log.Println("[Behavior] Pausing before parsing")
		stealth.RandomDelay(800, 1500)

		// Parse profiles
		profiles := search.ParseVisibleResults(page)

		for _, p := range profiles {
			if err := store.EnsureProfile(p.URL); err != nil {
				log.Printf("[State] Failed to persist %s: %v", p.URL, err)
				continue
			}
			log.Printf("[State] Profile persisted: %s", p.URL)
		}

		log.Printf(
			"[Phase 4] Page %d complete — profiles processed: %d",
			pageNum,
			len(profiles),
		)
	}

	log.Println("\n=== Stage 4 Complete ===")
	log.Println("All discovered profiles are now persisted in SQLite state store")
}