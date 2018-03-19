package backend

import (
	"github.com/tonychee7000/Arremi/consts"
	midi "github.com/youpy/go-coremidi"
)

// MIDIDevice implies a Writer interface.
type MIDIDevice struct {
	client midi.Client
	source midi.Source
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
	midiDev.client, err = midi.NewClient(consts.ClientName)
	if err != nil {
		return err
	}

	midiDev.source, err = midi.NewSource(midiDev.client, consts.SourceName)
	return err
}

func (midiDev *MIDIDevice) Write(p []byte) (int, error) {
	var pack = midi.NewPacket(p, 0)
	midiDev.Signal <- 1
	err := pack.Received(&(midiDev.source))
	return len(p), err
}

// AllNoteOff I don't want panic!
func (midiDev *MIDIDevice) AllNoteOff() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 128; j++ {
			midiDev.Write([]byte{byte(0x90 + i), byte(j), 0})
		}
	}
}
