package uinput

import (
	"os"
	"fmt"
	"syscall"
	"errors"
	"bytes"
	"encoding/binary"
	"time"
)

// types needed from uinput.h
const (
	uinputMaxNameSize = 80
	uiDevCreate       = 0x5501
	uiDevDestroy      = 0x5502
	uiSetEvBit        = 0x40045564
	uiSetKeyBit       = 0x40045565
	busUsb            = 0x03
)

// input event codes as specified in input-event-codes.h
const (
	evSyn     = 0x00
	evKey     = 0x01
	evRel     = 0x02
	evAbs     = 0x03
	synReport = 0
)

const (
	btnStateReleased = 0
	btnStatePressed  = 1
)

type inputID struct {
	Bustype uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

// translated to go from uinput.h
type uinputUserDev struct {
	Name       [uinputMaxNameSize]byte
	ID         inputID
	EffectsMax uint32
	Absmax     [64]int32
	Absmin     [64]int32
	Absfuzz    [64]int32
	Absflat    [64]int32
}

// translated to go from input.h
type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func closeDevice(deviceID *os.File) (err error) {
	err = releaseDevice(deviceID)
	if err != nil {
		return fmt.Errorf("failed to close device: %v", err)
	}
	return nil
}

func initKeyboardDevice(path string, name []byte) (deviceFile *os.File, err error) {
	deviceFile, err = os.OpenFile(path, syscall.O_WRONLY|syscall.O_NONBLOCK|syscall.O_NDELAY, 0666)
	if err != nil {
		return nil, errors.New("could not open device file")
	}

	// register device
	err = ioctl(deviceFile, uiSetEvBit, uintptr(evKey))
	if err != nil {
		err = releaseDevice(deviceFile)
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to close device: %v", err)
		}
		deviceFile.Close()
		return nil, fmt.Errorf("invalid file handle returned from ioctl: %v", err)
	}

	// register key events
	for i := 0; i < keyMax; i++ {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(i))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register key number %d: %v", i, err)
		}
	}
	var fixedSizeName [uinputMaxNameSize]byte
	copy(fixedSizeName[:], name)

	uidev :=
		uinputUserDev{
			Name: fixedSizeName,
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x4711,
				Product: 0x0815,
				Version: 1}}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, uidev)
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to write user device buffer: %v", err)
	}
	_, err = deviceFile.Write(buf.Bytes())
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to write uidev struct to device file: %v", err)
	}

	err = ioctl(deviceFile, uiDevCreate, uintptr(0))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to create device: %v", err)
	}
	time.Sleep(time.Millisecond * 200)
	return deviceFile, err
}

func releaseDevice(deviceFile *os.File) (err error) {
	return ioctl(deviceFile, uiDevDestroy, uintptr(0))
}

// original function taken from: https://github.com/tianon/debian-golang-pty/blob/master/ioctl.go
func ioctl(deviceFile *os.File, cmd, ptr uintptr) error {
	_, _, errorCode := syscall.Syscall(syscall.SYS_IOCTL, deviceFile.Fd(), cmd, ptr)
	if errorCode != 0 {
		return errorCode
	}
	return nil
}
func createVKeyboardDevice(path string, name []byte) (fd *os.File, err error) {
	fd, err = initKeyboardDevice(path, name)
	if err != nil {
		return nil, fmt.Errorf("error during initialization of keyboard '%s' at %s: %v", name, path, err)
	}

	return fd, nil
}

func sendBtnEvent(deviceFile *os.File, key int, btnState int) (err error) {
	if key < 1 || key > keyMax {
		return fmt.Errorf("could not send key event: invalid keycode '%d'", key);
	}
	buf, err := inputEventToBuffer(evKey, uint16(key), int32(btnState))
	if err != nil {
		return fmt.Errorf("key event could not be set: %v", err)
	}
	_, err = deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("writing btnEvent structure to the device file failed: %v", err)
	}
	err = syncEvents(deviceFile)
	if err != nil && err != syscall.EINVAL {
		return fmt.Errorf("sync to device file failed: %v", err)
	}
	return nil
}

func syncEvents(deviceFile *os.File) (err error) {
	buf, err := inputEventToBuffer(evSyn, 0, int32(synReport))
	if err != nil {
		return fmt.Errorf("writing sync event failed: %v", err)
	}
	_, err = deviceFile.Write(buf)
	return err
}

func inputEventToBuffer(evType uint16, evCode uint16, evValue int32)(buffer []byte, err error) {
	iev := inputEvent{
		Time: syscall.Timeval{0,0},
		Type: evType,
		Code: evCode,
		Value: evValue}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, iev)
	if err != nil {
		return nil, fmt.Errorf("failed to write input event to buffer: %v", err)
	}
	return buf.Bytes(), nil
}
