package models

import (
	"database/sql"
	"time"
)

type Bin struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type BinModel struct {
	DB *sql.DB
}

func (m *BinModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *BinModel) Get(id int) (Bin, error) {
	return Bin{}, nil
}

func (m *BinModel) Latest() ([]Bin, error) {
	return nil, nil
}
