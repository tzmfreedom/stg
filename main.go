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
	"time"
)

var timer int
var x, y float64
var score int
var phase int
var initPlayerX, initPlayerY float64
var speed float64
var shotLevel = 1

var currentStage []*EnemyPattern
var currentEnemyPattern []*EnemyPattern
var currentItemPatterns []*ItemPattern

const (
	screenWidth  = 320
	screenHeight = 240
	playerSizeX  = 20
	playerSizeY  = 5
)

const (
	PHASE_GAMESTART = iota
	PHASE_GAMEOVER
)

const (
	STATE_HIDDEN = iota
	STATE_DISPLAY_IN
	STATE_DEAD
	STATE_DISPLAY_OUT
)

const (
	ITEM_POWERUP = iota
)

type Enemy struct {
	Id    int
	X     float64
	Y     float64
	SizeX float64
	SizeY float64
	Score int
	HP    int
}

type Shot struct {
	Id        int
	X         float64
	Y         float64
	Vx        float64
	Vy        float64
	SizeX     float64
	SizeY     float64
	Power     int
	MoveCount int
}

type Block struct {
	X     float64
	Y     float64
	SizeX float64
	SizeY float64
}

type EnemyPattern struct {
	Id            int
	displayIn     int
	displayPeriod int
	displayTime   int
	State         int
	Enemy         Enemy
	X             float64
	Y             float64
	movePattern   func(p *EnemyPattern)
	shotPattern   func(p *EnemyPattern) *Shot
}

type Item struct {
	SizeX float64
	SizeY float64
}

type ItemPattern struct {
	Id   int
	X    float64
	Y    float64
	Item *Item
}

const FRAME_1SEC = 60

var shotIndex = 1
var itemIndex = 1
var shots = []*Shot{}
var enemyShots = []*Shot{}
var blocks = []*Block{}
var currentItemId = 0

var powerUpItem = &Item{
	SizeX: 20,
	SizeY: 20,
}

func update(screen *ebiten.Image) error {
	if phase != PHASE_GAMEOVER {
		timer++
		for _, pattern := range currentEnemyPattern {
			pattern.displayTime++
		}
		handleInput()
	}
	if phase == PHASE_GAMEOVER {
		ebitenutil.DebugPrint(screen, "GAME OVER!: "+strconv.Itoa(score))
	} else {
		ebitenutil.DebugPrint(screen, strconv.Itoa(score))

		itemGet()
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

func itemGet() {
	pattern := conflictItem()
	if pattern != nil {
		shotLevel++
		newItemPatterns := []*ItemPattern{}
		for _, p := range currentItemPatterns {
			if p.Id != pattern.Id {
				newItemPatterns = append(newItemPatterns, pattern)
			}
		}
		currentItemPatterns = newItemPatterns
	}
}

func emenyAction() {
	for _, pattern := range stage1 {
		if timer == pattern.displayIn {
			pattern.State = STATE_DISPLAY_IN
			currentEnemyPattern = append(currentEnemyPattern, pattern)
		}
	}
	newPattern := []*EnemyPattern{}
	for _, pattern := range currentEnemyPattern {
		if timer > pattern.displayIn+pattern.displayPeriod {
			pattern.State = STATE_DISPLAY_OUT
			continue
		}
		pattern.movePattern(pattern)
		shot := pattern.shotPattern(pattern)
		if shot != nil {
			enemyShots = append(enemyShots, shot)
		}
		newPattern = append(newPattern, pattern)
	}

	currentEnemyPattern = newPattern
}

func moveShot() {
	for _, shot := range enemyShots {
		if timer%shot.MoveCount == 0 {
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
		pattern := conflictEnemy(shot.X, shot.Y, shot.SizeX, shot.SizeY)
		if pattern != nil {
			pattern.Enemy.HP -= shot.Power
			if pattern.Enemy.HP <= 0 {
				score += pattern.Enemy.Score
				removeEnemy(pattern)
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

func removeEnemy(search *EnemyPattern) {
	newEnemyPattern := []*EnemyPattern{}
	for _, e := range currentEnemyPattern {
		if e.Id != search.Id {
			newEnemyPattern = append(newEnemyPattern, e)
		}
	}
	currentEnemyPattern = newEnemyPattern
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
			x -= speed
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
			Id:    shotIndex,
			X:     x,
			Y:     y,
			SizeX: 5,
			SizeY: 5,
			Vx:    5,
			Vy:    0,
			Power: 10,
		})
		shotIndex++

		if shotLevel > 1 && (v == 1 || v%30 == 0) {
			shots = append(shots, &Shot{
				Id:    shotIndex,
				X:     x,
				Y:     y,
				SizeX: 5,
				SizeY: 5,
				Vx:    3,
				Vy:    4,
				Power: 5,
			})
			shotIndex++
		}
		if shotLevel > 2 && (v == 1 || v%30 == 0) {
			shots = append(shots, &Shot{
				Id:    shotIndex,
				X:     x,
				Y:     y,
				SizeX: 5,
				SizeY: 5,
				Vx:    3,
				Vy:    -4,
				Power: 5,
			})
			shotIndex++
		}
	}
	if conflictEnemy(x, y, playerSizeX, playerSizeY) != nil {
		phase = PHASE_GAMEOVER
	}
}

func draw(screen *ebiten.Image) {
	img, _ := ebiten.NewImage(playerSizeX, playerSizeY, 0)
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

	for _, enemy := range currentEnemyPattern {
		img, _ := ebiten.NewImage(int(enemy.Enemy.SizeX), int(enemy.Enemy.SizeY), 0)
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

	for _, pattern := range currentItemPatterns {
		img, _ := ebiten.NewImage(int(pattern.Item.SizeX), int(pattern.Item.SizeY), 0)
		img.Fill(color.RGBA{0x00, 0x00, 0xff, 0xff})
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(pattern.X), float64(pattern.Y))
		screen.DrawImage(img, options)
	}
	// text.Draw(screen, string(score), scoreFont, scoreX, scoreY, color.White)
}

func conflictItem() *ItemPattern {
	for _, pattern := range currentItemPatterns {
		if pattern.Y < y+playerSizeY && pattern.Y+pattern.Item.SizeY > y && pattern.X+pattern.Item.SizeX > x && pattern.X < x+playerSizeX {
			return pattern
		}
	}
	return nil
}

func conflictEnemy(x, y, sizeX, sizeY float64) *EnemyPattern {
	for _, pattern := range currentEnemyPattern {
		if pattern.Y < y+sizeY && pattern.Y+pattern.Enemy.SizeY > y && pattern.X+pattern.Enemy.SizeX > x && pattern.X < x+sizeX {
			return pattern
		}
	}
	return nil
}

func conflictBlock(x, y, sizeX, sizeY float64) bool {
	for _, block := range blocks {
		if block.Y < y+sizeY && block.Y+block.SizeY > y && block.X+block.SizeX > x && block.X < x+sizeX {
			return true
		}
	}
	return false
}

func conflictPlayer(shot *Shot) bool {
	if shot.Y < y+playerSizeY && shot.Y+shot.SizeY > y && shot.X+shot.SizeX > x && shot.X < x+playerSizeX {
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
	currentStage = stage1
	currentItemPatterns = []*ItemPattern{
		{
			Id:   newItemId(),
			X:    100,
			Y:    100,
			Item: powerUpItem,
		},
	}
	rand.Seed(time.Now().UnixNano())
	blocks = append(blocks, &Block{
		X:     20 + float64(rand.Intn(300)),
		Y:     float64(rand.Intn(240)),
		SizeX: 5,
		SizeY: 50,
	})
}

func newItemId() int {
	currentItemId++
	return currentItemId
}

func debug(args ...interface{}) {
	pp.Println(args...)
}
