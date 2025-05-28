package service

import (
	"log"
	"snake_ladder/intf"
	"snake_ladder/packets"
	"snake_ladder/utils"
	"time"

	"github.com/google/uuid"
)

type GameService struct {
	GameRepo    intf.GameRepositoryIntf
	UserService intf.UserServiceIntf
}

func NewGameService(gameRepo intf.GameRepositoryIntf, userServiec intf.UserServiceIntf) intf.GameServiceIntf {
	return &GameService{
		GameRepo:    gameRepo,
		UserService: userServiec,
	}
}

// func (gs *GameService) CreateGame(userID string, dicetype int) *packets.UpdatePayloadGameStatus {
// 	gameID := uuid.New().String()
// 	gs.GameRepo.CreateGame(gameID, userID, dicetype)
// 	return gs.gameStatusPlayload(gameID)
// }

// func (gs *GameService) JoinGameByGameID(gameID string, userID string) *packets.UpdatePayloadGameStatus {
// 	gs.GameRepo.JoinGameByGameID(gameID, userID)
// 	status := gs.gameStatusPlayload(gameID)
// 	log.Println(status)
// 	return status
// }

func (gs *GameService) PlayTurn(gameID string, userID string) *packets.UpdatePayloadGameStatus {
	
	played := gs.GameRepo.PlayTurn(gameID, userID)
	if !played {
		return nil
	}
	status := gs.gameStatusPlayload(gameID)
	log.Println(status)
	return status
}

func (gs *GameService) BroadCastGameUpdate(gameID string, payload interface{}, packet_type string) {

	game := gs.GameRepo.GetGame(gameID)
	for _, val := range game.Players {
		msg := utils.MakePacket(val.UserID, packet_type, payload)
		gs.UserService.SendMessageToUser(val.UserID, msg)
	}
}

func (gs *GameService) CreateandJoin(userID1 string, userID2 string, dicetype int) *packets.UpdatePayloadGameStatus {
	gameID := uuid.New().String()
	conn1, err1 := gs.UserService.GetUserConn(userID1)
	conn2, err2 := gs.UserService.GetUserConn(userID2)
	if err1 != nil || err2 != nil {
		return nil
	}

	conn1.ReadDisconnect()
	gs.GameRepo.CreateandJoinTwoPlayer(userID1, userID2, gameID, dicetype)
	go gs.waitAndCheckConnection(gameID ,userID1,conn1.ReadDisconnect())
	go gs.waitAndCheckConnection(gameID,userID2,conn2.ReadDisconnect())
	gameStatus := gs.gameStatusPlayload(gameID)
	return gameStatus
}

func (gs *GameService) gameStatusPlayload(gameID string) *packets.UpdatePayloadGameStatus {
	games := gs.GameRepo.GetGame(gameID)
	var gameStatus packets.UpdatePayloadGameStatus
	gameStatus.GameID = gameID
	gameStatus.Start = games.Start
	gameStatus.Running = games.Running
	gameStatus.End = games.End
	gameStatus.WonBy = games.WonBy
	gameStatus.UserTurn = games.WhooseTurn
	gameStatus.SnakeAndLadder = games.SnakeAndLadder
	for _, player := range games.Players {
		gameStatus.Players = append(gameStatus.Players, packets.Players{UserID: player.UserID, Name: player.Name, Location: player.Location})
	}
	return &gameStatus
}

func (gs *GameService) EndGame(gameID string)*packets.UpdatePayloadGameStatus {
	gs.GameRepo.LeaveGame(gameID)
	gameStatus := gs.gameStatusPlayload(gameID)
	return gameStatus
}

func (gs *GameService) IfUserIsAlreadyPartOfSomeGameJoinHimThere(userID string) *packets.UpdatePayloadGameStatus {
	ok, gameID := gs.GameRepo.GetGameByUserID(userID)
	if !ok {
		return nil
	}
	if(gs.GameRepo.GetGame(gameID).End){
		return nil
	}

	gs.GameRepo.JoinGameByGameID(gameID, userID)
	gameStatus := gs.gameStatusPlayload(gameID)
	return gameStatus
}

func (gs *GameService) waitAndCheckConnection(gameID string, userID string,disconnect chan struct{}) {
	select {
	case <-disconnect:
		log.Println("disconnection signal received for ",userID)
		gs.GameRepo.GetGame(gameID).PlayerMap[userID].Connected=false
		time.Sleep(30*time.Second)
		if(!gs.GameRepo.GetGame(gameID).PlayerMap[userID].Connected){
			gs.EndGame(gameID)
			return
		}
	}
}
