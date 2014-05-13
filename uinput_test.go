package uinput

import (
	"testing"
)

// This test will create a basic VKeyboard, send a key command and then close the keyboard device
func TestBasicVKeyboard(t *testing.T) {
	devId, err := CreateVKeyboardDevice("/dev/uinput")

	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = SendBtnEvent(devId, KEY_1, 1)

	if err != nil {
		t.Fatalf("Failed to send key event. Last error was: %s\n", err)
	}

	err = CloseDevice(devId)

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}
