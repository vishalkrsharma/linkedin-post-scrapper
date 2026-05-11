# linkedin-hunter — Project Context for Claude Code

> Paste this file into Claude Code at the start of your session to resume planning/implementation.

---

## Project Goal

Build a **personal LinkedIn hiring-post scraper** that:
- Searches LinkedIn for hiring/job posts based on keywords
- Extracts structured data from results
- Stores data in a database
- Exports to CSV
- Sends Telegram alerts for new matches
- Displays results in a lightweight dashboard
- Runs both **manually** and on a **cron schedule**

---

## Key Decisions Already Made

| Decision | Choice | Reason |
|---|---|---|
| Scale | Personal tool (just me) | Low volume, simpler auth |
| LinkedIn access method | **Playwright browser automation** with real session cookies | Free, reliable at low volume, avoids API restrictions |
| Backend language | **Go (Fiber framework)** | Developer is actively learning Go |
| Scraper language | **Python + Playwright + FastAPI** | Go's Playwright bindings are immature; Python is battle-tested here |
| Database | **PostgreSQL** | Via GORM in Go |
| Scheduler | **gocron** (Go library) | Built into the Go backend |
| Alerts | **Telegram Bot** via `gopkg.in/telebot.v3` | Simple, free |
| Dashboard | **HTMX + Go Templ** | No React overhead for a personal tool |

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Your Laptop / VPS                  │
│                                                       │
│  ┌──────────────┐     ┌───────────────────────────┐  │
│  │  Go Backend  │────▶│  Python Scraper Service   │  │
│  │  (main app)  │     │  (Playwright + FastAPI)   │  │
│  └──────┬───────┘     └───────────────────────────┘  │
│         │                                             │
│  ┌──────▼───────┐     ┌───────────────────────────┐  │
│  │  PostgreSQL  │     │   Cron Scheduler (Go)     │  │
│  │  + CSV export│     │   or system cron          │  │
│  └──────┬───────┘     └───────────────────────────┘  │
│         │                                             │
│  ┌──────▼───────┐     ┌───────────────────────────┐  │
│  │  Go REST API │     │  Telegram Bot Alerts      │  │
│  │  + Dashboard │     │  (Go telebot library)     │  │
│  │  (HTMX/Templ)│     └───────────────────────────┘  │
│  └──────────────┘                                     │
└─────────────────────────────────────────────────────┘
```

---

## Project Structure

```
linkedin-hunter/
├── cmd/
│   └── server/main.go          # entry point
├── internal/
│   ├── scraper/client.go       # HTTP client → Python service
│   ├── jobs/repository.go      # GORM DB layer
│   ├── jobs/service.go         # business logic
│   ├── scheduler/cron.go       # gocron setup
│   ├── alerts/telegram.go      # bot alerts
│   └── api/handlers.go         # Fiber route handlers
├── templates/                  # Templ HTML templates
├── scraper-service/            # Python Playwright service
│   ├── main.py                 # FastAPI app
│   ├── linkedin.py             # Playwright scraping logic
│   └── requirements.txt
├── docker-compose.yml          # PostgreSQL + both services
└── .env
```

---

## Database Schema (planned)

```sql
CREATE TABLE job_posts (
  id          SERIAL PRIMARY KEY,
  title       TEXT,
  company     TEXT,
  location    TEXT,
  posted_at   TIMESTAMP,
  url         TEXT UNIQUE,
  raw_text    TEXT,
  scraped_at  TIMESTAMP DEFAULT NOW()
);
```

---

## Implementation Phases

| Phase | What to build | Status |
|---|---|---|
| **1** | Python scraper + FastAPI (LinkedIn → JSON) | ⬜ TODO |
| **2** | Go backend + PostgreSQL + GORM | ⬜ TODO |
| **3** | Manual trigger via CLI / API (end-to-end flow) | ⬜ TODO |
| **4** | Cron scheduler (gocron) | ⬜ TODO |
| **5** | Telegram alerts | ⬜ TODO |
| **6** | HTMX dashboard + CSV export | ⬜ TODO |

---

## Important Notes & Gotchas

- **Start with Phase 1 first** — if LinkedIn blocks Playwright, nothing else matters
- **Use saved LinkedIn session cookies** (not login automation) — far more reliable
- **Add 3–5s random delays** between Playwright requests to avoid rate limiting
- **Schedule cron at off-peak hours** (e.g., 3 AM IST) to reduce detection risk
- LinkedIn's DOM selectors change frequently — use `data-` attributes where possible, add fallback selectors

---

## Developer Background

- **Languages known**: Java, Spring Boot, JavaScript
- **Currently learning**: Go (Golang) — familiar with Fiber, Echo, GORM from research
- **Python experience**: LangChain/LangGraph, FastAPI, uv package manager
- **OS**: Arch Linux, zsh, Vim
- **Machine**: Lenovo ThinkPad E14 Gen 6 (Ryzen 7, 16GB RAM)

---

## Where to Resume

Start implementation at **Phase 1**:

> "Let's implement Phase 1 — the Python Playwright scraper service. Create `scraper-service/linkedin.py` and `scraper-service/main.py` with a FastAPI endpoint `POST /scrape` that accepts `{keyword, location}` and returns a list of job posts using Playwright with saved session cookies."
