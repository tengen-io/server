package providers

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IdentityProvider struct {
	db         *sqlx.DB
	BcryptCost int
}

func NewIdentityProvider(db *sqlx.DB, bcryptCost int) *IdentityProvider {
	return &IdentityProvider{
		db,
		bcryptCost,
	}
}

func (p *IdentityProvider) CreateIdentity(email string, password string, name string) (*models.Identity, error) {
	// TODO(eac): re-add validation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), p.BcryptCost)
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Identity
	ts := pq.FormatTimestamp(time.Now().UTC())

	// TODO(eac): do a precondition check for duplicate users to save autoincrement IDs
	identity := tx.QueryRowx("INSERT INTO identities (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, email", email, passwordHash, ts, ts)
	err = identity.Scan(&rv.Id, &rv.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := tx.QueryRowx("INSERT INTO users (identity_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name", rv.Id, name, ts, ts)
	err = user.Scan(&rv.User.Id, &rv.User.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

// TODO(eac): switch to sqlx.get
func (p *IdentityProvider) GetIdentityById(id int32) (*models.Identity, error) {
	var identity models.Identity
	row := p.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND i.id = $1", id)
	err := row.Scan(&identity.Id, &identity.Email, &identity.User.Id, &identity.User.Name)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}
