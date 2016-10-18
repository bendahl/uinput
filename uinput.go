/*
Package uinput provides access to the userland input device driver uinput on linux systems.
For now, only the creation of a virtual keyboard is supported. The keycodes, that are available
and can be used to trigger key press events, are part of this package ("KEY_1" for number 1, for
example).

In order to use the virtual keyboard, you will need to follow these three steps:

	1. Initialize the device
		Example: vk, err := CreateKeyboard("/dev/uinput", "Virtual Keyboard")

	2. Send Button events to the device
		Example: err = vk.SendKeyPress(uinput.KEY_D)
				 err = vk.SendKeyRelease(uinput.KEY_D)

	3. Close the device
		Example: err = vk.Close()
*/
package uinput

/*
#include "uinputwrapper.h"
*/
import "C"
import (
	"errors"
	"io"
	"unsafe"
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
	name string
	fd   int
}

func CreateKeyboard(path, name string) (Keyboard, error) {
	fd, err := createVKeyboardDevice(path, name)
	if err != nil {
		return nil, err
	}

	return vKeyboard{name, fd}, nil
}

func (vk vKeyboard) SendKeyPress(key int) error {
	return sendBtnEvent(vk.fd, key, 1)
}

func (vk vKeyboard) SendKeyRelease(key int) error {
	return sendBtnEvent(vk.fd, key, 0)
}

// Close will close the device and free resources.
// It's usually a good idea to use defer to call this function.
func (vk vKeyboard) Close() error {
	return closeDevice(vk.fd)
}

func createVKeyboardDevice(path, name string) (deviceId int, err error) {
	uinputDevice := C.CString(path)
	defer C.free(unsafe.Pointer(uinputDevice))

	if name == "" {
		name = "uinput_default_vkeyboard"
	}
	virtDeviceName := C.CString(name)
	defer C.free(unsafe.Pointer(virtDeviceName))

	var fd C.int
	fd = C.initVKeyboardDevice(uinputDevice, virtDeviceName)
	if fd < 0 {
		// TODO: Map ErrValues into more specific Errors
		return 0, errors.New("Could not initialize device.")
	}

	return int(fd), nil
}

func sendBtnEvent(deviceId int, key int, btnState int) (err error) {
	if C.sendBtnEvent(C.int(deviceId), C.int(key), C.int(btnState)) < 0 {
		return errors.New("Sending keypress failed.")
	} else {
		return nil
	}
}

func closeDevice(deviceId int) (err error) {
	if int(C.releaseDevice(C.int(deviceId))) < 0 {
		return errors.New("Closing device failed.")
	} else {
		return nil
	}
}
