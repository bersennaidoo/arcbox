package mysql

import "database/sql"

type SnipsRepository struct {
	dbc *sql.DB
}

func New(dbc *sql.DB) *SnipsRepository {
	return &SnipsRepository{
		dbc: dbc,
	}
}
