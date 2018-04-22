package uinput

import "testing"

func TestValidateDevicePathEmptyPathPanics(t *testing.T) {
	expected := "device path must not be empty"
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual )
			}
		}
	}()
	validateDevicePath("")
	t.Fatalf("Empty path did not yield a panic")
}

func TestValidateDevicePathInvalidPathPanics(t *testing.T) {
	path := "/some/bogus/path"
	expected := "device path '" + path + "' does not exist"
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual )
			}
		}
	}()
	validateDevicePath(path)
	t.Fatalf("Invalid path did not yield a panic")
}
