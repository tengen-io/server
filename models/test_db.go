package models

type TestDB struct{}

func (db *TestDB) CheckPw(username, password string) (*User, error) {
	if username == "dude" && password == "dudedude" {
		return db.GetUser("dude")
	}
	return nil, invalidLoginError{}
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
		if id, ok := identifier.(int); ok {
			user.Id = id
		}
	case string:
		if id, ok := identifier.(string); ok {
			if user.Username != id {
				user.Id = len(id)
			}
			user.Username = id
		}
	}

	return user, nil
}

func (db *TestDB) GetGame(gameId interface{}) (*Game, error) {
	game := &Game{
		Id:           1,
		Status:       "active",
		PlayerTurnId: 1,
		BoardSize:    SmallBoardSize,
	}

	if gameId.(string) == "2" {
		game.Status = "complete"
	}

	if gameId.(string) == "3" {
		game.PlayerTurnId++
		game.Players = []Player{
			Player{Id: 1, UserId: 1},
			Player{UserId: 2},
		}
	}

	if gameId.(string) == "4" {
		game.Players = []Player{
			Player{Id: 1, UserId: 1},
			Player{UserId: 2},
		}
	}

	return game, nil
}

func (db *TestDB) GetGames(userId interface{}) ([]*Game, error) {
	game := &Game{
		Id:           1,
		Status:       "active",
		PlayerTurnId: 1,
	}
	return []*Game{game}, nil
}

func (db *TestDB) CreateGame(userId int, opponent *User, size int) (*Game, error) {
	return &Game{
		Id:           1,
		Status:       "active",
		PlayerTurnId: userId,
		BoardSize:    size,
	}, nil
}

func (db *TestDB) UpdateGame(userId int, game *Game, stone Stone, toRemove []Stone) error {
	game.Stones = []Stone{stone}

	return nil
}

func (db *TestDB) Pass(userId int, game *Game) (*Game, error) {
	return nil, nil
}
