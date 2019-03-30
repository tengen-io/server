package gql

import (
	"context"
	"errors"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/pubsub"
	"github.com/tengen-io/server/repository"
	"log"
)

type Resolver struct {
	repo   *repository.Repository
	auth auth
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
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type gameResolver struct{ *Resolver }

func (r *gameResolver) Users(ctx context.Context, obj *models.Game) ([]models.GameUserEdge, error) {
	return r.repo.GetUsersForGame(obj.Id)
}

type mutationResolver struct{ *Resolver }

func (m mutationResolver) CreateMatchmakingRequest(ctx context.Context, input models.CreateMatchmakingRequestInput) (*models.CreateMatchmakingRequestPayload, error) {
	identity, ok := ctx.Value(IdentityContextKey).(models.Identity)
	if !ok {
		return nil, errors.New("invalid user")
	}

	var rv models.CreateMatchmakingRequestPayload
	err := m.repo.WithTx(func(r *repository.Repository) error {
		req, err := r.CreateMatchmakingRequest(identity.User, input.Delta)
		rv.Request = req

		if err != nil {
			log.Println("unable to create mm request", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Println("unable to commit mm txn", err)
		return nil, err
	}

	return &rv, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id *string, name *string) (*models.User, error) {
	if id != nil {
		user, err := r.repo.GetUserById(*id)
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

func (r *queryResolver) Viewer(ctx context.Context) (*models.Identity, error) {
	identity, ok := ctx.Value(IdentityContextKey).(models.Identity)
	if !ok {
		// TODO(eac): this is asserted already by @hasAuth. Should I just ignore the error?
		return nil, errors.New("invalid user")
	}

	return &identity, nil
}

func (r *queryResolver) Game(ctx context.Context, id *string) (*models.Game, error) {
	if id != nil {
		game, err := r.repo.GetGameById(*id)
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
		return r.repo.GetGamesByIds(ids)
	}

	if len(states) > 0 {
		return r.repo.GetGamesByState(states)
	}

	panic("not implemented")
}

func (r *queryResolver) MatchmakingRequests(ctx context.Context) ([]models.MatchmakingRequest, error) {
	id, err := r.auth.authForContext(ctx)
	if err != nil {
		return nil, err
	}

	requests, err := r.repo.GetMatchmakingRequestsForUser(id.User)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MatchmakingRequestCompletions(ctx context.Context) (<-chan *models.MatchmakingRequestCompletionPayload, error) {
	id, err := r.auth.authForContext(ctx)
	if err != nil {
		return nil, err
	}

	rv := make(chan *models.MatchmakingRequestCompletionPayload, 5)
	c := r.repo.Subscribe(pubsub.MkTopic(pubsub.TopicCategoryMatchmakeRequests, id.User.Id))

	go func() {
		for event := range c {
			gameId, ok := event.Payload["game"].(string)
			if !ok {
				// TODO: log
				return
			}
			game, err := r.repo.GetGameById(gameId)

			if err != nil {
				// TODO: log
				return
			}

			payload := &models.MatchmakingRequestCompletionPayload{
				Game: *game,
			}

			rv <- payload
		}
	}()

	return rv, nil
}
