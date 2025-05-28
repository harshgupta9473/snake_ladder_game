package intf

import (
	"snake_ladder/models"
	"snake_ladder/packets"
	"snake_ladder/transport"
)

type UserRepositoryIntf interface {
	Connect(string, string, transport.Connection)
	Disconnect(string)
	GetUser(userID string) (*models.User, error)
}

type UserServiceIntf interface {
	Connect(userID string, name string, conn transport.Connection)
	Disconnect(string)
	SendMessageToUser(userID string, msg *packets.PacketResponse)
	GetUserConn(userID string)(transport.Connection,error)
}

type GameRepositoryIntf interface {
	PlayTurn(gameID string, userID string) bool
	GetGame(gameID string) *models.Game
	CreateandJoinTwoPlayer(userID1 string, userID2 string, gameID string, dicetype int) 
	GetGameByUserID(useID string) (bool, string)
	JoinGameByGameID(gameID string, userID string) bool
	LeaveGame(gameID string)
}

type GameServiceIntf interface {
	PlayTurn(gameID string, userID string) *packets.UpdatePayloadGameStatus
	CreateandJoin(userID1 string, userID2 string, dicetype int) *packets.UpdatePayloadGameStatus 
	BroadCastGameUpdate(gameID string, payload interface{}, packet_type string)
	IfUserIsAlreadyPartOfSomeGameJoinHimThere(userID string) *packets.UpdatePayloadGameStatus
	EndGame(gameID string)*packets.UpdatePayloadGameStatus
}

type MatchMakingServiceIntf interface {
	StartMatch(userID string, dicetype int) (bool, *packets.UpdatePayloadGameStatus)
	AnyPreviousMatch(userID string)(*packets.UpdatePayloadGameStatus)
}
