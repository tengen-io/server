package providers

import (
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/models"
	"strconv"
)

type UserProvider struct {
	db *sqlx.DB
}

func NewUserProvider(db *sqlx.DB) *UserProvider {
	return &UserProvider{
		db,
	}
}

func (p *UserProvider) GetUserById(id string) (*models.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = p.db.Select(&user, "SELECT name FROM users WHERE id = $1", idInt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
