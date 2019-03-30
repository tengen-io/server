package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/pubsub"
)

type Repository struct {
	db *sqlx.DB
	tx *sqlx.Tx
	pubsub *pubsub.DbPubSub
}

func NewRepository(db *sqlx.DB, pubsub *pubsub.DbPubSub) *Repository {
	return &Repository{
		db: db,
		pubsub: pubsub,
	}
}

func (r *Repository) WithTx(f func(r *Repository) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	rTx := &Repository{
		db: r.db,
		tx: tx,
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

func (r *Repository) Publish(topic pubsub.TopicCategory, payload pubsub.Event) error {
	var tx *sqlx.Tx
	if r.tx == nil {
		itx, err := r.db.Beginx()
		if err != nil {
			return err
		}

		tx = itx
		defer tx.Rollback()
	} else {
		tx = r.tx
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = tx.Exec("SELECT pg_notify($1, $2)", string(topic), string(payloadBytes))
	if err != nil {
		return err
	}

	if r.tx == nil {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Subscribe(topic pubsub.Topic) <-chan pubsub.Event {
	return r.pubsub.Subscribe(topic)
}
