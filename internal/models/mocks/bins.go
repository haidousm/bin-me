package mocks

import (
	"time"

	"binme.haido.us/internal/models"
)

var mockBin = models.Bin{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type BinModel struct{}

func (m *BinModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *BinModel) Get(id int) (models.Bin, error) {
	switch id {
	case 1:
		return mockBin, nil
	default:
		return models.Bin{}, models.ErrNoRecord
	}
}

func (m *BinModel) Latest() ([]models.Bin, error) {
	return []models.Bin{mockBin}, nil
}
