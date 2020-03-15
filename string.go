package lcd

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/jjjabc/lcd/wbimage"
	_ "golang.org/x/image/font"
	"image"
	"log"
)

func Strings(strs []string, size float64, spacing float64, f *truetype.Font) {
		//for i := -64; i < 65; i++ {
			//s := time.Now()
			disPic := stringToPic(strs, size, spacing, f, 0, 0)
			Picture(disPic)
			//time.Sleep(200*time.Millisecond - time.Now().Sub(s))
		//}
}
func stringToPic(strs []string, size float64, spacing float64, f *truetype.Font, x, y int) (disPic *wbimage.WB) {
	disPic = wbimage.NewWB(image.Rect(0, 0, 128, 64))
	fg := image.White
	c := freetype.NewContext()
	c.SetDPI(50.8)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(disPic.Bounds())
	c.SetDst(disPic)
	c.SetSrc(fg)
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6)-1)

	for _, str := range strs {
		_, err := c.DrawString(str, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}
	return
}
