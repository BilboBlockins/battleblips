package utils

import (
	"fmt"
	"time"
	"math/rand"
	"battleblips/ships"
	"battleblips/ocean"
	"battleblips/game"
)

//Gets random seed int in range max - min
func RandIntRange(min int, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min) + min
}

//Checks if a click is withing an ocean bounds
func ClickOnBoard(mx int, my int, o ocean.Ocean) bool {
	return mx > o.Offsetx && my > o.Offsety && mx < o.Rcorner && my < o.Bcorner
}

//Returns the normalized screen position of a mouse click and its game grid coordinates
func MapEvCoord(dim int, mx int, my int, offsetx int, offsety int) (int, int, int, int) {
	row := mx
	col := my
	rowStart := offsetx + 2
	colStart := offsety + 1
	if row % 4 != 0 {
		row = (mx / 4) * 4 + 2
	} else {
		row += 2
	}
	if col % 2 != 0 {
		col = (my / 2) * 2 + 1
	} else {
		col += 1
	}
	coordx := (row - rowStart) / 4
	coordy := (col - colStart) / 2
	return row, col, coordx, coordy
}

//Returns returns the screen position of the center of a game grid coordinate
func GetCoordPosition(x int, y int, offsetx int, offsety int) (int, int) {
	rowStart := offsetx + 2
	colStart := offsety + 1
	xpos :=  4 * x + rowStart
	ypos :=  2 * y + colStart
	return xpos, ypos
}

//Checks if the ship attached to the cursor is within the player board bounds
func CursorShipOnBoard(g *game.Game) bool {
	coordx, coordy := g.Cursor.X, g.Cursor.Y
	shipVal := g.Cursor.HasShip.GetShip()
	return ShipOnBoard(coordx, coordy, shipVal, g.PlBoard)
}

//Checks if a ship is within an ocean bounds
func ShipOnBoard(coordx int, coordy int, shipVal ships.Ship, o ocean.Ocean) bool {
	if shipVal.Vertical {
		return coordy <= o.Dim - shipVal.Length
	} else {
		return coordx <= o.Dim - shipVal.Length
	}
}

//Checks if the ship attached to the cursor is overlapping a set player ship
func CursorShipOverlapping(g *game.Game) bool {
	coordx, coordy := g.Cursor.X, g.Cursor.Y
	shipVal := g.Cursor.HasShip.GetShip()
	return ShipOverlapping(coordx, coordy, shipVal, g.PlBoard)
}

//Checks if a ship is overlapping another ship in an ocean
func ShipOverlapping(coordx int, coordy int, shipVal ships.Ship, o ocean.Ocean) bool {
	for i := 0; i < shipVal.Length; i++ {
		if shipVal.Vertical {
			if o.Grid[coordy + i][coordx] > 0 {
				return true
			}
		} else {
			if o.Grid[coordy][coordx + i] > 0 {
				return true
			}

		}
	}
	return false
}

// good - not sure about mplyr
//Saves a ship position in an ocean
func SetShipOnBoard(coordx int, coordy int, shipVal ships.Ship, o ocean.Ocean) {
	for i:= 0; i<shipVal.Length; i++ {
		if shipVal.Vertical {
			o.Grid[coordy + i][coordx] = shipVal.Type
		} else {
			o.Grid[coordy][coordx + i] = shipVal.Type
		}
	}
}

//todo: could split out - multi-player would only need to check incoming...
//probably don't need the hit bool since this should only be for 

//Checks if a ship has been hit and decrements its health if it has, returns the hit ship
func ShipHit(coordx int, coordy int, board int, g *game.Game) (ships.Ship, bool) {
	fleet := g.OpShips 
	if board == 0 { fleet = g.PlShips }
	for si, ship := range fleet {
		shipVal := ship.GetShip()
		x := shipVal.X
		y := shipVal.Y
		for i := 0; i<shipVal.Length; i++ {
			if x == coordx && y == coordy {
				shipVal.Health -= 1
				//Sunk ships should always be visible
				if shipVal.Health == 0 { shipVal.Visible = true }
				ship = ship.SetShip(shipVal)
				
				if g.Cursor.Board == 1 {
					g.OpShips[si] = ship
				} else {
					g.PlShips[si] = ship
				}
				return shipVal, true 	
			}
			if shipVal.Vertical {
				y++
			} else {
				x++
			}
		}
	}
	return ships.Ship{}, false
}

//todo: clean up, unlikely that we need bother with the board stuff
//should break out interior logic into a more general "RecordShot" function
//should also return named values to make it more clear...
//but this is fine for now while I work out the menu, game, and computer play

//Checks for a hit or miss in an open coord and records the shot in the incoming or firing grids
func CheckShot(coordx int, coordy int, g *game.Game) (board int, hit bool, open bool) {
	if g.Cursor.Board == 0 {
		if g.Incoming.Grid[coordy][coordx] == 0 {
			if g.PlBoard.Grid[coordy][coordx] > 0 {
				g.Incoming.Grid[coordy][coordx] = 2
				return 0, true, true
			} else {
				g.Incoming.Grid[coordy][coordx] = 1
				return 0, false, true
			}
		} else {
			return 0, false, false
		}
	} else {
		if g.Fired.Grid[coordy][coordx] == 0 {
			if g.OpBoard.Grid[coordy][coordx] > 0 {
				g.Fired.Grid[coordy][coordx] = 2
				return 1, true, true
			} else {
				g.Fired.Grid[coordy][coordx] = 1
				return 1, false, true
			}
		} else {
			return 1, false, false
		}
	}
}

//Issue: hits between ships can make only 1 space openings
//this alg only makes sense if there are no ships hit that are not sunk...
//need to do something different if there are already hits on the board or just treat squares with adjacent hits differently
func CheckSurrounding(coordx int, coordy int, threshold int, o ocean.Ocean) bool {
	openCount := 1
	for t := coordy-1; t>=0; t-- {
		if o.Grid[t][coordx] == 0 {
			openCount++
		}
		if openCount == threshold {
			return true
		}
	}
	openCount = 1
	fmt.Println("		top blocked")
	for b := coordy+1; b<o.Dim; b++ {
		if o.Grid[b][coordx] == 0 {
			openCount++
		}
		if openCount == threshold {
			return true
		}
	}
	openCount = 1
	fmt.Println("		bottom blocked")
	for l := coordx-1; l>=0; l-- {
		if o.Grid[coordy][l] == 0 {
			openCount++
		}
		if openCount == threshold {
			return true
		}
	}
	openCount = 1
	fmt.Println("		left blocked")
	for r := coordx+1; r<o.Dim; r++ {
		if o.Grid[coordy][r] == 0 {
			openCount++
		}
		if openCount == threshold {
			return true
		}
	}
	fmt.Println("		all blocked")
	return false
}