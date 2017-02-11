/*
Package uinput provides access to the userland input device driver uinput on linux systems.
For now, only the creation of a virtual keyboard is supported. The keycodes, that are available
and can be used to trigger key press events, are part of this package ("Key1" for number 1, for
example).

In order to use the virtual keyboard, you will need to follow these three steps:

	1. Initialize the device
		Example: vk, err := CreateKeyboard("/dev/uinput", "Virtual Keyboard")

	2. Send Button events to the device
		Example: err = vk.SendKeyPress(uinput.KeyD)
				 err = vk.SendKeyRelease(uinput.KeyD)

	3. Close the device
		Example: err = vk.Close()
*/
package uinput

import (
	"io"
	"os"
	"errors"
	"fmt"
)

// A Keyboard is an key event output device. It is used to
// enable a program to simulate HID keyboard input events.
type Keyboard interface {
	// SendKeyPress will send a keypress event to an existing keyboard device.
	// The key can be any of the predefined keycodes from uinputdefs.
	SendKeyPress(key int) error

	// SendKeyRelease will send a keyrelease event to an existing keyboard device.
	// The key can be any of the predefined keycodes from uinputdefs.
	SendKeyRelease(key int) error

	io.Closer
}

type vKeyboard struct {
	name       []byte
	deviceFile *os.File
}

// CreateKeyboard will create a new keyboard using the given uinput
// device path of the uinput device.
func CreateKeyboard(path string, name []byte) (Keyboard, error) {
	if path == "" {
		return nil, errors.New("device path must not be empty")
	}
	if len(name) > uinputMaxNameSize {
		return nil, fmt.Errorf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	}

	fd, err := createVKeyboardDevice(path, name)
	if err != nil {
		return nil, err
	}

	return vKeyboard{name: name, deviceFile: fd}, nil
}

// SendKeyPress will send the key code passed (see uinputdefs.go for available keycodes). Note that unless a key release
// event is sent to the device, the key will remain pressed and therefore input will continuously be generated. Therefore,
// do not forget to call "SendKeyRelease" afterwards.
func (vk vKeyboard) SendKeyPress(key int) error {
	return sendBtnEvent(vk.deviceFile, key, btnStatePressed)
}

// SendKeyRelease will release the given key passed as a parameter (see uinputdefs.go for available keycodes). In most
// cases it is recommended to call this function immediately after the "SendKeyPress" function in order to only issue a
// singel key press.
func (vk vKeyboard) SendKeyRelease(key int) error {
	return sendBtnEvent(vk.deviceFile, key, btnStateReleased)
}

// Close will close the device and free resources.
// It's usually a good idea to use defer to call this function.
func (vk vKeyboard) Close() error {
	return closeDevice(vk.deviceFile)
}
