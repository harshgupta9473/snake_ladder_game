package repository

import (
	"fmt"
	"snake_ladder/intf"
	"snake_ladder/models"
	"snake_ladder/transport"
)



type UserRepository struct {
    users map[string]models.User
}

func NewUserRepo()intf.UserRepositoryIntf{
	return &UserRepository{
		users: make(map[string]models.User),
	}
}


func (u *UserRepository) Connect(userID string,name string,conn transport.Connection){
	u.users[userID]=models.User{ID: userID,Name:name,Conn:conn }
}

func (u *UserRepository)Disconnect(userID string){
	delete(u.users,userID)
}



func (u *UserRepository)GetUser(userID string)(*models.User,error){
	user,ok:=u.users[userID]
	if !ok{
		return nil,fmt.Errorf("user not found")
	}
	return &user,nil
}