package backend

import (
	alsa "github.com/tonychee7000/Arremi/backend/alsa_wrapper"
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
	var err error
	midiDev.Signal = make(chan int, 4096)
	err = alsa.NewClient(consts.ClientName)
	if err != nil {
		return err
	}

	err = alsa.NewPort(consts.SourceName)
	if err != nil {
		return err
	}

	return nil
}

func (midiDev *MIDIDevice) Write(p []byte) (int, error) {
	midiDev.Signal <- 1
	alsa.SendData(p)
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
