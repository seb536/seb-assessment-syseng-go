package und6iobt

import "github.com/vanti-dev/techass-syseng-go/comm"

// Driver implements interrogation and control of a Q-SYS unD6IO-BT device.
type Driver struct {
	Comm comm.Endpoint
}
