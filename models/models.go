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
	Index    int
	Location int
	Name     string
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
	Ended          bool
	DiceType       int // 0 for normal and 1 for crooked
	SnakeAndLadder map[int]int
}
