package resolvers

import (
	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/rules"
	"github.com/graphql-go/graphql"
	"reflect"
)

// GetGame retrieves a Game by integer id.
func GetGame(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	return db.GetGame(p.Args["id"].(int))
}

// GetGames retrieves Games for a given User id.
func GetGames(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	return db.GetGames(p.Args["userId"].(int))
}

// CreateGame starts a new Game, with the current User and a provided opponent
// as Players.
func CreateGame(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, invalidTokenError{}
	}

	opponentId := p.Args["opponentId"].(int)

	claims, err := ValidateToken(token)

	if err != nil {
		return nil, err
	}

	return db.CreateGame(claims.UserId, opponentId)
}

// Pass executes a pass maneuver for the Game with the current User.
func Pass(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, invalidTokenError{}
	}
	gameId := p.Args["gameId"].(int)

	claims, err := ValidateToken(token)

	if err != nil {
		return nil, err
	}

	user, _ := db.GetUser(claims.UserId)
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

	return db.Pass(user.Id, game)
}

// AddStone executes a move by the current User to add a Stone
// to the Board, with given X and Y coordinates.
func AddStone(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, invalidTokenError{}
	}
	gameId := p.Args["gameId"].(int)

	claims, err := ValidateToken(token)

	if err != nil {
		return nil, err
	}

	user, _ := db.GetUser(claims.UserId)
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

	var currentPlayer models.Player
	for i, player := range game.Players {
		if player.UserId == user.Id {
			currentPlayer = game.Players[i]
		}
	}

	x := p.Args["x"].(int)
	y := p.Args["y"].(int)

	stone := models.Stone{x, y, currentPlayer.Color}

	// TODO: update game depending on rules evaluation
	toRemove, err := rules.Run(game.Board, stone)

	if err != nil {
		return nil, err
	}

	var flattenedToRemove []models.Stone
	for _, str := range toRemove {
		for _, s := range str {
			flattenedToRemove = append(flattenedToRemove, s)
		}
	}

	newStones := make([]models.Stone, 0)
	for _, s := range game.Board.Stones {
		if !contains(flattenedToRemove, s) {
			newStones = append(newStones, s)
		}
	}

	game.Board.Stones = newStones

	db.UpdateBoard(game)

	return game, nil
}

func validateGame(game *models.Game) error {
	if game.Status == "complete" {
		return gameCompleteError{}
	} else if game.Status == "not-started" {
		return gameNotStartedError{}
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
		if reflect.DeepEqual(s, stone) {
			return true
		}
	}

	return false
}
