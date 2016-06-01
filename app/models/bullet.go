package models

import (
)

type bullet struct {
	mapObject
}

func (bullet *bullet) Iterate() {
	bullet.mapObject.Iterate()
	for _,mapObject := range bullet.playground.GetMapObjects() {
		if (mapObject.GetTypeName() != "tank") {
			continue
		}
		if (mapObject.GetPlayer().GetId() == bullet.player.GetId()) {
			continue
		}
		distance := mapObject.GetPos()
		distance.Sub(bullet.pos)
		if (distance.Length() < 32) {
			go mapObject.GetParent().Destroy()
			go bullet.Destroy()
			break
		}
	}
}

func NewBullet(player *player) *bullet {
	bullet := &bullet{}
	bullet.initMapObject(bullet, player)

	return bullet
}
