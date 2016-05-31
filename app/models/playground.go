package models

import (
	"time"
	"sync"
	"github.com/revel/revel"
)

type Playground struct {
	players		[]*player
	mapObjects	[]MapObjectI
	runned		 bool
	mutex		*sync.Mutex
}

// GetMapObjects

func (p *Playground) GetMapObjects() ([]MapObjectI) {
	return p.mapObjects
}

// forRender

func (p *Playground) forRender() (map[string]interface{}) {
	result  := make(map[string]interface{})

	//revel.TRACE.Printf("p.mapObjects == %v", p.mapObjects)

	var objects []map[string]interface{}
	for _,mapObject := range p.mapObjects {
		object  := map[string]interface{}{
			"type":		mapObject.GetTypeName(),
			"playerId":	mapObject.GetPlayer().GetId(),
			"pos":		mapObject.GetPos(),
			"direction":	mapObject.GetMomentum(),
		}

		objects  = append(objects, object)
	}

	result["objects"] = objects

	return result
}


// iterate

func (p *Playground) iterate() {
	defer func(){p.Unlock()}()
	p.Lock()

	for _,mapObject := range p.mapObjects {
		mapObject.Iterate()
	}
	playgroundStatus := map[string]interface{}{"playground": p.forRender()}
	for _,player    := range p.players {
		playgroundStatus["playerId"] = player.GetId()
		player.Notify(playgroundStatus)
	}
}

// Start/Stop

func (p *Playground) Start() {
	p.runned = true
	p.mutex  = &sync.Mutex{}
	go func() {
		for p.runned {
			go p.iterate()
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (p *Playground) Stop() {
	p.runned = false
}

// Lock/Unlock

func (p *Playground) Lock() {
	p.mutex.Lock();
}
func (p *Playground) Unlock() {
	p.mutex.Unlock();
}

// players

func (p *Playground) AddPlayer(player *player) {
	defer func(){p.Unlock()}()
	p.Lock()

	p.players = append(p.players, player)

	player.SetPlayground(p)
}

func (p *Playground) RemovePlayer(removingPlayer *player) {
	defer func(){p.Unlock()}()
	p.Lock()

	for i,player := range p.players {
		if (player == removingPlayer) {
			p.players = append(p.players[:i], p.players[i+1:]...)
			return
		}
	}
}

// mapObjects

func (p *Playground) AddMapObject(mapObject MapObjectI) {
	defer func(){p.Unlock()}()
	p.Lock()
	revel.TRACE.Printf("AddMapObject()")

	mapObject.SetPlayground(p)
	p.mapObjects = append(p.mapObjects, mapObject)
}

func (p *Playground) RemoveMapObject(removingMapObject MapObjectI) {
	defer func(){p.Unlock()}()
	p.Lock()

	for i,mapObject := range p.mapObjects {
		if (mapObject == removingMapObject) {
			p.mapObjects = append(p.mapObjects[:i], p.mapObjects[i+1:]...)
			return
		}
	}
}

