package server

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/tengen-io/server/models"
)

func hasAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	_, ok := ctx.Value(IdentityContextKey).(models.Identity)
	if !ok {
		return nil, fmt.Errorf("access denied")
	}

	return next(ctx)
}

func Directives() DirectiveRoot {
	return DirectiveRoot{
		HasAuth: hasAuth,
	}
}
