package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestBasicTouchPadMoves(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer func(absDev TouchPad) {
		err := absDev.Close()
		if err != nil {
			t.Fatalf("Failed to close device. Last error was: %s\n", err)
		}
	}(absDev)

	err = absDev.MoveTo(0, 0)
	if err != nil {
		t.Fatalf("Failed to move cursor to initial position. Last error was: %s\n", err)
	}

	err = absDev.MoveTo(100, 200)
	if err != nil {
		t.Fatalf("Failed to move cursor to position x:100, y:200. Last error was: %s\n", err)
	}

}

func TestTouchPadClicks(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer func(absDev TouchPad) {
		err := absDev.Close()
		if err != nil {
			t.Fatalf("Failed to close device. Last error was: %s\n", err)
		}
	}(absDev)

	err = absDev.RightClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = absDev.LeftClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

}

func TestTouchPadButtonPresses(t *testing.T) {
	absDev, err := CreateTouchPad("/dev/uinput", []byte("Test TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer func(absDev TouchPad) {
		err := absDev.Close()
		if err != nil {
			t.Fatalf("Failed to close device. Last error was: %s\n", err)
		}
	}(absDev)

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
}

func TestTouchPadCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateTouchPad("", []byte("TouchDevice"), 0, 1024, 0, 768)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestTouchPadCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateTouchPad(path, []byte("TouchDevice"), 0, 1024, 0, 768)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestTouchPadCreationFailsOnWrongPathName(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "uinput-touchpad-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	expected := "failed to register key device: failed to close device: inappropriate ioctl for device"
	_, err = CreateTouchPad(file.Name(), []byte("TouchDevice"), 0, 1024, 0, 768)
	if err == nil || !(expected == err.Error()) {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestTouchPadCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateTouchPad("/dev/uinput", []byte(name), 0, 1024, 0, 768)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
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

func TestSingleTouchEvent(t *testing.T) {
	dev, err := CreateTouchPad("/dev/uinput", []byte("touchpad"), 0, 200, 0, 100)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer dev.Close()

	err = dev.TouchDown()
	if err != nil {
		t.Fatalf("Failed to issue touch down event: %v", err)
	}

	err = dev.TouchUp()
	if err != nil {
		t.Fatalf("Failed to issue touch up event: %v", err)
	}

}

func TestTouchPadSyspath(t *testing.T) {
	dev, err := CreateTouchPad("/dev/uinput", []byte("TouchPad"), 0, 1024, 0, 768)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}

	sysPath, err := dev.FetchSyspath()
	if err != nil {
		t.Fatalf("Failed to fetch syspath. Last error was: %s\n", err)
	}

	if sysPath[:32] != "/sys/devices/virtual/input/input" {
		t.Fatalf("Expected syspath to start with /sys/devices/virtual/input/input, but got %s", sysPath)
	}

	t.Logf("Syspath: %s", sysPath)
}
