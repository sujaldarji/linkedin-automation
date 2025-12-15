package auth

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Login performs a demo-safe LinkedIn login.
// It allows manual CAPTCHA resolution and fails gracefully if login
// does not complete within a limited number of retries.
func Login(page *rod.Page) bool {
	email := os.Getenv("LINKEDIN_EMAIL")
	password := os.Getenv("LINKEDIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("Missing LinkedIn credentials in env")
		return false
	}

	// 1. Navigate to login page
	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitLoad()

	// 2. Fill credentials
	emailInput := page.Timeout(15 * time.Second).
		MustElement(`input[name="session_key"]`)
	passInput := page.MustElement(`input[name="session_password"]`)

	emailInput.MustClick()
	emailInput.MustInput(email)
	time.Sleep(700 * time.Millisecond)

	passInput.MustClick()
	passInput.MustInput(password)
	time.Sleep(700 * time.Millisecond)

	// 3. Submit login form
	page.MustElement(`button[type="submit"]`).MustClick()

	// 4. Allow LinkedIn to redirect
	time.Sleep(3 * time.Second)

	// 5. Graceful post-login check loop (max 5 attempts)
	const maxChecks = 6
	for i := 1; i <= maxChecks; i++ {
		info, err := page.Info()
		if err != nil {
			log.Println("Page no longer available (possibly closed by user)")
			return false
		}

		url := info.URL
		log.Printf("Post-login check %d/%d: %s\n", i, maxChecks, url)

		// Successful login
		if strings.Contains(url, "/feed") {
			return true
		}

		// Checkpoint / CAPTCHA flow
		if strings.Contains(url, "/checkpoint") {
			log.Println("Checkpoint/CAPTCHA detected. Please solve it manually...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Redirected back to login or error page
		if strings.Contains(url, "/login") {
			log.Println("Redirected back to login page. Login not completed.")
			return false
		}

		// Unknown state, wait once more
		time.Sleep(5 * time.Second)
	}

	log.Println("Login not completed after maximum wait time")
	return false
}
