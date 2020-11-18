package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
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

func TestDialCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateDial("", []byte("DialDevice"))
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestDialCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateDial(path, []byte("DialDevice"))
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestDialCreationFailsOnWrongPathName(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "uinput-dial-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	expected := "failed to register dial input device: failed to close device: inappropriate ioctl for device"
	_, err = CreateDial(file.Name(), []byte("DialDevice"))
	if err == nil || !(expected == err.Error()) {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestDialCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateDial("/dev/uinput", []byte(name))
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}
