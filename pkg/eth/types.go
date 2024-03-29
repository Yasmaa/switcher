// Package eth implements a software ethernet switch.
package eth

import (
	"fmt"
	"io"
	"net"
)

// STUDENTS: do not edit this file

// MACAddress is a ethernet MAC address
type MACAddress [6]byte

func (addr MACAddress) String() string {
	return net.HardwareAddr(addr[:]).String()
}

// BroadcastAddress is the MAC address to use when the desire it to broadcast to all nodes on the LAN.
var BroadcastAddress = MACAddress{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

// Frame is an ethernet frame
// see https://en.wikipedia.org/wiki/Ethernet_frame
type Frame struct {
	Source, Destination MACAddress
	Data                []byte
}

// String implements fmt.Stringer
func (f *Frame) String() string {
	if len(f.Data) <= 2 {
		return fmt.Sprintf("%s->%s data %v", f.Source, f.Destination, f.Data)
	}
	// large packet
	return fmt.Sprintf("%s->%s length %d", f.Source, f.Destination, len(f.Data))

}

// Port represents a physical ethernet port on a switch
type Port interface {
	io.ReadWriteCloser
}
