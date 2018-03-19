package serialPort

import (
	"io/ioutil"
	"regexp"
)

// GetSerialPorts is to list all serial port for Arduino device.
func GetSerialPorts() ([]string, error) {
	f, err := ioutil.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range f {
		if regexp.MustCompile("^ttyACM([0-9]+)|^cu.usbmodem").MatchString(file.Name()) {
			fileList = append(fileList, file.Name())
		}
	}

	return fileList, nil
}
