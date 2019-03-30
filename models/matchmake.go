package models

type MatchmakingRequest struct {
	NodeFields
	Queue  string
	User User
	Rank   int
	Delta  int
}

func (MatchmakingRequest) IsNode() {}
