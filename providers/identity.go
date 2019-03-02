package providers

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IdentityProvider struct {
	db *sqlx.DB
	BcryptCost int
}

func NewIdentityProvider(db *sqlx.DB, bcryptCost int) *IdentityProvider {
	return &IdentityProvider{
		db,
		bcryptCost,
	}
}

func (p *IdentityProvider) CreateIdentity(input models.CreateIdentityInput) (*models.Identity, error) {
	// TODO(eac): re-add validation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), p.BcryptCost)
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Identity
	ts := pq.FormatTimestamp(time.Now().UTC())

	identity := tx.QueryRowx("INSERT INTO identities (email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?) RETURNING id, email", input.Email, passwordHash, ts, ts)
	err = identity.StructScan(&rv)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := tx.QueryRowx("INSERT INTO users (identity_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)", rv.Id, input.Name, ts, ts)
	err = user.StructScan(&rv.User)
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

func (p *IdentityProvider) GetIdentityById(id int32) (*models.Identity, error) {
	var identity models.Identity
	row := p.db.QueryRowx("SELECT i.*, u.* FROM identities i JOIN users u USING (i.id, u.identity_id)")
	err := row.StructScan(&identity)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}
