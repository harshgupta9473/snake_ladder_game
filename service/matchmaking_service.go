package service

import (
	"snake_ladder/intf"
	"snake_ladder/packets"
)

type MatchMakingService struct {
	GameService intf.GameServiceIntf
	waitingUserIDForDiceType1 string
	waitingUserIDForDiceType2 string
}



func NewMatchMakingService(gs intf.GameServiceIntf)*MatchMakingService{
	return &MatchMakingService{
		GameService: gs,
		waitingUserIDForDiceType1: "",
		waitingUserIDForDiceType2: "",
	}
}

func (mm *MatchMakingService)StartMatch(userID string,dicetype int)(bool,*packets.UpdatePayloadGameStatus){
	if(dicetype==0){
		if(mm.waitingUserIDForDiceType1==""){
			mm.waitingUserIDForDiceType1=userID
			return false,nil
		//return waiting for more player to join
	}else{
		status:=mm.GameService.CreateandJoin(mm.waitingUserIDForDiceType1,userID,0)
		mm.waitingUserIDForDiceType1=""
		return true,status
	}
	}else{
		if(mm.waitingUserIDForDiceType2==""){

		}else{
			stauts:=mm.GameService.CreateandJoin(mm.waitingUserIDForDiceType2,userID,1)
			mm.waitingUserIDForDiceType2=""
			return true,stauts
		}
	}

	
	return false ,nil
}

