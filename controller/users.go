package controllers

import (
	"encoding/json"
	"fmt"
	"log"
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

func WebsocketHandler(svc intf.UserServiceIntf, mm intf.MatchMakingServiceIntf, gs intf.GameServiceIntf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UserConnRequest
		userID := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")

		if userID == "" || name == "" {
			http.Error(w, "wrong format", http.StatusBadRequest)
			log.Println("Missing id or name")
			return
		}
		req.ID = userID
		req.Name = name
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading: ", err)
			log.Println(err)
			return
		}
		// defer ws.Close()

		conn := transport.NewConnection(ws)

		svc.Connect(req.ID, req.Name, conn)
		conn.Start()
		status:=mm.AnyPreviousMatch(userID)
		if(status!=nil){
			gs.BroadCastGameUpdate(status.GameID,status,"game_played")
		}

		go func() {

			for {
				select {
				case msg := <-conn.ReadMsg():

					switch msg.Header.RequestType {
					case "jg":
						JoinGameHandler(msg, mm, gs)
					case "pt":
						PlayGameHandler(msg, gs)
					}
					// case: for done channel or graceful shutdown
				}
			}
		}()

		go func(userID string){
			for{
				select{
				case <-conn.ReadDisconnect():
					log.Println("disconnection signal received for ",userID)
					 WaitForReconnectionHandler(userID,gs)
				}
			}
		}(userID )
	}
}

func JoinGameHandler(packet *packets.Packet, mm intf.MatchMakingServiceIntf, gs intf.GameServiceIntf) {

	payloadBytes, err := json.Marshal(packet.Payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	var temp map[string]interface{}
	err = json.Unmarshal(payloadBytes, &temp)
	if err != nil {
		fmt.Println("Error unmarshaling payload to map:", err)
		return
	}

	if _, ok := temp["dice_type"]; !ok {
		fmt.Println("DiceType is missing in payload")
		return
	}

	var join packets.JoinGame
	err = json.Unmarshal(payloadBytes, &join)
	if err != nil {
		fmt.Println("Error unmarshaling to JoinGame:", err)
		return
	}

	joined, status := mm.StartMatch(packet.Header.UserId, join.DiceType)
	if joined {
		gs.BroadCastGameUpdate(status.GameID, status, "game_started")
	}
}

func PlayGameHandler(packet *packets.Packet, gs intf.GameServiceIntf) {
	payloadBytes, err := json.Marshal(packet.Payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	var temp map[string]interface{}
	err = json.Unmarshal(payloadBytes, &temp)
	if err != nil {
		fmt.Println("Error unmarshaling payload to map:", err)
		return
	}

	if _, ok := temp["game_id"]; !ok {
		fmt.Println("GameID is missing in payload")
		return
	}

	var playgame packets.PlayTurn
	err = json.Unmarshal(payloadBytes, &playgame)
	if err != nil {
		fmt.Println("Error unmarshaling to JoinGame:", err)
		return
	}

	status := gs.PlayTurn(playgame.GameID, packet.Header.UserId)
	if(status==nil){
		return
	}
	gs.BroadCastGameUpdate(status.GameID, status, "game_played")
}

func WaitForReconnectionHandler(userID string,gs intf.GameServiceIntf){
	ok:=gs.WaitAndCheckConnection(userID)
	if ok{
		log.Println("reconnected")
	}
}


// {
//     "header":{
//         "user_id":"harsh9473",
//         "request_type":"jg"
//     },
//     "payload":{
//         "dice_type":0
//     }
// }