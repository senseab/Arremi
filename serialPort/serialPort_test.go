package serialPort

import (
	"testing"
)

func TestGetSerialPorts(T *testing.T) {
	portList, err := GetSerialPorts()
	if err != nil {
		T.Error("ERROR: ", err, "\n")
	}
	T.Log("Show all serial ports")
	for _, p := range portList {
		T.Log("\t", p)
	}
}

