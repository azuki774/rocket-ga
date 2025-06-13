package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"rocket-ga/internal/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/spf13/cobra"
)

var ErrUsualEnd = errors.New("usual end")

type Game struct {
	T         int
	RocketImg *ebiten.Image
	EarthImg  *ebiten.Image
	MoonImg   *ebiten.Image
	Rocket    *model.Object
	Earth     *model.Object
	Moon      *model.Object
}

func (g *Game) Update() error {
	g.T += 1
	// 状態計算
	thrustcmds := []model.ThrustCommand{
		{
			StartTime: 0,
			Duration:  70,
			Angle:     math.Pi / 4,
			Power:     0.08,
		},
	}
	nr := g.Rocket.EmulateNextBy2(float64(g.T), *g.Earth, *g.Moon, thrustcmds)
	// 状態更新
	g.Rocket = nr
	fmt.Printf("t=%d, X=%f, Y=%f, vX=%f, vY=%f, m=%f\n", g.T, g.Rocket.Pos.X, g.Rocket.Pos.Y, g.Rocket.Vel.X, g.Rocket.Vel.Y, g.Rocket.Mass)

	landCondition := g.Rocket.IsCollision(*g.Moon)
	if landCondition == model.ColisionClash {
		fmt.Println("CLASH!")
		return ErrUsualEnd
	} else if landCondition == model.ColisionLand {
		fmt.Println("LAND!")
		return ErrUsualEnd
	}

	landConditionEarth := g.Rocket.IsCollision(*g.Earth)
	if landConditionEarth == model.ColisionClash || landConditionEarth == model.ColisionLand {
		fmt.Println("EARTH CLASH!")
		return ErrUsualEnd
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// エミュレート座標 0, 0 = 画面表示 500, 500
	op.GeoM.Translate(g.Rocket.Pos.X+500, g.Rocket.Pos.Y+500)
	screen.DrawImage(g.RocketImg, op)

	op = &ebiten.DrawImageOptions{}
	// エミュレート座標 0, 0 = 画面表示 500, 500
	op.GeoM.Translate(g.Earth.Pos.X+500, g.Earth.Pos.Y+500)
	screen.DrawImage(g.EarthImg, op)

	op = &ebiten.DrawImageOptions{}
	// エミュレート座標 0, 0 = 画面表示 500, 500
	op.GeoM.Translate(g.Moon.Pos.X+500, g.Moon.Pos.Y+500)
	screen.DrawImage(g.MoonImg, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("t=%d, X=%f, Y=%f, vX=%f, vY=%f, m=%f\n", g.T, g.Rocket.Pos.X, g.Rocket.Pos.Y, g.Rocket.Vel.X, g.Rocket.Vel.Y, g.Rocket.Mass))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rocket-ga",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := startEmulate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func startEmulate() error {
	ebiten.SetWindowTitle("Rocket GA")
	ebiten.SetWindowSize(1000, 1000)

	g := &Game{}

	img1, _, err := ebitenutil.NewImageFromFile("rocket.png")
	if err != nil {
		return err
	}
	g.RocketImg = img1
	img2, _, err := ebitenutil.NewImageFromFile("earth.png")
	if err != nil {
		return err
	}
	g.EarthImg = img2

	img3, _, err := ebitenutil.NewImageFromFile("moon.png")
	if err != nil {
		return err
	}
	g.MoonImg = img3

	// 初期位置
	g.Rocket = &model.Object{
		Mass: model.RocketMass,
		Pos:  model.Vector{X: model.InitRocketPosX, Y: model.InitRocketPosY},
		Vel:  model.Vector{X: model.InitRocketVelX, Y: model.InitRocketVelY},
	}

	g.Earth = &model.Object{
		Mass:   model.EarthMass,
		Pos:    model.Vector{X: model.InitEarthPosX, Y: model.InitEarthPosY},
		Vel:    model.Vector{X: 0, Y: 0},
		Radius: model.EarthRadius,
	}

	g.Moon = &model.Object{
		Mass:   model.MoonMass,
		Pos:    model.Vector{X: model.InitMoonPosX, Y: model.InitMoonPosY},
		Vel:    model.Vector{X: 0, Y: 0},
		Radius: model.MoonRadius,
	}

	if err := ebiten.RunGame(g); err != nil {
		return err
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rocket-ga.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
