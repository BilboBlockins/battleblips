package draw

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"battleblips/ocean"
	"battleblips/ships"
	"battleblips/game"
	"battleblips/utils"
)

func Update(s tcell.Screen, g *game.Game) {
	s.Clear()
	drawGrid(s, g.Options.DefStyle, g.Options.AltStyle, g.PlBoard)
	drawGrid(s, g.Options.DefStyle, g.Options.AltStyle, g.OpBoard)
	drawShips(s, g.PlShips, g.PlBoard)
	drawShips(s, g.OpShips, g.OpBoard)

	if g.Cursor.Board == 0 {
		xpos, ypos := 
		utils.GetCoordPosition(g.Cursor.X, g.Cursor.Y, g.PlBoard.Offsetx, g.PlBoard.Offsety)
		DrawTarget(s, g.Options.DefStyle, xpos, ypos)
		if g.Cursor.PlacingShip {
			currentShip := g.Cursor.HasShip
			shipVal := currentShip.GetShip()
			style := shipVal.Style
			if !utils.CursorShipOnBoard(g) ||
				utils.CursorShipOverlapping(g) {
				style = g.Options.HStyle
			}
			if shipVal.Vertical {
				currentShip.SetV(s, xpos, ypos, style)
			} else {
				currentShip.SetH(s, xpos, ypos, style)
			}
		}
	}
	if g.Cursor.Board == 1 {
		xpos, ypos := 
		utils.GetCoordPosition(g.Cursor.X, g.Cursor.Y, g.OpBoard.Offsetx, g.OpBoard.Offsety)
		DrawTarget(s, g.Options.DefStyle, xpos, ypos)
	}

	drawBoardMarkers(s, g.Options.HStyle, g.Options.MStyle, g.Fired)
	drawBoardMarkers(s, g.Options.HStyle, g.Options.MStyle, g.Incoming)
	
	s.Show()
}

func drawShips(s tcell.Screen, plShips []ships.Shiper, o ocean.Ocean) {
	for _, ship := range plShips {
		shipVal := ship.GetShip()
		if shipVal.X != -1 && shipVal.Y != -1 && shipVal.Visible {
			xpos, ypos := 
			utils.GetCoordPosition(shipVal.X, shipVal.Y, o.Offsetx, o.Offsety)
			style := shipVal.Style
			if shipVal.Health == 0 {
				style = shipVal.SunkStyle
			}
			if shipVal.Vertical {
				ship.SetV(s, xpos, ypos, style)
			} else {
				ship.SetH(s, xpos, ypos, style)
			}
		}
	}
}

func drawBoardMarkers(s tcell.Screen, Hstyle tcell.Style, Mstyle tcell.Style, o ocean.Ocean) {
	for i, r := range o.Grid {
		for j, v := range r {
			if v > 0 {
				xpos, ypos := utils.GetCoordPosition(j,i, o.Offsetx, o.Offsety)
				if v == 1 {
					EmitStr(s, xpos, ypos, Mstyle, "●")
				}
				if v == 2 {
					EmitStr(s, xpos, ypos, Hstyle, "●")
				}
			}

		}
	}
}

func DrawTarget(s tcell.Screen, style tcell.Style, x int, y int) {
	EmitStr(s, x-2, y-1, style, "┌")
	EmitStr(s, x-1, y-1, style, "─")
	EmitStr(s, x+2, y-1, style, "┐")
	EmitStr(s, x+1, y-1, style, "─")
	EmitStr(s, x-2, y+1, style, "└")
	EmitStr(s, x-1, y+1, style, "─")
	EmitStr(s, x+2, y+1, style, "┘")
	EmitStr(s, x+1, y+1, style, "─")
	EmitStr(s, x, y, style, "+")
}


func drawGrid(s tcell.Screen, style tcell.Style, lineStyle tcell.Style, o ocean.Ocean) {
	hlen := o.Rcorner - o.Offsetx
	vlen := o.Bcorner - o.Offsety

	s.SetContent(o.Offsetx, o.Offsety, tcell.RuneULCorner, nil, style)
	s.SetContent(o.Rcorner, o.Offsety, tcell.RuneURCorner, nil, style)
	s.SetContent(o.Offsetx, o.Bcorner, tcell.RuneLLCorner, nil, style)
	s.SetContent(o.Rcorner, o.Bcorner, tcell.RuneLRCorner, nil, style)

	for r:=0; r<o.Dim+1; r++ {
		ax := o.Offsetx + r * o.Cellx + 2
		ay := o.Offsety - 1
		for i:=0; i<hlen-1; i++ {
			x := 1 + i + o.Offsetx
			y := o.Celly * r + o.Offsety
			if r == 0 || r == o.Dim {
				s.SetContent(x, y, tcell.RuneHLine, nil, style)
			} else {
				s.SetContent(x, y, tcell.RuneHLine, nil, lineStyle)
			}
		}
		if r < o.Dim {
			s.SetContent(ax, ay, rune(65 + r), nil, style)
		}
	}

	for c:=0; c<o.Dim+1; c++ {
		d := c % 10
		td := c / 10 
		nx := o.Offsetx - 2
		ny := o.Offsety + c * o.Celly + 1
		for i:=0; i<vlen-1; i++ {
			x := o.Cellx * c + o.Offsetx
			y := 1 + i + o.Offsety
			if c == 0 || c == o.Dim {
				s.SetContent(x, y, tcell.RuneVLine, nil, style)
			} else {
				s.SetContent(x, y, tcell.RuneVLine, nil, lineStyle)
			}
		}
		if c < o.Dim {
			s.SetContent(nx, ny, rune(48+d), nil, style)
			if td > 0 {
				s.SetContent(nx-1, ny, rune(48+td), nil, style)
			}
		}
	}
}

func EmitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
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