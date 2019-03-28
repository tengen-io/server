package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/db"
	"github.com/tengen-io/server/pubsub"
)

type Repository struct {
	db *sqlx.DB
	h  db.Handle
	pubsub *pubsub.DbPubSub
}

func NewRepository(db *sqlx.DB, pubsub *pubsub.DbPubSub) *Repository {
	return &Repository{
		db: db,
		h:  db,
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

func (d *Repository) Publish(topic pubsub.TopicCategory, payload pubsub.Event) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = d.h.Exec("SELECT pg_notify($1, $2)", string(topic), string(payloadBytes))
	if err != nil {
		return err
	}

	return nil
}

func (d *Repository) Subscribe(topic pubsub.Topic) <-chan pubsub.Event {
	return d.pubsub.Subscribe(topic)
}
