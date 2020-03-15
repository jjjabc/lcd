package wbimage

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"
)

func TestWBColor_RGBA(t *testing.T) {
	src, err := imaging.Open("./exp.png")
	if err != nil {
		t.Error(err)
	}
	wb := NewWB(src.Bounds())
	for y := 0; y < wb.Bounds().Dy(); y++ {
		for x := 0; x < wb.Bounds().Dx(); x++ {
			wb.Set(x, y, src.At(x, y))
		}
	}
	dstFile, err := os.Create("./dst.jpg")
	if err != nil {
		t.Error(err)
	}
	err = jpeg.Encode(dstFile, wb, &jpeg.Options{Quality: 90})
	if err != nil {
		t.Error(err)
	}
}

func TestWB_WBAt(t *testing.T) {
	i, err := imaging.Open("./exp.png")
	if err != nil {
		t.Error(err)
	}
	srcWidth:=i.Bounds().Max.Sub(i.Bounds().Min).X
	srcHeight:=i.Bounds().Max.Sub(i.Bounds().Min).Y
	var resizedImg image.Image
	if srcHeight>64||srcWidth>128{
		resizedImg=imaging.Resize(i,128,64,imaging.Lanczos)
		log.Printf("resize")
	}else{
		resizedImg=i
		log.Printf("Crop")
	}
/*	for y := 0; y < 64; y++ {
		for x := 0; x < 128; x++ {
			_,g,_,_:=resizedImg.At(x,y).RGBA()
			if g>0x8000{
				fmt.Printf("1")
			}else{
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
	return*/
	var bitmap [128 *8]byte
	index:=0
	for y := 0; y < 64; y++ {
		for x := 0; x < 128/8; x++ {
			var b byte
			for bitIndex:=0;bitIndex<8;bitIndex++ {
				_, g, _, _ := Default.ColorModel().Convert(resizedImg.At(x*8+bitIndex, y)).RGBA()
				if g > 0xf000 {
					b = b | (0x80 >> bitIndex)
				}
			}
			bitmap[index]=b
			fmt.Printf("%8b",b)
			index++
		}
		fmt.Printf("\n")
	}
}