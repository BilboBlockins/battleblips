package game

import (
	"fmt"
	"os"
	// "os/exec"
	"runtime"
	// "math/rand"
	// "time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"battleblips/ships"
	"battleblips/ocean"
)

type GameOptions struct {
	Players int
	Difficulty int
	GridDim int
	DefStyle tcell.Style
	AltStyle tcell.Style
	HStyle tcell.Style
	MStyle tcell.Style
	PlStyle tcell.Style
	OpStyle tcell.Style
}

type Game struct {
	Options GameOptions
	PlBoard ocean.Ocean
	OpBoard ocean.Ocean
	Fired ocean.Ocean
	Incoming ocean.Ocean
	PlShips []ships.Shiper
	OpShips []ships.Shiper
	Cursor Cursor
}

type Cursor struct {
	X int
	Y int
	Board int
	HasShip ships.Shiper
	PlacingShip bool
}

func Init(opts GameOptions) (*Game) {
	bo := ocean.NewOcean(opts.GridDim, 4, 2)
	ob := ocean.NewOcean(opts.GridDim, bo.Rcorner + 8, bo.Offsety)
	fi := ocean.NewOcean(opts.GridDim, bo.Rcorner + 8, bo.Offsety)
	in := ocean.NewOcean(opts.GridDim, 4, 2)
	g := &Game{ 
		Options: opts, 
		PlBoard: bo, 
		OpBoard: ob, 
		Fired: fi, 
		Incoming: in,
	}
	return g
}

func InitShips(g *Game, shipList []int) (*Game) {
	plShips := []ships.Shiper{}
	opShips := []ships.Shiper{}

	for _, T := range shipList {
		var typedShip ships.Shiper
		newShip := ships.Ship{
			X: -1, 
			Y: -1, 
			Visible: true,
			Vertical: true,
			Style: g.Options.PlStyle,
			SunkStyle: g.Options.AltStyle,
		}
		
		switch T {
		case 5:
			newShip.Name = "Carrier"
			newShip.Type, newShip.Health, newShip.Length = 5,5,5
			newShip.Coords = make([]ships.Coord, 5)
			typedShip = ships.Carrier{ Ship: newShip }
		case 4:
			newShip.Name = "Battleship"
			newShip.Type, newShip.Health, newShip.Length = 4,4,4
			newShip.Coords = make([]ships.Coord, 4)
			typedShip = ships.Battleship{ Ship: newShip }
		case 3:
			newShip.Name = "Cruiser"
			newShip.Type, newShip.Health, newShip.Length = 3,3,3
			newShip.Coords = make([]ships.Coord, 3)
			typedShip = ships.Cruiser{ Ship: newShip }
		case 2:
			newShip.Name = "Submarine"
			newShip.Type, newShip.Health, newShip.Length = 2,3,3
			newShip.Coords = make([]ships.Coord, 3)
			typedShip = ships.Submarine{ Ship: newShip }
		case 1:	
			newShip.Name = "Destroyer"
			newShip.Type, newShip.Health, newShip.Length = 1,2,2
			newShip.Coords = make([]ships.Coord, 2)
			typedShip = ships.Destroyer{ Ship: newShip }
		}
		plShips = append(plShips, typedShip)
	}

	for _, ship := range plShips {
		typedShip := ship
		newShip := ship.GetShip()
		newShip.Visible = false
		newShip.Style = g.Options.OpStyle
		typedShip = typedShip.SetShip(newShip)
		opShips = append(opShips, typedShip)
	}

	g.PlShips = plShips
	g.OpShips = opShips

	return g
}

func InitScreen() (string, tcell.Screen) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "CMD.EXE"
		} else {
			shell = "/bin/sh"
		}
	}

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	return shell, s
}