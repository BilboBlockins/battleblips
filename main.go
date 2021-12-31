package main

import (
	"fmt"
	"os"
	"os/exec"
	// "runtime"
	// "math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	// "github.com/gdamore/tcell/v2/encoding"
	// "github.com/mattn/go-runewidth"
	// "battleblips/ships"
	"battleblips/utils"
	"battleblips/game"
	"battleblips/draw"
	"battleblips/player"
)

func main() {

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorGreen)

	altStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorSeaGreen)

	hStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorRed)

	mStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	opts:= game.GameOptions{
		Players: 1,
		Difficulty: 0,
		GridDim: 10,
		DefStyle: defStyle,
		AltStyle: altStyle,
		HStyle: hStyle,
		MStyle: mStyle,
		PlStyle: defStyle,
		OpStyle: defStyle,
	}

	shipList := []int{1,2,3,4,5}

	g := game.Init(opts)
	g = game.InitShips(g, shipList)

	shell, s := game.InitScreen()

	s.SetStyle(g.Options.DefStyle)
	s.EnableMouse()

	draw.Update(s, g)

	player.PlaceShips(s, g)
	player.AiPlaceShips(s, g)



	// gridScreen(s)
	bn := 0
	xcoord, ycoord := g.Cursor.X, g.Cursor.Y
	aiTurn := false


// Event Crap
//==============================
	for {
		g.Cursor.X = xcoord
		g.Cursor.Y = ycoord
		g.Cursor.Board = 1
		draw.Update(s, g)
		draw.EmitStr(s, 1, 30, g.Options.DefStyle, "Awaiting orders...")
		s.Show()
		if !aiTurn {
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
				s.Sync()
			case *tcell.EventKey:

				if ev.Rune() == 'Z' || ev.Rune() == 'z' {
					// CtrlZ doesn't really suspend the process, but we use it to execute a subshell.
					if err := s.Suspend(); err == nil {
						fmt.Printf("Executing shell (%s -l)...\n", shell)
						fmt.Printf("Battleblips Paused: press [CTRL+D] to exit subshell and return to game.\n")
						c := exec.Command(shell, "-l" ) // NB: -l works .exe too (ignored)
						c.Stdin = os.Stdin
						c.Stdout = os.Stdout
						c.Stderr = os.Stderr
						c.Run()
			
						if err := s.Resume(); err != nil {
							panic("failed to resume: " + err.Error())
						}
						draw.Update(s,g)
						s.Sync()
					}
				}	
				// switch ev.Rune() {
				// case 'P':
				// 	draw.EmitStr(s, 4, 28, g.Options.DefStyle, "P or p was pressed")
				// 	s.Show()
				// case 'D':
				// 	draw.EmitStr(s, 4, 28, g.Options.DefStyle, "D or d was pressed")
				// 	s.Show()
				// case 'p':
				// 	draw.EmitStr(s, 4, 28, g.Options.DefStyle, "P or p was pressed")
				// 	s.Show()
				// case 'd':
				// 	draw.EmitStr(s, 4, 28, g.Options.DefStyle, "D or d was pressed")
				// 	s.Show()
				// }

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
					
					if board, hit, open := utils.CheckShot(xcoord, ycoord, g); hit && open {
						if board == 0 {
							//should allow self shooting - just for testing
							hitShip, _ := utils.ShipHit(xcoord, ycoord, board, g)
							draw.Update(s, g)
							draw.EmitStr(s, 1, 30, g.Options.DefStyle, fmt.Sprintf("Hit %v, %v", hitShip))
							draw.EmitStr(s, 1, 31, g.Options.DefStyle, "You should fire on yourselves - cap'em!")
						} else {
							hitShip, _ := utils.ShipHit(xcoord, ycoord, board, g)
							draw.Update(s, g)
							if hitShip.Health == 0 {
								draw.EmitStr(s, 1, 30, g.Options.DefStyle, fmt.Sprintf("Enemy %v sunk Cap'n!", hitShip.Name))	
							} else {
								draw.EmitStr(s, 1, 30, g.Options.DefStyle, "Hit Cap'n!")
							}
						}
					} else {
						draw.Update(s, g)
						if open {
							draw.EmitStr(s, 1, 30, g.Options.DefStyle, "Missed Cap'n")
						} else {
							draw.EmitStr(s, 1, 30, g.Options.DefStyle, "Already fired at those coordinates Cap'n")
						}
					}
					
					draw.EmitStr(s, 4, 27, g.Options.DefStyle, fmt.Sprintf("Fired on b%v %v %v", g.Cursor.Board, string(rune(65 + xcoord)), ycoord))
					s.Show()
					aiTurn = true

					time.Sleep(1400 * time.Millisecond)

				case tcell.KeyEscape:
					s.Fini()

					g.PlBoard.Print()
					fmt.Println("")
					for _, ship := range g.PlShips {
						fmt.Println(ship)
					}
					fmt.Println("")
					g.OpBoard.Print()
					fmt.Println("")
					for _, ship := range g.OpShips {
						fmt.Println(ship)
					}
				
					os.Exit(0)
				}		
			}
		} else {
			player.AiFire(s,g)
			aiTurn = false

		}
	}

}







// func inputEvent()


// func gridScreen(s tcell.Screen) {

// 	// cruiser := ships.Cruiser{
// 	// 	ships.Ship{
// 	// 		Name: "Cruiser",
// 	// 		Health: 3,
// 	// 	},
// 	// }

// 	sub := ships.Battleship{
// 		ships.Ship{
// 			Name: "Sub",
// 			Health: 3,
// 		},
// 	}
// 	// w, h := s.Size()
// 	s.Clear()
// 	// drawGrid(s, defStyle, 10, 4, 2)
	

// 	gcx, gcy, gcx1, gcy1 := drawGrid(s, defStyle, 10, 4, 2)
// 	g2cx, g2cy, g2cx1, g2cy1 := drawGrid(s, defStyle, 10, gcx1 + 8, 2) // 8 offset to keep cell line divisible by 4

// 	s.Show()

// 	defStyle := tcell.StyleDefault.
// 		Background(tcell.ColorBlack).
// 		Foreground(tcell.ColorGreen)



// 	s.SetStyle(defStyle)

// 	s.EnableMouse()

// 	bn := 0
// 	mx := -1
// 	my := -1
// 	xpos := -1
// 	ypos := -1
// 	xcoord := 0
// 	ycoord := 0 

// 	// drawIndex := 0
	

// 	for {
// 		switch ev := s.PollEvent().(type) {
// 			case *tcell.EventMouse:
// 				mx, my = ev.Position()
// 				if clickOnBoard(mx, my, gcx, gcy, gcx1, gcy1) || 
// 					clickOnBoard(mx, my, g2cx, g2cy, g2cx1, g2cy1) {
// 					switch ev.Buttons() {
// 					case tcell.ButtonNone:
// 						if bn > 0 {
// 							switch bn {
// 							case 1:
// 								xpos, ypos, xcoord, ycoord = 
// 								mapEvCoord(10,mx,my,4,2)
// 								emitStr(s, 4, 25, defStyle, fmt.Sprintf("Hit on %v %v", string(rune(65 + xcoord)), ycoord))
// 								s.Clear()
								
// 								drawGrid(s, defStyle, 10, 4, 2)
// 								drawCursor(s, xpos, ypos, defStyle)
								
// 								sub.SetV(s, 14, 11, defStyle)
// 								sub.SetH(s, 10, 7, defStyle)
// 								drawBoardMarkers(s, g.Board.Grid, 4, 2)
								
// 								// sub.SetH(s, xpos, ypos, defStyle)
// 							case 2:
// 								xpos, ypos, xcoord, ycoord = 
// 								mapEvCoord(10,mx,my,4,2)
// 								sub.SetV(s, xpos, ypos, defStyle)
// 								emitStr(s, 4, 25, defStyle, fmt.Sprintf("Hit on %v %v", string(rune(65 + xcoord)), ycoord))
// 							case 3:	
// 							}
// 						}
// 						bn = 0
// 					case tcell.Button1:
// 						bn = 1
// 					case tcell.Button2:
// 						bn = 2
// 					case tcell.Button3:
// 						bn = 3
// 					}
					
// 					emitStr(s, 4, 26, defStyle, fmt.Sprintf("Got click on board @ - x: %v, y: %v, button: %v", mx, my, bn))
// 					s.Show()
// 				}
// 			case *tcell.EventResize:
// 				s.Clear()
// 				drawGrid(s, defStyle, 10, 4, 2)
// 				drawGrid(s, defStyle, 10, gcx1 + 8, 2)
// 				s.Sync()
// 			case *tcell.EventKey:

// 				switch ev.Key() {
// 					case tcell.KeyUp:
// 						ycoord--
// 						xpos, ypos = getCoordPosition(xcoord, ycoord, 4,2)


// 						s.Clear()
								
// 						drawGrid(s, defStyle, 10, 4, 2)
// 						drawCursor(s, xpos, ypos, defStyle)
						
// 						sub.SetV(s, 14, 11, defStyle)
// 						sub.SetH(s, 10, 7, defStyle)
// 						drawBoardMarkers(s, g.Board.Grid, 4, 2)

// 						emitStr(s, 4, 32, defStyle, fmt.Sprintf("UP           "))

// 					case tcell.KeyDown:
						
// 						ycoord++
// 						xpos, ypos = getCoordPosition(xcoord, ycoord, 4,2)
// 						s.Clear()
								
// 						drawGrid(s, defStyle, 10, 4, 2)
// 						drawCursor(s, xpos, ypos, defStyle)
						
// 						sub.SetV(s, 14, 11, defStyle)
// 						sub.SetH(s, 10, 7, defStyle)
// 						drawBoardMarkers(s, g.Board.Grid, 4, 2)

// 						emitStr(s, 4, 32, defStyle, fmt.Sprintf("DOWN           "))

// 					case tcell.KeyRight:
// 						xcoord++
// 						xpos, ypos = getCoordPosition(xcoord, ycoord, 4,2)
						
// 						s.Clear()
								
// 						drawGrid(s, defStyle, 10, 4, 2)
// 						drawCursor(s, xpos, ypos, defStyle)
						
// 						sub.SetV(s, 14, 11, defStyle)
// 						sub.SetH(s, 10, 7, defStyle)
// 						drawBoardMarkers(s, g.Board.Grid, 4, 2)

// 						emitStr(s, 4, 32, defStyle, fmt.Sprintf("RIGHT           "))

// 					case tcell.KeyLeft:
// 						xcoord--
// 						xpos, ypos = getCoordPosition(xcoord, ycoord, 4,2)

// 						s.Clear()
								
// 						drawGrid(s, defStyle, 10, 4, 2)
// 						drawCursor(s, xpos, ypos, defStyle)
						
// 						sub.SetV(s, 14, 11, defStyle)
// 						sub.SetH(s, 10, 7, defStyle)
// 						drawBoardMarkers(s, g.Board.Grid, 4, 2)

// 						emitStr(s, 4, 32, defStyle, fmt.Sprintf("LEFT           "))

// 					case tcell.KeyEnter:
// 						g.Board.Grid[xcoord][ycoord] = 1
// 						drawBoardMarkers(s, g.Board.Grid, 4, 2)
// 						// emitStr(s, xpos, ypos, hitStyle, "●")
// 						emitStr(s, 4, 32, defStyle, fmt.Sprintf("ENTER           "))
// 						// drawIndex++
// 					case tcell.KeyEscape:
// 						s.Fini()
// 						os.Exit(0)

// 				}

// 				emitStr(s, 4, 29, defStyle, fmt.Sprintf("got key event", ev.Name() ))
// 				s.Show()
// 		}
// 	}

// }
