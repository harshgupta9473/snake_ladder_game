package service

import (
	"log"
	"snake_ladder/intf"
	"snake_ladder/packets"
	"snake_ladder/utils"

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

func (gs *GameService) CreateGame(userID string, dicetype int) *packets.UpdatePayloadGameStatus {
	gameID := uuid.New().String()
	gs.GameRepo.CreateGame(gameID, userID, dicetype)
	return gs.gameStatusPlayload(gameID)
}

func (gs *GameService) JoinGameByGameID(gameID string, userID string) *packets.UpdatePayloadGameStatus {
	gs.GameRepo.JoinGameByGameID(gameID, userID)
	status := gs.gameStatusPlayload(gameID)
	log.Println(status)
	return status
}

func (gs *GameService) PlayTurn(gameID string, userID string) *packets.UpdatePayloadGameStatus {
	gs.GameRepo.PlayTurn(gameID, userID)
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
	gs.GameRepo.CreateandJoinTwoPlayer(userID1, userID2, gameID, dicetype)
	gameStatus := gs.gameStatusPlayload(gameID)
	return gameStatus
}

// JoinGameByGameID(string,string)
// 	PlayTurn(string,string)

func (gs *GameService) gameStatusPlayload(gameID string) *packets.UpdatePayloadGameStatus {
	games := gs.GameRepo.GetGame(gameID)
	var gameStatus packets.UpdatePayloadGameStatus
	gameStatus.GameID = gameID
	gameStatus.Start = games.Start
	gameStatus.Running = games.Running
	gameStatus.End = games.End
	gameStatus.WonBy = games.WonBy
	gameStatus.UserTurn = games.WhooseTurn
	for _, player := range games.Players {
		gameStatus.Players = append(gameStatus.Players, packets.Players{UserID: player.UserID, Name: player.Name, Location: player.Location})
	}
	return &gameStatus
}
