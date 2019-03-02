package models

import "time"

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Identity struct {
	Id    int  `json:"id",db:"identity_id"`
	Email string `json:"email"`
	User
}

type User struct {
	Id int `json:"id",db:"user_id"`
	Name       string `json:"name"`
}

