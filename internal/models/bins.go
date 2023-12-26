package models

import (
	"database/sql"
	"errors"
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
	stmt := `SELECT id, title, content, created, expires FROM bins
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var bin Bin

	err := row.Scan(&bin.ID, &bin.Title, &bin.Content, &bin.Created, &bin.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Bin{}, ErrNoRecord
		} else {
			return Bin{}, err
		}
	}

	return bin, nil
}

func (m *BinModel) Latest() ([]Bin, error) {
	stmt := `SELECT id, title, content, created, expires FROM bins
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bins []Bin
	for rows.Next() {
		var bin Bin
		err = rows.Scan(&bin.ID, &bin.Title, &bin.Content, &bin.Created, &bin.Expires)
		if err != nil {
			return nil, err
		}
		bins = append(bins, bin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bins, nil
}
