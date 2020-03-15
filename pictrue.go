package lcd

import (
	"github.com/disintegration/imaging"
	"github.com/jjjabc/lcd/wbimage"
	"image"
)

func Picture(i image.Image) {
	var resizedImg image.Image
	if i.Bounds().Dy() > 64 || i.Bounds().Dx() > 128 {
		resizedImg = imaging.Resize(i, 128, 64, imaging.Lanczos)
	} else {
		resizedImg = imaging.CropAnchor(i, 128, 64, imaging.Center)
	}
	var bitmap [128 * 8]byte
	index := 0
	for y := 0; y < 64; y++ {
		for x := 0; x < 128/8; x++ {
			var b byte
			for bitIndex := 0; bitIndex < 8; bitIndex++ {
				_, g, _, _ := wbimage.Default.ColorModel().Convert(resizedImg.At(8*x+bitIndex, y)).RGBA()
				if g > 0xf000 {
					b = b | (0x80 >> bitIndex)
				}
			}
			bitmap[index] = b
			index++
		}
	}
	pic(&bitmap)
}
