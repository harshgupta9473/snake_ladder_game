package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snake_ladder/intf"
	"snake_ladder/models"
	"snake_ladder/packets"

	"snake_ladder/transport"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketHandler(svc intf.UserServiceIntf, mm intf.MatchMakingServiceIntf,gs intf.GameServiceIntf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UserConnRequest
		err:=json.NewDecoder(r.Body).Decode(&req)
		if err!=nil{
			http.Error(w,"wrong format",http.StatusBadRequest)
			return
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading: ", err)
			return
		}
		defer ws.Close()
		
		
		conn := transport.NewConnection(ws)
		
		svc.Connect(req.ID, req.Name, conn)
		conn.Start()

		go func() {

			for {
				select {
				case msg := <-conn.ReadMsg():

					switch msg.Header.RequestType {
					case "jg":
						JoinGameHandler(msg,mm,gs)
					case "pt":
						PlayGameHandler(msg,gs)	
					}
					// case: for done channel or graceful shutdown
				}
			}
		}()
	}
}

func JoinGameHandler(packet *packets.Packet,mm intf.MatchMakingServiceIntf,gs intf.GameServiceIntf) {
	joined,status:=mm.StartMatch(packet.Header.UserId, packet.Payload.(packets.JoinGame).DiceType)
	if joined{
		gs.BroadCastGameUpdate(status.GameID,status,"game_started")
	}
}

func PlayGameHandler(packet *packets.Packet,gs intf.GameServiceIntf){
	status:=gs.PlayTurn(packet.Payload.(packets.PlayTurn).GameID,packet.Header.UserId)
	gs.BroadCastGameUpdate(status.GameID,status,"game_played")
}

