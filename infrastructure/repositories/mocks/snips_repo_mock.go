package mocks

import (
	"time"

	"github.com/bersennaidoo/arcbox/domain/models"
)

var mockSnip = &models.Snip{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnipsMockRepository struct{}

func (m *SnipsMockRepository) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}
func (m *SnipsMockRepository) Get(id int) (*models.Snip, error) {
	switch id {
	case 1:
		return mockSnip, nil
	default:
	}
	return nil, models.ErrNoRecord
}

func (m *SnipsMockRepository) Latest() ([]*models.Snip, error) {
	return []*models.Snip{mockSnip}, nil
}
