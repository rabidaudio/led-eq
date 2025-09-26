package colorlight

import (
	"testing"
)

// func failIfErr(t *testing.T, err error) {
// 	assert.NoError(t, err)
// 	if err != nil {
// 		t.Fail()
// 	}
// }

func TestDetect(t *testing.T) {
	// ifi, err := net.InterfaceByName("en4")
	// failIfErr(t, err)

	// open a raw socket. here for macos AF_PACKET isn't supported so we're using an ip4 socket...
	// s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
	// failIfErr(t, err)
	// ...and specifying we're using our own headers
	// err = syscall.SetsockoptInt(s, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	// failIfErr(t, err)
	// also specify the device to use TODO: linux only??
	// err = syscall.SetsockoptString(s, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ifi.Name)
	// failIfErr(t, err)
	// setsockopt(sockfd, SOL_SOCKET, SO_BINDTODEVICE, ifname, strlen(ifname))
	// addr := syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}} // localhost (doesn't matter really)

	// s, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	// failIfErr(t, err)
	// err = syscall.SetsockoptString(s, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ifi.Name)
	// failIfErr(t, err)

	// err = syscall.Sendto(s, b)
	// failIfErr(t, err)

	// f := ethernet.Frame{
	// 	Source:      SrcAddr,
	// 	Destination: DestAddr,
	// 	EtherType:   0x0700,
	// 	Payload:     make([]byte, 270),
	// }
	// b, err := f.MarshalBinary()
	// failIfErr(t, err)
	// fmt.Printf("%x\n", b)

	// con, err := raw.ListenPacket(ifi, uint16(f.EtherType), &raw.Config{BPFDirection: 1})
	// failIfErr(t, err)
	// defer con.Close()

	// // err = con.SetPromiscuous(true)
	// // failIfErr(t, err)

	// t.Logf("local: %v", con.LocalAddr())

	// v, err := con.WriteTo(b, &raw.Addr{HardwareAddr: DestAddr})
	// failIfErr(t, err)
	// t.Logf("v: %v", v)
}
