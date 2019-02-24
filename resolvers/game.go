package resolvers

import (
	// "fmt"
	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/rules"
	"github.com/graphql-go/graphql"
	// "reflect"
)

// GetGame retrieves a Game by integer id.
func (r *Resolvers) GetGame(p graphql.ResolveParams) (interface{}, error) {
	return r.db.GetGame(p.Args["id"].(string))
}

// GetGames retrieves Games for a given User id.
func (r *Resolvers) GetGames(p graphql.ResolveParams) (interface{}, error) {
	return r.db.GetGames(p.Args["userId"].(string))
}

func (r *Resolvers) GetLobby(p graphql.ResolveParams) (interface{}, error) {
	return r.db.GetGames(nil)
}

// CreateGame starts a new Game, with the current User and a provided opponent
// as Players.
func (r *Resolvers) CreateGame(p graphql.ResolveParams) (interface{}, error) {
	currentUser, ok := p.Context.Value("currentUser").(*models.User)
	if !ok {
		return nil, invalidTokenError{}
	}

	opponentUsername := p.Args["opponentUsername"].(string)
	opponent, err := r.db.GetUser(opponentUsername)
	if err != nil {
		return nil, err
	}

	if opponent.Id == currentUser.Id {
		return nil, sameUserError{}
	}

	return r.db.CreateGame(currentUser.Id, opponent)
}

// Pass executes a pass maneuver for the Game with the current User.
func (r *Resolvers) Pass(p graphql.ResolveParams) (interface{}, error) {
	currentUser, ok := p.Context.Value("currentUser").(*models.User)
	if !ok {
		return nil, invalidTokenError{}
	}

	gameId := p.Args["gameId"].(string)
	game, err := r.db.GetGame(gameId)

	if err != nil {
		return nil, err
	}

	if err := validateGame(game); err != nil {
		return nil, err
	}

	if err := validateTurn(game, currentUser.Id); err != nil {
		return nil, err
	}

	return r.db.Pass(currentUser.Id, game)
}

// AddStone executes a move by the current User to add a Stone
// to the Board, with given X and Y coordinates.
func (r Resolvers) AddStone(p graphql.ResolveParams) (interface{}, error) {
	currentUser, ok := p.Context.Value("currentUser").(*models.User)
	if !ok {
		return nil, invalidTokenError{}
	}
	gameId := p.Args["gameId"].(string)
	game, err := r.db.GetGame(gameId)

	if err != nil {
		return nil, err
	}

	if err := validateGame(game); err != nil {
		return nil, err
	}

	if err := validateTurn(game, currentUser.Id); err != nil {
		return nil, err
	}

	currentPlayer, _ := game.CurrentPlayer(currentUser.Id)

	x := p.Args["x"].(int)
	y := p.Args["y"].(int)

	stone := models.Stone{X: x, Y: y, Color: currentPlayer.Color}

	if !contains(game.Stones, stone) {
		game.Stones = append(game.Stones, stone)
	} else {
		return nil, stoneExistsError{}
	}

	stringsToRemove, err := rules.Run(game.BoardSize, game.Stones, stone)
	toRemove := flatten(stringsToRemove)

	if err != nil {
		return nil, err
	}

	if err := validateKo(game, toRemove); err != nil {
		return nil, err
	}

	if len(toRemove) == 1 {
		game.LastTaker = &stone
	} else {
		game.LastTaker = nil
	}

	newStones := make([]models.Stone, 0)
	for _, s := range game.Stones {
		if !contains(toRemove, s) {
			newStones = append(newStones, s)
		}
	}

	err = r.db.UpdateGame(currentUser.Id, game, stone, toRemove)

	if err != nil {
		return nil, err
	}

	for _, s := range game.Stones {
		if s.X == stone.X && s.Y == stone.Y {
			stone = s
		}
	}

	stone.Game = game

	return stone, nil
}

func validateGame(game *models.Game) error {
	if game.Status == "complete" {
		return gameCompleteError{}
		// } else if game.Status == "not-started" {
		// 	return gameNotStartedError{}
	} else {
		return nil
	}
}

func validateTurn(game *models.Game, userId int) error {
	var player *models.Player
	for i, p := range game.Players {
		if p.UserId == userId {
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

func contains(stones []models.Stone, stone models.Stone) bool {
	for _, s := range stones {
		if s.X == stone.X && s.Y == stone.Y {
			return true
		}
	}

	return false
}

func validateKo(game *models.Game, toRemove []models.Stone) error {
	if len(toRemove) != 1 {
		return nil
	}
	if game.LastTaker.X == toRemove[0].X && game.LastTaker.Y == toRemove[0].Y {
		return koViolationError{}
	}
	return nil
}

func flatten(strings []rules.String) []models.Stone {
	var flattened []models.Stone
	for _, str := range strings {
		for _, s := range str {
			flattened = append(flattened, s)
		}
	}
	return flattened
}
