package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"strconv"
	"testing"
)

func Test_transitionGameState(t *testing.T) {
	testCases := []struct {
		name string
		state models.GameState
		gameType models.GameType
		numUsers int
		newState models.GameState
	}{
		{
			"standard game invitation with one user",
			models.GameStateInvitation, models.GameTypeStandard,
			1,
			models.GameStateInvitation,
		},
		{
			"standard game invitation with two users",
			models.GameStateInvitation, models.GameTypeStandard,
			2,
			models.GameStateInProgress,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			g := models.Game{
				State: testCase.state,
				Type: testCase.gameType,
				BoardSize: 19,
				NodeFields: models.NodeFields{
					Id: "1",
				},
			}

			users := make([]models.GameUserEdge, 0)
			for i := 0; i < testCase.numUsers; i++ {
				user := models.GameUserEdge{
					Index: i,
					Type: models.GameUserEdgeTypePlayer,
					User: models.User{
						Name: "test user",
						NodeFields: models.NodeFields{
							Id: strconv.Itoa(i),
						},
					},
				}

				users = append(users, user)
			}

			newState, err := transitionGameState(g, users)
			assert.NoError(t, err)
			assert.Equal(t, newState, testCase.newState)
		})
	}
}
