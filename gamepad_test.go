package uinput

import (
	"testing"
	"time"
	"uinput"
)

// This test inputs the konami code
func TestInfiniteKonami(t *testing.T) {
	vg, err := uinput.CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 3; j++ {
			err = vg.ButtonPress(uinput.ButtonDpadUp)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

			err = vg.ButtonPress(uinput.ButtonDpadDown)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

		}

		for j := 0; j < 3; j++ {
			err = vg.ButtonPress(uinput.ButtonDpadLeft)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

			err = vg.ButtonPress(uinput.ButtonDpadRight)
			if err != nil {
				t.Fatalf("Failed to send button press. Last error was: %s\n", err)
			}

		}

		err = vg.ButtonPress(uinput.ButtonSouth)
		if err != nil {
			t.Fatalf("Failed to send button press. Last error was: %s\n", err)
		}

		err = vg.ButtonPress(uinput.ButtonEast)
		if err != nil {
			t.Fatalf("Failed to send button press. Last error was: %s\n", err)
		}

		err = vg.ButtonPress(uinput.ButtonStart)
		if err != nil {
			t.Fatalf("Failed to send button press. Last error was: %s\n", err)
		}
	}
}

// This test moves the axes around a bit
func TestAxisMovement(t *testing.T) {
	vg, err := uinput.CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}

	err = vg.LeftStickMove(0.2, 1.0)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}

	err = vg.RightStickMove(0.2, 1.0)
	if err != nil {
		t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
	}
}
