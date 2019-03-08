package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tengen-io/server/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	game     *GameProvider
	identity *IdentityProvider
	user     *UserProvider
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Game() GameResolver {
	return &gameResolver{r}
}

type gameResolver struct{ *Resolver }

func (r *gameResolver) Users(ctx context.Context, obj *models.Game) ([]*models.GameUserEdge, error) {
	return r.game.GetUsersForGame(obj.Id)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateGameInvitation(ctx context.Context, input *models.CreateGameInvitationInput) (*models.Game, error) {
	identity, _ := ctx.Value("currentUser").(models.Identity)
	rv, err := r.game.CreateInvitation(identity, *input)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (r *mutationResolver) JoinGame(ctx context.Context, gameId string) (*models.JoinGamePayload, error) {
	identity, ok:= ctx.Value("currentUser").(models.Identity)
	if !ok {
		return nil, errors.New("invalid user")
	}

	rv, err := r.game.CreateGameUser(gameId, identity.Id, models.GameUserEdgeTypePlayer)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id *string, name *string) (*models.User, error) {
	if id != nil {
		user, err := r.user.GetUserById(*id)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context, ids []string, names []string) ([]*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) Game(ctx context.Context, id *string) (*models.Game, error) {
	if id != nil {
		game, err := r.game.GetGameById(*id)
		if err != nil {
			return nil, err
		}

		return game, nil
	}

	panic("not implemented")
}

func (r *queryResolver) Games(ctx context.Context, ids []string, states []models.GameState) ([]*models.Game, error) {
	if len(ids) > 0 && len(states) > 0 {
		return nil, errors.New("arguments are mutually exclusive")
	}

	if len(ids) > 0 {
		return r.game.GetGamesByIds(ids)
	}

	if len(states) > 0 {
		return r.game.GetGamesByState(states)
	}

	panic("not implemented")
}
