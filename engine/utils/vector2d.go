package utils

import (
	"nuvem/engine/coder"
)

type Vector2D struct {
	X int
	Y int
}

func NewVector2D(x, y int) *Vector2D {
	return &Vector2D{
		X: x,
		Y: y,
	}
}

func NewVector2DByMap(m map[string]interface{}) *Vector2D {
	return &Vector2D{
		X: int(ForceUint32("x", m)),
		Y: int(ForceUint32("y", m)),
	}
}

func (v2d *Vector2D) Equals(other *Vector2D) bool {
	if v2d.X == other.X && v2d.Y == other.Y {
		return true
	}
	return false
}

func (v2d *Vector2D) Clone(other *Vector2D) *Vector2D {
	v2d.X = other.X
	v2d.Y = other.Y
	return v2d
}

func (v2d *Vector2D) Incr(other *Vector2D) *Vector2D {
	v2d.X += other.X
	v2d.Y += other.Y
	return v2d
}

func (v2d *Vector2D) Decr(other *Vector2D) *Vector2D {
	v2d.X -= other.X
	v2d.Y -= other.Y
	return v2d
}

func (v2d *Vector2D) ToJSON() coder.JSON {
	return coder.JSON{
		"x": v2d.X,
		"y": v2d.Y,
	}
}
