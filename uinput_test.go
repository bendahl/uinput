package uinput

import (
	"testing"
)

// This test will create a basic VKeyboard, send a key command and then close the keyboard device
func TestBasicVKeyboard(t *testing.T) {
	vk := VKeyboard{}
	err := vk.Create("/dev/uinput")

	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.SendKeyPress(KEY_1)

	if err != nil {
		t.Fatalf("Failed to send key event. Last error was: %s\n", err)
	}

	err = vk.SendKeyRelease(KEY_1)

	if err != nil {
		t.Fatalf("Failed to send key event. Last error was: %s\n", err)
	}

	err = vk.Close()

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

// This test will confirm that a proper error code is returned if an invalid uinput path is
// passed to the library
func TestInvalidDevicePath(t *testing.T) {
	vk := VKeyboard{}
	err := vk.Create("/invalid/path")

	if err == nil {
		// this usually shouldn't happen, but if the device is created, we need to close it
		vk.Close()
		t.Fatalf("Expected error code,but received %s instead.\n", err)
	}
}

// This test will confirm that a proper error code is returned if an invalid keycode is
// passed to the library
func TestInvalidKeycode(t *testing.T) {
	vk := VKeyboard{}
	err := vk.Create("/dev/uinput")

	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.SendKeyPress(4711)

	if err == nil {
		t.Fatalf("Sending an invalid keycode did not trigger an error. Got: %d.\n")
	}

	vk.Close()
}
