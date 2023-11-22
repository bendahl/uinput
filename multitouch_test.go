package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestBasicMultiTouchMoves(t *testing.T) {
	absDev, err := CreateMultiTouch("/dev/uinput", []byte("Test MultiTouch"), 0, 1024, 0, 768, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual multi touch device. Last error was: %s\n", err)
	}
	defer func(absDev MultiTouch) {
		err := absDev.Close()
		if err != nil {
			t.Fatalf("Failed to close device. Last error was: %s\n", err)
		}
	}(absDev)

	contacts := absDev.GetContacts()

	if len(contacts) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contacts)))
	}

	err = contacts[0].TouchDownAt(0, 0)
	if err != nil {
		t.Fatalf("Failed to move contact 0 to initial position. Last error was: %s\n", err)
	}

	err = contacts[0].TouchDownAt(100, 200)
	if err != nil {
		t.Fatalf("Failed to move contact 0 to position x:100, y:200. Last error was: %s\n", err)
	}
}

func TestBasicMultiTouchGesture(t *testing.T) {
	absDev, err := CreateMultiTouch("/dev/uinput", []byte("Test MultiTouch"), 0, 1024, 0, 768, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual multi touch device. Last error was: %s\n", err)
	}
	defer func(absDev MultiTouch) {
		err := absDev.Close()
		if err != nil {
			t.Fatalf("Failed to close device. Last error was: %s\n", err)
		}
	}(absDev)

	contacts := absDev.GetContacts()

	if len(contacts) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contacts)))
	}

	for i := int32(0); i < 200; i++ {
		time.Sleep(3 * time.Millisecond)
		for n := int32(0); n < 3; n++ {
			y := 255 - i
			err := contacts[n].TouchDownAt(0, y)
			if err != nil {
				t.Fatalf("Failed to move contact %s at [0, %s]", fmt.Sprint(i), fmt.Sprint(y))
			}
		}
	}
}

func TestMultiTouchCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateMultiTouch("", []byte("TouchDevice"), 0, 1024, 0, 768, 3)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMultiTouchCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateMultiTouch(path, []byte("TouchDevice"), 0, 1024, 0, 768, 3)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestMultiTouchCreationFailsOnWrongPathName(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "uinput-MultiTouch-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	expected := "failed to register key device: failed to close device: inappropriate ioctl for device"
	_, err = CreateMultiTouch(file.Name(), []byte("TouchDevice"), 0, 1024, 0, 768, 3)
	if err == nil || !(expected == err.Error()) {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMultiTouchCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateMultiTouch("/dev/uinput", []byte(name), 0, 1024, 0, 768, 3)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMultiTouchMoveToFailsOnClosedDevice(t *testing.T) {
	absDev, err := CreateMultiTouch("/dev/uinput", []byte("Test MultiTouch"), 0, 1024, 0, 768, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}

	_ = absDev.Close()

	contacts := absDev.GetContacts()

	if len(contacts) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contacts)))
	}

	err = contacts[0].TouchDownAt(1, 1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMultipleMultiTouchsWithDifferentSizes(t *testing.T) {
	horizontal, err := CreateMultiTouch("/dev/uinput", []byte("horizontal_pad"), 0, 200, 0, 100, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer horizontal.Close()
	vertical, err := CreateMultiTouch("/dev/uinput", []byte("vertical_pad"), 0, 100, 0, 200, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer vertical.Close()

	contactsHorizontal := horizontal.GetContacts()

	if len(contactsHorizontal) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contactsHorizontal)))
	}

	err = contactsHorizontal[0].TouchDownAt(200, 100)
	if err != nil {
		t.Fatalf("Unable to move cursor on horizontal pad: %v", err)
	}

	contactsVertical := vertical.GetContacts()

	if len(contactsVertical) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contactsVertical)))
	}

	err = contactsHorizontal[0].TouchDownAt(100, 200)
	if err != nil {
		t.Fatalf("Unable to move cursor on horizontal pad: %v", err)
	}
}

func TestMultiTouchPositioningInUpperLeftCorner(t *testing.T) {
	dev, err := CreateMultiTouch("/dev/uinput", []byte("MultiTouch"), 0, 200, 0, 100, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer dev.Close()

	contacts := dev.GetContacts()

	if len(contacts) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contacts)))
	}

	err = contacts[0].TouchDownAt(0, 0)

	if err != nil {
		t.Fatalf("Failed to move cursor to upper left corner: %v", err)
	}
}

func TestMultiTouchSingleTouchEvent(t *testing.T) {
	dev, err := CreateMultiTouch("/dev/uinput", []byte("MultiTouch"), 0, 200, 0, 100, 3)
	if err != nil {
		t.Fatalf("Failed to create the virtual touch pad. Last error was: %s\n", err)
	}
	defer dev.Close()

	contacts := dev.GetContacts()

	if len(contacts) != 3 {
		t.Fatalf("Failed to create contacts, expected 3, got %s", fmt.Sprint(len(contacts)))
	}

	err = contacts[0].TouchDownAt(0, 0)
	if err != nil {
		t.Fatalf("Failed to issue touch down event at [0,0]: %v", err)
	}

	err = contacts[0].TouchUp()
	if err != nil {
		t.Fatalf("Failed to issue touch up event [0,0]: %v", err)
	}

}

func TestMultiTouchSyspath(t *testing.T) {
	dev, err := CreateMultiTouch("/dev/uinput", []byte("MultiTouch"), 0, 1024, 0, 768, 3)
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
