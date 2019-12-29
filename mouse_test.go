package uinput

import (
	"fmt"
	"os"
	"testing"
)

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

	err = relDev.Wheel(false, 1)
	if err != nil {
		t.Fatalf("Failed to perform wheel movement. Last error was: %s\n", err)
	}

	err = relDev.Wheel(true, 1)
	if err != nil {
		t.Fatalf("Failed to perform horizontal wheel movement. Last error was: %s\n", err)
	}

	err = relDev.Close()
	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestMouseCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateMouse("", []byte("MouseDevice"))
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMouseCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateMouse(path, []byte("MouseDevice"))
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestMouseCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateMouse("/dev/uinput", []byte(name))
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMouseLeftClickFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.LeftClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseLeftPressFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.LeftPress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseLeftReleaseFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.LeftRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseRightClickFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.RightClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseRightPressFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.RightPress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestVMouse_RightReleaseFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.RightRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseMoveUpFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MoveUp(1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseMoveDownFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MoveDown(1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseMoveLeftFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MoveLeft(1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseMoveRightFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MoveRight(1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseWheelFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouse("/dev/uinput", []byte("Test Basic Mouse"))
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.Wheel(false, 1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}
