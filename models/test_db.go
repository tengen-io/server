package models

type TestDB struct{}

func (db *TestDB) CheckPw(username, password string) (*User, error) {
	return nil, nil
}

func (db *TestDB) CreateUser(username, email, password, passwordConfirm string) (*User, error) {
	user := &User{
		Id:       1,
		Username: username,
		Email:    email,
	}
	return user, nil
}

func (db *TestDB) GetUser(identifier interface{}) (*User, error) {
	user := &User{
		Id:       1,
		Username: "dude",
		Email:    "dude@dude.dude",
	}

	switch identifier.(type) {
	case int:
		if identifier.(int) == 1 {
			return user, userNotFoundError{}
		}
	case string:
		if identifier.(string) == "dude" {
			return user, userNotFoundError{}
		}
	}

	return nil, userNotFoundError{}
}

func (db *TestDB) GetGame(gameId interface{}) (*Game, error) {
	return nil, nil
}

func (db *TestDB) GetGames(userId interface{}) ([]*Game, error) {
	return nil, nil
}

func (db *TestDB) CreateGame(userId int, opponent *User) (*Game, error) {
	return nil, nil
}

func (db *TestDB) UpdateBoard(userId int, game *Game) (*Game, error) {
	return nil, nil
}

func (db *TestDB) Pass(userId int, game *Game) (*Game, error) {
	return nil, nil
}
