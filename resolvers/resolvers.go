package resolvers

import (
	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/providers"
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
