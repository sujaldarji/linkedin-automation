# LinkedIn Automation

This repository contains a **Go-based proof of concept (PoC)** for LinkedIn browser automation using a real Chrome session powered by the **Rod** library.

---
## 🎥 Video Demonstration

A short walkthrough

**Video link:** https://youtu.be/kgsKSQly7JA

---

## 📌 Project Scope

This PoC focuses on:

- Persistent browser session handling
- Login flow with manual checkpoint / CAPTCHA support
- LinkedIn People Search navigation
- Parsing visible profile URLs from search results
- Pagination with controlled limits
- SQLite-based state tracking
- Read-only profile visits with human-like behavior
---

## ⚙️ How It Works

### 1. Browser & Login
- Launches a real Chrome instance with a persistent profile
- Reuses login sessions across runs
- Allows manual completion of LinkedIn checkpoints if triggered

### 2. Search & Parse
- Builds LinkedIn People Search URLs
- Navigates result pages
- Extracts visible profile URLs
- Stores unique profiles in SQLite

### 3. Profile Visit Loop
- Retrieves unvisited profiles from state store
- Opens profiles in new tabs
- Scrolls and pauses to simulate reading
- Marks profiles as visited

---

## 🗃 State Management

The application uses **SQLite (pure Go)** to track profile state across runs.

Stored information includes:
- Profile URL
- Visit status
- Creation timestamp
- Last action timestamp

This prevents duplicate processing and enables controlled iteration.

---

## 🔐 Environment Configuration

Create a `.env` file using the provided template.

### `.env.example`

```env
# LinkedIn login credentials
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password
```
## Running the Project
```
go mod tidy

go run cmd/bot/main.go
```
---

## Disclaimer

This project is intended solely for **educational and technical demonstration purposes**.  
Users are responsible for complying with LinkedIn’s terms of service and applicable laws.
