package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bersennaidoo/arcbox/domain/models"
	my "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository struct {
	dbu *sql.DB
}

func NewUsersRepository(dbu *sql.DB) *UsersRepository {
	return &UsersRepository{
		dbu: dbu,
	}
}

func (u *UsersRepository) Insert(name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = u.dbu.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *my.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (u *UsersRepository) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := u.dbu.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UsersRepository) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := u.dbu.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
