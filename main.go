package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/k0kubun/pp"
	"image/color"
	"log"
	"math/rand"
	"strconv"
)

var timer int
var x, y int
var score int
var phase int
var initPlayerX, initPlayerY int

const (
	screenWidth = 320
	screenHeight = 240
	playerSize = 5
)

const (
	PHASE_GAMESTART = iota
	PHASE_GAMEOVER
)

type Enemy struct {
	X int
	Y int
	Vx int
	Vy int
	SizeX int
	SizeY int
	Counter int
	VCounter int
}

var enemies = []*Enemy{}

func update(screen *ebiten.Image) error {
	if phase != PHASE_GAMEOVER {
		timer++
		handleInput()
	}
	if phase == PHASE_GAMEOVER {
		ebitenutil.DebugPrint(screen, "GAME OVER!: " + strconv.Itoa(score))
	} else {
		ebitenutil.DebugPrint(screen, strconv.Itoa(score))

		emenyAction()
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	draw(screen)
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Tetris!"); err != nil {
		log.Fatal(err)
	}
}

func emenyAction() {
	for _, enemy := range enemies {
		enemy.Counter++
		if enemy.Counter == 1 {
			enemy.X += enemy.Vx
			enemy.Y += enemy.Vy
			enemy.Counter = 0
		}
	}
}

func handleInput() {
	if v := inpututil.KeyPressDuration(ebiten.KeyLeft); v > 0 && v%2 == 0 {
		if x > 0 {
			x-=5
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyRight); v > 0 && v%2 == 0 {
		if x < screenWidth {
			x += 5
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyDown); v > 0 && v%2 == 0 {
		if y < screenHeight {
			y += 5
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyUp); v > 0 && v%2 == 0 {
		if y > 0 {
			y -= 5
		}
	}
	if isConflict() {
		phase = PHASE_GAMEOVER
	}
}

func draw(screen *ebiten.Image) {
	img, _ := ebiten.NewImage(playerSize, playerSize, 0)
	img.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, options)

	for _, enemy := range enemies {
		img, _ := ebiten.NewImage(playerSize, playerSize, 0)
		img.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(enemy.X), float64(enemy.Y))
		screen.DrawImage(img, options)
	}
	// text.Draw(screen, string(score), scoreFont, scoreX, scoreY, color.White)
}

func isConflict() bool {
	for _, enemy := range enemies {
		if enemy.Y < y + playerSize && enemy.Y + enemy.SizeY > y && enemy.X + enemy.SizeX > x && enemy.X < x + playerSize {
			return true
		}
	}
	return false
}

func init() {
	timer = 0
	score = 0
	x = initPlayerX
	y = initPlayerY
	phase = PHASE_GAMESTART
	n := 30 + rand.Intn(5)
	for i := 0; i < n; i++ {
		enemy := &Enemy{
			X: 20 + rand.Intn(300),
			Y: rand.Intn(240),
			Vx: -10 + rand.Intn(20),
			Vy: -10 + rand.Intn(20),
			SizeX: 5 + rand.Intn(5),
			SizeY: 5 + rand.Intn(5),
			Counter: 0,
		}
		enemies = append(enemies, enemy)
	}
}

func debug(args ...interface{}) {
	pp.Println(args...)
}


