package main

import (
	"fmt"
	"log"
	"os"

	"linkedin-hunter/internal/alerts"
	"linkedin-hunter/internal/api"
	"linkedin-hunter/internal/jobs"
	"linkedin-hunter/internal/scraper"
	"linkedin-hunter/internal/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "linkedin_hunter"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&jobs.JobPost{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	repo := jobs.NewRepository(db)
	service := jobs.NewService(repo)
	scraperClient := scraper.NewClient(getEnv("SCRAPER_URL", "http://localhost:8000"))

	telegramToken := getEnv("TELEGRAM_BOT_TOKEN", "")
	var telegramBot *alerts.TelegramBot
	if telegramToken != "" {
		telegramBot = alerts.NewBot(telegramToken, getEnv("TELEGRAM_CHAT_ID", ""))
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	api.SetupRoutes(app, service, scraperClient, telegramBot)

	schedulerKeywords := getEnv("SCHEDULE_KEYWORDS", "software engineer,devops")
	schedulerLocation := getEnv("SCHEDULE_LOCATION", "India")
	scheduler.Start(repo, scraperClient, telegramBot, schedulerKeywords, schedulerLocation)

	port := getEnv("PORT", "3000")
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}