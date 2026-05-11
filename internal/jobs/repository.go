package jobs

import (
	"time"

	"gorm.io/gorm"
)

type JobPost struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:text" json:"title"`
	Company   string    `gorm:"type:text" json:"company"`
	Location  string    `gorm:"type:text" json:"location"`
	PostedAt  string    `gorm:"type:text" json:"posted_at"`
	URL       string    `gorm:"type:text;uniqueIndex" json:"url"`
	RawText   string    `gorm:"type:text" json:"raw_text"`
	ScrapedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"scraped_at"`
}

type Repository interface {
	Create(job *JobPost) error
	FindAll() ([]JobPost, error)
	FindByURL(url string) (*JobPost, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(job *JobPost) error {
	return r.db.Create(job).Error
}

func (r *repository) FindAll() ([]JobPost, error) {
	var jobs []JobPost
	err := r.db.Order("scraped_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *repository) FindByURL(url string) (*JobPost, error) {
	var job JobPost
	err := r.db.Where("url = ?", url).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&JobPost{}, id).Error
}