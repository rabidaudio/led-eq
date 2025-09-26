package colorlight

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mdlayher/ethernet"
)

var srcAddr = must(net.ParseMAC("22:22:33:44:55:66"))
var destAddr = must(net.ParseMAC("11:22:33:44:55:66"))
var bcastAddr = must(net.ParseMAC("FF:FF:FF:FF:FF:FF"))

// https://hkubota.wordpress.com/2022/01/31/winter-project-colorlight-5a-75b-protocol/

type MsgType ethernet.EtherType

const (
	DisplayFrame       MsgType = 0x0107
	SetBrightness      MsgType = 0x0a00 // + brightness
	PixelRow           MsgType = 0x5500 // or 5501? + MSB of row number
	DetectCard         MsgType = 0x0700
	DetectCardResponse MsgType = 0x0805 // 11:22:33:44:55:66 -> ff:ff:ff:ff:ff:ff
)

type Message interface {
	Frame() ethernet.Frame
}

type BrightnessMessage struct {
	Brightness Brightness
}

func (bm BrightnessMessage) Frame() ethernet.Frame {
	pl := make([]byte, 63)
	pl[0] = byte(bm.Brightness)
	pl[1] = byte(bm.Brightness)
	pl[2] = 0xff
	return ethernet.Frame{
		Source:      srcAddr,
		Destination: destAddr,
		EtherType:   ethernet.EtherType(SetBrightness) + ethernet.EtherType(bm.Brightness),
		Payload:     pl,
	}
}

type SetRowMessage struct {
	// Columns int
	RowIndex    uint16
	Data        []Pixel // BGR, pixels 0-n
	PixelOffset uint16
}

func (rm SetRowMessage) Frame() ethernet.Frame {
	et := ethernet.EtherType(PixelRow)
	if rm.RowIndex >= 256 && rm.RowIndex <= 511 {
		et += 1
	} else {
		et += ethernet.EtherType(rm.RowIndex >> 8)
	}
	pl := make([]byte, len(rm.Data)+7)
	pl[0] = byte(rm.RowIndex & 0xFF)
	ncol := uint16(len(rm.Data) / 3)
	binary.BigEndian.PutUint16(pl[1:2], rm.PixelOffset) // TODO: LE?
	binary.BigEndian.PutUint16(pl[3:4], ncol)
	pl[5] = 0x08
	pl[6] = 0x88
	for i, p := range rm.Data {
		pl[i*3+7+0] = p.Blue
		pl[i*3+7+1] = p.Green
		pl[i*3+7+2] = p.Red
	}
	return ethernet.Frame{
		Source:      srcAddr,
		Destination: destAddr,
		EtherType:   et,
		Payload:     pl,
	}
}

type DisplayMessage struct {
	Brightness Brightness
	ColorTemp  ColorTemp
}

func (dm DisplayMessage) Frame() ethernet.Frame {
	pl := make([]byte, 98)
	pl[21] = byte(dm.Brightness)
	pl[22] = 0x05
	pl[24] = byte(dm.ColorTemp.Red)
	pl[25] = byte(dm.ColorTemp.Green)
	pl[26] = byte(dm.ColorTemp.Blue)
	return ethernet.Frame{
		Source:      srcAddr,
		Destination: destAddr,
		EtherType:   ethernet.EtherType(DisplayFrame),
		Payload:     pl,
	}
}

type DetectCardMessage struct{}

func (DetectCardMessage) Frame() ethernet.Frame {
	return ethernet.Frame{
		Source:      srcAddr,
		Destination: destAddr,
		EtherType:   ethernet.EtherType(DetectCard),
		Payload:     make([]byte, 270),
	}
}

type DetectCardResponseAckMessage struct{ ControllerID uint8 }

func (dm DetectCardResponseAckMessage) Frame() ethernet.Frame {
	pl := make([]byte, 270)
	pl[2] = dm.ControllerID
	return ethernet.Frame{
		Source:      srcAddr,
		Destination: destAddr,
		EtherType:   ethernet.EtherType(DetectCard),
		Payload:     pl,
	}
}

type CardInfo struct {
	Version      uint8
	VersionMajor uint8
	VersionMinor uint8
	NumColumns   uint16
	NumRows      uint16
	FrameCounter uint16
	RuntimeMs    uint32
	ControllerID uint8
}

func ParseCardInfo(f ethernet.Frame) (ci CardInfo, err error) {
	if f.EtherType != ethernet.EtherType(DetectCardResponse) {
		err = fmt.Errorf("colorlight: unable to parse card info: wrong type, expected %x but was %x", DetectCardResponse, f.EtherType)
		return
	}
	if !bytes.Equal(f.Source, srcAddr) {
		err = fmt.Errorf("colorlight: unable to parse card info: wrong src addr, expected %v but was %v", srcAddr, f.Source)
		return
	}
	if !bytes.Equal(f.Destination, bcastAddr) {
		err = fmt.Errorf("colorlight: unable to parse card info: wrong src addr, expected %v but was %v", bcastAddr, f.Destination)
		return
	}
	if len(f.Payload) < 1056 {
		err = fmt.Errorf("colorlight: unable to parse card info: incomplete payload, expected 1056 bytes but was %d", len(f.Payload))
		return
	}
	ci.Version = f.Payload[0]
	ci.VersionMajor = f.Payload[1]
	ci.VersionMinor = f.Payload[2]
	ci.NumColumns = binary.BigEndian.Uint16(f.Payload[20:21])
	ci.NumRows = binary.BigEndian.Uint16(f.Payload[22:23])
	ci.FrameCounter = binary.BigEndian.Uint16(f.Payload[39:40])
	ci.RuntimeMs = binary.BigEndian.Uint32(f.Payload[45:48])
	ci.ControllerID = f.Payload[62]
	return
}

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
