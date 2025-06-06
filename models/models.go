package models

import "snake_ladder/transport"

type UserConnRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   string
	Name string
	Conn transport.Connection
}

type Player struct {
	UserID   string
	Location int
	Name     string
	Connected bool
}

type Game struct {
	ID             string
	Players        []*Player
	Start          bool
	Running        bool
	End            bool
	Turn           int
	WonBy          string
	WhooseTurn     string
	PlayerMap      map[string]*Player
	DiceType       int // 0 for normal, 1 for crooked
	SnakeAndLadder map[int]int
}
