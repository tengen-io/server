package models

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("with invalid email", createUserInvalidEmail)
	t.Run("with short password", createUserInvalidPassword)
	t.Run("with mismatched password", createUserPasswordMismatch)
	t.Run("with proper arguments", createUserSucceed)
	t.Run("with existing username", createUserUsernameTaken)
	t.Run("with existing email", createUserEmailTaken)

	teardown()
}

func createUserSucceed(t *testing.T) {
	_, err := db.CreateUser("dude", "dude@dude.dude", "dudedude", "dudedude")

	if err != nil {
		t.Error(err)
	}
}

func createUserInvalidEmail(t *testing.T) {
	_, err := db.CreateUser("dude", "bad-email", "dudedude", "dudedude")

	expectErr(t, invalidEmailError{}, err)
}

func createUserInvalidPassword(t *testing.T) {
	_, err := db.CreateUser("dude", "dude@dude.dude", "shortpw", "shortpw")

	expectErr(t, passwordTooShortError{}, err)
}

func createUserPasswordMismatch(t *testing.T) {
	_, err := db.CreateUser("dude", "dude@dude.dude", "dudefirst", "dudeother")

	expectErr(t, passwordMismatchError{}, err)
}

func createUserUsernameTaken(t *testing.T) {
	_, err := db.CreateUser("dude", "dude@dude.dude", "dudedude", "dudedude")

	expectErr(t, usernameTakenError{}, err)
}

func createUserEmailTaken(t *testing.T) {
	_, err := db.CreateUser("duder", "dude@dude.dude", "dudedude", "dudedude")

	expectErr(t, emailTakenError{}, err)
}

func TestGetUser(t *testing.T) {
	t.Run("with id", getUserById)
	t.Run("with username", getUserByUsername)
	t.Run("not found", getUserNotFound)

	teardown()
}

func getUserById(t *testing.T) {
	user, _ := db.CreateUser("getuser", "getuser@dude.dude", "dudedude", "dudedude")
	_, err := db.GetUser(user.Id)

	if err != nil {
		t.Error(err)
	}
}

func getUserByUsername(t *testing.T) {
	_, err := db.GetUser("getuser")

	if err != nil {
		t.Error(err)
	}
}

func getUserNotFound(t *testing.T) {
	_, err := db.GetUser("notfound")

	expectErr(t, userNotFoundError{}, err)
}
