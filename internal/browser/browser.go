package browser

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"linkedin-automation/internal/stealth"
)

// Browser represents one human-like browser session
type Browser struct {
	Instance *rod.Browser
	Stealth  *stealth.Config
}

// New initializes Chrome, connects Rod, and prepares stealth config
func New() *Browser {
	u := launcher.New().
		Bin(`C:\Program Files\Google\Chrome\Application\chrome.exe`).
		Headless(false).
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
	page := b.Instance.MustPage(url)

	stealth.Apply(page, b.Stealth)

	page.MustWaitLoad()
	return page
}

// Close shuts down the browser
func (b *Browser) Close() {
	_ = b.Instance.Close()
}
