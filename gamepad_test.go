package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// This test inputs the konami code
func TestInfiniteKonami(t *testing.T) {
	vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 3; j++ {
			err = vg.ButtonPress(ButtonDpadUp)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

			err = vg.ButtonPress(ButtonDpadDown)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

		}

		for j := 0; j < 3; j++ {
			err = vg.ButtonPress(ButtonDpadLeft)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

			err = vg.ButtonPress(ButtonDpadRight)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

		}

		err = vg.ButtonPress(ButtonSouth)
		if err != nil {
			t.Fatalf("Failed to send button press. Last error was: %s\n", err)
		}

		err = vg.ButtonPress(ButtonEast)
		if err != nil {
			t.Fatalf("Failed to send button press. Last error was: %s\n", err)
		}

		err = vg.ButtonPress(ButtonStart)
		if err != nil {
			t.Fatalf("Failed to send button press. Last error was: %s\n", err)
		}
	}
}

// This test moves the axes around a bit
func TestAxisMovement(t *testing.T) {
	vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}

	err = vg.LeftStickMove(0.2, 1.0)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}

	err = vg.LeftStickMoveX(0.2)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}

	err = vg.LeftStickMoveY(1)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}

	err = vg.RightStickMove(0.2, 1.0)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}

	err = vg.RightStickMoveX(0.2)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}

	err = vg.RightStickMoveY(1)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}
}

func TestHatMovement(t *testing.T) {
	vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}

	err = vg.HatPress(HatDirection(HatUp))
	if err != nil {
		t.Fatalf("Falied to move hat up")
	}
	err = vg.HatRelease(HatDirection(HatUp))
	if err != nil {
		t.Fatalf("Failed to release hat")
	}
	err = vg.HatPress(HatDirection(HatRight))
	if err != nil {
		t.Fatalf("Falied to move hat right")
	}
	err = vg.HatRelease(HatDirection(HatRight))
	if err != nil {
		t.Fatalf("Failed to release hat")
	}
	err = vg.HatPress(HatDirection(HatDown))
	if err != nil {
		t.Fatalf("Falied to move hat down")
	}
	err = vg.HatRelease(HatDirection(HatDown))
	if err != nil {
		t.Fatalf("Failed to release hat")
	}
	err = vg.HatPress(HatDirection(HatLeft))
	if err != nil {
		t.Fatalf("Falied to move hat left")
	}
	err = vg.HatRelease(HatDirection(HatLeft))
	if err != nil {
		t.Fatalf("Failed to release hat")
	}
}

func TestGamepadCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateGamepad("", []byte("Gamepad"), 0xDEAD, 0xBEEF)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestGamepadCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateGamepad(path, []byte("Gamepad"), 0xDEAD, 0xBEEF)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestGamepadCreationFailsOnWrongPathName(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "uinput-gamepad-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	expected := "failed to register virtual gamepad device: failed to close device: inappropriate ioctl for device"
	_, err = CreateGamepad(file.Name(), []byte("GamepadDevice"), 0xDEAD, 0xBEEF)
	if err == nil || !(expected == err.Error()) {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestGamepadCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateGamepad("/dev/uinput", []byte(name), 0xDEAD, 0xBEEF)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestGamepadMoveToFailsOnClosedDevice(t *testing.T) {
	gamepad, err := CreateGamepad("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
	_ = gamepad.Close()
	err = gamepad.LeftStickMoveX(1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}
