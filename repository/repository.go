package repository

import "github.com/jmoiron/sqlx"

type dbHandle interface {
	sqlx.Queryer
	sqlx.Execer
	sqlx.Preparer

	Rebind(string) string
}

type Repository struct {
	db *sqlx.DB
	h  dbHandle
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		db: db,
		h:  db,
	}
}

func (m *Repository) WithTx(f func(r *Repository) error) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	r := &Repository{
		db: m.db,
		h:  tx,
	}

	err = f(r)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
