package gql

import (
	"errors"
	"github.com/tengen-io/server/models"
)

func transitionGameState(game models.Game, players []models.GameUserEdge) (models.GameState, error) {
	switch game.Type {
	case models.GameTypeStandard:
		switch game.State {
		case models.GameStateInvitation:
			var numPlayers = 0
			for _, player := range players {
				if player.Type == models.GameUserEdgeTypeOwner || player.Type == models.GameUserEdgeTypePlayer {
					numPlayers += 1
				}
			}

			if numPlayers > 2 {
				return -1, errors.New("invalid game state")
			}

			if numPlayers == 2 {
				return models.GameStateInProgress, nil
			}

			return game.State, nil
		default:
			return -1, errors.New("unknown game state")
		}
	default:
		return -1, errors.New("unknown game type")
	}
}
