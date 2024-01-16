package contracts

import "github.com/bersennaidoo/arcbox/domain/models"

type SnipRepositoryInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*models.Snip, error)
	Latest() ([]*models.Snip, error)
}

type UserRepositoryInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}
