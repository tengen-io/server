package main

func HandleError(e error) error {
	if e == nil {
		return e
	}

	switch e.Error() {
	case "pq: duplicate key value violates unique constraint \"users_username_key\"":
		return usernameTakenError{e}
	case "pq: duplicate key value violates unique constraint \"users_email_key\"":
		return emailTakenError{e}
	case "crypto/bcrypt: hashedPassword is not the hash of the given password":
		return invalidLoginError{e}
	case "pq: insert or update on table \"players\" violates foreign key constraint \"players_user_id_fkey\"":
		return userNotFoundError{}
	default:
		return e
	}
}

type passwordMismatchError struct {
	Err error
}
type usernameTakenError struct {
	Err error
}
type emailTakenError struct {
	Err error
}
type invalidEmailError struct {
	Err error
}
type invalidLoginError struct {
	Err error
}
type userNotFoundError struct{}
type gameNotFoundError struct{}
type userNotInGameError struct{}
type wrongTurnError struct{}

func (e passwordMismatchError) Error() string {
	return "Passwords do not match"
}

func (e usernameTakenError) Error() string {
	return "Username is already taken"
}

func (e emailTakenError) Error() string {
	return "Email is already taken"
}

func (e invalidEmailError) Error() string {
	return "Email format is invalid"
}

func (e invalidLoginError) Error() string {
	return "Invalid login"
}

func (e userNotFoundError) Error() string {
	return "User not found"
}

func (e gameNotFoundError) Error() string {
	return "Game not found"
}

func (e userNotInGameError) Error() string {
	return "User not in game"
}

func (e wrongTurnError) Error() string {
	return "User must wait for turn"
}
