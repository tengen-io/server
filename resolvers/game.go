package resolvers

import (
    "github.com/camirmas/go_stop/models"
    "github.com/graphql-go/graphql"
)

func GetGame(p graphql.ResolveParams) (interface{}, error) {
    db := p.Context.Value("db").(models.Database)
    return db.GetGame(p.Args["id"].(int))
}

func GetGames(p graphql.ResolveParams) (interface{}, error) {
    db := p.Context.Value("db").(models.Database)
    return db.GetGames(p.Args["userId"].(int))
}

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

    return db.Pass(claims.UserId, game)
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
