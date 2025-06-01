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
	//Logger      logger.ZapLogger
}

func NewGameService(gameRepo intf.GameRepositoryIntf, userServiec intf.UserServiceIntf) intf.GameServiceIntf {
	return &GameService{
		GameRepo:    gameRepo,
		UserService: userServiec,
		//Logger:      logger.NewLogger("service", "game-service"),
	}
}

func (gs *GameService) CreateandJoin(userID1 string, userID2 string, dicetype int) *packets.UpdatePayloadGameStatus {
	gameID := uuid.New().String()
	gs.GameRepo.CreateandJoinTwoPlayer(userID1, userID2, gameID, dicetype)
	gameStatus := gs.gameStatusPlayload(gameID)
	//gs.Logger.LogInfo("game played", zap.String("game status","played"))
	return gameStatus
}

func (gs *GameService) PlayTurn(gameID string, userID string) *packets.UpdatePayloadGameStatus {

	played := gs.GameRepo.PlayTurn(gameID, userID)
	if !played {
		return nil
	}
	status := gs.gameStatusPlayload(gameID)
	log.Println(status)
	//gs.Logger.LogInfo("game played", zap.String("game status","played turn"))
	return status
}

func (gs *GameService) IfUserIsAlreadyPartOfSomeGameJoinHimThere(userID string) *packets.UpdatePayloadGameStatus {
	ok, gameID := gs.GameRepo.GetGameByUserID(userID)
	if !ok {
		return nil
	}
	if gs.GameRepo.GetGame(gameID).End {
		return nil
	}

	gs.GameRepo.JoinGameByGameID(gameID, userID)
	gameStatus := gs.gameStatusPlayload(gameID)
		//gs.Logger.LogInfo("game joined again", zap.String("game status","user joined left game"))
	return gameStatus
}

func (gs *GameService) BroadCastGameUpdate(gameID string, payload interface{}, packet_type string) {

	game := gs.GameRepo.GetGame(gameID)
	for _, val := range game.Players {
		msg := utils.MakePacket(val.UserID, packet_type, payload)
		gs.UserService.SendMessageToUser(val.UserID, msg)
			//.Logger.LogInfo("broadcasting game status", zap.String("user",val.UserID))
	}
}

func (gs *GameService) EndGame(gameID string) *packets.UpdatePayloadGameStatus {
	gs.GameRepo.LeaveGame(gameID)
	gameStatus := gs.gameStatusPlayload(gameID)
	return gameStatus
}

func (gs *GameService) WaitAndCheckConnection(userID string) bool {
	ok, gameID := gs.GameRepo.GetGameByUserID(userID)
	if !ok {
		return false
	}
	if !gs.GameRepo.GetGame(gameID).PlayerMap[userID].Connected {
		return false
	}
	gs.GameRepo.GetGame(gameID).PlayerMap[userID].Connected = false
	time.Sleep(30 * time.Second)
	if !gs.GameRepo.GetGame(gameID).PlayerMap[userID].Connected {
		gs.EndGame(gameID)
		return true
	}
	return true
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
