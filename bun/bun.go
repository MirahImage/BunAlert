package bun

import (
	"image"

	"github.com/kellydun/golang-geo"
)

// Bun struct containing bun information
type Bun struct {
	Location    *geo.Point
	Size        int
	Picture     *image.RGBA
	Description string
}

// LogBun logs a bun, creating a new Bun
func (bun *Bun) LogBun(size int, description string) {
	bun.Location = geo.NewPoint(0, 0)
	bun.Size = size
	bun.Description = description
}
