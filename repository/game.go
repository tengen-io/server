package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"strconv"
	"time"
)

// TODO(eac): Add validation
// TODO(eac): Switch to sqlx binding
func (r *Repository) CreateGame(gameType models.GameType, boardSize int, gameState models.GameState, users []models.User) (*models.Game, error) {
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

	var rv models.Game
	ts := pq.FormatTimestamp(time.Now().UTC())

	game := tx.QueryRow("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, type, state, board_size", gameType, gameState, boardSize, ts, ts)
	err := game.Scan(&rv.Id, &rv.Type, &rv.State, &rv.BoardSize)
	if err != nil {
		return nil, err
	}

	insertStmt, err := tx.Prepare("INSERT INTO game_user (game_id, user_id, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		_, err = insertStmt.Exec(rv.Id, user.Id, models.GameUserEdgeTypePlayer, ts, ts)
		if err != nil {
			return nil, err
		}
	}

	/*	payload := pubsub.Event{
			Event: "create",
			Payload: &rv,
		}

		r.pubsub.Publish(payload, "game", "game_"+rv.State.String()) */

	if r.tx == nil {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return &rv, nil
}

// TODO(eac): Validation? here or in the resolver
// TODO(eac): Need this anymore?
func (r *Repository) CreateGameUser(gameId string, userId string, edgeType models.GameUserEdgeType) (*models.Game, error) {
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

	rows, err := tx.Query("SELECT user_id, type FROM game_user WHERE game_id = $1", gameId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	gameUsers := make([]models.GameUserEdge, 0)
	for rows.Next() {
		var gameUser models.GameUserEdge
		err := rows.Scan(&gameUser.User.Id, &gameUser.Type)
		if err != nil {
			return nil, err
		}

		gameUsers = append(gameUsers, gameUser)
	}

	var rv models.Game
	ts := pq.FormatTimestamp(time.Now().UTC())
	_, err = tx.Exec("INSERT INTO game_user (game_id, user_id, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", gameId, userId, edgeType, ts, ts)
	if err != nil {
		return nil, err
	}

	newEdge := models.GameUserEdge{
		Type: edgeType,
		User: models.User{
			NodeFields: models.NodeFields{
				Id: userId,
			},
		},
	}
	gameUsers = append(gameUsers, newEdge)

	game := tx.QueryRowx("SELECT id, type, state FROM games WHERE id = $1", gameId)
	err = game.Scan(&rv.Id, &rv.Type, &rv.State)
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

func (r *Repository) GetGameById(id string) (*models.Game, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var game models.Game
	row := r.db.QueryRowx("SELECT * FROM games WHERE id = $1", idInt)
	err = row.StructScan(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *Repository) GetUsersForGame(id string) ([]models.GameUserEdge, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query("SELECT type, user_id, name FROM game_user gu, users u WHERE game_id = $1 AND gu.user_id = u.id", idInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rv = make([]models.GameUserEdge, 0)
	for rows.Next() {
		var i models.GameUserEdge
		err := rows.Scan(&i.Type, &i.User.Id, &i.User.Name)
		if err != nil {
			return nil, err
		}

		rv = append(rv, i)
	}

	return rv, nil
}

func (r *Repository) GetGamesByIds(ids []string) ([]*models.Game, error) {
	idInts := make([]int, len(ids))
	for i, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		idInts[i] = idInt
	}

	query, fragArgs, err := sqlx.In("SELECT * FROM games WHERE id IN (?)", idInts)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, fragArgs...)

	query = r.db.Rebind(query)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rv := make([]*models.Game, 0)
	for rows.Next() {
		var i models.Game
		err := rows.StructScan(&i)
		if err != nil {
			return nil, err
		}
		rv = append(rv, &i)
	}

	return rv, nil
}

func (r *Repository) GetGamesByState(states []models.GameState) ([]*models.Game, error) {
	query, fragArgs, err := sqlx.In("SELECT * FROM games WHERE state IN (?)", states)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	args = append(args, fragArgs...)

	query = r.db.Rebind(query)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rv := make([]*models.Game, 0)
	for rows.Next() {
		var i models.Game
		err := rows.StructScan(&i)
		if err != nil {
			return nil, err
		}
		rv = append(rv, &i)
	}

	return rv, nil
}
