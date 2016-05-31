package models

import (
	"github.com/joonazan/vec2"
	"sync"
	"reflect"
	"github.com/revel/revel"
)

type MapObjectI interface {
	SetPlayground(*Playground)
	GetPlayground() *Playground
	GetPos() vec2.Vector
	SetPos(pos vec2.Vector)
	GetMomentum() vec2.Vector
	SetMomentum(momentum vec2.Vector)
	GetTypeName() string
	GetPlayer() *player
	GetParent() MapObjectI
	Iterate()
	Destroy()
	OutBoundary()
}

type mapObject struct {
	parent		 MapObjectI

	playground	*Playground

	player		*player

	pos		 vec2.Vector
	momentum	 vec2.Vector

	mutex		*sync.Mutex

	typeName	string
}

func (mapObject *mapObject) GetParent() MapObjectI {
	return mapObject.parent
}

func (mapObject *mapObject) OutBoundary() {
	go mapObject.parent.Destroy()
}

func (mapObject *mapObject) Iterate() {
	if (mapObject.momentum.Length() == 0) {
		return
	}

	mapObject.pos.Add ( mapObject.momentum )

	if (mapObject.pos.X < 0 || mapObject.pos.Y < 0 || mapObject.pos.X > 1000 || mapObject.pos.Y > 600) {
		mapObject.parent.OutBoundary()
	}
}

func (mapObject *mapObject) GetPos() vec2.Vector {
	return mapObject.pos
}

func (mapObject *mapObject) SetPos(pos vec2.Vector) {
	mapObject.pos = pos
}

func (mapObject *mapObject) GetMomentum() vec2.Vector {
	return mapObject.momentum
}

func (mapObject *mapObject) SetMomentum(momentum vec2.Vector) {
	mapObject.momentum = momentum
}

func (mapObject *mapObject) GetTypeName() string {
	return mapObject.typeName
}

func (mapObject *mapObject) GetPlayer() *player {
	return mapObject.player
}

func (mapObject *mapObject) SetPlayground(playground *Playground) {
	if (mapObject.playground == playground) {
		return
	}
	if (mapObject.playground != nil) {
		mapObject.playground.RemoveMapObject(mapObject);
	}
	mapObject.playground = playground;
}
func (mapObject *mapObject) GetPlayground() *Playground {
	return mapObject.playground
}

func (mapObject *mapObject) initMapObject(parent MapObjectI, player *player) {
	mapObject.player   = player
	mapObject.parent   = parent
	mapObject.typeName = reflect.ValueOf(parent).Elem().Type().Name()

	revel.TRACE.Printf("initMapObject")
	player.GetPlayground().AddMapObject(parent)
}


func (mapObject *mapObject) deinitMapObject() {
	mapObject.playground.RemoveMapObject(mapObject.parent)
}

func (mapObject *mapObject) Destroy() {
	mapObject.deinitMapObject()
}

