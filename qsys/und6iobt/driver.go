package und6iobt

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/vanti-dev/assessment-syseng-go/bluetooth"
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

func (d Driver) ConnectionChanged(last bluetooth.Connection) (bluetooth.Connection, error) {
	fmt.Println("ConnectionChanged Function Called")

	// call Bluetooth Status command (13.20)
	// command: BTS<CR>
	// Example Response: ACK BTS 2<CR>

	fmt.Println("Announce Function Called")

	err:= d.Comm.Connect()
	if (err != nil) {
		return bluetooth.ConnectionUnknown, errors.New("connect failed")
	}
	new := last

	for last != new {
		n, err := d.Comm.Write([]byte("BTS\r"))
		if err != nil {
			return bluetooth.ConnectionUnknown, errors.New("write failed")
		}
		fmt.Println("wrote ", n, " bytes")
		
		//Create a byte array
		buf := make([]byte, 8)
		for {
			n, err := d.Comm.Read(buf)
			if err != nil {
				return bluetooth.ConnectionUnknown, errors.New("read failed")
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
		if len(split) == 3 {
			if(split[0] != "ACK") {
				return bluetooth.ConnectionUnknown, errors.New("nack received")
			}
			if(split[1] != "BTS\r") {
				return bluetooth.ConnectionUnknown, errors.New("incorrect response received")
			}
			
			n , err := strconv.Atoi(split[2])

			//check if error occured
			if err != nil{
			  //executes if there is any error
			  fmt.Println(err)
			}else{
				new = bluetooth.Connection(n)
			}
		} else {
			return bluetooth.ConnectionUnknown, errors.New("response length incorrect")
		}
	}

	err = d.Comm.Close()
	if (err != nil) {
		return bluetooth.ConnectionUnknown, errors.New("close failed")
	}

	return bluetooth.ConnectionUnknown, nil
}