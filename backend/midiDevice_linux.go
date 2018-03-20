package backend

/*
#ifndef ARREMI
#define ARREMI
#cgo LDFLAGS: -lasound
#include <stdlib.h>
#include "alsa_wrapper/alsa_wrapper.h"
#endif
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/tonychee7000/Arremi/consts"
)

// MIDIDevice implies a Writer interface.
type MIDIDevice struct {
	Signal chan int
}

// NewMIDIDevice construction func
func NewMIDIDevice() (*MIDIDevice, error) {
	var mididev = new(MIDIDevice)
	err := mididev.Init()
	return mididev, err
}

// Init the client and source
func (midiDev *MIDIDevice) Init() error {
	var status int

	midiDev.Signal = make(chan int, 4096)

	cDevName := C.CString(consts.ClientName)
	defer C.free(unsafe.Pointer(cDevName))
	status = int(C.new_client(cDevName))
	if status != 0 {
		stage, errCode := resolveErrCode(status)
		switch stage {
		case 1:
			return fmt.Errorf("Error while opening sequencer. %d", errCode)
		case 2:
			return fmt.Errorf("Error while getting sequencer id. %d", errCode)
		case 3:
			return fmt.Errorf("Error while setting sequencer name. %d", errCode)
		}
	}

	cPortName := C.CString(consts.SourceName)
	defer C.free(unsafe.Pointer(cPortName))
	status = int(C.new_port(cPortName))
	if status != 0 {
		return fmt.Errorf("Error while createing sequencer port. %d", status)
	}

	return nil
}

func (midiDev *MIDIDevice) Write(p []byte) (int, error) {
	midiDev.Signal <- 1
	cData := (*C.CChar)(unsafe.Pointer(&p[0]))
	defer C.free(unsafe.Pointer(cData))
	var status = int(C.send_data(cData, C.int(len(p))))
	if status != 0 {
		stage, errCode := resolveErrCode(status)
		switch stage {
		case 1:
			return 0, fmt.Errorf("Error while creating MIDI event. %d", errCode)
		case 2:
			return 0, fmt.Errorf("Error while encoding MIDI event. %d", errCode)
		case 3:
			return 0, fmt.Errorf("Error while sending data. %d", errCode)
		}
	}
	return len(p), nil
}

// AllNoteOff I don't want panic!
func (midiDev *MIDIDevice) AllNoteOff() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 128; j++ {
			midiDev.Write([]byte{byte(0x90 + i), byte(j), 0})
		}
	}
}

func resolveErrCode(int code) (int, int) {
	return code >> 16, code & 0xffff
}
