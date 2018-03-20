package alsa_wrapper

/*
#ifndef ARREMI
#define ARREMI
#cgo LDFLAGS: -lasound
#include <stdlib.h>
#include "alsa_wrapper.h"
#endif
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// NewClient to register new ALSA client.
func NewClient(name string) error {
	cDevName := C.CString(name)
	defer C.free(unsafe.Pointer(cDevName))
	var status = int(C.new_client(cDevName))
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
}

// NewPort to reigister new MIDI port.
func NewPort(name string) error {
	cPortName := C.CString(name)
	defer C.free(unsafe.Pointer(cPortName))
	var status = int(C.new_port(cPortName))
	if status != 0 {
		return fmt.Errorf("Error while createing sequencer port. %d", status)
	}
}

// SendData to ALSA
func SendData(p []byte) error {
	cData := (*C.char)(unsafe.Pointer(&p[0]))
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
}

func resolveErrCode(code int) (int, int) {
	return code >> 16, code & 0xffff
}
