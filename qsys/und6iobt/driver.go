package und6iobt

import (
	"errors"
	"fmt"
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

// list of Bluetooth Status that can be returned by the und6iobt device
// see section 13.20 of the 'QSC unIFY 3rd Party Control Software API document'
const (
	Idle                          = 0
	Discoverable                  = 1
	ConnectedUnknownAVRCPsupport  = 2
	ConnectedAVRCPNotSupported    = 3
	ConnectedAVRCPSupported       = 4
	ConnectedAVRCPAndPDUSupported = 5
)

// function used to get the response, do some basic checks on it and 
// returns an array of the response's fields
func (d Driver) getResponse() ([]string, error) {

	// Create a buffer of 100 bytes to store the response
	buf := make([]byte, 100)
	_, err := d.Comm.Read(buf)
	if err != nil {
		return nil, errors.New("response read failed")
	}

	// convert the array of bytes to a string
	str := string(buf[:])

	// ensure the response contains the carriage return <CR>
	carrRetIndex := strings.Index(str, "\r")
	if carrRetIndex == -1 {
		return nil, errors.New("incorrect response found (no CR)")
	}

	// extract all fields from the response except the carriage return field
	str = str[:carrRetIndex]

	// split the response into an array of string using the 'space' separator
	split := strings.Split(str, " ")
	return split, nil
}

// Implementation of the HasAnnonce interface used to activate the Bluetooth 
// pairing process on the unD6IO-BT
// The pairing/connect mode can be remotely activated by using the 'Activate Pairing' command.
// see section 13.21 of the 'QSC unIFY 3rd Party Control Software API document'
// Example command: BTB<CR>
// Example response: ACK BTB<CR>
func (d Driver) Announce() error {

	// Open the Connection
	err := d.Comm.Connect()
	if err != nil {
		return errors.New("connect failed")
	}

	// Request the activation of the Bluetooth pairing process
	_ ,err = d.Comm.Write([]byte("BTB\r"))
	if err != nil {
		return errors.New("BTB request failed")
	}

	// Get response
	resp, err := d.getResponse()
	if err != nil {
		return err
	}

	// Check response length
	if len(resp) != 2 {
		err = fmt.Errorf("BTB response length incorrect - received %d fields (expected 2)", len(resp))
		return err
	}

	// Check ACK was received
	if resp[0] != "ACK" {
		err = fmt.Errorf("BTB ACK not received - %s received instead", resp[0])
		return err
	}

	// Check response command field
	if resp[1] != "BTB" {
		err = fmt.Errorf("BTB response incorrect command received - received %s", resp[1])
		return err
	}

	// Close the connection
	err = d.Comm.Close()
	if err != nil {
		return errors.New("close failed")
	}

	return nil
}

// Implementation of the HasConnection interface used to report the Bluetooth 
// connection status of the unD6IO-BT
// The Bluetooth status can be remotely returned by using the 'Bluetooth Status' command.
// see section 13.20 of the 'QSC unIFY 3rd Party Control Software API document'
// Example command: BTS<CR>
// Example response: ACK BTS 2<CR>
// Note: The ConnectionChanged blocks until the connection status is different from last.
func (d Driver) ConnectionChanged(last bluetooth.Connection) (bluetooth.Connection, error) {

	// Open the Connection
	err := d.Comm.Connect()
	if err != nil {
		return bluetooth.ConnectionUnknown, errors.New("connect failed")
	}

	// initialise new connection status as current status
	new := last

	// loop until status are same
	for last != new {

		// Request the current bluetooth connection status
		_, err = d.Comm.Write([]byte("BTS\r"))
		if err != nil {
			return bluetooth.ConnectionUnknown, errors.New("BTS command failed")
		}

		// Get response
		resp, err := d.getResponse()
		if err != nil {
			return bluetooth.ConnectionUnknown, err
		}

		// Check response length
		if len(resp) != 2 {
			err = fmt.Errorf("BTS response length incorrect - received %d fields (expected 2)", len(resp))
			return bluetooth.ConnectionUnknown, err
		}

		// Check ACK was received
		if resp[0] != "ACK" {
			err = fmt.Errorf("BTS ACK not received - %s received instead", resp[0])
			return bluetooth.ConnectionUnknown, err
		}

		// Check response command field
		if resp[1] != "BTS" {
			err = fmt.Errorf("BTS response incorrect command received - received %s", resp[1])
			return bluetooth.ConnectionUnknown, err
		}

		// get connection status
		state, err := strconv.Atoi(resp[2])
		if err == nil {
			return bluetooth.ConnectionUnknown, err
		}

		// Check if received status is connected or disconnected
		if state == Idle || state == Discoverable {
			new = bluetooth.ConnectionNotConnected
		} else if state == ConnectedUnknownAVRCPsupport ||
			state == ConnectedAVRCPNotSupported ||
			state == ConnectedAVRCPSupported ||
			state == ConnectedAVRCPAndPDUSupported {
			new = bluetooth.ConnectionConnected
		} else {
			err = fmt.Errorf("BTS incorrect status received %d - expected [0-5])", state)
			return bluetooth.ConnectionUnknown, err
		}
	}

	// Close the connection
	err = d.Comm.Close()
	if err != nil {
		return bluetooth.ConnectionUnknown, errors.New("close failed")
	}

	return new, nil
}

// Implementation of the HasName interface used to get the Bluetooth 
// Friendly Name on the unD6IO-BT
// The friendly name can be remotely read by using the 'Get Friendly Name' command.
// see section 13.17 of the 'QSC unIFY 3rd Party Control Software API document'
// Example command: BTN<CR>
// Example response: ACK BTN unD6IO-BT-010203<CR>
func (d Driver) Name() (string, error) {

	// Open the Connection
	err := d.Comm.Connect()
	if err != nil {
		return "", errors.New("connect failed")
	}

	// Request the bluetooth friendly name
	_, err = d.Comm.Write([]byte("BTN\r"))
	if err != nil {
		return "", errors.New("BTN request failed")
	}

	// Get response
	resp, err := d.getResponse()
	if err != nil {
		return "", err
	}

	// Check response length
	if len(resp) != 3 {
		err = fmt.Errorf("BTN response length incorrect - received %d fields (expected 3)", len(resp))
		return "", err
	}

	// Check ACK was received
	if resp[0] != "ACK" {
		err = fmt.Errorf("BTN ACK not received - %s received instead", resp[0])
		return "", err
	}

	// Check response command field
	if resp[1] != "BTN" {
		err = fmt.Errorf("BTN response incorrect command received - received %s", resp[1])
		return "", err
	}

	// retrieve bluetooth friendly name
	friendly_name := resp[2]

	// Close the connection
	err = d.Comm.Close()
	if err != nil {
		return "", errors.New("close failed")
	}

	return friendly_name, nil
}
