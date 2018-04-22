package uinput

import (
	"fmt"
	"io"
	"os"
)

// A Keyboard is an key event output device. It is used to
// enable a program to simulate HID keyboard input events.
type Keyboard interface {
	// KeyPress will cause the key to be pressed and immediately released.
	KeyPress(key int) error

	// KeyDown will send a keypress event to an existing keyboard device.
	// The key can be any of the predefined keycodes from uinputdefs.
	// Note that the key will be "held down" until "KeyUp" is called.
	KeyDown(key int) error

	// KeyUp will send a keyrelease event to an existing keyboard device.
	// The key can be any of the predefined keycodes from uinputdefs.
	KeyUp(key int) error

	io.Closer
}

type vKeyboard struct {
	name       []byte
	deviceFile *os.File
}

// CreateKeyboard will create a new keyboard using the given uinput
// device path of the uinput device.
func CreateKeyboard(path string, name []byte) (Keyboard, error) {
	validateDevicePath(path)
	validateUinputName(name)

	fd, err := createVKeyboardDevice(path, name)
	if err != nil {
		return nil, err
	}

	return vKeyboard{name: name, deviceFile: fd}, nil
}

// KeyPress will issue a single key press (push down a key and then immediately release it).
func (vk vKeyboard) KeyPress(key int) error {
	err := sendBtnEvent(vk.deviceFile, key, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyDown event: %v", err)
	}

	err = sendBtnEvent(vk.deviceFile, key, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyUp event: %v", err)
	}

	err = syncEvents(vk.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// KeyDown will send the key code passed (see uinputdefs.go for available keycodes). Note that unless a key release
// event is sent to the device, the key will remain pressed and therefore input will continuously be generated. Therefore,
// do not forget to call "KeyUp" afterwards.
func (vk vKeyboard) KeyDown(key int) error {
	err := sendBtnEvent(vk.deviceFile, key, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyDown event: %v", err)
	}

	err = syncEvents(vk.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// KeyUp will release the given key passed as a parameter (see uinputdefs.go for available keycodes). In most
// cases it is recommended to call this function immediately after the "KeyDown" function in order to only issue a
// single key press.
func (vk vKeyboard) KeyUp(key int) error {
	err := sendBtnEvent(vk.deviceFile, key, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyUp event: %v", err)
	}

	err = syncEvents(vk.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// Close will close the device and free resources.
// It's usually a good idea to use defer to call this function.
func (vk vKeyboard) Close() error {
	return closeDevice(vk.deviceFile)
}
