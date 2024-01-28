package diagram

import (
	"image/color"

	"github.com/jwendel/downsampling/core"
)

type Config struct {
	Title string
	Name  string
	Data  []core.Point[float64, float64]
	Color color.RGBA
}
