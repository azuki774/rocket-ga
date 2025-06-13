package cmd

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/spf13/cobra"
)

type Game struct {
	RocketImg  *ebiten.Image
	EarthImg   *ebiten.Image
	MoonImg    *ebiten.Image
	RocketImgX float64
	RocketImgY float64
}

func (g *Game) Update() error {
	g.RocketImgX += 1
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.RocketImgX, 0)
	screen.DrawImage(g.RocketImg, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(g.EarthImg, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(1600-80, 1200-80)
	screen.DrawImage(g.MoonImg, op)

	ebitenutil.DebugPrint(screen, "Rocket GA")
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
		// todo
		fmt.Println("Pass")
		ebiten.SetWindowTitle("Ebitengine 入門")
		ebiten.SetWindowSize(1600, 1200)

		g := &Game{}

		img, _, err := ebitenutil.NewImageFromFile("earth.png")
		if err != nil {
			panic(err)
		}
		g.EarthImg = img

		img, _, err = ebitenutil.NewImageFromFile("moon.png")
		if err != nil {
			panic(err)
		}
		g.MoonImg = img
		img, _, err = ebitenutil.NewImageFromFile("rocket.png")
		if err != nil {
			panic(err)
		}
		g.RocketImg = img

		if err := ebiten.RunGame(g); err != nil {
			panic(err)
		}
	},
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
