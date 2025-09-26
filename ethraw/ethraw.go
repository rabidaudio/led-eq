package ethraw

import (
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/packet"
	"github.com/rabidaudio/led-eq/colorlight"
)

type PacketInterface struct {
	iface *net.Interface
	conn  *packet.Conn
}

var _ colorlight.EthInterface = (*PacketInterface)(nil)

func New(deviceName string) (*PacketInterface, error) {
	iface, err := net.InterfaceByName(deviceName)
	if err != nil {
		return nil, err
	}
	conn, err := packet.Listen(iface, packet.Raw, 0, nil)
	if err != nil {
		return nil, err
	}
	pi := PacketInterface{
		iface: iface,
		conn:  conn,
	}
	return &pi, nil
}

func (i *PacketInterface) Send(f ethernet.Frame) error {
	b, err := f.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = i.conn.WriteTo(b, &packet.Addr{HardwareAddr: f.Destination})
	return err
}

func (i *PacketInterface) ReceiveBroadcast() (f ethernet.Frame, err error) {
	b := make([]byte, 1518) // TODO: should we support 9K jumbo frames?
	n, _, err := i.conn.ReadFrom(b)
	if err != nil {
		return
	}
	err = f.UnmarshalBinary(b[:n])
	return
}

func (i *PacketInterface) Close() error {
	if i.conn != nil {
		defer func() { i.conn = nil }()
		return i.conn.Close()
	}
	return nil
}
