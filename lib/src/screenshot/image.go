package screenshot

import "C"
import (
	"image"
	"image/color"
	"unsafe"
)

var ImageCache *image.YCbCr

func RGBToYCbCr(r, g, b uint8) (y, cb, cr uint8) {
	ret := C.RGBToYCbCr((C.uchar)(r), (C.uchar)(g), (C.uchar)(b))
	return uint8(ret.Y), uint8(ret.Cb), uint8(ret.Cr)
}

func CRGBToYCbCr444Linux(data, y, cb, cr []byte) {
	C.ImageRGBToYCbCr4442((*C.uchar)(unsafe.Pointer(&data[0])),
		C.int32_t(len(data)),
		(*C.uchar)(unsafe.Pointer(&y[0])),
		(*C.uchar)(unsafe.Pointer(&cb[0])),
		(*C.uchar)(unsafe.Pointer(&cr[0])))
}

func CRGBToYCbCr444Windows(data, y, cb, cr []byte) {
	C.ImageRGBToYCbCr4442((*C.uchar)(unsafe.Pointer(&data[0])),
		C.int32_t(len(data)),
		(*C.uchar)(unsafe.Pointer(&y[0])),
		(*C.uchar)(unsafe.Pointer(&cb[0])),
		(*C.uchar)(unsafe.Pointer(&cr[0])))
}

func ImageToRGBALinux(data []byte) {
	C.ImageToRGBALinux((*C.uchar)(unsafe.Pointer(&data[0])),
		C.int32_t(len(data)))
}

func ImageToRGBAWindows(data, bytes []byte) {
	C.ImageToRGBAWindows((*C.uchar)(unsafe.Pointer(&data[0])),
		C.int32_t(len(data)),
		(*C.uchar)(unsafe.Pointer(&bytes[0])))
}

func RGBAToYCbCr444(img *image.RGBA) *image.YCbCr {
	new_img := image.NewYCbCr(img.Rect, image.YCbCrSubsampleRatio444)
	new_img.Y = make([]uint8, len(img.Pix)/4)
	new_img.Cb = make([]uint8, len(img.Pix)/4)
	new_img.Cr = make([]uint8, len(img.Pix)/4)

	n := 0
	for i := 0; i < len(img.Pix); i += 4 {
		y, cb, cr := color.RGBToYCbCr(img.Pix[i], img.Pix[i+1], img.Pix[i+2])
		new_img.Y[n] = y
		new_img.Cb[n] = cb
		new_img.Cr[n] = cr
		n += 1
	}

	return new_img
}

func RGBAToYCbCr420(img *image.RGBA) *image.YCbCr {
	new_img := image.NewYCbCr(img.Rect, image.YCbCrSubsampleRatio420)
	new_img.Y = make([]uint8, len(img.Pix)/4)
	new_img.Cb = make([]uint8, len(img.Pix)/16)
	new_img.Cr = make([]uint8, len(img.Pix)/16)

	cn := 0
	for y := 0; y < img.Rect.Dy()/2; y += 1 {
		for x := 0; x < img.Rect.Dx()/2; x += 1 {
			x0, y0 := x*2, y*2
			x1, y1 := x*2+1, y*2
			x2, y2 := x*2, y*2+1
			x3, y3 := x*2+1, y*2+1

			co0 := img.RGBAAt(x0, y0)
			cy0, cb0, cr0 := color.RGBToYCbCr(co0.R, co0.G, co0.B)
			co1 := img.RGBAAt(x1, y1)
			cy1, cb1, cr1 := color.RGBToYCbCr(co1.R, co1.G, co1.B)
			co2 := img.RGBAAt(x2, y2)
			cy2, cb2, cr2 := color.RGBToYCbCr(co2.R, co2.G, co2.B)
			co3 := img.RGBAAt(x3, y3)
			cy3, cb3, cr3 := color.RGBToYCbCr(co3.R, co3.G, co3.B)

			new_img.Y[x0+y0*img.Rect.Dx()] = cy0
			new_img.Y[x1+y1*img.Rect.Dx()] = cy1
			new_img.Y[x2+y2*img.Rect.Dx()] = cy2
			new_img.Y[x3+y3*img.Rect.Dx()] = cy3
			new_img.Cb[cn] = cb0/4 + cb1/4 + cb2/4 + cb3/4
			new_img.Cr[cn] = cr0/4 + cr1/4 + cr2/4 + cr3/4
			cn += 1
		}
	}

	return new_img
}
