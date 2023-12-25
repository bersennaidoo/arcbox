package mysql

import (
	"database/sql"
	"errors"

	"github.com/bersennaidoo/arcbox/domain/models"
)

type SnipsRepository struct {
	dbc *sql.DB
}

func New(dbc *sql.DB) *SnipsRepository {
	return &SnipsRepository{
		dbc: dbc,
	}
}

func (s *SnipsRepository) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snips (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := s.dbc.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SnipsRepository) Get(id int) (*models.Snip, error) {
	stmt := `SELECT id, title, content, created, expires FROM snips
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := s.dbc.QueryRow(stmt, id)

	snip := models.Snip{}

	err := row.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &snip, nil
}

func (s *SnipsRepository) Latest() ([]*models.Snip, error) {

	stmt := `SELECT id, title, content, created, expires FROM snips
WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := s.dbc.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snips := []*models.Snip{}

	for rows.Next() {
		snip := &models.Snip{}

		err = rows.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
		if err != nil {
			return nil, err
		}

		snips = append(snips, snip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snips, nil
}
