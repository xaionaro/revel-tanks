package models

import (
	  "encoding/json"
	  "github.com/revel/revel"
	  "golang.org/x/net/websocket"
	  "github.com/joonazan/vec2"
	  "sync/atomic"
)

type playerMessage struct {
	Action	string
	Args	map[string]interface{}
}

type player struct {
	id		 int64
	nickname	 string
	websocket	*websocket.Conn
	tank		*tank
	playground	*Playground
}

var playerCounter int64

func (player *player) GetId() int64 {
	return player.id
}

func (player *player) GetNickname() string {
	return player.nickname
}

func (player *player) SetPlayground(playground *Playground) {
	player.playground = playground
}

func (player *player) GetPlayground() *Playground {
	return player.playground
}

func (player *player) CreateTank() {
	if (player.tank != nil) {
		player.DestroyTank()
	}
	player.tank = NewTank(player);
}

func (player *player) GetTank() *tank {
	return player.tank
}

func (player *player) DestroyTank() {
	if (player.tank != nil) {
		player.tank.Destroy()
		player.tank = nil
	}
}

func (player *player) Destroy() {
	revel.TRACE.Printf("player.Destroy()")
	player.DestroyTank()
	player.playground.RemovePlayer(player)
}

func (player *player) Notify(message map[string]interface{}) {

	// Serializing the message

	marshalized,err := json.Marshal(message)
	if (err != nil) {
		revel.ERROR.Printf("player.Notify(): Cannot marshalize message \"%v\": %v", message, err.Error())
		return
	}

	// Sending the message

	_,err = player.websocket.Write(marshalized)
	if (err != nil) {
		revel.ERROR.Printf("player.Notify(): Cannot send message \"%v\": %v", message, err.Error())
		player.Destroy()
	}
}

func (player *player) considerMessage(playerMessage playerMessage) {
	revel.TRACE.Printf("player.considerMessage(%v)", playerMessage)
	switch (playerMessage.Action) {
		case "fire":
			if (player.tank == nil) {
				player.CreateTank()
			}
			player.tank.Fire()
			break
		case "setDirection":
			if (player.tank == nil) {
				player.CreateTank()
			}

			directionData := playerMessage.Args["direction"].(map[string]interface{})

			direction     := vec2.Vector{
					X: directionData["x"].(float64),
					Y: directionData["y"].(float64),
				}
			player.tank.SetDirection(direction)
			break
	}
}

func (player *player) Play() {
	revel.TRACE.Printf("player.Play()")

	for {
		var playerMessageBytes []byte
		{	// Receiving a massage
			err := websocket.Message.Receive(player.websocket, &playerMessageBytes)
			if (err != nil) {
				revel.TRACE.Printf("player.Play(): WS is closed: %v", err.Error())
				return
			}
		}
		{	// Unseriazlizing the message
			var playerMessage playerMessage
			err := json.Unmarshal(playerMessageBytes, &playerMessage)
			if (err != nil) {
				revel.TRACE.Printf("player.Play(): Cannot unmarshal the message \"%v\": %v", string(playerMessageBytes), err.Error())
				return
			}

			// Considering the message
			player.considerMessage(playerMessage)
		}
	}
}

func NewPlayer(ws *websocket.Conn, nickname string) *player {
	return &player{id: atomic.AddInt64(&playerCounter, 1), websocket: ws, nickname: nickname}
}


