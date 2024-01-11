package mysql

import "database/sql"

type UsersRepository struct {
	dbu *sql.DB
}

func NewUsersRepository(dbu *sql.DB) *UsersRepository {
	return &UsersRepository{
		dbu: dbu,
	}
}

func (u *UsersRepository) Insert(name, email, password string) error {
	return nil
}

func (u *UsersRepository) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (u *UsersRepository) Exists(id int) (bool, error) {
	return false, nil
}
