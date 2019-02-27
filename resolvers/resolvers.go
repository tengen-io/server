package resolvers

import (
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/providers"
)

type Resolvers struct {
	db   models.DB
	auth *providers.Auth
}

func NewResolvers(db models.DB, auth *providers.Auth) *Resolvers {
	return &Resolvers{
		db,
		auth,
	}
}
