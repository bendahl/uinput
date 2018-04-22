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
	"fmt"
	"os"
)

func validateDevicePath(path string) {
	if path == "" {
		panic("device path must not be empty")
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("device path '%s' does not exist", path))
	}
}

func validateUinputName(name []byte) {
	if len(name) > uinputMaxNameSize {
		panic(fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize))
	}
}
