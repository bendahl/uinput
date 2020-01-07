package uinput

import (
	"testing"
)

func TestDialWheel(t *testing.T) {
	relDev, err := CreateDial("/dev/uinput", []byte("Test Dial"))
	if err != nil {
		t.Fatalf("Failed to create the virtual dial. Last error was: %s\n", err)
	}

	err = relDev.Turn(1)
	if err != nil {
		t.Fatalf("Failed to perform wheel movement. Last error was: %s\n", err)
	}

	err = relDev.Close()
	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestDialFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateDial("/dev/uinput", []byte("Test Dial"))
	if err != nil {
		t.Fatalf("Failed to create the virtual dial. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.Turn(1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}
