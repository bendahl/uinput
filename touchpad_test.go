package uinput

import "testing"

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
