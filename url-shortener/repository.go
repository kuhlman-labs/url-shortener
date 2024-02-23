package urlshortener

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type URLSchema struct {
	gorm.Model
	Slug     string `gorm:"type:varchar(100);unique_index"`
	ShortUrl string `gorm:"type:varchar(100);unique_index"`
	LongUrl  string `gorm:"type:varchar(100);unique_index"`
}

type SQLURLRepository struct {
	db *gorm.DB
}

type URLRepository interface {
	CreateURL(u *URLSchema) error
	ReadURL(slug string) (*URLSchema, error)
	UpdateURL(slug string, newLongURL string) error
	DeleteURL(slug string) error
}

func NewSQLURLRepository() (*SQLURLRepository, error) {
	db, err := gorm.Open("sqlite3", "url.db")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&URLSchema{})

	return &SQLURLRepository{
		db: db,
	}, nil
}

func (s *SQLURLRepository) CreateURL(u *URLSchema) error {
	if err := s.db.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (s *SQLURLRepository) ReadURL(longURL string) (*URLSchema, error) {
	var url URLSchema
	if err := s.db.Where("long_url = ?", longURL).First(&url).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil // no error, just no record found
		}
		return nil, err
	}
	return &url, nil
}

func (s *SQLURLRepository) UpdateURL(longURL string, newLongURL string) error {
	var url URLSchema
	if err := s.db.Model(&url).Where("long_url = ?", longURL).Update("long_url", newLongURL).Error; err != nil {
		return err
	}
	return nil
}

func (s *SQLURLRepository) DeleteURL(longURL string) error {
	var url URLSchema
	if err := s.db.Where("long_url = ?", longURL).Delete(&url).Error; err != nil {
		return err
	}
	return nil
}
