package main

import (
	"context"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/providers"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	identity *providers.IdentityProvider
	user     *providers.UserProvider
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateIdentity(ctx context.Context, input *models.CreateIdentityInput) (*models.Identity, error) {
	rv, err := r.identity.CreateIdentity(*input)
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
