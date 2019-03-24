package matchmake

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type inMemoryPool struct {
	pool    []request
	matches []struct{ i, j request }
}

func (p *inMemoryPool) requests() []request {
	return p.pool
}

func (p *inMemoryPool) match(i, j request) {
	p.matches = append(p.matches, struct{ i, j request }{i, j})
}

func Test_tick(t *testing.T) {
	testCases := []struct {
		name     string
		requests []request
		expected []struct{ i, j int64 }
	}{
		{
			"within delta",
			[]request{r(1, 1, 5, 3), r(2, 2, 3, 3)},
			[]struct{ i, j int64 }{{1, 2}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			pool := inMemoryPool{
				pool:    testCase.requests,
				matches: make([]struct{ i, j request }, 0),
			}

			matchmaker := newMatchmaker(&pool, time.Duration(1)*time.Second)
			matches := matchmaker.tick()

			found := 0
		expected:
			for _, expected := range testCase.expected {
				for _, match := range matches {
					if (match.i.requestId == expected.i || match.i.requestId == expected.j) &&
						(match.j.requestId == expected.i || match.j.requestId == expected.j) {
						found++
						continue expected
					}
				}
			}

			assert.Len(t, testCase.expected, found)
		})
	}
}

func r(id, user int64, rank, delta int) request {
	return request{
		id, user, rank, delta,
	}
}
