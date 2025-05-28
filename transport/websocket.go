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
	ReadDisconnect() chan struct{}
}

type WSConnection struct {
	conn         *websocket.Conn
	readchan     chan *packets.Packet
	writechan    chan *packets.PacketResponse
	quit         chan struct{}
	disconnected chan struct{}
}

func NewConnection(conn *websocket.Conn) Connection {
	return &WSConnection{
		conn:         conn,
		readchan:     make(chan *packets.Packet),
		writechan:    make(chan *packets.PacketResponse),
		quit:         make(chan struct{}),
		disconnected: make(chan struct{}),
	}
}

func (ws *WSConnection) readLoop() {
	defer ws.Close()
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			ws.disconnected <- struct{}{}
			return
		}
		if len(msg) == 0 {
			log.Println("empty message")
			continue
		}
		var packet *packets.Packet
		err = json.Unmarshal(msg, &packet)
		if err != nil {
			log.Println("invalid json packet", err)
			continue
		}

		ws.readchan <- packet
	}
}

func (ws *WSConnection) writeLoop() {
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
				ws.disconnected <- struct{}{}
				return
			}
		case <-ws.quit:
			log.Println("Quit signal received, closing write loop")
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

func (ws *WSConnection) ReadDisconnect() chan struct{} {
	return ws.disconnected
}
