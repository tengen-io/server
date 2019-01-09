package main

type User struct {
	Id                int    `json:id`
	Username          string `json:username`
	Email             string `json:email`
	encryptedPassword string
	Timestamps
}

type AuthUser struct {
	Jwt  string `json:token`
	User *User  `json:user`
}
