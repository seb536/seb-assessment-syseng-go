package main

import (
	"fmt"

	"github.com/vanti-dev/assessment-syseng-go/bluetooth"
	"github.com/vanti-dev/assessment-syseng-go/qsys/und6iobt"
)

func main() {
	d := und6iobt.Driver{}
	
	d.Announce()

	mystring, error := d.Name()
	if error == nil {
		fmt.Printf(" result = %s\n", mystring)
	}

	connection, error := d.ConnectionChanged(bluetooth.ConnectionUnknown)
	if error == nil {
		fmt.Printf(" result = %d\n", connection)
	}
}
