package matchmake

import (
	"github.com/jmoiron/sqlx"
	"github.com/tengen-io/server/db"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/repository"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type pool interface {
	requests() ([]models.MatchmakingRequest, error)
	match(models.MatchmakingRequest, models.MatchmakingRequest) error
}

type DbPool struct {
	repo repository.Repository
}

func (p DbPool) requests() ([]models.MatchmakingRequest, error) {
	return p.repo.GetMatchmakingRequests()
}

func (p DbPool) match(i models.MatchmakingRequest, j models.MatchmakingRequest) error {
	users := []models.User{
		{NodeFields: models.NodeFields{Id: i.Id}}, {NodeFields: models.NodeFields{Id: j.Id}},
	}

	err := p.repo.WithTx(func(r *repository.Repository) error {
		_, err := r.CreateGame(models.GameTypeStandard, 19, models.GameStateNegotiation, users)
		if err != nil {
			return err
		}

		err = r.DeleteMatchmakingRequest(i)
		if err != nil {
			return err
		}

		err = r.DeleteMatchmakingRequest(j)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type match struct {
	i models.MatchmakingRequest
	j models.MatchmakingRequest
}

type matchmaker struct {
	pool         pool
	tickInterval time.Duration
}

func newMatchmaker(pool pool, tick time.Duration) *matchmaker {
	return &matchmaker{
		pool:         pool,
		tickInterval: tick,
	}
}

func (m *matchmaker) run() {
	for true {
		log.Println("starting matchmaking run...")
		start := time.Now()
		matches, err := m.tick()
		if err != nil {
			log.Printf("error in matchmaking run: %s", err)
		} else {
			log.Printf("matchmaking done. matches: %d", len(matches))
		}
		duration := time.Since(start)
		log.Printf("elapsed: %dms", duration / time.Millisecond)

		for _, match := range matches {
			m.pool.match(match.i, match.j)
		}

		time.Sleep(m.tickInterval)
	}
}

func (m *matchmaker) tick() ([]match, error) {
	requests, err := m.pool.requests()
	if err != nil {
		return nil, err
	}

	removed := make([]bool, len(requests))
	matches := make([]match, 0)

	for i := 0; i < len(requests); i++ {
		if removed[i] {
			continue
		}

		r1 := requests[i]
		pairs := make([]models.MatchmakingRequest, 0)

		for j := 0; j < len(requests); j++ {
			if i == j || removed[j] {
				continue
			}

			r2 := requests[j]
			delta := abs(r1.Rank - r2.Rank)
			if delta < r1.Delta && delta < r2.Delta {
				pairs = append(pairs, r2)
			}
		}

		if len(pairs) > 0 {
			sort.Slice(pairs, func(i, j int) bool {
				di := abs(r1.Rank - pairs[i].Rank)
				dj := abs(r1.Rank - pairs[j].Rank)
				return di < dj
			})

			removed[i] = true
			// TODO(eac): bring j along
			for j, r2 := range requests {
				if r2 == pairs[0] {
					removed[j] = true
				}
			}

			matches = append(matches, match{r1, pairs[0]})
		}
	}

	return matches, nil
}

func Start() {
	db := makeDb()
	pool := DbPool{
		repo: repository.NewRepository(db),

	}

	log.Printf("starting matchmaker")
	tickTimeStr := os.Getenv("TENGEN_MATCHMAKE_TICK_TIME_MS")
	tickTimeMs, err := strconv.Atoi(tickTimeStr)
	tickTime := time.Duration(tickTimeMs) * time.Millisecond

	if err != nil {
		log.Fatalf("Could not parse TENGEN_MATCHMAKE_TICK_TIME_MS")
	}

	matchmaker := newMatchmaker(pool, tickTime)
	matchmaker.run()
}

func abs(n int) int {
	if n < 0 {
		return -1 * n
	}
	return n
}

func makeDb() *sqlx.DB {
	port, err := strconv.Atoi(os.Getenv("TENGEN_DB_PORT"))
	if err != nil {
		log.Fatal("Cold not parse TENGEN_DB_PORT")
	}

	config := db.PostgresDBConfig{
		Host:     os.Getenv("TENGEN_DB_HOST"),
		Port:     port,
		User:     os.Getenv("TENGEN_DB_USER"),
		Database: os.Getenv("TENGEN_DB_DATABASE"),
		Password: os.Getenv("TENGEN_DB_PASSWORD"),
	}

	db, err := db.NewPostgresDb(config)
	if err != nil {
		log.Fatal("Unable to connect to DB.", err)
	}

	return db
}

