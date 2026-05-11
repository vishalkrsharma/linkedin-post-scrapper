package api

import (
	"encoding/csv"
	"fmt"
	"strconv"

	"linkedin-hunter/internal/alerts"
	"linkedin-hunter/internal/jobs"
	"linkedin-hunter/internal/scraper"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	service jobs.Service,
	scraperClient scraper.Client,
	telegramBot *alerts.TelegramBot,
) {
	api := app.Group("/api")

	api.Get("/jobs", handleGetJobs(service))
	api.Delete("/jobs/:id", handleDeleteJob(service))
	api.Post("/scrape", handleScrape(service, scraperClient, telegramBot))
	api.Get("/export", handleExport(service))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("LinkedIn Hunter API - Go to /api/jobs")
	})
}

func handleGetJobs(service jobs.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobs, err := service.GetAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(jobs)
	}
}

func handleDeleteJob(service jobs.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
		}

		if err := service.Delete(uint(id)); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "deleted"})
	}
}

func handleScrape(
	service jobs.Service,
	scraperClient scraper.Client,
	telegramBot *alerts.TelegramBot,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type ScrapeReq struct {
			Keyword  string `json:"keyword"`
			Location string `json:"location"`
		}

		var req ScrapeReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		if req.Keyword == "" {
			req.Keyword = "software engineer"
		}
		if req.Location == "" {
			req.Location = "India"
		}

		result, err := scraperClient.Scrape(req.Keyword, req.Location)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		savedJobs, err := service.SaveFromScraper(result)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if telegramBot != nil && len(savedJobs) > 0 {
			msg := jobs.FormatJobMessage(savedJobs)
			telegramBot.SendMessage(msg)
		}

		return c.JSON(fiber.Map{
			"new_jobs":  len(savedJobs),
			"total":     result.Count,
			"jobs":      savedJobs,
		})
	}
}

func handleExport(service jobs.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobs, err := service.GetAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "text/csv")
		c.Set("Content-Disposition", "attachment; filename=jobs.csv")

		writer := csv.NewWriter(c.Response().BodyWriter())
		writer.Write([]string{"ID", "Title", "Company", "Location", "URL", "Posted At", "Scraped At"})
		for _, job := range jobs {
			writer.Write([]string{
				fmt.Sprintf("%d", job.ID),
				job.Title,
				job.Company,
				job.Location,
				job.URL,
				job.PostedAt,
				job.ScrapedAt.String(),
			})
		}
		writer.Flush()
		return nil
	}
}