package board

import (
	"image/color"
	"math"

	"github.com/facuellarg/tetris/piece"
)

// type Board [][]int
// type BoardColor [][]color.Color

type Board struct {
	plates [][]int
	colors [][]color.Color
}

const (
	EMPTY = "_"
	BLOCK = "X"
)

var (
	//gray color for empty blocks
	EMPTY_COLOR = color.RGBA{192, 192, 192, 255}
)

func CreateBoard(width int, height int) *Board {
	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}
	boardColor := make([][]color.Color, height)
	for i := range boardColor {
		boardColor[i] = make([]color.Color, width)
	}
	for i := range boardColor {
		for j := range boardColor[i] {
			boardColor[i][j] = EMPTY_COLOR
		}
	}
	return &Board{board, boardColor}
}

func (b *Board) GetBoardColor() [][]color.Color {
	return b.colors
}

func (b *Board) GetWidth() int {
	return len(b.plates[0])
}

func (b *Board) GetHeight() int {
	return len(b.plates)
}

func (b *Board) IsEmpty(x, y int) bool {
	return !(x < 0 || y < 0 || x >= b.GetHeight() || y >= b.GetWidth() || b.plates[x][y] == 1)
}

func (b *Board) SetBlock(x, y int) {
	b.plates[x][y] = 1
}

func (b *Board) Clear() {
	for i := range b.plates {
		for j := range b.plates[i] {
			b.plates[i][j] = 0
		}
	}
}

func (b *Board) IsDownCollision(x, y int, shape piece.Shape) bool {
	if y+shape.GetHeight() >= b.GetHeight() {
		return true
	}

	for c := 0; c < shape.GetWidth(); c++ {
		for r := 0; r < shape.GetHeight(); r++ {
			if shape[r][c] == 1 && !b.IsEmpty(y+r, c+x) {
				return true
			}
		}
	}
	return false
}

func (b *Board) IsLeftCollision(x, y int, shape piece.Shape) bool {
	if x < 0 {
		return true
	}
	for r := 0; r < shape.GetHeight(); r++ {
		if shape[r][0] == 1 && !b.IsEmpty(r+y, x) {
			return true
		}
	}
	return false
}

func (b *Board) IsRightCollision(x, y int, shape piece.Shape) bool {
	if x+shape.GetWidth() > b.GetWidth() {
		return true
	}
	lastColumn := shape.GetWidth() - 1
	for r := 0; r < shape.GetHeight(); r++ {
		if shape[r][lastColumn] == 1 && !b.IsEmpty(r+y, x+lastColumn) {
			return true
		}
	}
	return false
}

func (b *Board) AddPieceToBoard(p piece.Piece) {
	for i := 0; i < p.GetHeight(); i++ {
		for j := 0; j < p.GetWidth(); j++ {
			if p.Shape[i][j] == 1 {
				b.SetBlock(p.Position.Y+i, p.Position.X+j)
				b.colors[p.Position.Y+i][p.Position.X+j] = makeLighter(p.Color, .5)

			}
		}
	}
}

func (b *Board) RemoveLine(line int) {
	width := b.GetWidth()
	b.plates = append(b.plates[:line], b.plates[line+1:]...)
	b.plates = append(make([][]int, 1), b.plates...)
	b.plates[0] = make([]int, width)

	b.colors = append(b.colors[:line], b.colors[line+1:]...)
	b.colors = append(make([][]color.Color, 1), b.colors...)
	b.colors[0] = make([]color.Color, width)
	for i := 0; i < width; i++ {
		b.colors[0][i] = EMPTY_COLOR
	}
}

func (b *Board) CheckLines(start, end int) int {
	lines := 0
	for i := start; i < end; i++ {
		if b.IsLineComplete(i) {
			b.RemoveLine(i)
			lines++
		}
	}
	return lines
}

func (b *Board) IsLineComplete(line int) bool {
	for i := 0; i < b.GetWidth(); i++ {
		if b.plates[line][i] == 0 {
			return false
		}
	}
	return true
}

func (b *Board) GetColor(x, y int) color.Color {
	return b.colors[x][y]
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
