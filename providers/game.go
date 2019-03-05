package providers

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"strconv"
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
// TODO(eac): Switch to sqlx binding
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

func (p *GameProvider) GetGameById(id string) (*models.Game, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var game models.Game
	err = p.db.Get(&game, "SELECT * FROM games WHERE id = $1", idInt)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (p *GameProvider) GetUsersForGame(id string) ([]*models.GameUserEdge, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var rv = make([]*models.GameUserEdge, 1)
	err = p.db.Select(&rv, "SELECT * FROM game_user WHERE game_id = $1", idInt)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (p *GameProvider) GetGamesByIdAndState(ids []string, states []models.GameState) ([]models.Game, error) {
	query := "SELECT * FROM games WHERE "
	args := make([]interface{}, 0)

	if len(ids) > 0 {
		idInts := make([]int, len(ids))
		for i, id := range ids {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			idInts[i] = idInt
		}

		fragment, fragArgs, err := sqlx.In("id IN (?)", idInts)
		if err != nil {
			return nil, err
		}

		query += fragment
		args = append(args, fragArgs...)
	}

	if len(states) > 0 {
		query += "id IN (?)"
		fragment, fragArgs, err := sqlx.In("state IN (?)", states)
		if err != nil {
			return nil, err
		}
		if len(ids) > 0 {
			query += " AND"
		}
		query += fragment
		args = append(args, fragArgs...)
	}

	query = p.db.Rebind(query)
	rows, err := p.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rv := make([]models.Game, 0)
	for rows.Next() {
		var i models.Game
		err := rows.StructScan(&i)
		if err != nil {
			return nil, err
		}
		rv = append(rv, i)
	}

	return rv, nil
}
