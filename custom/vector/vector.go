package vector

import (
	"custom/point"
)

type Vector struct {
	X float64
	Y float64
}

func New(X float64, Y float64) Vector {
	e := Vector{X, Y}
	return e
}

func (vector Vector) Normalized() Vector{
	length := point.New(0,0).DistanceTo( point.New(vector.X,vector.Y))
	return Vector{ vector.X / length, vector.Y / length}
}
