/*
Package uinput provides access to the userland input device driver uinput on linux systems.
For now, only the creation of a virtual keyboard is supported. The keycodes, that are available
and can be used to trigger key press events, are part of this package ("KEY_1" for number 1, for
example).

In order to use the virtual keyboard, you will need to follow these three steps:

	1. Initialize the device
		Example: vk := VKeyboard{}
	             err := vk.Create("/dev/uinput")


	2. Send Button events to the device
		Example: err = vk.SendKeyPress(uinput.KEY_D)

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
	"unsafe"
)

// VKeyboard represents a virtual keyboard device. There are several
// methods available to work with this virtual device. Devices can be
// created, receive events, and closed.
type VKeyboard struct {
	// The Name of the uinput device. Will be trimmed to a Max Length of 80 bytes.
	// If left blank the device will have a default name.
	Name string

	id int
}

// Create creates a new virtual keyboard device.
// Make sure to pass the correct path to the current system's
// uinput device (usually either "/dev/uinput" or "/dev/input/uinput".
func (vk *VKeyboard) Create(path string) (err error) {
	vk.id = -1
	var ret error

	vk.id, ret = createVKeyboardDevice(path, vk.Name)
	return ret
}

// SendKeyPress will send a keypress event to an existing keyboard device.
// The key can be any of the predefined keycodes from uinputdefs.
func (vk *VKeyboard) SendKeyPress(key int) (err error) {
	if vk.id < 0 {
		return errors.New("Keyboard not initialized. Sending keypress event failed.")
	}

	return sendBtnEvent(vk.id, key, 1)
}

// SendKeyRelease will send a keyrelease event to an existing keyboard device.
// The key can be any of the predefined keycodes from uinputdefs.
func (vk *VKeyboard) SendKeyRelease(key int) (err error) {
	if vk.id < 0 {
		return errors.New("Keyboard not initialized. Sending keyrelease event failed.")
	}

	return sendBtnEvent(vk.id, key, 0)
}

// Close will close the device and free resources.
// It's usually a good idea to use defer to call this function.
func (vk *VKeyboard) Close() (err error) {
	if vk.id < 0 {
		return errors.New("Keyboard not initialized. Closing device failed.")
	}
	return closeDevice(vk.id)
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
