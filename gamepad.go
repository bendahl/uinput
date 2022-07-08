package uinput

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

const MaximumAxisValue = 32767

// These are the directions and values for hat events.
type HatDirection int

const (
	HatUp int = iota + 1
	HatDown
	HatLeft
	HatRight

	ReleaseHatUp
	ReleaseHatDown
	ReleaseHatLeft
	ReleaseHatRight
)

// A Gamepad is a hybrid key / absolute change event output device.
// Is is used to enable a priogram to simulate gamepad input events.
type Gamepad interface {
	// KeyPress will cause the key to be pressed and immediately released.
	ButtonPress(key int) error

	// ButtonDown will send a keypress event to an existing gamepad device.
	// The key can be any of the predefined keycodes from keycodes.go.
	// Note that the key will be "held down" until "KeyUp" is called.
	ButtonDown(key int) error

	// ButtonUp will send a keyrelease event to an existing gamepad device.
	// The key can be any of the predefined keycodes from keycodes.go.
	ButtonUp(key int) error

	////// The following stick events take in normalized values (-1.0:1.0)
	// These methods will move the left stick's position to the given value.
	LeftStickMoveX(value float32) error
	LeftStickMoveY(value float32) error

	// These methods will move the right stick's position to the given value.
	RightStickMoveX(value float32) error
	RightStickMoveY(value float32) error

	// These methods will move the stick's position to the given X and Y positions.
	LeftStickMove(x, y float32) error
	RightStickMove(x, y float32) error

	// These are alternative methods to send dpad events.
	HatPress(direction HatDirection) error
	HatRelease(direction HatDirection) error

	// Sends out a SYN event.
	syncEvents() error

	io.Closer
}

type vGamepad struct {
	name       []byte
	deviceFile *os.File
}

// CreateGamepad will create a new gamepad using the given uinput
// device path of the uinput device.
func CreateGamepad(path string, name []byte, vendor uint16, product uint16) (Gamepad, error) { // TODO: Consider moving this to a generic function that works for all devices
	err := validateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = validateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createVGamepadDevice(path, name, vendor, product)
	if err != nil {
		return nil, err
	}

	return vGamepad{name: name, deviceFile: fd}, nil
}

func (vg vGamepad) ButtonPress(key int) error {
	err := vg.ButtonDown(key)
	if err != nil {
		return err
	}
	err = vg.ButtonUp(key)
	if err != nil {
		return err
	}
	return nil
}

func (vg vGamepad) ButtonDown(key int) error {
	return sendBtnEvent(vg.deviceFile, []int{key}, btnStatePressed)
}

func (vg vGamepad) ButtonUp(key int) error {
	return sendBtnEvent(vg.deviceFile, []int{key}, btnStateReleased)
}

func (vg vGamepad) LeftStickMoveX(value float32) error {
	return vg.sendStickAxisEvent(absX, value)
}

func (vg vGamepad) LeftStickMoveY(value float32) error {
	return vg.sendStickAxisEvent(absY, value)
}

func (vg vGamepad) RightStickMoveX(value float32) error {
	return vg.sendStickAxisEvent(absRX, value)
}

func (vg vGamepad) RightStickMoveY(value float32) error {
	return vg.sendStickAxisEvent(absRY, value)
}

func (vg vGamepad) RightStickMove(x, y float32) error {
	values := map[uint16]float32{}
	values[absRX] = x
	values[absRY] = y

	return vg.sendStickEvent(values)
}

func (vg vGamepad) LeftStickMove(x, y float32) error {
	values := map[uint16]float32{}
	values[absX] = x
	values[absY] = y

	return vg.sendStickEvent(values)
}

func (vg vGamepad) HatPress(direction HatDirection) error {
	return vg.sendHatEvent(direction)
}

func (vg vGamepad) HatRelease(direction HatDirection) error {
	return vg.sendHatEvent(direction)
}

func (vg vGamepad) syncEvents() error {
	buf, err := inputEventToBuffer(inputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  evSyn,
		Code:  uint16(synReport),
		Value: 0})
	if err != nil {
		return fmt.Errorf("writing sync event failed: %v", err)
	}
	_, err = vg.deviceFile.Write(buf)
	return err
}

func (vg vGamepad) sendStickAxisEvent(absCode uint16, value float32) error {
	ev := inputEvent{
		Type:  evAbs,
		Code:  absCode,
		Value: denormalizeInput(value),
	}

	buf, err := inputEventToBuffer(ev)
	if err != nil {
		return fmt.Errorf("writing abs stick event failed: %v", err)
	}

	_, err = vg.deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write abs stick event to device file: %v", err)
	}

	return vg.syncEvents()
}

func (vg vGamepad) sendStickEvent(values map[uint16]float32) error {
	for code, value := range values {
		ev := inputEvent{
			Type:  evAbs,
			Code:  code,
			Value: denormalizeInput(value),
		}

		buf, err := inputEventToBuffer(ev)
		if err != nil {
			return fmt.Errorf("writing abs stick event failed: %v", err)
		}

		_, err = vg.deviceFile.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write abs stick event to device file: %v", err)
		}
	}

	return vg.syncEvents()
}

func (vg vGamepad) sendHatEvent(direction HatDirection) error {
	var event uint16
	var value int32

	// TODO: This is a -questionable- (terrible) way to handle this, I'm open to other ideas.
	switch int(direction) {
	case HatUp:
		{
			event = absHat0Y
			value = -1
		}
	case HatDown:
		{
			event = absHat0Y
			value = 1
		}
	case HatLeft:
		{
			event = absHat0Y
			value = -1
		}
	case HatRight:
		{
			event = absHat0Y
			value = 1
		}
	case ReleaseHatUp:
		{
			event = absHat0Y
			value = 0
		}
	case ReleaseHatDown:
		{
			event = absHat0Y
			value = 0
		}
	case ReleaseHatLeft:
		{
			event = absHat0Y
			value = 0
		}
	case ReleaseHatRight:
		{
			event = absHat0Y
			value = 0
		}
	default:
		{
			return errors.New("Failed to parse input direction")
		}
	}

	ev := inputEvent{
		Type:  evAbs,
		Code:  event,
		Value: value,
	}

	buf, err := inputEventToBuffer(ev)
	if err != nil {
		return fmt.Errorf("writing abs stick event failed: %v", err)
	}

	_, err = vg.deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write abs stick event to device file: %v", err)
	}

	return vg.syncEvents()
}

func (vg vGamepad) Close() error {
	return closeDevice(vg.deviceFile)
}

func createVGamepadDevice(path string, name []byte, vendor uint16, product uint16) (fd *os.File, err error) {
	// This array is needed to register the event keys for the gamepad device.
	keys := []uint16{
		ButtonGamepad,

		ButtonSouth,
		ButtonEast,
		ButtonNorth,
		ButtonWest,

		ButtonBumperLeft,
		ButtonBumperRight,
		ButtonTriggerLeft,
		ButtonTriggerRight,
		ButtonThumbLeft,
		ButtonThumbRight,

		ButtonSelect,
		ButtonStart,

		ButtonDpadUp,    // * * *
		ButtonDpadDown,  // * These buttons can be used instead of the hat events.
		ButtonDpadLeft,  // *
		ButtonDpadRight, // * * *

		ButtonMode,
	}

	// This array is for the absolute events for the gamepad device.
	abs_events := []uint16{
		absX,
		absY,
		absZ,
		absRX,
		absRY,
		absRZ,
		absHat0X,
		absHat0Y,
	}

	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual gamepad device: %v", err)
	}

	// register button events
	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register virtual gamepad device: %v", err)
	}

	for _, code := range keys {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(code))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register key number %d: %v", code, err)
		}
	}

	// register absolute events
	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute event input device: %v", err)
	}

	for _, event := range abs_events {
		err = ioctl(deviceFile, uiSetAbsBit, uintptr(event))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register absolute event %v: %v", event, err)
		}
	}

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  vendor,
				Product: product,
				Version: 1}})
}

// Takes in a normalized value (-1.0:1.0) and return an event value
func denormalizeInput(value float32) int32 {
	return int32(value * MaximumAxisValue)
}
