package repository

import (
	"github.com/tengen-io/server/models"
)

func (r *Repository) GetMatchmakingRequests() ([]models.MatchmakingRequest, error) {
	rows, err := r.h.Query("SELECT id, user_id, rank, rank_delta, created_at, updated_at FROM matchmake_pool")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	requests := make([]models.MatchmakingRequest, 0)
	for rows.Next() {
		var request models.MatchmakingRequest
		err := rows.Scan(&request.Id, &request.UserId, &request.Rank, &request.Delta, &request.CreatedAt, &request.UpdatedAt)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *Repository) DeleteMatchmakingRequest(request models.MatchmakingRequest) error {
	_, err := r.h.Exec("DELETE FROM matchmake_pool WHERE id = $1", request.Id)
	return err
}
