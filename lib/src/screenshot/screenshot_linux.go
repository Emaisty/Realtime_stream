package screenshot

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"image"
	"time"
)

var Conn *xgb.Conn

type POS struct {
	X, Y int
}

type SIZE struct {
	W, H int
}

type RESIZE struct {
	W, H int
}

type CAPTURE struct {
	W, H int
	B    *[]byte
}

func InitConn() {
	var err error
	Conn, err = xgb.NewConn()
	if err != nil {
		panic(err)
	}
}

func CloseConn() {
	Conn.Close()
}

func ScreenRect() (image.Rectangle, error) {
	c := Conn

	screen := xproto.Setup(c).DefaultScreen(c)
	x := screen.WidthInPixels
	y := screen.HeightInPixels

	return image.Rect(0, 0, int(x), int(y)), nil
}

func CaptureWindowMust(pos *POS, size *SIZE, resize *RESIZE, toSBS bool, cursor bool) *image.RGBA {
	img, err := CaptureWindow(pos, size, resize, toSBS, cursor)
	for err != nil {
		img, err = CaptureWindow(pos, size, resize, toSBS, cursor)
		time.Sleep(10 * time.Millisecond)
	}
	return img
}

func CaptureScreenYCbCrMust(toSBS bool, cursor bool) *image.YCbCr {
	img, err := CaptureScreenYCbCr444()
	for err != nil {
		img, err = CaptureScreenYCbCr444()
		time.Sleep(10 * time.Millisecond)
	}
	return img
}

func CaptureScreen() (*image.RGBA, error) {
	r, e := ScreenRect()
	if e != nil {
		return nil, e
	}
	return CaptureRect(r)
}

func CaptureScreenYCbCr444() (*image.YCbCr, error) {
	r, e := ScreenRect() // 20us
	if e != nil {
		return nil, e
	}
	return CaptureRectYCbCr444(r)
}

func CaptureRect(rect image.Rectangle) (*image.RGBA, error) {
	c := Conn

	screen := xproto.Setup(c).DefaultScreen(c)
	x, y := rect.Dx(), rect.Dy()
	xImg, err := xproto.GetImage(c, xproto.ImageFormatZPixmap, xproto.Drawable(screen.Root), int16(rect.Min.X), int16(rect.Min.Y), uint16(x), uint16(y), 0xffffffff).Reply()
	if err != nil {
		return nil, err
	}

	data := xImg.Data
	for i := 0; i < len(data); i += 4 {
		data[i], data[i+2], data[i+3] = data[i+2], data[i], 255
	}

	img := &image.RGBA{data, 4 * x, image.Rect(0, 0, x, y)}
	return img, nil
}

func CaptureRectYCbCr444(rect image.Rectangle) (*image.YCbCr, error) {
	c := Conn

	t := time.Now()
	screen := xproto.Setup(c).DefaultScreen(c)
	x, y := rect.Dx(), rect.Dy()
	xImg, err := xproto.GetImage(c, xproto.ImageFormatZPixmap, xproto.Drawable(screen.Root), int16(rect.Min.X), int16(rect.Min.Y), uint16(x), uint16(y), 0xffffffff).Reply()
	if err != nil {
		return nil, err
	}

	data := xImg.Data
	dataLen := len(data)

	tt := time.Now()
	if ImageCache == nil {
		ImageCache = image.NewYCbCr(image.Rect(0, 0, x, y), image.YCbCrSubsampleRatio444)
	}
	ttt := time.Now()

	CRGBToYCbCr444Linux(data, ImageCache.Y, ImageCache.Cb, ImageCache.Cr)
	tttt := time.Now()
	log.Println(fmt.Sprintf("Shot: %v, Create: %v, Convert: %v", tt.Sub(t), ttt.Sub(tt), tttt.Sub(ttt)))
	// Shot: 14.734765ms, Create: 108ns, Convert: 9.515677ms

	return ImageCache, nil
}

func CaptureWindow(pos *POS, size *SIZE, resize *RESIZE, toSBS bool, cursor bool) (*image.RGBA, error) {

	c := Conn

	screen := xproto.Setup(c).DefaultScreen(c)

	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(c, true, uint16(len(aname)), aname).Reply()
	if err != nil {
		return nil, fmt.Errorf("error occurred, when xproto.InternAtom 0 err:%v.\n", err)
	}

	reply, err := xproto.GetProperty(c, false, screen.Root, activeAtom.Atom, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return nil, fmt.Errorf("error occurred, when xproto.GetProperty 0 err:%v.\n", err)
	}
	windowId := xproto.Window(xgb.Get32(reply.Value))

	ginfo, err := xproto.GetGeometry(c, xproto.Drawable(windowId)).Reply()
	if err != nil {
		return nil, err
	}

	width := int(ginfo.Width) - pos.X
	height := int(ginfo.Height) - pos.Y
	if size.W != 0 && size.H != 0 {
		width = size.W
		height = size.H
	}

	xImg, err := xproto.GetImage(c, xproto.ImageFormatZPixmap, xproto.Drawable(windowId), int16(pos.X), int16(pos.Y), uint16(width), uint16(height), 0xffffffff).Reply()
	if err != nil {
		return nil, err
	}

	data := xImg.Data
	ImageToRGBALinux(data)
	var img *image.RGBA
	if toSBS {
		img = &image.RGBA{append(data, data...), 4 * width, image.Rect(pos.X, pos.Y, width*2-2, height-1)}
	} else {
		img = &image.RGBA{data, 4 * width, image.Rect(pos.X, pos.Y, width-1, height-1)}
	}
	return img, nil
}
