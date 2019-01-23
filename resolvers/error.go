package resolvers

type userNotInGameError struct{}
type wrongTurnError struct{}
type gameCompleteError struct{}
type gameNotStartedError struct{}
type koViolationError struct{}
type stoneExistsError struct{}
type sameUserError struct{}

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

func (e koViolationError) Error() string {
	return "Invalid Ko move"
}

func (e stoneExistsError) Error() string {
	return "Stone already exists at that location"
}

func (e sameUserError) Error() string {
	return "Cannot choose self as opponent"
}
