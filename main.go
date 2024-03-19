package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/facuellarg/tetris/board"
	"github.com/facuellarg/tetris/piece"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	b             *board.Board
	SquareSize    int
	logicalHeight int
	logicalWidth  int
	lastUpdate    time.Time
	dropDownRate  time.Duration
	score         int
	gameOver      bool
}

var p = piece.CreatePiece(piece.PIECES[0], piece.Position{X: 1, Y: 0}, piece.COLORS[0])

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if !g.b.IsLeftCollision(p.Position.X-1, p.Position.Y, p.Shape) {
			p.Position.X--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if !g.b.IsRightCollision(p.Position.X+1, p.Position.Y, p.Shape) {
			p.Position.X++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		UpdateYPosition(g, p)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		pieceRotate := piece.Rotate(p.Shape)
		if !pieceIsColliding(g, pieceRotate, p.Position.X, p.Position.Y) {
			p.Shape = pieceRotate
		}
	}

	if g.gameOver && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		*g = *StartGame()

		newPieceIndex := rand.Intn(len(piece.PIECES))
		startX := rand.Intn(g.b.GetWidth()-piece.PIECES[newPieceIndex].GetWidth()-2) + 2

		*p = *piece.CreatePiece(
			piece.PIECES[newPieceIndex],
			piece.Position{X: startX, Y: 1},
			piece.COLORS[newPieceIndex],
		)
		if pieceIsColliding(g, p.Shape, p.Position.X, p.Position.Y) {
			// Game over
			g.gameOver = true

		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over", g.logicalWidth/2-40, g.logicalHeight/2)
		ebitenutil.DebugPrintAt(screen, "Press Enter to restart", g.logicalWidth/2-80, g.logicalHeight/2+20)
		// print score
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), g.logicalWidth/2-40, g.logicalHeight/2+40)
		return
	}

	//print last time the board was updated
	if time.Since(g.lastUpdate) > g.dropDownRate {
		UpdateYPosition(g, p)
		g.lastUpdate = time.Now()
	}
	for i := 0; i < g.b.GetHeight(); i++ {
		for j := 0; j < g.b.GetWidth(); j++ {
			vector.DrawFilledRect(screen, float32(j*g.SquareSize), float32(i*g.SquareSize), float32(g.SquareSize), float32(g.SquareSize), makeLighter(
				g.b.GetColor(i, j), .5,
			), true)
			vector.DrawFilledRect(screen, float32(j*g.SquareSize), float32(i*g.SquareSize), float32(g.SquareSize-1), float32(g.SquareSize-1), g.b.GetColor(i, j), false)
		}
	}

	// Draw the piece
	for i := 0; i < len(p.Shape); i++ {
		for j := 0; j < len(p.Shape[i]); j++ {
			if p.Shape[i][j] == 1 {
				vector.DrawFilledRect(screen, float32((j+p.Position.X)*g.SquareSize), float32((i+p.Position.Y)*g.SquareSize), float32(g.SquareSize), float32(g.SquareSize), makeLighter(
					p.Color, .5,
				), true)
				vector.DrawFilledRect(screen, float32((j+p.Position.X)*g.SquareSize), float32((i+p.Position.Y)*g.SquareSize), float32(g.SquareSize-2), float32(g.SquareSize-2), p.Color, true)
			}
		}
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f FPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 0, 20)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.logicalWidth, g.logicalHeight
}

func main() {

	g := StartGame()
	ebiten.SetWindowSize(
		g.logicalWidth*3,
		g.logicalHeight*3,
	)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func UpdateYPosition(g *Game, p *piece.Piece) {
	if g.b.IsDownCollision(p.Position.X, p.Position.Y+1, p.Shape) {
		g.b.AddPieceToBoard(*p)
		g.score += g.b.CheckLines(p.Position.Y, p.Position.Y+p.GetHeight()) * 100

		*p = *CreateRandomPiece(g.b)
		if pieceIsColliding(g, p.Shape, p.Position.X, p.Position.Y) {
			// Game over
			g.gameOver = true

		}
	} else {
		p.Position.Y++
	}
}

func pieceIsColliding(g *Game, shape piece.Shape, x, y int) bool {
	return g.b.IsDownCollision(x, y, shape) || g.b.IsLeftCollision(x, y, shape) || g.b.IsRightCollision(x, y, shape)
}

func StartGame() *Game {
	g := &Game{
		b:            board.CreateBoard(10, 15),
		SquareSize:   30,
		lastUpdate:   time.Now(),
		dropDownRate: time.Second / 5,
	}

	g.logicalHeight = g.b.GetHeight() * g.SquareSize
	g.logicalWidth = g.b.GetWidth() * g.SquareSize

	return g
}

func CreateRandomPiece(b *board.Board) *piece.Piece {
	newPieceIndex := rand.Intn(len(piece.PIECES))
	fmt.Println("argument rand ", b.GetWidth()-piece.PIECES[newPieceIndex].GetWidth()-2)
	startX := rand.Intn(b.GetWidth()-piece.PIECES[newPieceIndex].GetWidth()-2) + 2
	return piece.CreatePiece(
		piece.PIECES[newPieceIndex],
		piece.Position{X: startX, Y: 0},
		piece.COLORS[newPieceIndex],
	)
}

func makeLighter(c color.Color, factor float64) color.RGBA {
	// if factor < 1 {
	// 	factor = 1 // Ensure factor is not less than 1, which would make the color darker instead of lighter.
	// }

	// Convert color.Color to RGBA, which we can work with.
	r, g, b, a := c.RGBA()

	// Since RGBA returns color components in the range [0, 65535], we first scale them to [0, 255].
	r, g, b = r>>8, g>>8, b>>8

	// Lighten the color by increasing the RGB values towards 255 by the specified factor.
	// We use math.Min to ensure the values do not go beyond 255.
	newR := math.Min(255, float64(r)*factor)
	newG := math.Min(255, float64(g)*factor)
	newB := math.Min(255, float64(b)*factor)

	return color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a >> 8)}
}
