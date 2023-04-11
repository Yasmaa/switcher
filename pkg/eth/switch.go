package eth

import "io"

// STUDENTS: edit this file

// EthernetSwitch implements the ethernet switch functionality
type EthernetSwitch struct {
}

// NewEthernetSwitch create a new switch with sendQueueSize for the size of the sending buffers and the provided ports.
func NewEthernetSwitch(sendQueueSize int, ports ...Port) *EthernetSwitch {
	// STUDENTS: implement this
	return nil
}

// Run the ethernet switch.
// Blocks until all the io.Readers of the ports are closed (return io.EOF).
// Returns any unrecoverable error from reading (other than io.EOF) or writing to the ports.
// Before returning, this closes the writer for each port.
func (sw *EthernetSwitch) Run() error {
	// STUDENTS: implement this
	return nil
}

// RunSize returns the number of elements of the MAC table
// This may only be called while Run() is called.
func (sw *EthernetSwitch) RunSize() int {
	// STUDENTS: implement this
	return 9999
}

// ReadFrame reads a single frame from r.
// If the frame is not valid, return a nil Frame and a nil error.
func ReadFrame(r io.Reader) (*Frame, error) {
	// STUDENTS: implement this
	return nil, nil
}

// WriteFrame write the ethernet frame to w.
func WriteFrame(w io.Writer, frame Frame) (int, error) {
	// STUDENTS: implement this
	return 0, nil
}
