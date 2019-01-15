package models

import (
    "database/sql"
    _ "fmt"
    "github.com/lib/pq"
    "time"
)


type Game struct {
	Id           int    `json:id`
	Status       string `json:status`
	PlayerTurnId int    `json:playerTurnId`
	Players      []Player
	Timestamps
}

func (db *DB) GetGame(gameId int) (*Game, error) {
    rows, _ := db.Query("SELECT * from games where id = $1", gameId)
    games, _ := parseGameRows(rows)
    if len(games) == 0 {
        return nil, gameNotFoundError{}
    }
    game := &games[0]

    if err := buildGame(db, game); err != nil {
        return nil, err
    }

    return game, nil
}

func (db *DB) GetGames(userId int) ([]*Game, error) {
    rows, _ := db.Query("SELECT DISTINCT games.* FROM games JOIN players P ON P.user_id = $1", userId)
    games, _ := parseGameRows(rows)

    gameRefs := make([]*Game, 0)
    for _, game := range games {
        buildGame(db, &game)
        gameRefs = append(gameRefs, &game)
    }

    return gameRefs, nil
}

// CreateGame builds all the necessary information to start a game, including
// associated Player entries.
func (db *DB) CreateGame(userId, opponentId int) (*Game, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }
    time := pq.FormatTimestamp(time.Now())

    // Create Game
    rows, err := tx.Query("INSERT INTO games VALUES (nextval('games_id_seq'), $1, $2, $3, $4) RETURNING *", "not-started", nil, time, time)
    if err != nil {
        _ = tx.Rollback()
        return nil, err
    }
    games, _ := parseGameRows(rows)
    game := &games[0]

    // Create Player 1 (inviter)
    player1, err := createPlayer(tx, userId, game.Id, "active", "black", time)
    if err != nil {
        _ = tx.Rollback()
        return nil, err
    }
    user, _ := db.GetUser(userId)
    player1.User = user

    // Create Player 2 (invitee)
    player2, err := createPlayer(tx, opponentId, game.Id, "user-pending", "white", time)
    if err != nil {
        _ = tx.Rollback()
        return nil, err
    }
    user, _ = db.GetUser(opponentId)
    player2.User = user

    game.Players = []Player{*player1, *player2}

    if err != nil {
        _ = tx.Rollback()
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    rows, _ = db.Query("UPDATE games SET player_turn_id = $1 WHERE id = $2 RETURNING *", player1.Id, game.Id)

    games, _ = parseGameRows(rows)

    return &games[0], nil
}

// Pass is a game action where a player decides that they cannot make a
// move on their turn. If both players pass, the game ends.
func (db *DB) Pass(userId, gameId int) (*Game, error) {
    user, _ := db.GetUser(userId)
    game, err := db.GetGame(gameId)

    if err != nil {
        return nil, err
    }

    if err := validateGame(game); err != nil {
        return nil, err
    }

    if err := validateTurn(game, user.Id); err != nil {
        return nil, err
    }

    var otherPlayer Player
    for _, player := range game.Players {
        if player.userId != user.Id {
            otherPlayer = player
        }
    }

    time := pq.FormatTimestamp(time.Now())

    var rows *sql.Rows
    if otherPlayer.HasPassed {
        rows, err = db.Query("UPDATE games SET (status, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", "complete", time, game.Id)
    } else {
        tx, _ := db.Begin()
        rows, err = tx.Query("UPDATE players SET (has_passed, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", true, time, game.Id)
        if err != nil {
            return nil, err
        }
        rows, err = tx.Query("UPDATE games SET (player_turn_id, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", otherPlayer.Id, time, game.Id)
        if err != nil {
            return nil, err
        }

        if err := tx.Commit(); err != nil {
            return nil, err
        }
    }

    games, _ := parseGameRows(rows)
    game = &games[0]

    buildGame(db, game)

    return game, nil
}

func validateGame(game *Game) error {
    if game.Status == "complete" {
        return gameCompleteError{}
    } else if game.Status == "not-started" {
        return gameNotStartedError{}
    } else {
        return nil
    }
}

func validateTurn(game *Game, userId int) error {
    var player *Player
    for i, p := range game.Players {
        if p.userId == userId {
            player = &game.Players[i]
        }
    }

    if player == nil {
        return userNotInGameError{}
    }

    if player.Id != game.PlayerTurnId {
        return wrongTurnError{}
    }

    return nil
}

func buildGame(db *DB, game *Game) error {
    rows, _ := db.Query("SELECT * from players where game_id = $1", game.Id)
    players, _ := parsePlayerRows(rows)

    // We save an extra query by getting both users
    rows, _ = db.Query("SELECT DISTINCT users.* FROM users JOIN players P ON (P.game_id = $1)", game.Id)
    users, _ := parseUserRows(rows)

    for i, player := range players {
        for j, user := range users {
            if player.userId == user.Id {
                players[i].User = &users[j]
            }
        }
    }

    game.Players = players

    return nil
}

func parseGameRows(rows *sql.Rows) ([]Game, error) {
    defer rows.Close()

    games := make([]Game, 0)
    for rows.Next() {
        var game Game
        rows.Scan(
            &game.Id,
            &game.Status,
            &game.PlayerTurnId,
            &game.InsertedAt,
            &game.UpdatedAt,
        )
        games = append(games, game)
    }

    return games, rows.Err()
}
