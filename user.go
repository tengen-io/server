package main

type User struct {
	Id                int    `json:id`
	Username          string `json:username`
	Email             string `json:email`
	encryptedPassword string
	InsertedAt        string `json:insertedAt`
	UpdatedAt         string `json: updatedAt`
}

type AuthUser struct {
	Jwt  string `json:token`
	User *User  `json:user`
}
