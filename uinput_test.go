package uinput

import (
	"os"
	"testing"
)

func TestValidateDevicePathEmptyPathPanics(t *testing.T) {
	expected := "device path must not be empty"
	err := validateDevicePath("")
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestValidateDevicePathInvalidPathPanics(t *testing.T) {
	path := "/some/bogus/path"
	err := validateDevicePath(path)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestValidateUinputNameEmptyNamePanics(t *testing.T) {
	expected := "device name may not be empty"
	err := validateUinputName(nil)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}
