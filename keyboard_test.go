package uinput

import "testing"

// This test will confirm that basic key events are working.
// Note that only Key1 is used here, as the purpose of this test is to ensure that the event handling for
// keyboard devices is working. All other keys, defined in uinputdefs should work as well if this test passes.
// Another thing to keep in mind is that there are certain key codes that might not be great candidates for
// unit testing, as they may create unwanted side effects, like logging out the current user, etc...
func TestBasicKeyboard(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.KeyPress(Key1)
	if err != nil {
		t.Fatalf("Failed to send key press. Last error was: %s\n", err)
	}

	err = vk.KeyDown(Key1)
	if err != nil {
		t.Fatalf("Failed to send key down event. Last error was: %s\n", err)
	}

	err = vk.KeyUp(Key1)
	if err != nil {
		t.Fatalf("Failed to send key up event. Last error was: %s\n", err)
	}

	err = vk.Close()

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}
