package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/k0kubun/pp"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var timer int
var x, y float64
var score int
var phase int
var initPlayerX, initPlayerY float64
var speed float64

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
	Type int
	Id int
	X float64
	Y float64
	Vx float64
	Vy float64
	SizeX float64
	SizeY float64
	Score int
	HP int
}

type Shot struct {
	Id int
	X float64
	Y float64
	Vx float64
	Vy float64
	SizeX float64
	SizeY float64
	Power int
	MoveCount int
}

type Block struct {
	X float64
	Y float64
	SizeX float64
	SizeY float64
}

var shotIndex = 1
var enemies = []*Enemy{}
var shots = []*Shot{}
var enemyShots = []*Shot{}
var blocks = []*Block{}

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
		moveShot()
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	draw(screen)
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "STG"); err != nil {
		log.Fatal(err)
	}
}

func emenyAction() {
	for _, enemy := range enemies {
		enemy.X += enemy.Vx
		enemy.Y += enemy.Vy

		if enemy.Type == 1 && timer % 60 == 0 {
			v := math.Sqrt(math.Pow(x - enemy.X, 2) + math.Pow(y - enemy.Y, 2))
			enemyShots = append(enemyShots, &Shot{
				Id: shotIndex,
				X: enemy.X,
				Y: enemy.Y,
				Vx: (x - enemy.X)*(1 + float64(rand.Intn(5)))/v,
				Vy: (y - enemy.Y)*(1 + float64(rand.Intn(5)))/v,
				SizeX: 2,
				SizeY: 2,
				MoveCount: 4,
			})
			shotIndex++
		}
	}
}

func moveShot() {
	for _, shot := range enemyShots {
		if timer % shot.MoveCount == 0 {
			shot.X += shot.Vx
			shot.Y += shot.Vy
		}
		if conflictPlayer(shot) {
			phase = PHASE_GAMEOVER
			return
		}
		if conflictBlock(shot.X, shot.Y, shot.SizeX, shot.SizeY) {
			removeEnemyShot(shot)
		}
		if shot.X < 0 || shot.Y < 0 || shot.X > screenWidth || shot.Y > screenHeight {
			removeEnemyShot(shot)
		}
	}
	for _, shot := range shots {
		shot.X += shot.Vx
		shot.Y += shot.Vy
		enemy := conflictEnemy(shot.X, shot.Y, shot.SizeX, shot.SizeY)
		if enemy != nil {
			enemy.HP -= shot.Power
			if enemy.HP <= 0 {
				score += enemy.Score
				removeEnemy(enemy)
				removeShot(shot)
			}
		}
		if conflictBlock(shot.X, shot.Y, shot.SizeX, shot.SizeY) {
			removeShot(shot)
		}
		if shot.X < 0 || shot.Y < 0 || shot.X > screenWidth || shot.Y > screenHeight {
			removeShot(shot)
		}
	}
}

func removeEnemy(search *Enemy) {
	newEnemies := []*Enemy{}
	for _, e := range enemies {
		if e.Id != search.Id {
			newEnemies = append(newEnemies, e)
		}
	}
	enemies = newEnemies
}

func removeShot(search *Shot) {
	newShots := []*Shot{}
	for _, s := range shots {
		if s.Id != search.Id {
			newShots = append(newShots, s)
		}
	}
	shots = newShots
}

func removeEnemyShot(search *Shot) {
	newShots := []*Shot{}
	for _, s := range enemyShots {
		if s.Id != search.Id {
			newShots = append(newShots, s)
		}
	}
	enemyShots = newShots
}


func handleInput() {
	if v := inpututil.KeyPressDuration(ebiten.KeyLeft); v > 0 {
		if x > 0 {
			x-=speed
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyRight); v > 0 {
		if x < screenWidth {
			x += speed
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyDown); v > 0 {
		if y < screenHeight {
			y += speed
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyUp); v > 0 {
		if y > 0 {
			y -= speed
		}
	}
	if v := inpututil.KeyPressDuration(ebiten.KeyS); v == 1 || (v > 0 && v%10 == 0) {
		shots = append(shots, &Shot{
			Id: shotIndex,
			X: x+5,
			Y: y+5,
			SizeX: 5,
			SizeY: 5,
			Vx: 1,
			Vy: 0,
			Power: 10,
		})
		shotIndex++
	}
	if conflictEnemy(x, y, playerSize, playerSize) != nil {
		phase = PHASE_GAMEOVER
	}
}

func draw(screen *ebiten.Image) {
	img, _ := ebiten.NewImage(playerSize, playerSize, 0)
	img.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, options)

	for _, shot := range shots {
		img, _ := ebiten.NewImage(int(shot.SizeX), int(shot.SizeY), 0)
		img.Fill(color.RGBA{0x00, 0x00, 0xff, 0xff})
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(shot.X), float64(shot.Y))
		screen.DrawImage(img, options)
	}

	for _, shot := range enemyShots {
		img, _ := ebiten.NewImage(int(shot.SizeX), int(shot.SizeY), 0)
		img.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(shot.X), float64(shot.Y))
		screen.DrawImage(img, options)
	}

	for _, enemy := range enemies {
		img, _ := ebiten.NewImage(int(enemy.SizeX), int(enemy.SizeY), 0)
		img.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(enemy.X), float64(enemy.Y))
		screen.DrawImage(img, options)
	}

	for _, block := range blocks {
		img, _ := ebiten.NewImage(int(block.SizeX), int(block.SizeY), 0)
		img.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(block.X), float64(block.Y))
		screen.DrawImage(img, options)
	}
	// text.Draw(screen, string(score), scoreFont, scoreX, scoreY, color.White)
}

func conflictEnemy(x, y, sizeX, sizeY float64) *Enemy {
	for _, enemy := range enemies {
		if enemy.Y < y + sizeY && enemy.Y + enemy.SizeY > y && enemy.X + enemy.SizeX > x && enemy.X < x + sizeX {
			return enemy
		}
	}
	return nil
}

func conflictBlock(x, y, sizeX, sizeY float64) bool {
	for _, block := range blocks {
		if block.Y < y + sizeY && block.Y + block.SizeY > y && block.X + block.SizeX > x && block.X < x + sizeX {
			return true
		}
	}
	return false
}

func conflictPlayer(shot *Shot) bool {
	if shot.Y < y + playerSize && shot.Y + shot.SizeY > y && shot.X + shot.SizeX > x && shot.X < x + playerSize {
		return true
	}
	return false
}

func init() {
	timer = 0
	score = 0
	x = initPlayerX
	y = initPlayerY
	speed = 1
	phase = PHASE_GAMESTART
	rand.Seed(time.Now().UnixNano())
	n := 30 + rand.Intn(5)
	for i := 0; i < n; i++ {
		enemy := &Enemy{
			Id: i,
			X: 20 + float64(rand.Intn(300)),
			Y: float64(rand.Intn(240)),
			Vx: 0,
			Vy: 0,
			//Vx: -5 + rand.Intn(10),
			//Vy: -5 + rand.Intn(10),
			SizeX: 5 + float64(rand.Intn(5)),
			SizeY: 5 + float64(rand.Intn(5)),
			Score: 100 + rand.Intn(50),
			Type: rand.Intn(2),
			HP: 10 + rand.Intn(100),
		}
		enemies = append(enemies, enemy)
	}
	blocks = append(blocks, &Block{
		X: 20 + float64(rand.Intn(300)),
		Y: float64(rand.Intn(240)),
		SizeX: 5,
		SizeY: 50,
	})
}

func debug(args ...interface{}) {
	pp.Println(args...)
}


