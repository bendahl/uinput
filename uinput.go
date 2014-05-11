/*
Package uinput provides access to the userland input device driver uinput on linux systems.
For now, only the creation of a virtual keyboard is supported. The keycodes, that are available
and can be used to trigger key press events, are part of this package ("KEY_1" for number 1, for
example).

In order to use the virtual keyboard, you will need to follow these three steps:

	1. Initialize the device (don't forget to store the deviceId)
		Example: devId, err := uinput.CreateVKeyboardDevice("/dev/uinput")

	2. Send Button events to the device
		Example: err = uinput.SendBtnEvent(devId, uinput.KEY_D, 1)

	3. Close the device
		Example: err = uinput.CloseDevice(devId)
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

// CreateVKeyboardDevice creates a new uinput device.
// Make sure to pass the correct path to the current system's
// uinput device (usually either "/dev/uinput" or "/dev/input/uinput".
func CreateVKeyboardDevice(path string) (deviceId int, err error) {
	var fd C.int
	var deviceName = C.CString(path)
	defer C.free(unsafe.Pointer(deviceName))

	fd = C.initVKeyboardDevice(deviceName)
	if fd < 0 {
		return 0, errors.New("Could not initialize device")
	}

	return int(fd), nil
}

// SendBtnEvent will send a button event to the newly created device.
// The id refers to the deviceId returned by the create method. The key
// can be any of the predefined keycodes. The btnState can either be 
// pressed (1) or released (0).
func SendBtnEvent(deviceId int, key int, btnState int) (err error) {
	if C.sendBtnEvent(C.int(deviceId), C.int(key), C.int(btnState)) < 0 {
		return errors.New("Sending keypress failed")
	} else {
		return nil
	}
}

// CloseDevice will close the device and free resources. 
// It's usually a good idea to use defer to call this function.
func CloseDevice(deviceId int) (err error) {
	if int(C.releaseDevice(C.int(deviceId))) < 0 {
		return errors.New("Closing device failed")
	} else {
		return nil
	}
}
