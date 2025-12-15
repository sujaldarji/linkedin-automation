package browser

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"linkedin-automation/internal/stealth"
)

// Browser represents one persistent human browser session
type Browser struct {
	Instance *rod.Browser
	Stealth  *stealth.Config
}

// New initializes Chrome with a persistent user profile
func New() *Browser {
	// Create persistent Chrome profile directory
	profileDir := filepath.Join(".", "chrome-profile")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		log.Fatalf("failed to create chrome profile dir: %v", err)
	}

	log.Printf("Using persistent profile: %s", profileDir)

	u := launcher.New().
		Bin(`C:\Program Files\Google\Chrome\Application\chrome.exe`).
		UserDataDir(profileDir). // 🔑 PERSISTENT SESSION
		Headless(false).
		Set("disable-blink-features", "AutomationControlled").         // REQUIRED for CAPTCHA/manual login
		MustLaunch()

	browser := rod.New().ControlURL(u)

	if err := browser.Connect(); err != nil {
		log.Fatalf("failed to connect browser: %v", err)
	}

	return &Browser{
		Instance: browser,
		Stealth:  stealth.NewConfig(),
	}
}

// NewPage creates a new page with stealth applied
func (b *Browser) NewPage(url string) *rod.Page {
	page := b.Instance.MustPage()

	stealth.Apply(page, b.Stealth)

	page.MustNavigate(url)
	page.MustWaitLoad()
	return page
}

// Close shuts down the browser
func (b *Browser) Close() {
	_ = b.Instance.Close()
}
