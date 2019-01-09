package main

type Game struct {
	Id           int    `json:id`
	Status       string `json:status`
	PlayerTurnId int    `json:playerTurnId`
	Players      *[]Player
	Timestamps
}
