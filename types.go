package main
/*

import (
	"github.com/graphql-go/graphql"
	"github.com/tengen-io/server/models"
)

type Objects struct {
	User     *graphql.Object
	AuthUser *graphql.Object
	Player   *graphql.Object
	Game     *graphql.Object
	Stone    *graphql.Object
}

func buildObjects() *Objects {
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*models.User); ok {
						return user.Id, nil
					}
					return nil, nil
				},
			},
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*models.User); ok {
						return user.Username, nil
					}
					return nil, nil
				},
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(*models.User); ok {
						return user.Email, nil
					}
					return nil, nil
				},
			},
		},
	})

	playerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Player",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(models.Player); ok {
						return player.Id, nil
					}
					return nil, nil
				},
			},
			"status": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(models.Player); ok {
						return player.Status, nil
					}
					return nil, nil
				},
			},
			"color": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(models.Player); ok {
						return player.Color, nil
					}
					return nil, nil
				},
			},
			"hasPassed": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(models.Player); ok {
						return player.HasPassed, nil
					}
					return nil, nil
				},
			},
			"user": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if player, ok := p.Source.(models.Player); ok {
						return player.User, nil
					}
					return nil, nil
				},
			},
		},
	})

	stoneType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Stone",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if stone, ok := p.Source.(models.Stone); ok {
						return stone.Id, nil
					}
					return nil, nil
				},
			},
			"x": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if stone, ok := p.Source.(models.Stone); ok {
						return stone.X, nil
					}
					return nil, nil
				},
			},
			"y": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if stone, ok := p.Source.(models.Stone); ok {
						return stone.Y, nil
					}
					return nil, nil
				},
			},
			"color": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if stone, ok := p.Source.(models.Stone); ok {
						return stone.Color, nil
					}
					return nil, nil
				},
			},
		},
	})

	gameType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Game",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.Id, nil
					}
					return nil, nil
				},
			},
			"status": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.Status, nil
					}
					return nil, nil
				},
			},
			"playerTurnId": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.PlayerTurnId, nil
					}
					return nil, nil
				},
			},
			"players": &graphql.Field{
				Type: graphql.NewList(playerType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.Players, nil
					}
					return nil, nil
				},
			},
			"boardSize": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.BoardSize, nil
					}
					return nil, nil
				},
			},
			"lastTaker": &graphql.Field{
				Type: stoneType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.LastTaker, nil
					}
					return nil, nil
				},
			},
			"stones": &graphql.Field{
				Type: graphql.NewList(stoneType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if game, ok := p.Source.(*models.Game); ok {
						return game.Stones, nil
					}
					return nil, nil
				},
			},
		},
	})

	stoneType.AddFieldConfig(
		"game",
		&graphql.Field{
			Type: gameType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if stone, ok := p.Source.(models.Stone); ok {
					return stone.Game, nil
				}
				return nil, nil
			},
		},
	)

	tokenType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthUser",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if authUser, ok := p.Source.(*models.AuthUser); ok {
						return authUser.User, nil
					}
					return nil, nil
				},
			},
			"token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if authUser, ok := p.Source.(*models.AuthUser); ok {
						return authUser.Jwt, nil
					}
					return nil, nil
				},
			},
		},
	})

	return &Objects{
		userType,
		tokenType,
		playerType,
		gameType,
		stoneType,
	}
} */
