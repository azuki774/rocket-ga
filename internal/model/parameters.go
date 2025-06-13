package model

var G float64 = 1 // 重力定数
var EarthMass float64 = 10000
var MoonMass float64 = 123
var RocketMass float64 = 1

type Vector struct {
	X float64
	Y float64
}

func NewVector(X float64, Y float64) *Vector {
	return &Vector{X: X, Y: Y}
}
