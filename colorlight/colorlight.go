package colorlight

import (
	"io"
	"log"

	"github.com/mdlayher/ethernet"
)

type EthInterface interface {
	io.Closer
	Send(ethernet.Frame) error
	ReceiveBroadcast() (ethernet.Frame, error)
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
	iface      EthInterface
	info       CardInfo
}

func New() Colorlight {
	return Colorlight{
		Brightness: 255,
		ColorTemp:  BalancedTemp,
	}
}

func (cl *Colorlight) Open(iface EthInterface) error {
	cl.iface = iface
	ci, err := cl.detectCard()
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
	defer func() { cl.iface = nil }()
	return cl.iface.Close()
}

func (cl *Colorlight) send(m Message) error {
	return cl.iface.Send(m.Frame())
}

func (cl *Colorlight) receiveMatching(matcher func(ethernet.Frame) bool) (f ethernet.Frame, err error) {
	for {
		f, err = cl.iface.ReceiveBroadcast()
		if err != nil {
			return
		}
		if matcher(f) {
			return
		} else {
			log.Printf("colorlight: ignoring broadcast (no match): %v", f)
		}
	}
}

func (cl *Colorlight) detectCard() (ci CardInfo, err error) {
	err = cl.send(DetectCardMessage{})
	if err != nil {
		return
	}
	f, err := cl.receiveMatching(func(f ethernet.Frame) bool {
		_, err := ParseCardInfo(f)
		return err == nil
	})
	if err != nil {
		return
	}
	ci, _ = ParseCardInfo(f)
	err = cl.send(DetectCardResponseAckMessage{ControllerID: ci.ControllerID})
	return
}
