/*
Package uinput is a pure go package that provides access to the userland input device driver uinput on linux systems.
Virtual keyboard devices as well as virtual mouse input devices may be created using this package.
The keycodes and other event definitions, that are available and can be used to trigger input events,
are part of this package ("Key1" for number 1, for example).

In order to use the virtual keyboard, you will need to follow these three steps:

	1. Initialize the device
		Example: vk, err := CreateKeyboard("/dev/uinput", "Virtual Keyboard")

	2. Send Button events to the device
		Example (print a single D):
			err = vk.KeyPress(uinput.KeyD)

		Example (keep moving right by holding down right arrow key):
				 err = vk.KeyDown(uinput.KeyRight)

		Example (stop moving right by releasing the right arrow key):
				 err = vk.KeyUp(uinput.KeyRight)

	3. Close the device
		Example: err = vk.Close()

A virtual mouse input device is just as easy to create and use:

	1. Initialize the device:
		Example: vm, err := CreateMouse("/dev/uinput", "DangerMouse")

	2. Move the cursor around and issue click events
		Example (move mouse right):
			err = vm.MoveRight(42)

		Example (move mouse left):
			err = vm.MoveLeft(42)

		Example (move mouse up):
			err = vm.MoveUp(42)

		Example (move mouse down):
			err = vm.MoveDown(42)

		Example (trigger a left click):
			err = vm.LeftClick()

		Example (trigger a right click):
			err = vm.RightClick()

	3. Close the device
		Example: err = vm.Close()


If you'd like to use absolute input events (move the cursor to specific positions on screen), use the touch pad.
Note that you'll need to specify the size of the screen area you want to use when you initialize the
device. Here are a few examples of how to use the virtual touch pad:

	1. Initialize the device:
		Example: vt, err := CreateTouchPad("/dev/uinput", "DontTouchThis", 0, 1024, 0, 768)

	2. Move the cursor around and issue click events
		Example (move cursor to the top left corner of the screen):
			err = vt.MoveTo(0, 0)

		Example (move cursor to the position x: 100, y: 250):
			err = vt.MoveTo(100, 250)

		Example (trigger a left click):
			err = vt.LeftClick()

		Example (trigger a right click):
			err = vt.RightClick()

	3. Close the device
		Example: err = vt.Close()

*/
package uinput

import (
	"errors"
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

// A Mouse is a device that will trigger an absolute change event.
// For details see: https://www.kernel.org/doc/Documentation/input/event-codes.txt
type Mouse interface {
	// MoveLeft will move the mouse cursor left by the given number of pixel.
	MoveLeft(pixel int32) error

	// MoveRight will move the mouse cursor right by the given number of pixel.
	MoveRight(pixel int32) error

	// MoveUp will move the mouse cursor up by the given number of pixel.
	MoveUp(pixel int32) error

	// MoveDown will move the mouse cursor down by the given number of pixel.
	MoveDown(pixel int32) error

	// LeftClick will issue a single left click.
	LeftClick() error

	// RightClick will issue a right click.
	RightClick() error

	io.Closer
}

type vMouse struct {
	name       []byte
	deviceFile *os.File
}

// A TouchPad is an input device that uses absolute axis events, meaning that you can specify
// the exact position the cursor should move to. Therefore, it is necessary to define the size
// of the rectangle in which the cursor may move upon creation of the device.
type TouchPad interface {
	// MoveTo will move the cursor to the specified position on the screen
	MoveTo(x int32, y int32) error

	// LeftClick will issue a single left click.
	LeftClick() error

	// RightClick will issue a right click.
	RightClick() error

	io.Closer
}

type vTouchPad struct {
	name       []byte
	deviceFile *os.File
}

// CreateTouchPad will create a new touch pad device. note that you will need to define the x and y axis boundaries
// (min and max) within which the cursor maybe moved around.
func CreateTouchPad(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (TouchPad, error) {
	if path == "" {
		return nil, errors.New("device path must not be empty")
	}
	if len(name) > uinputMaxNameSize {
		return nil, fmt.Errorf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	}

	fd, err := createTouchPad(path, name, minX, maxX, minY, maxY)
	if err != nil {
		return nil, err
	}

	return vTouchPad{name: name, deviceFile: fd}, nil
}

func (vTouch vTouchPad) MoveTo(x int32, y int32) error {
	return sendAbsEvent(vTouch.deviceFile, x, y)
}

func (vTouch vTouchPad) LeftClick() error {
	err := sendBtnEvent(vTouch.deviceFile, evBtnLeft, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the LeftClick event: %v", err)
	}

	err = sendBtnEvent(vTouch.deviceFile, evBtnLeft, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyUp event: %v", err)
	}

	err = syncEvents(vTouch.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

func (vTouch vTouchPad) RightClick() error {
	err := sendBtnEvent(vTouch.deviceFile, evBtnRight, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the RightClick event: %v", err)
	}

	err = sendBtnEvent(vTouch.deviceFile, evBtnRight, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyUp event: %v", err)
	}

	err = syncEvents(vTouch.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

func (vTouch vTouchPad) Close() error {
	return closeDevice(vTouch.deviceFile)
}

// CreateMouse will create a new mouse input device. A mouse is a device that allows relative input.
// Relative input means that all changes to the x and y coordinates of the mouse pointer will be
func CreateMouse(path string, name []byte) (Mouse, error) {
	if path == "" {
		return nil, errors.New("device path must not be empty")
	}
	if len(name) > uinputMaxNameSize {
		return nil, fmt.Errorf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	}

	fd, err := createMouse(path, name)
	if err != nil {
		return nil, err
	}

	return vMouse{name: name, deviceFile: fd}, nil
}

// MoveLeft will move the cursor left by the number of pixel specified.
func (vRel vMouse) MoveLeft(pixel int32) error {
	return sendRelEvent(vRel.deviceFile, relX, -pixel)
}

// MoveRight will move the cursor right by the number of pixel specified.
func (vRel vMouse) MoveRight(pixel int32) error {
	return sendRelEvent(vRel.deviceFile, relX, pixel)
}

// MoveUp will move the cursor up by the number of pixel specified.
func (vRel vMouse) MoveUp(pixel int32) error {
	return sendRelEvent(vRel.deviceFile, relY, -pixel)
}

// MoveDown will move the cursor down by the number of pixel specified.
func (vRel vMouse) MoveDown(pixel int32) error {
	return sendRelEvent(vRel.deviceFile, relY, pixel)
}

// LeftClick will issue a LeftClick.
func (vRel vMouse) LeftClick() error {
	err := sendBtnEvent(vRel.deviceFile, evBtnLeft, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the LeftClick event: %v", err)
	}

	err = sendBtnEvent(vRel.deviceFile, evBtnLeft, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyUp event: %v", err)
	}

	err = syncEvents(vRel.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// RightClick will issue a RightClick
func (vRel vMouse) RightClick() error {
	err := sendBtnEvent(vRel.deviceFile, evBtnRight, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the RightClick event: %v", err)
	}

	err = sendBtnEvent(vRel.deviceFile, evBtnRight, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to issue the KeyUp event: %v", err)
	}

	err = syncEvents(vRel.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// Close closes the device and releases the device.
func (vRel vMouse) Close() error {
	return closeDevice(vRel.deviceFile)
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
