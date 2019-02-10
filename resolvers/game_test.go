package resolvers

import (
	"context"
	"github.com/camirmas/go_stop/models"
	"testing"
)

func TestCreateGame(t *testing.T) {
	t.Run("with missing token", createGameInvalidToken)
	t.Run("with self as opponent", createGameSelf)
	t.Run("with non-self opponent", createGame)
}

func createGameInvalidToken(t *testing.T) {
	r, params := setup()
	params.Context = context.WithValue(params.Context, "token", nil)

	_, err := r.CreateGame(params)

	expectErr(t, invalidTokenError{}, err)
}

func createGameSelf(t *testing.T) {
	r, params := setupAuth()
	params.Args["opponentUsername"] = "dude"

	_, err := r.CreateGame(params)

	expectErr(t, sameUserError{}, err)
}

func createGame(t *testing.T) {
	r, params := setupAuth()
	params.Args["opponentUsername"] = "saitama"

	_, err := r.CreateGame(params)

	if err != nil {
		t.Errorf("Expected new game, got error: %s", err.Error())
	}
}

func TestGetGame(t *testing.T) {
	r, params := setup()
	params.Args["id"] = "1"

	game, err := r.GetGame(params)

	if err != nil {
		t.Errorf("Expected Game, got error: %s", err.Error())
	}

	if _, ok := game.(*models.Game); !ok {
		t.Errorf("Expected Game, got %v", game)
	}
}

func TestGetGames(t *testing.T) {
	r, params := setup()
	params.Args["userId"] = "1"

	games, err := r.GetGames(params)

	if err != nil {
		t.Error("Expected Games, got error")
	}

	if _, ok := games.([]*models.Game); !ok {
		t.Errorf("Expected Games, got %v", games)
	}
}

func TestGetLobby(t *testing.T) {
	r, params := setup()

	games, err := r.GetLobby(params)

	if err != nil {
		t.Error("Expected Games, got error")
	}

	if _, ok := games.([]*models.Game); !ok {
		t.Errorf("Expected Games, got %v", games)
	}
}

func TestPass(t *testing.T) {
	t.Run("with missing token", passInvalidToken)
	t.Run("with already completed game", passComplete)
	t.Run("with invalid turn", passInvalidTurn)
	t.Run("when not in game", passNotInGame)
	t.Run("with valid turn", pass)
}

func passInvalidToken(t *testing.T) {
	r, params := setup()
	params.Context = context.WithValue(params.Context, "token", nil)

	_, err := r.Pass(params)

	expectErr(t, invalidTokenError{}, err)
}

func passComplete(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "2"

	_, err := r.Pass(params)

	expectErr(t, gameCompleteError{}, err)
}

func passInvalidTurn(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "3"

	_, err := r.Pass(params)

	expectErr(t, wrongTurnError{}, err)
}

func passNotInGame(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "10"

	_, err := r.Pass(params)

	expectErr(t, userNotInGameError{}, err)
}

func pass(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "4"

	game, err := r.Pass(params)

	if err != nil {
		t.Errorf("Expected Game, got error: %s", err.Error())
	}

	if _, ok := game.(*models.Game); !ok {
		t.Errorf("Expected Game, got %v", game)
	}
}

func TestAddStone(t *testing.T) {
	t.Run("with missing token", addStoneInvalidToken)
	t.Run("with already completed game", addStoneComplete)
	t.Run("with invalid turn", addStoneInvalidTurn)
	t.Run("when not in game", addStoneNotInGame)
	t.Run("with valid turn", addStone)
}

func addStoneInvalidToken(t *testing.T) {
	r, params := setup()
	params.Context = context.WithValue(params.Context, "token", nil)

	_, err := r.AddStone(params)

	expectErr(t, invalidTokenError{}, err)
}

func addStoneComplete(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "2"

	_, err := r.AddStone(params)

	expectErr(t, gameCompleteError{}, err)
}

func addStoneInvalidTurn(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "3"

	_, err := r.AddStone(params)

	expectErr(t, wrongTurnError{}, err)
}

func addStoneNotInGame(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "10"

	_, err := r.AddStone(params)

	expectErr(t, userNotInGameError{}, err)
}

func addStone(t *testing.T) {
	r, params := setupAuth()
	params.Args["gameId"] = "4"
	params.Args["x"] = 3
	params.Args["y"] = 3

	stone, err := r.AddStone(params)

	if err != nil {
		t.Errorf("Expected Stone, got error: %s", err.Error())
	}

	if _, ok := stone.(models.Stone); !ok {
		t.Errorf("Expected Stone, got %v", stone)
	}
}
