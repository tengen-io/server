package main

import (
	"time"
)

type Player struct {
	Id        int    `json:id`
	userId    int    `json:userId`
	gameId    int    `json:gameId`
	Status    string `json:status`
	Color     string `json:color`
	Stats     Stats  `json:stats`
	HasPassed bool   `json:hasPassed`
	User      *User
	Game      *Game
	Timestamps
}

type Timestamps struct {
	InsertedAt time.Time `json:insertedAt`
	UpdatedAt  time.Time `json:updatedAt`
}

type Stats struct{}
