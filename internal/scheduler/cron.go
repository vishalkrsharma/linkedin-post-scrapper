package scheduler

import (
	"log"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"linkedin-hunter/internal/alerts"
	"linkedin-hunter/internal/jobs"
	"linkedin-hunter/internal/scraper"
)

func Start(
	repo jobs.Repository,
	scraperClient scraper.Client,
	telegramBot *alerts.TelegramBot,
	keywords string,
	location string,
) {
	c := cron.New()

	_, err := c.AddFunc("55 15 * * *", func() {
		runScrape(repo, scraperClient, telegramBot, keywords, location)
	})
	if err != nil {
		log.Printf("Failed to add cron job: %v", err)
	}

	c.Start()
	log.Println("Scheduler started - running at 9:25 PM IST daily")
}

func runScrape(
	repo jobs.Repository,
	scraperClient scraper.Client,
	telegramBot *alerts.TelegramBot,
	keywords string,
	location string,
) {
	log.Println("Starting scheduled scrape...")
	keywordList := strings.Split(keywords, ",")

	for _, keyword := range keywordList {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}

		log.Printf("Scraping: %s in %s", keyword, location)

		result, err := scraperClient.Scrape(keyword, location)
		if err != nil {
			log.Printf("Scraper error for %s: %v", keyword, err)
			continue
		}

		savedJobs, err := jobs.SaveJobsFromScraper(repo, result)
		if err != nil {
			log.Printf("Failed to save jobs: %v", err)
			continue
		}

		if telegramBot != nil && len(savedJobs) > 0 {
			msg := jobs.FormatJobMessage(savedJobs)
			if err := telegramBot.SendMessage(msg); err != nil {
				log.Printf("Telegram send error: %v", err)
			}
		}

		time.Sleep(5 * time.Second)
	}

	log.Println("Scheduled scrape completed")
}