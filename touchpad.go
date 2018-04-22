package uinput

import (
	"fmt"
	"io"
	"os"
)

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

	// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
	// LeftRelease is invoked.
	LeftPress() error

	// LeftRelease will simulate the release of the left mouse button.
	LeftRelease() error

	// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
	// RightRelease is invoked.
	RightPress() error

	// RightRelease will simulate the release of the right mouse button.
	RightRelease() error

	io.Closer
}

type vTouchPad struct {
	name       []byte
	deviceFile *os.File
}

// CreateTouchPad will create a new touch pad device. note that you will need to define the x and y axis boundaries
// (min and max) within which the cursor maybe moved around.
func CreateTouchPad(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (TouchPad, error) {
	validateDevicePath(path)
	validateUinputName(name)

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

// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
// LeftRelease is invoked.
func (vTouch vTouchPad) LeftPress() error {
	err := sendBtnEvent(vTouch.deviceFile, evBtnLeft, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed press the left mouse button: %v", err)
	}
	err = syncEvents(vTouch.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// LeftRelease will simulate the release of the left mouse button.
func (vTouch vTouchPad) LeftRelease() error {
	err := sendBtnEvent(vTouch.deviceFile, evBtnLeft, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to release the left mouse button: %v", err)
	}
	err = syncEvents(vTouch.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
// RightRelease is invoked.
func (vTouch vTouchPad) RightPress() error {
	err := sendBtnEvent(vTouch.deviceFile, evBtnRight, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to press the right mouse button: %v", err)
	}
	err = syncEvents(vTouch.deviceFile)
	if err != nil {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

// RightRelease will simulate the release of the right mouse button.
func (vTouch vTouchPad) RightRelease() error {
	err := sendBtnEvent(vTouch.deviceFile, evBtnRight, btnStateReleased)
	if err != nil {
		return fmt.Errorf("Failed to release the right mouse button: %v", err)
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
