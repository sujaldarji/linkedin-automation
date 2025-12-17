package search

import (
	"log"
	"strings"

	"github.com/go-rod/rod"
)

type Profile struct {
	URL string
}

// ParseProfiles extracts visible profile URLs from LinkedIn search results
func ParseProfiles(page *rod.Page) ([]Profile, error) {
	var profiles []Profile

	// Debug: confirm page is loaded
	html, _ := page.HTML()
	log.Printf("[DEBUG] Page HTML length: %d", len(html))

	// Enumerate frames
	frames := page.Frames()
	log.Printf("[DEBUG] Number of frames: %d", len(frames))

	for i, f := range frames {
		frameHTML, _ := f.HTML()
		log.Printf("[DEBUG] Frame %d HTML length: %d", i, len(frameHTML))

		anchors, err := f.Elements(`a[href*="/in/"]`)
		if err != nil {
			continue
		}

		log.Printf("[DEBUG] Frame %d profile links: %d", i, len(anchors))

		for _, a := range anchors {
			href, err := a.Attribute("href")
			if err != nil || href == nil {
				continue
			}

			url := normalizeProfileURL(*href)
			if url == "" {
				continue
			}

			profiles = append(profiles, Profile{
				URL: url,
			})
		}
	}

	return profiles, nil
}

// normalizeProfileURL ensures clean LinkedIn profile URLs
func normalizeProfileURL(raw string) string {
	if !strings.Contains(raw, "/in/") {
		return ""
	}

	// Strip query params
	if idx := strings.Index(raw, "?"); idx != -1 {
		raw = raw[:idx]
	}

	// Ensure absolute URL
	if strings.HasPrefix(raw, "/") {
		raw = "https://www.linkedin.com" + raw
	}

	return raw
}
