package repository

import (
	"github.com/lib/pq"
	"github.com/tengen-io/server/models"
	"strconv"
	"time"
)

func (r *Repository) GetMatchmakingRequests() ([]models.MatchmakingRequest, error) {
	rows, err := r.h.Query("SELECT id, user_id, rank, rank_delta, created_at, updated_at FROM matchmake_requests")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	requests := make([]models.MatchmakingRequest, 0)
	for rows.Next() {
		var request models.MatchmakingRequest
		err := rows.Scan(&request.Id, &request.User.Id, &request.Rank, &request.Delta, &request.CreatedAt, &request.UpdatedAt)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *Repository) GetMatchmakingRequestsForUser(user models.User) ([]models.MatchmakingRequest, error) {
	rows, err := r.h.Query("SELECT id, user_id, rank, rank_delta, created_at, updated_at FROM matchmake_requests WHERE user_id = $1", user.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	requests := make([]models.MatchmakingRequest, 0)
	for rows.Next() {
		var request models.MatchmakingRequest
		err := rows.Scan(&request.Id, &request.User.Id, &request.Rank, &request.Delta, &request.CreatedAt, &request.UpdatedAt)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *Repository) CreateMatchmakingRequest(user models.User, delta int) (*models.MatchmakingRequest, error) {
	now := time.Now().UTC()
	ts := pq.FormatTimestamp(now)

	// TODO(eac): add real ranks and queues
	row := r.h.QueryRowx("INSERT INTO matchmake_requests (queue, user_id, rank, rank_delta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", "FIXME", user.Id, 10, delta, ts, ts)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.MatchmakingRequest{
		NodeFields: models.NodeFields{
			Id: strconv.FormatInt(id, 10),
			CreatedAt: now,
		},
		User: user,
		Delta: delta,
		Rank: 10,
		Queue: "FIXME",
	}, nil

}

func (r *Repository) DeleteMatchmakingRequest(request models.MatchmakingRequest) error {
	_, err := r.h.Exec("DELETE FROM matchmake_requests WHERE id = $1", request.Id)
	return err
}
