package transport

import (
	"encoding/json"
	"log"
	"snake_ladder/packets"

	"github.com/gorilla/websocket"
)

type Connection interface {
	Start()
	ReadMsg() <-chan *packets.Packet
	WriteMsg(msg *packets.PacketResponse)
	Close() error
}

type WSConnection struct {
	conn      *websocket.Conn
	readchan  chan *packets.Packet
	writechan chan *packets.PacketResponse
	quit      chan struct{}
}

func NewConnection(conn *websocket.Conn) Connection {
	return &WSConnection{
		conn:      conn,
		readchan:  make(chan *packets.Packet),
		writechan: make(chan *packets.PacketResponse),
		quit:      make(chan struct{}),
	}
}

func (ws *WSConnection) readLoop() {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			//error
			break
		}
		var packet *packets.Packet
		err = json.Unmarshal(msg, &packet)
		if err != nil {
			//
		}

		ws.readchan <- packet
	}
}

func (ws *WSConnection) writeLoop() {
	defer ws.Close()
	for {
		select {
		case msg := <-ws.writechan:
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling the data")
				continue
			}
			err = ws.conn.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Println("Error sending data")
				return
			}
		case <-ws.quit:
			return
		}
	}
}

func (ws *WSConnection) Start() {
	go ws.readLoop()
	go ws.writeLoop()
}

func (ws *WSConnection) ReadMsg() <-chan *packets.Packet {
	return ws.readchan
}

func (ws *WSConnection) WriteMsg(msg *packets.PacketResponse) {
	ws.writechan <- msg
}

func (ws *WSConnection) Close() error {
	close(ws.quit)
	return ws.conn.Close()
}
