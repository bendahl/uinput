package uinput

import "testing"

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
