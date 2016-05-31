package controllers

import (
	  "github.com/revel/revel"
	  "golang.org/x/net/websocket"
	a "tanks/app"
	m "tanks/app/models"
)

type Game struct {
	*revel.Controller
}

func (c Game) MainWindow() revel.Result {
	return c.Render()
}

func (c Game) WebSocket(ws *websocket.Conn) revel.Result {
	player := m.NewPlayer(ws, "unnamed");

	a.Playground.AddPlayer(player)
	player.Play()
	player.Destroy();

	return c.Render()
}
