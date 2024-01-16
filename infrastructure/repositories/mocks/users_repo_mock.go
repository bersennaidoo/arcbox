package mocks

import "github.com/bersennaidoo/arcbox/domain/models"

type UsersMockRepository struct{}

func (m *UsersMockRepository) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UsersMockRepository) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}
func (m *UsersMockRepository) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
