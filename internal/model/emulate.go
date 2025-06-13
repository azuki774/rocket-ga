package model

import (
	"math"
)

type Object struct {
	Mass   float64
	Pos    Vector  // 位置
	Vel    Vector  // 速度
	Radius float64 // 半径
}

type Vector struct {
	X float64
	Y float64
}

// 2質点からの万有引力を受けて、1秒後の位置を計算
// 与えたオブジェクトは更新して、新たな座標を返す
func (o *Object) EmulateNextBy2(m1 Object, m2 Object) *Object {
	// 衝突している場合を判定
	if CollisionObject(*o, m1) || CollisionObject(*o, m2) {
		v := Vector{X: 0, Y: 0}
		p := Vector{X: o.Pos.X, Y: o.Pos.Y}
		no := Object{Mass: o.Mass, Pos: p, Vel: v}
		return &no
	}

	// 各天体からの引力を個別に計算
	f1 := computeGravityForce(m1, *o) // m1 が o を引く力
	f2 := computeGravityForce(m2, *o) // m2 が o を引く力

	// 力の合成
	f := Vector{X: f1.X + f2.X, Y: f1.Y + f2.Y}

	// 加速度計算 a = F/m
	a := Vector{X: f.X / o.Mass, Y: f.Y / o.Mass}

	// 1 秒後の「速度」を更新する v = v + a * 1.0
	v := Vector{X: o.Vel.X + a.X, Y: o.Vel.Y + a.Y}

	// 1 秒後の「位置」を更新する p = p + v * 1.0
	p := Vector{X: o.Pos.X + v.X, Y: o.Pos.Y + v.Y}

	no := Object{Mass: o.Mass, Pos: p, Vel: v}
	return &no
}

func CollisionObject(o1 Object, o2 Object) bool {
	dVec := Vector{X: o2.Pos.X - o1.Pos.X, Y: o2.Pos.Y - o1.Pos.Y}
	d := math.Sqrt(dVec.X*dVec.X + dVec.Y*dVec.Y)
	if d <= o1.Radius+o2.Radius {
		// 衝突
		return true
	}
	return false
}

// o1 が o2 を引く力を計算する
func computeGravityForce(o1 Object, o2 Object) (v Vector) {
	// o1 から o2 へのベクトル
	dEarth := Vector{X: o2.Pos.X - o1.Pos.X, Y: o2.Pos.Y - o1.Pos.Y}
	// o1 と o2 の距離
	d := math.Sqrt(dEarth.X*dEarth.X + dEarth.Y*dEarth.Y)
	d2 := dEarth.X*dEarth.X + dEarth.Y*dEarth.Y
	// 力の大きさを計算
	f := G * o1.Mass * o2.Mass / d2
	// 力のベクトルを計算
	v = Vector{X: -f * dEarth.X / d, Y: -f * dEarth.Y / d}
	return v
}
