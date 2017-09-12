package display

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Step struct {
	X float64
	Y float64
	Z float64
}

type DisplayTrajectory struct {
	pixelSize int
	size      float64
	image     *image.RGBA
	gc        *draw2dimg.GraphicContext
}

func NewDisplayTrajectory(pixelSize int, size float64) DisplayTrajectory {
	dt := DisplayTrajectory{
		pixelSize: pixelSize,
		size:      size,
	}
	dt.image = image.NewRGBA(image.Rect(0, 0, pixelSize, pixelSize))
	dt.gc = draw2dimg.NewGraphicContext(dt.image)

	// Set some properties
	dt.gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	dt.gc.SetLineWidth(1)

	return dt
}

func (dt DisplayTrajectory) DrawStepXY(position Step, newStep Step) {
	// Draw a closed shape
	dt.gc.MoveTo(float64(dt.pixelSize)/2.0+position.X/dt.size*float64(dt.pixelSize), float64(dt.pixelSize)/2.0+position.Y/dt.size*float64(dt.pixelSize))
	dt.gc.LineTo(float64(dt.pixelSize)/2.0+(position.X+newStep.X)/dt.size*float64(dt.pixelSize), float64(dt.pixelSize)/2.0+(position.Y+newStep.Y)/dt.size*float64(dt.pixelSize))
	dt.gc.FillStroke()
}

func (dt DisplayTrajectory) DrawStepXZ(position Step, newStep Step) {
	// Draw a closed shape
	dt.gc.MoveTo(float64(dt.pixelSize)/2.0+position.X/dt.size*float64(dt.pixelSize), position.Z/dt.size*float64(dt.pixelSize))
	dt.gc.LineTo(float64(dt.pixelSize)/2.0+(position.X+newStep.X)/dt.size*float64(dt.pixelSize), (position.Z+newStep.Z)/dt.size*float64(dt.pixelSize))
	dt.gc.FillStroke()
}

func (dt DisplayTrajectory) DrawStepYZ(position Step, newStep Step) {
	// Draw a closed shape
	dt.gc.MoveTo(float64(dt.pixelSize)/2.0+position.Y/dt.size*float64(dt.pixelSize), position.Z/dt.size*float64(dt.pixelSize))
	dt.gc.LineTo(float64(dt.pixelSize)/2.0+(position.Y+newStep.Y)/dt.size*float64(dt.pixelSize), (position.Z+newStep.Z)/dt.size*float64(dt.pixelSize))
	dt.gc.FillStroke()
}

func (dt DisplayTrajectory) SaveFile(filename string) {
	draw2dimg.SaveToPngFile(filename, dt.image)
}
