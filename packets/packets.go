package packets

type Packet struct {
	Header struct {
		UserId     string `json: "user_id"`
		RequestType string `json:"request_type"`
	} `json:"header"`
	Payload interface{} `json:"payload"`
}

type JoinGame struct {
	DiceType int `json:"dice_type"`
}

type PlayTurn struct {
	GameID string  `json:"game_id"`
}

type Players struct{
	UserID   string   `json:"user_id"`
		Name     string  `json:"name"`
		Location int   `json:"location"`
}

type UpdatePayloadGameStatus struct {
	GameID string `json:"game_id"`
	Start   bool `json:"start"`
	End     bool `json:"end"`
	Running bool `json:"running"`
	UserTurn    string `json:"user_turn"`
	WonBy   string
	Players []Players   `json:"players"`
}




type PacketResponse struct{
	Header struct{
		UserId string `json:"user_id"`
		PacketType string `json:"packet_type"`
	}  `json:"json"`
	Payload interface{} `json:"payload"`
}