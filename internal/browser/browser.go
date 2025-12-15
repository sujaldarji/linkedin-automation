package browser

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Browser wraps rod.Browser for clean architecture
type Browser struct {
	Instance *rod.Browser
}

// New launches a visible Chrome browser
func New() *Browser {
	u := launcher.New().
		Bin(`C:\Program Files\Google\Chrome\Application\chrome.exe`).
		Headless(false). // VERY IMPORTANT for demo
		MustLaunch()

	browser := rod.New().ControlURL(u)

	err := browser.Connect()
	if err != nil {
		log.Fatalf("failed to connect browser: %v", err)
	}

	return &Browser{
		Instance: browser,
	}
}

// Close shuts down the browser
func (b *Browser) Close() {
	_ = b.Instance.Close()
}
