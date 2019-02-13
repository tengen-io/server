package resolvers

import "github.com/camirmas/go_stop/models"

type Resolvers struct {
	db         models.DB
	signingKey []byte
}

func NewResolvers(db models.DB, signingKey []byte) *Resolvers {
	return &Resolvers{
		db,
		signingKey,
	}
}
