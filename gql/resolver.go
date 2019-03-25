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

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MatchmakingRequests(ctx context.Context, user string) (<-chan models.MatchmakingRequestsPayload, error) {
	rv := make(chan models.MatchmakingRequestsPayload, 5)
	c := r.repo.Subscribe(pubsub.MkTopic(pubsub.TopicCategoryMatchmakeRequests, user))

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

			payload := &models.MatchmakingRequestCompletePayload{
				Game: *game,
			}

			rv <- payload
		}
	}()

	return rv, nil
}

