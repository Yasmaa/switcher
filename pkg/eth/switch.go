package eth

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"sync"
)

// STUDENTS: edit this file

// EthernetSwitch implements the ethernet switch functionality
type EthernetSwitch struct {
	mu sync.Mutex

	sendQueueSize int
	ports         []Port
	macTable      map[MACAddress]int
}

// NewEthernetSwitch create a new switch with sendQueueSize for the size of the sending buffers and the provided ports.
func NewEthernetSwitch(sendQueueSize int, ports ...Port) *EthernetSwitch {
	// STUDENTS: implement this
	return &EthernetSwitch{
		sendQueueSize: sendQueueSize,
		ports:         ports,
		macTable:      make(map[MACAddress]int),
	}
}

// Run the ethernet switch.
// Blocks until all the io.Readers of the ports are closed (return io.EOF).
// Returns any unrecoverable error from reading (other than io.EOF) or writing to the ports.
// Before returning, this closes the writer for each port.
func (sw *EthernetSwitch) Run() error {
	// Create channels for each port's send queue
	sendChans := make([]chan *Frame, len(sw.ports))
	for i := range sendChans {
		sendChans[i] = make(chan *Frame, sw.sendQueueSize)
	}

	// Start goroutine to forward frames from each port's send queue to the appropriate port
	for i, sendChan := range sendChans {
		go func(portID int, sendChan <-chan *Frame) {
			for frame := range sendChan {
				if _, err := WriteFrame(sw.ports[portID], *frame); err != nil {
					break
				}
			}
			err := sw.ports[portID].Close()
			if err != nil {
				return
			}
		}(i, sendChan)
	}

	// the `main` goroutine sends frames to the channel `sendChans`
	for i, port := range sw.ports {
		go sw.pool(port, i, sendChans)

	}

	return nil
}

func (sw *EthernetSwitch) pool(port Port, i int, sendChans []chan *Frame) {
	defer func() { close(sendChans[i]) }()
	for {
		frame, err := ReadFrame(port)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			// break out of the loop when the port is closed
			break
		}
		if frame != nil {
			sw.forward(i, frame, sendChans)
		}
	}
}

func (sw *EthernetSwitch) forward(i int, frame *Frame, sendChans []chan *Frame) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.macTable[frame.Source] = i
	dstPort, ok := sw.macTable[frame.Destination]
	if !ok {
		// broadcast the frame to all ports except the source port
		for j := range sendChans {
			if j != i {
				select {
				case sendChans[j] <- frame:
				default: // if the send queue is full, drop the frame and move on
				}
			}
		}
	} else {
		select {
		case sendChans[dstPort] <- frame:
		default:
		}
	}

}

// RunSize returns the number of elements of the MAC table
// This may only be called while Run() is called.
func (sw *EthernetSwitch) RunSize() int {
	// STUDENTS: implement this
	sw.mu.Lock()
	size := len(sw.macTable)
	defer sw.mu.Unlock()
	return size
}

// ReadFrame reads a single frame from r.
// If the frame is not valid, return a nil Frame and a nil error.
func ReadFrame(r io.Reader) (*Frame, error) {
	// STUDENTS: implement this
	header := make([]byte, 14)

	if _, err := io.ReadFull(r, header); err != nil {
		if errors.Is(err, io.EOF) {

			return nil, io.EOF
		}
		return nil, errors.Unwrap(err)
	}

	size := binary.BigEndian.Uint16(header[12:14])

	if size > 1500 {

		return nil, fmt.Errorf("invalid Ethernet II frame")
	}

	data := make([]byte, size)

	if _, err := io.ReadFull(r, data); err != nil {
		return nil, errors.Unwrap(err)
	}

	var recievedCrc uint32

	if err := binary.Read(r, binary.BigEndian, &recievedCrc); err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}

	if _, err := buffer.Write(header); err != nil {
		return nil, errors.Unwrap(err)
	}

	if _, err := buffer.Write(data); err != nil {
		return nil, errors.Unwrap(err)
	}
	crc := crc32.ChecksumIEEE(buffer.Bytes())

	if crc != recievedCrc {

		return nil, nil
	}

	return &Frame{
		Destination: MACAddress{header[0], header[1], header[2], header[3], header[4], header[5]},
		Source:      MACAddress{header[6], header[7], header[8], header[9], header[10], header[11]},
		Data:        data,
	}, nil
}

// WriteFrame write the ethernet frame to w.
func WriteFrame(w io.Writer, frame Frame) (int, error) {
	// STUDENTS: implement this
	size := len(frame.Data)
	if size > 1500 {

		return 0, fmt.Errorf("invalid Ethernet II frame")
	}

	header := make([]byte, 14)

	copy(header[:6], frame.Destination[:])
	copy(header[6:12], frame.Source[:])

	binary.BigEndian.PutUint16(header[12:14], uint16(size))

	buffer := &bytes.Buffer{}

	n, err := buffer.Write(header)
	if err != nil {
		return n, errors.Unwrap(err)
	}

	m, err := buffer.Write(frame.Data)
	if err != nil {
		return n + m, errors.Unwrap(err)
	}
	crc := crc32.ChecksumIEEE(buffer.Bytes())

	if _, err := buffer.WriteTo(w); err != nil {
		return 0, errors.Unwrap(err)
	}

	if err := binary.Write(w, binary.BigEndian, crc); err != nil {
		return 0, err
	}

	return n + m + 4, nil
}
