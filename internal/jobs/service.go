package jobs

import (
	"fmt"
	"strings"
	"time"
)

type ScraperResult struct {
	Jobs  []JobPost `json:"jobs"`
	Count int       `json:"count"`
}

type Service interface {
	ScrapAndSave(keyword, location string) ([]JobPost, error)
	GetAll() ([]JobPost, error)
	Delete(id uint) error
	GetNewJobsCount() (int64, error)
	SaveFromScraper(result ScraperResult) ([]JobPost, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ScrapAndSave(keyword, location string) ([]JobPost, error) {
	return nil, nil
}

func (s *service) GetAll() ([]JobPost, error) {
	return s.repo.FindAll()
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) GetNewJobsCount() (int64, error) {
	return 0, nil
}

func (s *service) SaveFromScraper(result ScraperResult) ([]JobPost, error) {
	return SaveJobsFromScraper(s.repo, result)
}

func SaveJobsFromScraper(repo Repository, scraperResult ScraperResult) ([]JobPost, error) {
	var savedJobs []JobPost

	for _, job := range scraperResult.Jobs {
		existing, err := repo.FindByURL(job.URL)
		if err == nil && existing != nil {
			continue
		}

		job.ScrapedAt = time.Now()
		if err := repo.Create(&job); err != nil {
			fmt.Printf("Failed to save job: %v\n", err)
			continue
		}
		savedJobs = append(savedJobs, job)
	}

	return savedJobs, nil
}

func FormatJobMessage(jobs []JobPost) string {
	if len(jobs) == 0 {
		return "No new jobs found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d new job(s)\n\n", len(jobs)))

	for i, job := range jobs {
		if i >= 5 {
			sb.WriteString(fmt.Sprintf("\n...and %d more", len(jobs)-5))
			break
		}
		sb.WriteString(fmt.Sprintf("• %s at %s\n  📍 %s\n  🔗 %s\n\n",
			job.Title, job.Company, job.Location, job.URL))
	}

	return sb.String()
}