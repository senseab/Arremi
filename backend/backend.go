package backend

import (
	"io"

	"github.com/jacobsa/go-serial/serial"
	"github.com/tonychee7000/Arremi/consts"
)

var (
	// MidiDev is a global midi device.
	MidiDev *MIDIDevice

	// MIDIError check this. if not nil, go exit
	MIDIError error
)

func init() {
	MidiDev, MIDIError = NewMIDIDevice()
}

// Run called by frontend.
func Run(chSerialName chan string, ch chan int, errCh chan error) {
	serialName := <-chSerialName

	sPort, err := serial.Open(serial.OpenOptions{
		PortName:        "/dev/" + serialName,
		BaudRate:        consts.SerialBaudrate,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 3,
	})
	if err != nil {
		errCh <- err
		return
	}
	defer MidiDev.AllNoteOff()
	defer sPort.Close()

	go func() {
		_, err := io.Copy(MidiDev, sPort)
		if err != nil {
			errCh <- err
		}
		ch <- 1
	}()

	for {
		select {
		case <-ch:
			return
		}
	}
}
