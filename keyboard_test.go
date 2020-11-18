package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// This test will confirm that basic key events are working.
// Note that only Key1 is used here, as the purpose of this test is to ensure that the event handling for
// keyboard devices is working. All other keys, defined in keycodes.go should work as well if this test passes.
// Another thing to keep in mind is that there are certain key codes that might not be great candidates for
// unit testing, as they may create unwanted side effects, like logging out the current user, etc...
func TestKeysInValidRangeWork(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.KeyPress(keyReserved)
	if err != nil {
		t.Fatalf("Failed to send key press. Last error was: %s\n", err)
	}

	err = vk.KeyDown(keyReserved)
	if err != nil {
		t.Fatalf("Failed to send key down event. Last error was: %s\n", err)
	}

	err = vk.KeyUp(keyReserved)
	if err != nil {
		t.Fatalf("Failed to send key up event. Last error was: %s\n", err)
	}

	err = vk.KeyPress(keyMax)
	if err != nil {
		t.Fatalf("Failed to send key press. Last error was: %s\n", err)
	}

	err = vk.KeyDown(keyMax)
	if err != nil {
		t.Fatalf("Failed to send key down event. Last error was: %s\n", err)
	}

	err = vk.KeyUp(keyMax)
	if err != nil {
		t.Fatalf("Failed to send key up event. Last error was: %s\n", err)
	}

	err = vk.Close()

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestKeyboardCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateKeyboard("", []byte("KeyboardDevice"))
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestKeyboardCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateKeyboard(path, []byte("KeyboardDevice"))
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestKeyboardCreationFailsOnWrongPathName(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "uinput-keyboard-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	expected := "failed to register virtual keyboard device: failed to close device: inappropriate ioctl for device"
	_, err = CreateKeyboard(file.Name(), []byte("DialDevice"))
	if err == nil || !(expected == err.Error()) {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestKeyboardCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateKeyboard("/dev/uinput", []byte(name))
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestKeyOutsideOfRangeKeyPressFails(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	defer vk.Close()

	err = vk.KeyPress(keyMax + 1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

	err = vk.KeyPress(-1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

}
func TestKeyOutsideOfRangeKeyUpFails(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	defer vk.Close()

	err = vk.KeyUp(keyMax + 1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

	err = vk.KeyUp(-1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

}

func TestKeyOutsideOfRangeKeyDownFails(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	defer vk.Close()

	err = vk.KeyDown(keyMax + 1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

	err = vk.KeyDown(-1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

}

func TestKeyPressFailsIfDeviceIsClosed(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	vk.Close()

	err = vk.KeyPress(Key1)
	if err == nil {
		t.Fatalf("Expected KeyPress to fail, but no error was returned.")
	}
}

func TestKeyUpFailsIfDeviceIsClosed(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	vk.Close()

	err = vk.KeyUp(Key1)
	if err == nil {
		t.Fatalf("Expected KeyPress to fail, but no error was returned.")
	}
}

func TestKeyDownFailsIfDeviceIsClosed(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	vk.Close()

	err = vk.KeyDown(Key1)
	if err == nil {
		t.Fatalf("Expected KeyPress to fail, but no error was returned.")
	}
}
