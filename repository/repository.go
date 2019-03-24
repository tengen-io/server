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

func (r *Repository) WithTx(f func(r *Repository) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	rTx := &Repository{
		db: r.db,
		h:  tx,
	}

	err = f(rTx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
