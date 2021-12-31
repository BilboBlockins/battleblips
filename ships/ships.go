package ships

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Shiper interface {
	GetShip() Ship
	SetShip(Ship) Shiper
	SetH(s tcell.Screen, x int, y int, style tcell.Style)
	SetV(s tcell.Screen, x int, y int, style tcell.Style)
}

type Ship struct {
	Name string
	Type int
	X int
	Y int
	Coords []Coord
	Length int
	Health int
	Visible bool
	Vertical bool
	Style tcell.Style
	SunkStyle tcell.Style
}

type Destroyer struct {
	Ship
}

type Submarine struct {
	Ship
}

type Cruiser struct {
	Ship
}

type Battleship struct {
	Ship
}

type Carrier struct {
	Ship
}

type Coord struct {
	X int
	Y int
}

func (des Destroyer) SetH(s tcell.Screen, x int, y int, style tcell.Style) {
	setHorizontal(s, x-2, y, style, "<{⫼[H]⫼}")
}

func (des Destroyer) SetV(s tcell.Screen, x int, y int, style tcell.Style) {
	s.SetContent(x, y, '▤', nil, style)
	s.SetContent(x, y + 1, 'I', nil, style)
	s.SetContent(x, y + 2, '▤', nil, style)
	s.SetContent(x - 1, y, '⎧', nil, style)
	s.SetContent(x + 1, y, '⎫', nil, style)
	s.SetContent(x - 1, y + 1, '|', nil, style)
	s.SetContent(x + 1, y + 1, '|', nil, style)
	s.SetContent(x - 1, y + 2, '⎩', nil, style)
	s.SetContent(x + 1, y + 2, '⎭', nil, style)
	s.SetContent(x, y + 3, 'V', nil, style)
}

func (des Destroyer) GetShip() Ship {
	return des.Ship
}

func (des Destroyer) SetShip(s Ship) Shiper {
	des.Ship = s
	return des
}

func (sub Submarine) SetH(s tcell.Screen, x int, y int, style tcell.Style) {
	setHorizontal(s, x-1, y, style, "(⫼]⫼[H]⫼[⫼)")
}

func (sub Submarine) SetV(s tcell.Screen, x int, y int, style tcell.Style) {
	s.SetContent(x, y, '⫼', nil, style)
	s.SetContent(x, y + 1, '⫼', nil, style)
	s.SetContent(x, y + 2, 'I', nil, style)
	s.SetContent(x, y + 3, '⫼', nil, style)
	s.SetContent(x, y + 4, '⫼', nil, style)
	s.SetContent(x - 1, y, '⎧', nil, style)
	s.SetContent(x + 1, y, '⎫', nil, style)
	s.SetContent(x - 1, y + 1, '|', nil, style)
	s.SetContent(x + 1, y + 1, '|', nil, style)
	s.SetContent(x - 1, y + 2, '|', nil, style)
	s.SetContent(x + 1, y + 2, '|', nil, style)
	s.SetContent(x - 1, y + 3, '|', nil, style)
	s.SetContent(x + 1, y + 3, '|', nil, style)
	s.SetContent(x - 1, y + 4, '⎩', nil, style)
	s.SetContent(x + 1, y + 4, '⎭', nil, style)
}

func (sub Submarine) GetShip() Ship {
	return sub.Ship
}

func (sub Submarine) SetShip(s Ship) Shiper {
	sub.Ship = s
	return sub
}

func (cru Cruiser) SetH(s tcell.Screen, x int, y int, style tcell.Style) {
	setHorizontal(s, x-2, y, style, "<{⫼]⫼[H]⫼[⫼}")
}

func (cru Cruiser) SetV(s tcell.Screen, x int, y int, style tcell.Style) {
	s.SetContent(x, y, '▤', nil, style)
	s.SetContent(x, y + 1, '▤', nil, style)
	s.SetContent(x, y + 2, 'I', nil, style)
	s.SetContent(x, y + 3, '▤', nil, style)
	s.SetContent(x, y + 4, '▤', nil, style)
	s.SetContent(x - 1, y, '⎧', nil, style)
	s.SetContent(x + 1, y, '⎫', nil, style)
	s.SetContent(x - 1, y + 1, '|', nil, style)
	s.SetContent(x + 1, y + 1, '|', nil, style)
	s.SetContent(x - 1, y + 2, '|', nil, style)
	s.SetContent(x + 1, y + 2, '|', nil, style)
	s.SetContent(x - 1, y + 3, '|', nil, style)
	s.SetContent(x + 1, y + 3, '|', nil, style)
	s.SetContent(x - 1, y + 4, '⎩', nil, style)
	s.SetContent(x + 1, y + 4, '⎭', nil, style)
	s.SetContent(x, y + 5, 'V', nil, style)
}

func (cru Cruiser) GetShip() Ship {
	return cru.Ship
}

func (cru Cruiser) SetShip(s Ship) Shiper {
	cru.Ship = s
	return cru
}

func (bat Battleship) SetH(s tcell.Screen, x int, y int, style tcell.Style) {
	setHorizontal(s, x-2, y, style, "<{⫼]⫼[H]⫼[H]⫼[⫼}")
}

func (bat Battleship) SetV(s tcell.Screen, x int, y int, style tcell.Style) {
	s.SetContent(x, y, '▤', nil, style)
	s.SetContent(x, y + 1, '▤', nil, style)
	s.SetContent(x, y + 2, 'I', nil, style)
	s.SetContent(x, y + 3, '▤', nil, style)
	s.SetContent(x, y + 4, 'I', nil, style)
	s.SetContent(x, y + 5, '▤', nil, style)
	s.SetContent(x, y + 6, '▤', nil, style)
	s.SetContent(x - 1, y, '⎧', nil, style)
	s.SetContent(x + 1, y, '⎫', nil, style)
	s.SetContent(x - 1, y + 1, '|', nil, style)
	s.SetContent(x + 1, y + 1, '|', nil, style)
	s.SetContent(x - 1, y + 2, '|', nil, style)
	s.SetContent(x + 1, y + 2, '|', nil, style)
	s.SetContent(x - 1, y + 3, '|', nil, style)
	s.SetContent(x + 1, y + 3, '|', nil, style)
	s.SetContent(x - 1, y + 4, '|', nil, style)
	s.SetContent(x + 1, y + 4, '|', nil, style)
	s.SetContent(x - 1, y + 5, '|', nil, style)
	s.SetContent(x + 1, y + 5, '|', nil, style)
	s.SetContent(x - 1, y + 6, '⎩', nil, style)
	s.SetContent(x + 1, y + 6, '⎭', nil, style)
	s.SetContent(x, y + 7, 'V', nil, style)
}

func (bat Battleship) GetShip() Ship {
	return bat.Ship
}

func (bat Battleship) SetShip(s Ship) Shiper {
	bat.Ship = s
	return bat
}

func (car Carrier) SetH(s tcell.Screen, x int, y int, style tcell.Style) {
	setHorizontal(s, x-1, y, style, "[⫼]⫼[⫼]⫼[H]⫼[H]⫼[⫼]")
}

func (car Carrier) SetV(s tcell.Screen, x int, y int, style tcell.Style) {
	s.SetContent(x, y, '▤', nil, style)
	s.SetContent(x, y + 1, '▤', nil, style)
	s.SetContent(x, y + 2, 'I', nil, style)
	s.SetContent(x, y + 3, '▤', nil, style)
	s.SetContent(x, y + 4, '▤', nil, style)
	s.SetContent(x, y + 5, 'I', nil, style)
	s.SetContent(x, y + 6, '▤', nil, style)
	s.SetContent(x, y + 7, '▤', nil, style)
	s.SetContent(x, y + 8, '▤', nil, style)
	s.SetContent(x - 1, y, '⌈', nil, style)
	s.SetContent(x + 1, y, '⌉', nil, style)
	s.SetContent(x - 1, y + 1, '∣', nil, style)
	s.SetContent(x + 1, y + 1, '|', nil, style)
	s.SetContent(x - 1, y + 2, '|', nil, style)
	s.SetContent(x + 1, y + 2, '|', nil, style)
	s.SetContent(x - 1, y + 3, '|', nil, style)
	s.SetContent(x + 1, y + 3, '|', nil, style)
	s.SetContent(x - 1, y + 4, '|', nil, style)
	s.SetContent(x + 1, y + 4, '|', nil, style)
	s.SetContent(x - 1, y + 5, '|', nil, style)
	s.SetContent(x + 1, y + 5, '|', nil, style)
	s.SetContent(x - 1, y + 6, '|', nil, style)
	s.SetContent(x + 1, y + 6, '|', nil, style)
	s.SetContent(x - 1, y + 7, '|', nil, style)
	s.SetContent(x + 1, y + 7, '|', nil, style)
	s.SetContent(x - 1, y + 8, '⌊', nil, style)
	s.SetContent(x + 1, y + 8, '⌋', nil, style)
}

func (car Carrier) GetShip() Ship {
	return car.Ship
}

func (car Carrier) SetShip(s Ship) Shiper {
	car.Ship = s
	return car
}


func setHorizontal(s tcell.Screen, x int, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}