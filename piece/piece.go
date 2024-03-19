package piece

import "image/color"

type (
	Shape    [][]int
	Position struct {
		X, Y int
	}
	Piece struct {
		Shape    Shape
		Position Position
		Color    color.Color
	}
)

var (
	SQUARE = [][]int{
		{1, 1},
		{1, 1},
	}
	L = [][]int{
		{1, 0},
		{1, 0},
		{1, 1},
	}
	Z = [][]int{
		{1, 1, 0},
		{0, 1, 1},
	}
	T = [][]int{
		{1, 1, 1},
		{0, 1, 0},
	}
	I = [][]int{
		{1},
		{1},
		{1},
		{1},
	}

	PIECES = []Shape{SQUARE, L, Z, T, I}
	COLORS = []color.Color{color.RGBA{51, 153, 255, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{255, 255, 0, 255}, color.RGBA{0, 255, 255, 255}}
)

func CreatePiece(s Shape, p Position, c color.Color) *Piece {
	return &Piece{s, p, c}
}

func ListOfPieces() []Piece {
	pieces := make([]Piece, 0)
	for i, shape := range PIECES {
		pieces[i] = Piece{shape, Position{0, 0}, COLORS[i]}
	}
	return pieces
}

func (p *Piece) Rotate() {
	p.Shape = Rotate(p.Shape)
}

func Rotate(s Shape) Shape {
	rotated := make(Shape, len(s[0]))
	for i := range rotated {
		rotated[i] = make([]int, len(s))
	}

	for i := range s {
		for j := range s[i] {
			rotated[j][len(s)-1-i] = s[i][j]
		}
	}
	return rotated
}

func (p *Piece) GetWidth() int {
	return len(p.Shape[0])
}

func (p *Piece) GetHeight() int {
	return len(p.Shape)
}

func (p *Piece) GetShape() Shape {
	return p.Shape
}

func (s Shape) GetWidth() int {
	return len(s[0])
}

func (s Shape) GetHeight() int {
	return len(s)
}
