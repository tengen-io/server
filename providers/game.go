package providers

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"time"
)

type GameProvider struct {
	db *sqlx.DB
}

func NewGameProvider(db *sqlx.DB) *GameProvider {
	return &GameProvider{
		db,
	}
}

// TODO(eac): Add user from auth
// TODO(eac): Add validation
func (p *GameProvider) CreateInvitation(input models.CreateGameInvitationInput) (*models.Game, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	var rv models.Game
	ts := pq.FormatTimestamp(time.Now().UTC())

	game := tx.QueryRowx("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, type, state, board_size", input.Type, models.Invitation, input.BoardSize, ts, ts)
	err = game.Scan(&rv.Id, &rv.Type, &rv.State, &rv.BoardSize)
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
