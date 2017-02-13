package uinput

import (
	"testing"
)

// This test will create a basic VKeyboard, send a key command and then close the keyboard device
func TestBasicKeyboard(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.SendKeyPress(Key1)

	if err != nil {
		t.Fatalf("Failed to send key event. Last error was: %s\n", err)
	}

	err = vk.SendKeyRelease(Key1)

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
	vk, err := CreateKeyboard("/invalid/path", []byte("Invalid Device Path"))
	if err == nil {
		// this usually shouldn't happen, but if the device is created, we need to close it
		vk.Close()
		t.Fatalf("Expected error code,but received %s instead.\n", err)
	}
}

// This test will confirm that a proper error code is returned if an invalid keycode is
// passed to the library
func TestInvalidKeycode(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.SendKeyPress(4711)
	if err == nil {
		t.Fatalf("Sending an invalid keycode did not trigger an error.\n")
	}

	vk.Close()
}

// This test will create a basic Mouse, send an absolute axis event and close the device
func TestBasicMouse(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}

	err = relDev.MoveCursor(1000, 100)

	if err != nil {
		t.Fatalf("Failed to send key event. Last error was: %s\n", err)
	}

	err = relDev.Close()

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}
