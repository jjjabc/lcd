package wbimage

import (
	"image"
	"image/color"
)

type WB struct {
	// Pix holds the image's pixels, as gray values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []bool
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	Rect   image.Rectangle
}

func (p *WB) ColorModel() color.Model {
	return WBModel
}

func (p *WB) Bounds() image.Rectangle {
	return p.Rect
}

func (p *WB) At(x, y int) color.Color {
	return p.WBAt(x, y)
}

func (p *WB) WBAt(x, y int) WBColor {
	if !(image.Point{x, y}.In(p.Rect)) {
		return false
	}
	i := p.PixOffset(x, y)
	return WBColor(p.Pix[i])
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *WB) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}
func (p *WB) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = bool(p.ColorModel().Convert(c).(WBColor))
}

func NewWB(r image.Rectangle) *WB {
	w, h := r.Dx(), r.Dy()
	pix := make([]bool, 1*w*h)
	return &WB{pix, 1 * w, r}
}

var WBModel color.Model = color.ModelFunc(wbModel)

type WBColor bool

func (c WBColor) RGBA() (r, g, b, a uint32) {
	if c {
		return 0xffff, 0xffff, 0xffff, 0xffff
	}
	return 0x0000, 0x0000, 0x0000, 0xffff
}

func wbModel(c color.Color) color.Color {
	if _, ok := c.(WBColor); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	//
	// The 24 is 16 + 8. The 16 is the same as used in RGBToYCbCr. The 8 is
	// because the return value is 8 bit color, not 16 bit color.
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
	return WBColor(y > 127)
}

var Default  = NewWB(image.Rectangle{})