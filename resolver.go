package main

import (
	"context"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/providers"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	game     *providers.GameProvider
	identity *providers.IdentityProvider
	user     *providers.UserProvider
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
	panic("not implemented yet")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateGameInvitation(ctx context.Context, input *models.CreateGameInvitationInput) (*models.Game, error) {
	rv, err := r.game.CreateInvitation(*input)
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

func (r *queryResolver) Games(ctx context.Context, ids []string) ([]*models.Game, error) {
	panic("not implemented")
}
