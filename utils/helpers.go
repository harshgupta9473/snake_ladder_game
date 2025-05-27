package utils

import "snake_ladder/packets"

func MakePacket(userID string, packet_type string, payload interface{})*packets.PacketResponse{
	var packet packets.PacketResponse
	packet.Header.UserId=userID
	packet.Header.PacketType=packet_type
	packet.Payload=payload
	return &packet
}
