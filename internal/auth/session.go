package auth

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func RestoreSession(page *rod.Page) bool {
	if _, err := os.Stat("cookies.json"); err != nil {
		return false
	}

	if err := LoadCookies(page, "cookies.json"); err != nil {
		log.Println("Failed to load cookies:", err)
		return false
	}

	page.MustNavigate("https://www.linkedin.com/feed")
	page.MustWaitLoad()

	return strings.Contains(page.MustInfo().URL, "/feed")
}

// WaitForHumanCheckpoint pauses automation
func WaitForHumanCheckpoint(page *rod.Page) {
	log.Println("⚠️ Checkpoint/CAPTCHA detected")
	log.Println("👉 Please solve CAPTCHA manually in browser")
	log.Println("👉 Automation will resume after checkpoint")

	for {
		time.Sleep(3 * time.Second)
		url := page.MustInfo().URL

		if strings.Contains(url, "/feed") {
			log.Println("✅ Checkpoint resolved by human")
			return
		}
	}
}
