package main

import "math"

var currentEnemyId = 0

func NewPorn() Enemy {
	currentEnemyId++
	return Enemy{
		Id:    currentEnemyId,
		SizeX: 20,
		SizeY: 20,
		Score: 100,
		HP:    10,
	}
}

var p1 = func(p *EnemyPattern) {
	if p.displayTime < FRAME_1SEC*3 {
		p.X -= 1.4
	}
	if p.displayTime >= FRAME_1SEC*3 && p.displayTime < FRAME_1SEC*5 {
		p.X += 1
		p.Y += 1
	}
	if p.displayTime >= FRAME_1SEC*5 && p.displayTime < FRAME_1SEC*8 {
		p.X -= 1.4
	}
}

func NewShot(p *EnemyPattern) *Shot {
	v := math.Sqrt(math.Pow(x-p.X, 2) + math.Pow(y-p.Y, 2))
	return &Shot{
		Id:        newShotId(),
		X:         p.X,
		Y:         p.Y,
		Vx:        (x - p.X) * 1.5 / v,
		Vy:        (y - p.Y) * 1.5 / v,
		MoveCount: 1,
		SizeX:     3,
		SizeY:     3,
		Power:     10,
	}
}

var s1 = func(p *EnemyPattern) *Shot {
	if p.displayTime == FRAME_1SEC*3 {
		return NewShot(p)
	}

	if p.displayTime == FRAME_1SEC*6 {
		return NewShot(p)
	}
	return nil
}

var currentPatternId = 1

func newPatternId() int {
	currentPatternId++
	return currentPatternId
}

func newShotId() int {
	shotIndex++
	return shotIndex
}

var stage1 = []*EnemyPattern{
	{
		Id:            newPatternId(),
		displayIn:     1,
		displayPeriod: FRAME_1SEC * 10,
		Enemy:         NewPorn(),
		X:             320,
		Y:             100,
		movePattern:   p1,
		shotPattern:   s1,
	},
	{
		Id:            newPatternId(),
		displayIn:     21,
		displayPeriod: FRAME_1SEC * 10,
		Enemy:         NewPorn(),
		X:             320,
		Y:             100,
		movePattern:   p1,
		shotPattern:   s1,
	},
	{
		Id:            newPatternId(),
		displayIn:     41,
		displayPeriod: FRAME_1SEC * 10,
		Enemy:         NewPorn(),
		X:             320,
		Y:             100,
		movePattern:   p1,
		shotPattern:   s1,
	},
}
