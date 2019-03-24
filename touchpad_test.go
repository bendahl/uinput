package uinput

import (
	"fmt"
	"testing"
)

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

	err = absDev.LeftPress()
	if err != nil {
		t.Fatalf("Failed to perform left key press. Last error was: %s\n", err)
	}

	err = absDev.LeftRelease()
	if err != nil {
		t.Fatalf("Failed to perform left key release. Last error was: %s\n", err)
	}

	err = absDev.RightPress()
	if err != nil {
		t.Fatalf("Failed to perform right key press. Last error was: %s\n", err)
	}

	err = absDev.RightRelease()
	if err != nil {
		t.Fatalf("Failed to perform right key release. Last error was: %s\n", err)
	}

	err = absDev.Close()
	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestTouchPadCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual)
			}
		}
	}()
	_, _ = CreateTouchPad("", []byte("TouchDevice"), 0, 1024, 0, 768)
	t.Fatalf("Empty path did not yield a panic")
}

func TestTouchPadCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	expected := "device path '" + path + "' does not exist"
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual)
			}
		}
	}()
	_, _ = CreateTouchPad(path, []byte("TouchDevice"), 0, 1024, 0, 768)
	t.Fatalf("Invalid path did not yield a panic")
}

func TestTouchPadCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual)
			}
		}
	}()
	_, _ = CreateTouchPad("/dev/uinput", []byte(name), 0, 1024, 0, 768)
	t.Fatalf("Invalid name did not yield a panic")
}

func TestTouchPadMoveToFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.MoveTo(1, 1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestTouchPadLeftClickFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.LeftClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestTouchPadLeftPressFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.LeftPress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestTouchPadLeftReleaseFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.LeftRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestTouchPadRightClickFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.RightClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestTouchPadRightPressFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.RightPress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestTouchPadRightReleaseFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	_ = absDev.Close()
	err = absDev.RightRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMultipleTouchPadsWithDifferentSizes(t *testing.T) {
	horizontal, err := CreateTouchPad("/dev/uinput", []byte("horizontal_pad"), 0, 200, 0, 100)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer horizontal.Close()
	vertical, err := CreateTouchPad("/dev/uinput", []byte("vertical_pad"), 0, 100, 0, 200)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer vertical.Close()
	err = horizontal.MoveTo(200, 100)
	if err != nil {
		t.Fatalf("Unable to move cursor on horizontal pad: %v", err)
	}

	err = vertical.MoveTo(100, 200)
	if err != nil {
		t.Fatalf("Unable to move cursor on horizontal pad: %v", err)
	}

}

func TestPositioningInUpperLeftCorner(t *testing.T) {
	dev, err := CreateTouchPad("/dev/uinput", []byte("touchpad"), 0, 200, 0, 100)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer dev.Close()
	err = dev.MoveTo(0, 0)
	if err != nil {
		t.Fatalf("Failed to move cursor to upper left corner: %v", err)
	}
}
