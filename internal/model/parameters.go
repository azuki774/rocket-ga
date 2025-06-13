package model

var G float64 = 1 // 重力定数
var EarthMass float64 = 10000
var MoonMass float64 = 123

var RocketDryMass float64 = 0.1
var RocketFuelMass float64 = 0.9
var RocketMass float64 = RocketDryMass + RocketFuelMass // 0.1 -> 本体、0.9 -> 燃料

var InitRocketPosX float64 = 400
var InitRocketPosY float64 = -400
var InitRocketVelX float64 = -0.5
var InitRocketVelY float64 = 0.5

var InitEarthPosX float64 = 0
var InitEarthPosY float64 = 0
var InitMoonPosX float64 = 400
var InitMoonPosY float64 = 400

var EarthRadius float64 = 30

type ThrustCommand struct {
	StartTime float64
	Duration  float64
	Angle     float64
	Power     float64
}

var FuelConsumptionRate float64 = 0.01 // 燃料の消費レート

// 例：10回の噴射コマンドで1つの遺伝子とする
const NumCommands = 10

type Chromosome [NumCommands]ThrustCommand
