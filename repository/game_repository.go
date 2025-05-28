package repository

import (
	"log"
	"math/rand"
	"snake_ladder/intf"
	"snake_ladder/models"
	"time"
)

type GameRepo struct {
	Games map[string]*models.Game
	Board []int
	UserMap map[string]string
}

func NewGameRepo() intf.GameRepositoryIntf {
	return &GameRepo{
		Games: make(map[string]*models.Game),
		Board: createBoard(100),
		UserMap: make(map[string]string),
	}
}

// func (g *GameRepo) CreateGame(gameID string, userID string, dicetype int) {
// 	g.Games[gameID] = &models.Game{
// 		ID:       gameID,
// 		Players:  []*models.Player{},
// 		Turn:     0,
// 		PlayerMap: make(map[string]*models.Player),
// 		End:    false,
// 		DiceType: dicetype,
// 		Start: true,
// 		Running: true,
// 		SnakeAndLadder: g.generateSnakeAndLadder(100),
// 	}
// 	g.addIntoGame(gameID, userID)
// }

func (g *GameRepo) JoinGameByGameID(gameID string, userID string)bool {
	// check if game exists and if ended and if already joined
	if !g.ifGameExists(gameID) {
		//
		return false
	}

	if g.ifGameEnded(gameID) {
		//
		return false
	}

	// join
	g.addIntoGame(gameID, userID)
	return true
}

func (g *GameRepo) PlayTurn(gameID string, userID string) bool {
	// check if game is still live
	if !g.ifGameExists(gameID) {
		log.Println("Game does not exists")
		return false
	}

	if !g.ifUserIsPlayerInGame(gameID, userID) {
		log.Println("user is not the player of the game")
		return false
	}

	if g.ifGameEnded(gameID) {
		log.Println("Game Ended")
		return true
	}

	if g.whooseTurn(gameID) != userID {
		log.Println("Its not your turn")
		return false
	}
	loc := g.playTheGame(gameID, userID)
	if loc == 100 {
		g.Games[gameID].WonBy = userID
		g.Games[gameID].End = true
		g.Games[gameID].Running = false
		return true
	}

	return true
}

func (g *GameRepo) GetGame(gameID string) *models.Game {
	return g.Games[gameID]
}

func (g *GameRepo) CreateandJoinTwoPlayer(userID1 string, userID2 string, gameID string, dicetype int,disconnect1 *chan struct{},disconnect2 *chan struct{}) {
	g.Games[gameID] = &models.Game{
		ID:             gameID,
		Players:        []*models.Player{},
		Start:          true,
		End:            false,
		Running:        true,
		Turn:           0,
		PlayerMap:      make(map[string]*models.Player),
		DiceType:       dicetype,
		SnakeAndLadder: g.generateSnakeAndLadder(100),
	}
	g.Games[gameID].Players = append(g.Games[gameID].Players, &models.Player{UserID: userID1, Location: 0,Disconnected: disconnect1}, &models.Player{UserID: userID2, Location: 0,Disconnected: disconnect2})
	g.Games[gameID].PlayerMap[userID1] = g.Games[gameID].Players[0]
	g.Games[gameID].PlayerMap[userID2] = g.Games[gameID].Players[1]
	g.Games[gameID].WhooseTurn = g.whooseTurn(gameID)
	g.UserMap[userID1]=gameID
	g.UserMap[userID2]=gameID
}

func (g *GameRepo) LeaveGame(gameID string, userID string) {
	g.Games[gameID].End = true
}

func (g *GameRepo)GetGameByUserID(useID string)(bool,string){
	 gameID,ok:=g.UserMap[useID]
	 if !ok{
		return ok,""
	}
	return true,gameID
}

func (g *GameRepo) addIntoGame(gameID string, userID string) {
	g.Games[gameID].Players = append(g.Games[gameID].Players, &models.Player{UserID: userID})
}

func (g *GameRepo) ifGameExists(gameID string) bool {
	log.Println("inside game exists")
	_, ok := g.Games[gameID]
	return ok
}

func (g *GameRepo) ifGameEnded(gameID string) bool {
	log.Println("inside game ended")
	return g.Games[gameID].End
}

func (g *GameRepo) ifUserIsPlayerInGame(gameID string, userID string) bool {
	for _, player := range g.Games[gameID].Players {
		if player.UserID == userID {
			return true
		}
	}
	return false
}

func (g *GameRepo) whooseTurn(gameID string) string {
	log.Println("inside whose turn")
	k := (g.Games[gameID].Turn) % 2
	return g.Games[gameID].Players[k].UserID
	
}

func (g *GameRepo) playTheGame(gameID string, userID string) int {
	rand.Seed(time.Now().UnixNano())
	log.Println("play the game")
	var t int = g.Games[gameID].DiceType
	var dval int
	if t == 0 {
		// normal dice
		dval = rand.Intn(6) + 1
	} else {
		// crooked dice // odd numbered dice
		num := []int{1, 3, 5}
		dval = num[rand.Intn(len(num))]
	}
	g.Games[gameID].Turn = (g.Games[gameID].Turn + 1) % 2
	g.Games[gameID].WhooseTurn = g.whooseTurn(gameID)
	currLoc := g.Games[gameID].PlayerMap[userID].Location
	if currLoc+dval <= 100 {
		g.Games[gameID].PlayerMap[userID].Location = currLoc + dval
		if val, ok := g.Games[gameID].SnakeAndLadder[currLoc+dval]; ok {
			g.Games[gameID].PlayerMap[userID].Location = val
		}
		return g.Games[gameID].PlayerMap[userID].Location
	} else {
		return currLoc
	}
}

func (g *GameRepo) generateSnakeAndLadder(n int) map[int]int {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) {
		g.Board[i], g.Board[j] = g.Board[j], g.Board[i]
	})
	limit := n - (n % 2)
	mp := make(map[int]int)
	final := g.Board[:limit]
	for i := 0; i < limit; i = i + 2 {
		mp[final[i]] = final[i+1]
	}
	return mp
}

func createBoard(n int) []int {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	return nums
}

