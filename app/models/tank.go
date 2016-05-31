package models

import (
	"math/rand"
	"github.com/joonazan/vec2"
	"github.com/revel/revel"
)

type tank struct {
	mapObject
	velocity float64
}

func (tank *tank) Destroy() {
	tank.player.tank = nil
	tank.mapObject.Destroy()
}

func (tank *tank) SetDirection(direction vec2.Vector) {
	if (direction.Length() == 0) {
		tank.momentum = vec2.Vector{ 0, 0 }
		return
	}
	revel.TRACE.Printf("tank: %v", tank)
	tank.momentum = direction.Normalized()
	tank.momentum.Mul(tank.velocity)
}

func (tank *tank) Fire() {
	bullet := NewBullet(tank.player)
	bulletMomentum := tank.momentum
	bulletMomentum.Mul(5)
	bullet.SetMomentum(bulletMomentum)
	bullet.SetPos(tank.pos);
}

func NewTank(player *player) *tank {
	revel.TRACE.Printf("NewTank()")
	tank := &tank{velocity: 4}
	tank.initMapObject(tank, player)
	tank.SetPos(vec2.Vector{ rand.Float64()*1000, rand.Float64()*600 })
	return tank
}

func (tank *tank) OutBoundary() {
	if (tank.pos.X < 0) {
		tank.pos.X = 1000
	}
	if (tank.pos.X > 1000) {
		tank.pos.X = 0
	}
	if (tank.pos.Y < 0) {
		tank.pos.Y = 600
	}
	if (tank.pos.Y > 600) {
		tank.pos.Y = 0
	}
}

