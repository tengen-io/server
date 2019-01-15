package resolvers

type userNotInGameError struct{}
type wrongTurnError struct{}
type gameCompleteError struct{}
type gameNotStartedError struct{}

func (e userNotInGameError) Error() string {
    return "User not in game"
}

func (e wrongTurnError) Error() string {
    return "User must wait for turn"
}

func (e gameCompleteError) Error() string {
    return "Game has already finished"
}

func (e gameNotStartedError) Error() string {
    return "Game has not yet started"
}
