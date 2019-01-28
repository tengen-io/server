package resolvers

import (
	_ "fmt"
	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/rules"
	"github.com/graphql-go/graphql"
	"reflect"
)

// GetGame retrieves a Game by integer id.
func GetGame(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	return db.GetGame(p.Args["id"].(string))
}

// GetGames retrieves Games for a given User id.
func GetGames(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	return db.GetGames(p.Args["userId"].(string))
}

func GetLobby(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	return db.GetGames(nil)
}

// CreateGame starts a new Game, with the current User and a provided opponent
// as Players.
func CreateGame(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	signingKey := p.Context.Value("signingKey").([]byte)
	token, ok := p.Context.Value("token").(string)
	if !ok {
		return nil, invalidTokenError{}
	}

	opponentUsername := p.Args["opponentUsername"].(string)
	opponent, err := db.GetUser(opponentUsername)
	if err != nil {
		return nil, err
	}

	claims, err := ValidateToken(token, signingKey)
	if err != nil {
		return nil, err
	}

	if opponent.Id == claims.UserId {
		return nil, sameUserError{}
	}

	return db.CreateGame(claims.UserId, opponent)
}

// Pass executes a pass maneuver for the Game with the current User.
func Pass(p graphql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("db").(models.Database)
	signingKey := p.Context.Value("signingKey").([]byte)
	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, invalidTokenError{}
	}
	gameId := p.Args["gameId"].(string)

	claims, err := ValidateToken(token, signingKey)

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
	signingKey := p.Context.Value("signingKey").([]byte)
	token, ok := p.Context.Value("token").(string)

	if !ok {
		return nil, invalidTokenError{}
	}
	gameId := p.Args["gameId"].(string)

	claims, err := ValidateToken(token, signingKey)

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

	currentPlayer, _ := game.CurrentPlayer(user.Id)

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

	if len(toRemove) > 0 {
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

	game, err = db.UpdateBoard(user.Id, game, stone, toRemove)

	if err != nil {
		return nil, err
	}

	for _, s := range game.Stones {
		if s.X == stone.X && s.Y == stone.Y {
			stone = s
		}
	}

	return &stone, nil
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
	if reflect.DeepEqual(game.LastTaker, toRemove[0]) {
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
