package matchmake

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/models"
	"strconv"
	"testing"
	"time"
)

type inMemoryPool struct {
	pool    []models.MatchmakingRequest
	matches []struct{ i, j models.MatchmakingRequest }
}

func (p *inMemoryPool) requests() ([]models.MatchmakingRequest, error) {
	return p.pool, nil
}

func (p *inMemoryPool) match(i, j models.MatchmakingRequest) error {
	p.matches = append(p.matches, struct{ i, j models.MatchmakingRequest }{i, j})
	return nil
}

func Test_tick(t *testing.T) {
	testCases := []struct {
		name     string
		requests []models.MatchmakingRequest
		expected []struct{ i, j string }
	}{
		{
			"within delta",
			[]models.MatchmakingRequest{r(1, 1, 5, 3), r(2, 2, 3, 3)},
			[]struct{ i, j string }{{"1", "2"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			pool := inMemoryPool{
				pool:    testCase.requests,
				matches: make([]struct{ i, j models.MatchmakingRequest }, 0),
			}

			matchmaker := newMatchmaker(&pool, time.Duration(1)*time.Second)
			matches, err := matchmaker.tick()
			assert.NoError(t, err)

			found := 0
		expected:
			for _, expected := range testCase.expected {
				for _, match := range matches {
					if (match.i.Id == expected.i || match.i.Id == expected.j) &&
						(match.j.Id == expected.i || match.j.Id == expected.j) {
						found++
						continue expected
					}
				}
			}

			assert.Len(t, testCase.expected, found)
		})
	}
}

func r(id, user int64, rank, delta int) models.MatchmakingRequest {
	return models.MatchmakingRequest{
		NodeFields: models.NodeFields{Id: strconv.FormatInt(id, 10)},
		Queue: "",
		UserId: user,
		Rank: rank,
		Delta: delta,
	}
}
