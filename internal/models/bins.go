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

	stmt := `INSERT INTO bins (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BinModel) Get(id int) (Bin, error) {
	return Bin{}, nil
}

func (m *BinModel) Latest() ([]Bin, error) {
	return nil, nil
}
