package bun

import (
	"github.com/kellydun/golang-geo"
	"image"
)

type Bun struct {
	Location *geo.Point
	Size int
	Picture *image.RGBA
	Description string
}


func (bun *Bun) LogBun(size int, description string) {
	bun.Location = geo.NewPoint(0,0)
	bun.Size = size
	bun.Description = description
}

