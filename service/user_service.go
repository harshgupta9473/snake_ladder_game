package service

import (
	"log"
	"snake_ladder/intf"
	"snake_ladder/packets"
	"snake_ladder/transport"
)

type UserService struct {
	UserRepo intf.UserRepositoryIntf
}

func NewUserService(u intf.UserRepositoryIntf) intf.UserServiceIntf {
	return &UserService{
		UserRepo: u,
	}
}

func (s *UserService) Connect(userID string, name string, conn transport.Connection) {
	s.UserRepo.Connect(userID, name, conn)
}

func (s *UserService) Disconnect(userID string) {
	s.UserRepo.Disconnect(userID)
}

func (s *UserService) SendMessageToUser(userID string, msg *packets.PacketResponse) {
	user, err := s.UserRepo.GetUser(userID)
	if err != nil {
		log.Fatal("unable to send msg  to user")
	}
	user.Conn.WriteMsg(msg)
}
