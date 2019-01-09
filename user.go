package main

type User struct {
	Id                int    `json:id`
	Username          string `json:username`
	Email             string `json:email`
	encryptedPassword string
	Games             *[]Game
	Timestamps
}

type AuthUser struct {
	Jwt  string `json:token`
	User *User  `json:user`
}
