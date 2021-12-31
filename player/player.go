package player

import (
	"fmt"
	"os"
	// "math/rand"
	"time"
	"battleblips/draw"
	"battleblips/utils"
	"battleblips/ships"
	"battleblips/game"
	"github.com/gdamore/tcell/v2"
)


type Player struct {
	Name string
}

type Ai struct {
	Name string
	Color tcell.Style
	Difficulty int
}

func AiPlaceShips(s tcell.Screen, g *game.Game) {
	opShips := g.OpShips
	setShips := []ships.Shiper{}
	shipIndex := 0

	for len(opShips) != len(setShips) {
		currentShip := opShips[shipIndex]
		shipVal := currentShip.GetShip()
		randX := utils.RandIntRange(0, g.Options.GridDim)
		randY := utils.RandIntRange(0, g.Options.GridDim)
		vertVal := utils.RandIntRange(0,2)
		shipVal.Vertical = (vertVal != 0)

		//todo: could probably just pass ship val here too
		if utils.ShipOnBoard(randX, randY, shipVal, g.OpBoard) && 
			!utils.ShipOverlapping(randX, randY, shipVal, g.OpBoard) {
			//remove - just for testing
			time.Sleep(100 * time.Millisecond)
			shipVal.X = randX
			shipVal.Y = randY
			//make all enemy ships visible for testing
			shipVal.Visible = false
			currentShip = currentShip.SetShip(shipVal)
			setShips =  append(setShips, currentShip)
			g.OpShips = setShips
			//todo: just pass ship here?
			utils.SetShipOnBoard(shipVal.X, shipVal.Y, shipVal, g.OpBoard)
			shipIndex++
			draw.Update(s, g)
		}
	}
}

func AiFire(s tcell.Screen, g *game.Game) {
	difficulty := g.Options.Difficulty
	switch difficulty {
	case 0:
		AiRandomShot(s,g)
	case 1:
	case 2:	
	}
	return
}

func AiRandomShot(s tcell.Screen, g *game.Game) {
	foundNewCoord := false
	randX := utils.RandIntRange(0, g.Options.GridDim)
	randY := utils.RandIntRange(0, g.Options.GridDim)
	g.Cursor.Board = 0

	for !foundNewCoord {
		randX = utils.RandIntRange(0, g.Options.GridDim)
		randY = utils.RandIntRange(0, g.Options.GridDim)
		g.Cursor.X = randX 
		g.Cursor.Y = randY
		draw.Update(s, g)
		if g.Incoming.Grid[randY][randX] == 0 && utils.CheckSurrounding(randX, randY, 4, g.Incoming) {
			foundNewCoord = true
		}
		time.Sleep(300 * time.Millisecond)
		
	}

	if _, hit, open := utils.CheckShot(randX, randY, g); hit && open {
		hitShip, _ := utils.ShipHit(randX, randY, 0, g)
		draw.Update(s, g)
		if hitShip.Health == 0 {
			draw.EmitStr(s, 1, 30, g.Options.DefStyle, fmt.Sprintf("Our %v has been sunk Cap'n!", hitShip.Name))
		} else {
			draw.EmitStr(s, 1, 30, g.Options.DefStyle, fmt.Sprintf("Our %v has been hit!", hitShip.Name))
		}
	} else {
		draw.Update(s, g)
		draw.EmitStr(s, 1, 30, g.Options.DefStyle, "Enemy Missed Cap'n")
	}
	draw.EmitStr(s, 4, 27, g.Options.DefStyle, fmt.Sprintf("Incoming fired at b%v %v %v!", g.Cursor.Board, string(rune(65 + randX)), randY))
	s.Show()
	time.Sleep(1400 * time.Millisecond)
}


func PlaceShips(s tcell.Screen, g *game.Game) {
	bn := 0
	xcoord := g.Cursor.X
	ycoord := g.Cursor.Y
	plShips := g.PlShips
	setShips := []ships.Shiper{}
	shipIndex := 0
	currentShip := plShips[shipIndex]
	g.PlShips = setShips

	// Placing:
	for len(plShips) != len(setShips) {

		switch ev := s.PollEvent().(type) {
		case *tcell.EventMouse:
			mx, my := ev.Position()
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if bn > 0 {
					switch bn {
					case 1:
						if utils.ClickOnBoard(mx, my, g.PlBoard) {
							_, _, xcoord, ycoord = 
							utils.MapEvCoord(g.Options.GridDim, mx, my, g.PlBoard.Offsetx, g.PlBoard.Offsety)
							g.Cursor.Board = 0
						} 
						if utils.ClickOnBoard(mx, my, g.OpBoard) {
							_, _, xcoord, ycoord = 
							utils.MapEvCoord(g.Options.GridDim, mx, my, g.OpBoard.Offsetx, g.OpBoard.Offsety)
							g.Cursor.Board = 1
						}
						g.Cursor.X = xcoord
						g.Cursor.Y = ycoord
						draw.Update(s,g)
						ShowPlaceInstruct(s, g, plShips)
					case 2:

					case 3:	
					}
				}
				bn = 0
			case tcell.Button1:
				bn = 1
			case tcell.Button2:
				bn = 2
			case tcell.Button3:
				bn = 3
			}
		case *tcell.EventResize:
			draw.Update(s, g)
			ShowPlaceInstruct(s, g, plShips)
			s.Sync()
		case *tcell.EventKey:
			if ev.Rune() == 'p' || ev.Rune() == 'P' {
				g.Cursor.HasShip = currentShip
				g.Cursor.PlacingShip = true
				draw.Update(s,g)
			}

			if ev.Rune() == ' ' {

				if g.Cursor.PlacingShip {
					v := currentShip.GetShip()
					v.Vertical = !v.Vertical
					currentShip = currentShip.SetShip(v)
					g.Cursor.HasShip = currentShip
					draw.Update(s,g)
				}
			}

			switch ev.Key() {
			case tcell.KeyUp:
				if ycoord > 0 {
					ycoord--
					g.Cursor.Y = ycoord
				}
				draw.Update(s, g)
			case tcell.KeyDown:
				if ycoord < g.Options.GridDim - 1 {
					ycoord++
					g.Cursor.Y = ycoord
				}
				draw.Update(s, g)
			case tcell.KeyRight:
				if xcoord < g.Options.GridDim - 1 {
					xcoord++
					g.Cursor.X = xcoord
				} else {
					if g.Cursor.Board == 0 {
						g.Cursor.Board = 1
						xcoord = 0
						g.Cursor.X = xcoord
					}

				}
				draw.Update(s, g)
			case tcell.KeyLeft:
				if xcoord == 0 && g.Cursor.Board == 1 {
					g.Cursor.Board = 0
					xcoord = g.Options.GridDim
					g.Cursor.X = xcoord
				}
				if xcoord > 0 {
					xcoord--
					g.Cursor.X = xcoord
				}
				draw.Update(s, g)
			case tcell.KeyEnter:

				if g.Cursor.Board == 0 &&
					utils.CursorShipOnBoard(g) && 
					!utils.CursorShipOverlapping(g) {
					shipVal := currentShip.GetShip()
					shipVal.X = xcoord
					shipVal.Y = ycoord
					currentShip = currentShip.SetShip(shipVal)
					setShips =  append(setShips, currentShip)
					g.PlShips = setShips


					//todo: could just pass ship here
					utils.SetShipOnBoard(shipVal.X, shipVal.Y, shipVal, g.PlBoard)

					if shipIndex < len(plShips)-1 {
						shipIndex++
						currentShip = plShips[shipIndex]
						g.Cursor.HasShip = currentShip
					}

					if len(plShips) == len(setShips) {
						g.Cursor.PlacingShip = false
					}
					
				}

				draw.Update(s,g)

				draw.EmitStr(s, 4, 27, g.Options.DefStyle, fmt.Sprintf("Hit on b%v %v %v", g.Cursor.Board, string(rune(65 + xcoord)), ycoord))
				s.Show()

			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			}
			//todo: doesn't show instructions when placing destroyer
			ShowPlaceInstruct(s, g, plShips)
					
		}
	}
}

//probably should move to draw when fixed
func ShowPlaceInstruct(s tcell.Screen, g *game.Game, plShips []ships.Shiper) {
	if g.Cursor.PlacingShip && len(plShips) > 0 {
		draw.EmitStr(
			s, 4, 50, g.Options.DefStyle, 
			fmt.Sprintf("Placing %v | Click or Use Arrow Keys to Move | [SPACE] to Rotate | [ENTER] to Set", 
			g.Cursor.HasShip.GetShip().Name))
		
			s.Show()
	} else {
		draw.EmitStr(s, 4, 50, g.Options.DefStyle, fmt.Sprintf("Hit [pps] to Deploy Your Shit..."))
		s.Show()
	}
}