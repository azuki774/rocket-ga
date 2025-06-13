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

type ColisionCondition string

var safeLandingSpeed float64 = 10
var ColisionNone ColisionCondition = "COLISION_NONE"
var ColisionLand ColisionCondition = "COLISION_LAND"
var ColisionClash ColisionCondition = "COLISION_CLASH"

// o が m1 とぶつかっていない or 衝突 or 着陸かを判定する
func (o *Object) IsCollision(m1 Object) ColisionCondition {
	// 衝突している場合を判定
	if CollisionObject(*o, m1) {
		vabs := math.Sqrt(o.Vel.X*o.Vel.X + o.Vel.Y*o.Vel.Y)
		if vabs <= safeLandingSpeed {
			// 着陸
			return ColisionLand
		} else {
			// 衝突
			return ColisionClash
		}
	}

	return ColisionNone
}

// 2質点からの万有引力を受けて、1秒後の位置を計算
// 与えたオブジェクトは更新して、新たな座標を返す
func (o *Object) EmulateNextBy2(t float64, m1 Object, m2 Object, thurstCmds []ThrustCommand) *Object {
	// デフォルトの推力はゼロ
	thrustForce := Vector{X: 0, Y: 0}
	power := 0.0

	// 現在時刻で実行すべき噴射コマンドを探す
	for _, cmd := range thurstCmds {
		// 噴射期間中か？
		if t >= cmd.StartTime && t < cmd.StartTime+cmd.Duration {
			// 燃料が残っているか？
			if o.Mass > RocketDryMass { // o の質量が、空タン状態の o よりも大きければ燃料が残っている
				power = cmd.Power
				thrustForce.X = power * math.Cos(cmd.Angle)
				thrustForce.Y = power * math.Sin(cmd.Angle)
			}
			break // 同時刻に複数のコマンドは実行しないと仮定
		}
	}

	// 各天体からの引力を個別に計算
	f1 := computeGravityForce(m1, *o) // m1 が o を引く力
	f2 := computeGravityForce(m2, *o) // m2 が o を引く力

	// 力の合成
	f := Vector{X: f1.X + f2.X + thrustForce.X, Y: f1.Y + f2.Y + thrustForce.Y}

	// 加速度計算 a = F/m
	a := Vector{X: f.X / o.Mass, Y: f.Y / o.Mass}

	// 1 秒後の「速度」を更新する v = v + a * 1.0
	v := Vector{X: o.Vel.X + a.X, Y: o.Vel.Y + a.Y}

	// 1 秒後の「位置」を更新する p = p + v * 1.0
	p := Vector{X: o.Pos.X + v.X, Y: o.Pos.Y + v.Y}

	no := Object{Mass: o.Mass, Pos: p, Vel: v}
	if power > 0 {
		// もし燃料を噴射したなら、質量を減少させる
		fuelConsumed := FuelConsumptionRate * power * 1.0
		no.Mass -= fuelConsumed
		if no.Mass < RocketFuelMass {
			no.Mass = RocketFuelMass // ロケットの質量は空タンを下回らない
		}
	}
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
