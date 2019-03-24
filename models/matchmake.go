package models

type MatchmakingRequest struct {
	NodeFields
	Queue  string
	UserId int64
	Rank   int
	Delta  int
}
