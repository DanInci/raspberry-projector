package framebuffer

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"syscall"
	"unsafe"
)

const FBIOGET_VSCREENINFO = 0x4600
const FBIOPUT_VSCREENINFO = 0x4601
const FBIOGET_FSCREENINFO = 0x4602
const FB_TYPE_PACKED_PIXELS = 0
const FB_VISUAL_TRUECOLOR = 2

type FixScreenInfo struct {
	Id                               [16]byte
	Smem_start                       uintptr
	Smem_len, Type, Type_aux, Visual uint32
	Xpanstep, Ypanstep, Ywrapstep    uint16
	Line_length                      uint32
	Mmio_start                       uintptr
	Mmio_len, Accel                  uint32
	Capabilities                     uint16
	Reserved                         [2]uint16
}

type BitField struct {
	Offset, Length, Msb_right uint32
}

type VarScreenInfo struct {
	Xres, Yres,
	Xres_virtual, Yres_virtual,
	Xoffset, Yoffset,
	Bits_per_pixel, Grayscale uint32
	Red, Green, Blue, Transp BitField
	Nonstd, Activate,
	Height, Width,
	Accel_flags, Pixclock,
	Left_margin, Right_margin, Upper_margin, Lower_margin,
	Hsync_len, Vsync_len, Sync,
	Vmode, Rotate, Colorspace uint32
	Reserved [4]uint32
}

type FrameBuffer interface {
	draw.Image
	Close()
}

type BGR565 struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	File   *os.File
}

func (p *BGR565) Bounds() image.Rectangle { return p.Rect }
func (p *BGR565) ColorModel() color.Model { return color.NRGBAModel }
func (p *BGR565) PixOffset(x, y int) int  { return y*p.Stride + x*2 }
func (p *BGR565) Close()                  { syscall.Munmap(p.Pix); p.File.Close() }

func (p *BGR565) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.NRGBAModel.Convert(c).(color.NRGBA)
	p.Pix[i+0] = (c1.B >> 3) | ((c1.G >> 2) << 5)
	p.Pix[i+1] = (c1.G >> 5) | ((c1.R >> 3) << 3)
}

func (p *BGR565) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.NRGBA{}
	}
	i := p.PixOffset(x, y)
	return color.NRGBA{(p.Pix[i+1] >> 3) << 3, (p.Pix[i+1] << 5) | ((p.Pix[i+0] >> 5) << 2), p.Pix[i+0] << 3, 255}
}

type BGR struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	File   *os.File
}

func (p *BGR) Bounds() image.Rectangle { return p.Rect }
func (p *BGR) ColorModel() color.Model { return color.NRGBAModel }
func (p *BGR) PixOffset(x, y int) int  { return y*p.Stride + x*3 }
func (p *BGR) Close()                  { syscall.Munmap(p.Pix); p.File.Close() }

func (p *BGR) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.NRGBAModel.Convert(c).(color.NRGBA)
	p.Pix[i+0] = c1.B
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.R
}

func (p *BGR) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.NRGBA{}
	}
	i := p.PixOffset(x, y)
	return color.NRGBA{p.Pix[i+2], p.Pix[i+1], p.Pix[i+0], 255}
}

type BGR32 struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	File   *os.File
}

func (p *BGR32) Bounds() image.Rectangle { return p.Rect }
func (p *BGR32) ColorModel() color.Model { return color.NRGBAModel }
func (p *BGR32) PixOffset(x, y int) int  { return y*p.Stride + x*4 }
func (p *BGR32) Close()                  { syscall.Munmap(p.Pix); p.File.Close() }

func (p *BGR32) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.NRGBAModel.Convert(c).(color.NRGBA)
	p.Pix[i+0] = c1.B
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.R
}

func (p *BGR32) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.NRGBA{}
	}
	i := p.PixOffset(x, y)
	return color.NRGBA{p.Pix[i+2], p.Pix[i+1], p.Pix[i+0], 255}
}

type NBGRA struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	File   *os.File
}

func (p *NBGRA) Bounds() image.Rectangle { return p.Rect }
func (p *NBGRA) ColorModel() color.Model { return color.NRGBAModel }
func (p *NBGRA) PixOffset(x, y int) int  { return y*p.Stride + x*4 }
func (p *NBGRA) Close()                  { syscall.Munmap(p.Pix); p.File.Close() }

func (p *NBGRA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.NRGBAModel.Convert(c).(color.NRGBA)
	p.Pix[i+0] = c1.B
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.R
	p.Pix[i+3] = c1.A
}

func (p *NBGRA) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.NRGBA{}
	}
	i := p.PixOffset(x, y)
	return color.NRGBA{p.Pix[i+2], p.Pix[i+1], p.Pix[i+0], p.Pix[i+3]}
}

type UnsupportedError string

func (e UnsupportedError) Error() string { return "framebuffer: " + string(e) }

func Open(name string) (FrameBuffer, error) {
	file, err := os.OpenFile(name, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	var fixInfo FixScreenInfo
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), FBIOGET_FSCREENINFO, uintptr(unsafe.Pointer(&fixInfo))); errno != 0 {
		return nil, &os.SyscallError{"SYS_IOCTL", errno}
	}
	if fixInfo.Type != FB_TYPE_PACKED_PIXELS {
		return nil, UnsupportedError("fixInfo.Type != FB_TYPE_PACKED_PIXELS")
	}
	if fixInfo.Visual != FB_VISUAL_TRUECOLOR {
		return nil, UnsupportedError("fixInfo.Visual != FB_VISUAL_TRUECOLOR")
	}
	var varInfo VarScreenInfo
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), FBIOGET_VSCREENINFO, uintptr(unsafe.Pointer(&varInfo))); errno != 0 {
		return nil, &os.SyscallError{"SYS_IOCTL", errno}
	}

	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fixInfo.Smem_len), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	switch varInfo.Bits_per_pixel {
	case 32:
		if varInfo.Blue.Length != 8 {
			return nil, UnsupportedError("varInfo.Blue.Length != 8")
		}
		if varInfo.Blue.Offset != 0 {
			return nil, UnsupportedError("varInfo.Blue.Offset != 0")
		}
		if varInfo.Green.Length != 8 {
			return nil, UnsupportedError("varInfo.Green.Length != 8")
		}
		if varInfo.Green.Offset != 8 {
			return nil, UnsupportedError("varInfo.Green.Offset != 8")
		}
		if varInfo.Red.Length != 8 {
			return nil, UnsupportedError("varInfo.Red.Length != 8")
		}
		if varInfo.Red.Offset != 16 {
			return nil, UnsupportedError("varInfo.Red.Offset != 16")
		}
		if varInfo.Transp.Length == 0 {
			return &BGR32{mmap, int(fixInfo.Line_length), image.Rect(0, 0, int(varInfo.Xres), int(varInfo.Yres)).Add(image.Point{int(varInfo.Xoffset), int(varInfo.Yoffset)}), file}, nil
		} else if varInfo.Transp.Length == 8 && varInfo.Transp.Offset == 24 {
			return &NBGRA{mmap, int(fixInfo.Line_length), image.Rect(0, 0, int(varInfo.Xres), int(varInfo.Yres)).Add(image.Point{int(varInfo.Xoffset), int(varInfo.Yoffset)}), file}, nil
		}
	case 24:
		if varInfo.Blue.Length != 8 {
			return nil, UnsupportedError("varInfo.Blue.Length != 8")
		}
		if varInfo.Blue.Offset != 0 {
			return nil, UnsupportedError("varInfo.Blue.Offset != 0")
		}
		if varInfo.Green.Length != 8 {
			return nil, UnsupportedError("varInfo.Green.Length != 8")
		}
		if varInfo.Green.Offset != 8 {
			return nil, UnsupportedError("varInfo.Green.Offset != 8")
		}
		if varInfo.Red.Length != 8 {
			return nil, UnsupportedError("varInfo.Red.Length != 8")
		}
		if varInfo.Red.Offset != 16 {
			return nil, UnsupportedError("varInfo.Red.Offset != 16")
		}
		if varInfo.Transp.Length != 0 {
			return nil, UnsupportedError("varInfo.Transp.Length != 0")
		}
		return &BGR{mmap, int(fixInfo.Line_length), image.Rect(0, 0, int(varInfo.Xres), int(varInfo.Yres)).Add(image.Point{int(varInfo.Xoffset), int(varInfo.Yoffset)}), file}, nil
	case 16:
		if varInfo.Blue.Length != 5 {
			return nil, UnsupportedError("varInfo.Blue.Length != 5")
		}
		if varInfo.Blue.Offset != 0 {
			return nil, UnsupportedError("varInfo.Blue.Offset != 0")
		}
		if varInfo.Green.Length != 6 {
			return nil, UnsupportedError("varInfo.Green.Length != 6")
		}
		if varInfo.Green.Offset != 5 {
			return nil, UnsupportedError("varInfo.Green.Offset != 5")
		}
		if varInfo.Red.Length != 5 {
			return nil, UnsupportedError("varInfo.Red.Length != 5")
		}
		if varInfo.Red.Offset != 11 {
			return nil, UnsupportedError("varInfo.Red.Offset != 11")
		}
		if varInfo.Transp.Length != 0 {
			return nil, UnsupportedError("varInfo.Transp.Length != 0")
		}
		return &BGR565{mmap, int(fixInfo.Line_length), image.Rect(0, 0, int(varInfo.Xres), int(varInfo.Yres)).Add(image.Point{int(varInfo.Xoffset), int(varInfo.Yoffset)}), file}, nil
	}
	return nil, UnsupportedError("unsupported pixel format")
}
