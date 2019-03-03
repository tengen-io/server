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

// TODO(eac): switch to sqlx get
func (p *UserProvider) GetUserById(id string) (*models.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var rvId int
	var user models.User
	row := p.db.QueryRow("SELECT id, name FROM users WHERE id = $1", idInt)
	err = row.Scan(&rvId, &user.Name)
	user.Id = strconv.Itoa(rvId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
