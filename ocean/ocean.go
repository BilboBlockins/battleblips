package ocean

import (
	"fmt"
)

type Ocean struct {
	Dim int
	Offsetx int
	Offsety int
	Rcorner int
	Bcorner int
	Cellx int
	Celly int
	Grid [][]int
}

func NewOcean(dim int, offsetx int, offsety int) Ocean {
	o := Ocean{}
	o.Dim = dim
	o.Cellx = 4
	o.Celly = 2
	o.Rcorner = o.Dim * o.Cellx + offsetx
	o.Bcorner = o.Dim * o.Celly + offsety
	o.Offsetx = offsetx
	o.Offsety = offsety
	o.Grid = make_grid(dim)
	return o
}


func make_grid(dim int) [][]int {
	if dim < 5 { dim = 10 }
	if dim > 26 { dim = 26 }
	
	o := make([][]int, dim)
	for i := range o {
		for range o {
			o[i] = append(o[i],0)
		}
	}
	return o
}

func (o Ocean) Print() {
	for _, row := range o.Grid{
		fmt.Println(row)
	}
}
