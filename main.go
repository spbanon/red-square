package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	playerSize   = 50
	enemySize    = 60
	enemyCount   = 5
	speed        = 3
)

type Game struct {
	playerX, playerY float64
	enemies          [enemyCount]Enemy
}

type Enemy struct {
	x, y                   float64
	directionX, directionY float64
}

var start time.Time

func (g *Game) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	g.playerX = float64(cursorX) - playerSize/2
	g.playerY = float64(cursorY) - playerSize/2

	for i := range g.enemies {
		g.enemies[i].move(g.playerX, g.playerY)
	}

	return nil
}

func (e *Enemy) move(targetX, targetY float64) {
	e.x += e.directionX * speed
	e.y += e.directionY * speed

	directionX := targetX - e.x
	directionY := targetY - e.y
	length := math.Sqrt(directionX*directionX + directionY*directionY)
	if length < playerSize {
		fmt.Println("Gameover!\n", "With time: ", time.Since(start))
		os.Exit(0)
	}
	if rand.Float64() < 0.05 {
		e.directionX = rand.Float64()*2 - 1
		e.directionY = rand.Float64()*2 - 1
		length := math.Sqrt(e.directionX*e.directionX + e.directionY*e.directionY)
		e.directionX /= length
		e.directionY /= length
	}

	if e.x < 0 || e.x > screenWidth-enemySize {
		e.directionX *= -1
	}
	if e.y < 0 || e.y > screenHeight-enemySize {
		e.directionY *= -1
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	vector.DrawFilledRect(screen, float32(g.playerX), float32(g.playerY), float32(playerSize), float32(playerSize), color.RGBA{255, 0, 0, 255}, false)

	for _, enemy := range g.enemies {
		vector.DrawFilledRect(screen, float32(enemy.x), float32(enemy.y), float32(enemySize), float32(enemySize), color.RGBA{0, 0, 255, 255}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start = time.Now()
	game := &Game{}
	for i := range game.enemies {
		game.enemies[i] = Enemy{
			x:          float64(rand.Intn(screenWidth)),
			y:          float64(rand.Intn(screenHeight)),
			directionX: rand.Float64()*2 - 1,
			directionY: rand.Float64()*2 - 1,
		}
		length := math.Sqrt(game.enemies[i].directionX*game.enemies[i].directionX + game.enemies[i].directionY*game.enemies[i].directionY)
		game.enemies[i].directionX /= length
		game.enemies[i].directionY /= length
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Red Square Escape")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
