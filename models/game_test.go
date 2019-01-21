package models

import (
	_ "fmt"
	"testing"
)

func TestCreateGame(t *testing.T) {
	t.Run("with invalid user", createGameInvalidUser)
	t.Run("with invalid opponent", createGameInvalidOpponent)
	t.Run("with proper arguments", createGame)

	teardown()
}

func createGameInvalidUser(t *testing.T) {
	_, err := db.CreateGame(1234, "invalid")

	expectErr(t, userNotFoundError{}, err)
}

func createGameInvalidOpponent(t *testing.T) {
	user, _ := db.CreateUser("creategame", "creategame@dude.dude", "dudedude", "dudedude")
	_, err := db.CreateGame(user.Id, "invalid")

	expectErr(t, userNotFoundError{}, err)
}

func createGame(t *testing.T) {
	user1, _ := db.CreateUser("creategame1", "creategame1@dude.dude", "dudedude", "dudedude")
	user2, _ := db.CreateUser("creategame2", "creategame2@dude.dude", "dudedude", "dudedude")

	game, err := db.CreateGame(user1.Id, user2.Username)

	if err != nil {
		t.Error(err)
	}

	if len(game.Players) != 2 {
		t.Error("Expected 2 Players in game")
	}
	if game.Status != "active" {
		t.Error("Expected game status to be 'active'")
	}
}

func TestGetGames(t *testing.T) {
	t.Run("with user", getGamesByUser)
	t.Run("recent games", getGamesRecent)
	teardown()
}

func getGamesByUser(t *testing.T) {
	user1, _ := db.CreateUser("getgames1", "getgames1@dude.dude", "dudedude", "dudedude")
	user2, _ := db.CreateUser("getgames2", "getgames2@dude.dude", "dudedude", "dudedude")
	db.CreateGame(user1.Id, user2.Username)

	games, err := db.GetGames(user1.Id)

	if err != nil {
		t.Error(err)
	}

	if len(games) != 1 {
		t.Errorf("Expected 1 game, found %d", len(games))
	}
}

func getGamesRecent(t *testing.T) {
	games, err := db.GetGames(nil)

	if err != nil {
		t.Error(err)
	}

	if len(games) != 1 {
		t.Errorf("Expected 1 game, found %d", len(games))
	}
}

func TestPass(t *testing.T) {
	user1, _ := db.CreateUser("pass1", "pass1@dude.dude", "dudedude", "dudedude")
	user2, _ := db.CreateUser("pass2", "pass2@dude.dude", "dudedude", "dudedude")
	game, _ := db.CreateGame(user1.Id, user2.Username)

	turnId := game.PlayerTurnId

	game, err := db.Pass(user1.Id, game)

	if err != nil {
		t.Error(err)
	}

	if turnId == game.PlayerTurnId {
		t.Error("Expected player turn to change")
	}

	game, err = db.Pass(user2.Id, game)

	if err != nil {
		t.Error(err)
	}

	if game.Status != "complete" {
		t.Errorf("Expected game status to be 'complete', got '%s'", game.Status)
	}

	teardown()
}

func TestUpdateBoard(t *testing.T) {
	user1, _ := db.CreateUser("updateboard1", "updateboard1@dude.dude", "dudedude", "dudedude")
	user2, _ := db.CreateUser("updateboard2", "updateboard2@dude.dude", "dudedude", "dudedude")
	game, _ := db.CreateGame(user1.Id, user2.Username)

	stone := Stone{0, 0, "black"}
	game.Board.Stones = []Stone{stone}
	turnId := game.PlayerTurnId

	game, err := db.UpdateBoard(user1.Id, game)

	if err != nil {
		t.Error(err)
	}

	if turnId == game.PlayerTurnId {
		t.Error("Expected player turn to change")
	}

	if len(game.Board.Stones) != 1 {
		t.Error("Expected Board to have 1 Stone")
	}

	teardown()
}
