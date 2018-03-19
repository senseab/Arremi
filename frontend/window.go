package frontend

import (
	"fmt"
	"time"

	"regexp"

	"github.com/andlabs/ui"
	"github.com/tonychee7000/Arremi/consts"
	"github.com/tonychee7000/Arremi/serialPort"
	"github.com/tonychee7000/arremi/backend"
)

// WindowMain is to show the UI
func WindowMain() {
	// Let backend stop
	//var ch = make(chan int)

	var errCh = make(chan error, 0)
	var chCon = make(chan int, 1)
	var chSerialName = make(chan string, 1)
	var serialList []string

	window := ui.NewWindow(consts.WindowTitle, 300, 100, false)

	if backend.MIDIError != nil {
		ui.MsgBox(window, "MIDI Failure",
			fmt.Sprint("MIDI driver initialize failed: ", backend.MIDIError))
		ui.Quit()
	}

	labelSerialSelection := ui.NewLabel("Serial Port")
	serialSelection := ui.NewCombobox()
	buttonRefresh := ui.NewButton("Refresh")
	buttonRun := ui.NewCheckbox("Run")
	labelMidiActive := ui.NewLabel("MIDI Active  ")
	labelMidiSignal := ui.NewLabel(consts.MIDISignalOff)
	labelHint := ui.NewLabel(
		fmt.Sprint(
			"HINT:\nBaudrate of Arduino USB serial port should be set to ",
			consts.SerialBaudrate, " by using\n\n\tSerial.begin(", consts.SerialBaudrate, ")",
		),
	)

	serialList, err := loadSerial(serialSelection, buttonRun)

	if err != nil {
		ui.MsgBox(window, "Serial Failure",
			fmt.Sprint("Cannot list serial ports for reason: ", err))
		ui.Quit()
	}

	buttonRefresh.OnClicked(func(*ui.Button) {
		serialList, err = loadSerial(serialSelection, buttonRun)
		if err != nil {
			ui.MsgBox(window, "Serial Failure",
				fmt.Sprint("Cannot list serial ports for reason: ", err))
			ui.Quit()
		}
	})

	buttonRun.OnToggled(func(*ui.Checkbox) {
		if buttonRun.Checked() {
			serialSelection.Disable()
			buttonRefresh.Disable()
			chSerialName <- serialList[serialSelection.Selected()]
		} else {
			serialSelection.Enable()
			buttonRefresh.Enable()
			chCon <- 1
		}
	})

	mainBox := ui.NewVerticalBox()
	topBox := ui.NewHorizontalBox()
	buttomBox := ui.NewHorizontalBox()

	topBox.Append(serialSelection, true)
	topBox.Append(buttonRefresh, false)
	topBox.Append(ui.NewLabel("  "), false)
	topBox.Append(buttonRun, false)

	buttomBox.Append(labelSerialSelection, false)
	buttomBox.Append(ui.NewHorizontalSeparator(), true)
	buttomBox.Append(labelMidiActive, false)
	buttomBox.Append(labelMidiSignal, false)

	mainBox.Append(ui.NewHorizontalBox(), true)
	mainBox.Append(buttomBox, false)
	mainBox.Append(topBox, false)
	mainBox.Append(labelHint, true)
	mainBox.Append(ui.NewHorizontalBox(), true)

	window.SetMargined(true)
	window.SetChild(mainBox)
	window.OnClosing(func(*ui.Window) bool {
		if buttonRun.Checked() {
			return false
		}
		ui.Quit()
		return true
	})

	go func() {
		for {
			select {
			case err := <-errCh:
				ui.QueueMain(func() {
					if err != nil &&
						!regexp.MustCompile("file already closed").MatchString(err.Error()) {
						ui.MsgBox(window, "Backend Failure", fmt.Sprint(err))
					}
					serialSelection.Enable()
					buttonRefresh.Enable()
					buttonRun.SetChecked(false)
					labelMidiSignal.SetText(consts.MIDISignalOff)
				})
			case <-backend.MidiDev.Signal:
				ui.QueueMain(func() {
					labelMidiSignal.SetText(consts.MIDISignalOn)
				})
			case <-time.After(100 * time.Millisecond):
				ui.QueueMain(func() {
					labelMidiSignal.SetText(consts.MIDISignalOff)
				})
			}
		}
	}()

	go func() {
		for {
			backend.Run(chSerialName, chCon, errCh)
		}
	}()

	window.Show()
}

func loadSerial(serialSelection *ui.Combobox, buttonRun *ui.Checkbox) ([]string, error) {
	serialList, err := serialPort.GetSerialPorts()

	if len(serialList) == 0 {
		buttonRun.Disable()
		serialSelection.Append("-- No Arduino Found --")
		serialSelection.Disable()
	} else {
		for _, seiral := range serialList {
			serialSelection.Append(seiral)
		}
	}
	serialSelection.SetSelected(0)
	return serialList, err
}
