package lcd

import (
	_ "github.com/golang/freetype"
	"github.com/stianeikeland/go-rpio/v4"
	"golang.org/x/text/encoding/simplifiedchinese"
	"time"
)

const (
	ledW = 128
	ledH = 64
)

var (
	DPs = [8]rpio.Pin{
		rpio.Pin(0),
		rpio.Pin(5),
		rpio.Pin(6),
		rpio.Pin(13),
		rpio.Pin(19),
		rpio.Pin(26),
		rpio.Pin(12),
		rpio.Pin(16),
	}
	RS    = rpio.Pin(23)
	RW    = rpio.Pin(24)
	EN    = rpio.Pin(18)
	RESET = rpio.Pin(21)
)

func writeByte(b byte) {
	for i := 0; i < 8; i++ {
		switch (b >> i) & 0x01 {
		case 0:
			DPs[i].Low()
		case 1:
			DPs[i].High()
		}
	}
}
func chk_busy() {
	RS.Low()
	RW.High()
	EN.High()
	writeByte(0xFF)
	DPs[7].Input()
	for DPs[7].Read() == rpio.High {
	}
	DPs[7].Output()
	EN.Low()
}
func Cmd(cmd byte) {
	chk_busy()
	RS.Low()
	RW.Low()
	EN.High()
	writeByte(cmd)
	EN.Low()
}
func Data(b byte) {
	chk_busy()
	RS.High()
	RW.Low()
	EN.High()
	writeByte(b)
	EN.Low()
}
func sleep(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}
func Init() {
	for _, DP := range DPs {
		DP.Output()
	}

	RS.Output()
	RW.Output()
	EN.Output()
	reset()
	Cmd(0x38) //选择8bit数据流
	sleep(5)
	Cmd(0x01) //清除显示，并且设定地址指针为00H
	sleep(5)
}
func reset() {
	RESET.Output()
	RESET.Low()
	sleep(5)
	RESET.High()
}
func PrintString(s string) (err error) {
	ret, err := simplifiedchinese.HZGB2312.NewEncoder().Bytes([]byte(s))
	if err != nil {
		return
	}
	for i, r := range ret {
		switch i {
		case 0 * 16:
			Cmd(0x80)
		case 1 * 16:
			Cmd(0x90)
		case 2 * 16:
			Cmd(0x88)
		case 3 * 16:
			Cmd(0x98)
		}
		Data(r)
	}
	return nil
}
func pic(bitmap *[128 * 8]byte) {
	i := 0
	for y := 0x80; y < 0x80+32; y++ {
		Cmd(byte(y))
		Cmd(0x80)
		for x := 0; x < 16; x++ {
			Data(bitmap[i])
			i++
		}
	}
	for y := 0x80; y < 0x80+32; y++ {
		Cmd(byte(y))
		Cmd(0x88)
		for x := 0; x < 16; x++ {
			Data(bitmap[i])
			i++
		}
	}
}
func Clear() {
	Cmd(0x01)
	emptyMap := new([128 * 8]byte)
	pic(emptyMap)
}
func ImageMod() {
	Cmd(0x3e)
}
func TextMod()  {
	Cmd(0x0c)
}
