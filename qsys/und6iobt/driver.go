package und6iobt

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vanti-dev/assessment-syseng-go/comm"
)

// Driver implements interrogation and control of a Q-SYS unD6IO-BT device.
type Driver struct {
	// Comm abstracts the underlying transport communication with the device.
	// We can assume that the underlying implementation of Comm is compatible with our device model.
	Comm comm.Transport
}

func (d Driver) Announce() error {
	// call Activate Pairing command (13.21)
	// command: BTB<CR>
	// Expected Response: ACK BTB<CR>

	fmt.Println("Announce Function Called")

	err:= d.Comm.Connect()
	if (err != nil) {
		return errors.New("connect failed")
	}

	n, err := d.Comm.Write([]byte("BTB\r"))
    if err != nil {
        return errors.New("write failed")
    }
    fmt.Println("wrote ", n, " bytes")
	
	//Create a byte array
	buf := make([]byte, 8)
	for {
		n, err := d.Comm.Read(buf)
		if err != nil {
			return errors.New("read failed")
		}
		fmt.Printf("n = %v err = %v buf = %v\n", n, err, buf)
		fmt.Printf("buf[:n] = %q\n", buf[:n])
		if err == io.EOF {
			break
		}
	}
	str1 := string(buf[:])
	split := strings.Split(str1, " ")
	// need to remove carriage return character
	if len(split) == 2 {
		if(split[0] != "ACK") {
			return errors.New("nack received")
		}
		if(split[1] != "BTB\r") {
			return errors.New("incorrect response received")
		}
	} else {
		return errors.New("response length incorrect")
	}

	err = d.Comm.Close()
	if (err != nil) {
		return errors.New("close failed")
	}

	return nil
}