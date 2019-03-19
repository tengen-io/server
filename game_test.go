package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"testing"
)

func Test_transitionGameState(t *testing.T) {
	g := models.Game{
		State: models.GameStateInvitation,
		Type: models.GameTypeStandard,
		BoardSize: 19,
		NodeFields: models.NodeFields{
			Id: "1",
		},
	}

	users := []models.GameUserEdge{
		{
			Index: 0,
			Type: models.GameUserEdgeTypeOwner,
			User: models.User{
				Name: "test user",
				NodeFields: models.NodeFields{
					Id: "1",
				},
			},
		},
	}

	newState, err := transitionGameState(g, users)
	assert.NoError(t, err)
}
