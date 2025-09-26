package colorlight

import (
	"net"

	"github.com/mdlayher/ethernet"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

// Brightness of the display, from 0x00 (off) to 0xFF (full)
type Brightness uint8

// ColorTemp allows balancing the colors
type ColorTemp struct{ Red, Green, Blue Brightness }

var BalancedTemp ColorTemp = ColorTemp{255, 255, 255}

type Pixel struct{ Red, Green, Blue uint8 }
type Canvas [][]Pixel

type Colorlight struct {
	Brightness Brightness
	ColorTemp  ColorTemp
	iface      *net.Interface
	info       CardInfo
}

func New() Colorlight {
	return Colorlight{
		Brightness: 255,
		ColorTemp:  BalancedTemp,
	}
}

func (cl *Colorlight) Open(ifaceName string) error {
	ifi, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return err
	}
	cl.iface = ifi

	err = cl.send(DetectCardMessage{})
	if err != nil {
		return err
	}
	var f ethernet.Frame // TODO: receive
	ci, err := ParseCardInfo(f)
	if err != nil {
		cl.Close()
		return err
	}
	err = cl.send(DetectCardResponseAckMessage{ControllerID: ci.ControllerID})
	if err != nil {
		cl.Close()
		return err
	}
	cl.info = ci
	return nil
}

func (cl *Colorlight) Dims() (width, height uint16) {
	return cl.info.NumColumns, cl.info.NumRows
}

func (cl *Colorlight) send(m Message) error {
	// TODO
	return nil
}

func (cl *Colorlight) Draw(c Canvas) (err error) {
	err = cl.send(BrightnessMessage{cl.Brightness})
	if err != nil {
		return err
	}
	for x, row := range c {
		err = cl.send(SetRowMessage{RowIndex: uint16(x), Data: row, PixelOffset: 0})
		if err != nil {
			return err
		}
	}
	err = cl.send(DisplayMessage{cl.Brightness, cl.ColorTemp})
	return err
}

func (cl *Colorlight) Close() error {
	cl.iface = nil
	return nil
}
