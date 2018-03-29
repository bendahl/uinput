package uinput

import (
	"testing"
)

// This test will confirm that all basic key events are working
func TestBasicKeyboard(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.KeyPress(Key1)
	if err != nil {
		t.Fatalf("Failed to send key press. Last error was: %s\n", err)
	}

	err = vk.KeyDown(Key1)
	if err != nil {
		t.Fatalf("Failed to send key down event. Last error was: %s\n", err)
	}

	err = vk.KeyUp(Key1)
	if err != nil {
		t.Fatalf("Failed to send key up event. Last error was: %s\n", err)
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

// This test confirms that all basic mouse events are working as expected.
func TestBasicMouseMoves(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}

	err = relDev.MoveLeft(100)
	if err != nil {
		t.Fatalf("Failed to move mouse left. Last error was: %s\n", err)
	}

	err = relDev.MoveRight(150)
	if err != nil {
		t.Fatalf("Failed to move mouse right. Last error was: %s\n", err)
	}

	err = relDev.MoveUp(50)
	if err != nil {
		t.Fatalf("Failed to move mouse up. Last error was: %s\n", err)
	}

	err = relDev.MoveDown(100)
	if err != nil {
		t.Fatalf("Failed to move mouse down. Last error was: %s\n", err)
	}

	err = relDev.RightClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = relDev.LeftClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = relDev.LeftPress()
	if err != nil {
		t.Fatalf("Failed to perform left key press. Last error was: %s\n", err)
	}

	err = relDev.LeftRelease()
	if err != nil {
		t.Fatalf("Failed to perform left key release. Last error was: %s\n", err)
	}

	err = relDev.RightPress()
	if err != nil {
		t.Fatalf("Failed to perform right key press. Last error was: %s\n", err)
	}

	err = relDev.RightRelease()
	if err != nil {
		t.Fatalf("Failed to perform right key release. Last error was: %s\n", err)
	}

	err = relDev.Close()
	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestBasicTouchPadMoves(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}

	err = absDev.MoveTo(0, 0)
	if err != nil {
		t.Fatalf("Failed to move cursor to initial position. Last error was: %s\n", err)
	}

	err = absDev.MoveTo(100, 200)
	if err != nil {
		t.Fatalf("Failed to move cursor to position x:100, y:200. Last error was: %s\n", err)
	}
	err = absDev.RightClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = absDev.LeftClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = absDev.Close()
	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}
