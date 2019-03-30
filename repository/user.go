package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"strconv"
	"time"
)

func (r *Repository) CreateIdentity(email string, passwordHash []byte, name string) (*models.Identity, error) {
	var tx *sqlx.Tx
	if r.tx == nil {
		t, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}

		tx = t
		defer tx.Rollback()
	} else {
		tx = r.tx
	}
	// TODO(eac): re-add validation
	var rv models.Identity
	ts := pq.FormatTimestamp(time.Now().UTC())

	// TODO(eac): do a precondition check for duplicate users to save autoincrement IDs
	identity := tx.QueryRowx("INSERT INTO identities (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, email", email, passwordHash, ts, ts)
	err := identity.Scan(&rv.Id, &rv.Email)
	if err != nil {
		return nil, err
	}

	user := tx.QueryRowx("INSERT INTO users (identity_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name", rv.Id, name, ts, ts)
	err = user.Scan(&rv.User.Id, &rv.User.Name)
	if err != nil {
		return nil, err
	}

	if r.tx == nil {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return &rv, nil
}

func (r *Repository) GetIdentityById(id int32) (*models.Identity, error) {
	var identity models.Identity
	row := r.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND i.id = $1", id)
	err := row.Scan(&identity.Id, &identity.Email, &identity.User.Id, &identity.User.Name)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *Repository) GetIdentityByEmail(email string) (*models.Identity, error) {
	var identity models.Identity
	row := r.db.QueryRowx("SELECT i.id, i.email, u.id, u.name FROM identities i, users u WHERE i.id = u.identity_id AND i.email = $1", email)
	err := row.Scan(&identity.Id, &identity.Email, &identity.User.Id, &identity.User.Name)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *Repository) GetPwHashForEmail(email string) ([]byte, error) {
	row := r.db.QueryRowx("SELECT password_hash FROM identities WHERE email = $1", email)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return nil, err
	}

	return []byte(hash), nil
}

func (r *Repository) GetUserById(id string) (*models.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var rvId int
	var user models.User
	row := r.db.QueryRowx("SELECT id, name FROM users WHERE id = $1", idInt)
	err = row.Scan(&rvId, &user.Name)
	user.Id = strconv.Itoa(rvId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
