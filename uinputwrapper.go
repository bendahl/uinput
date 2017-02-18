package uinput

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

// types needed from uinput.h
const (
	uinputMaxNameSize = 80
	uiDevCreate       = 0x5501
	uiDevDestroy      = 0x5502
	uiSetEvBit        = 0x40045564
	uiSetKeyBit       = 0x40045565
	uiSetRelBit       = 0x40045566
	uiSetAbsBit       = 0x40045567
	busUsb            = 0x03
)

// input event codes as specified in input-event-codes.h
const (
	evSyn      = 0x00
	evKey      = 0x01
	evRel      = 0x02
	evAbs      = 0x03
	relX       = 0x0
	relY       = 0x1
	absX       = 0x0
	absY       = 0x1
	synReport  = 0
	evBtnLeft  = 0x110
	evBtnRight = 0x111
)

const (
	btnStateReleased = 0
	btnStatePressed  = 1
	absSize          = 64
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
	Absmax     [absSize]int32
	Absmin     [absSize]int32
	Absfuzz    [absSize]int32
	Absflat    [absSize]int32
}

// translated to go from input.h
type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func closeDevice(deviceFile *os.File) (err error) {
	err = releaseDevice(deviceFile)
	if err != nil {
		return fmt.Errorf("failed to close device: %v", err)
	}
	return deviceFile.Close()
}

func releaseDevice(deviceFile *os.File) (err error) {
	return ioctl(deviceFile, uiDevDestroy, uintptr(0))
}

func createDeviceFile(path string) (fd *os.File, err error) {
	deviceFile, err := os.OpenFile(path, syscall.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		return nil, errors.New("could not open device file")
	}
	return deviceFile, err
}

func registerDevice(deviceFile *os.File, evType uintptr) error {
	err := ioctl(deviceFile, uiSetEvBit, evType)
	if err != nil {
		err = releaseDevice(deviceFile)
		if err != nil {
			deviceFile.Close()
			return fmt.Errorf("failed to close device: %v", err)
		}
		deviceFile.Close()
		return fmt.Errorf("invalid file handle returned from ioctl: %v", err)
	}
	return nil
}

func createVKeyboardDevice(path string, name []byte) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual keyboard device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register virtual keyboard device: %v", err)
	}

	// register key events
	for i := 0; i < keyMax; i++ {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(i))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register key number %d: %v", i, err)
		}
	}

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x4711,
				Product: 0x0815,
				Version: 1}})
}

func toUinputName(name []byte) (uinputName [uinputMaxNameSize]byte) {
	var fixedSizeName [uinputMaxNameSize]byte
	copy(fixedSizeName[:], name)
	return fixedSizeName
}

func createMouse(path string, name []byte) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create relative axis input device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register virtual mouse device: %v", err)
	}
	// register button events (in order to enable left and right click)
	err = ioctl(deviceFile, uiSetKeyBit, uintptr(evBtnLeft))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register left click event: %v", err)
	}
	err = ioctl(deviceFile, uiSetKeyBit, uintptr(evBtnRight))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register right click event: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evRel))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register relative axis input device: %v", err)
	}

	// register x and y axis events
	err = ioctl(deviceFile, uiSetRelBit, uintptr(relX))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register relative x axis events: %v", err)
	}
	err = ioctl(deviceFile, uiSetRelBit, uintptr(relY))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register relative y axis events: %v", err)
	}

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x4711,
				Product: 0x0816,
				Version: 1}})
}

func createUsbDevice(deviceFile *os.File, dev uinputUserDev) (fd *os.File, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, dev)
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

func sendBtnEvent(deviceFile *os.File, key int, btnState int) (err error) {
	buf, err := inputEventToBuffer(inputEvent{
		Time:  syscall.Timeval{0, 0},
		Type:  evKey,
		Code:  uint16(key),
		Value: int32(btnState)})
	if err != nil {
		return fmt.Errorf("key event could not be set: %v", err)
	}
	_, err = deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("writing btnEvent structure to the device file failed: %v", err)
	}
	return err
}

func sendRelEvent(deviceFile *os.File, eventCode uint16, pixel int32) error {
	iev := inputEvent{
		Time:  syscall.Timeval{0, 0},
		Type:  evRel,
		Code:  eventCode,
		Value: pixel}

	buf, err := inputEventToBuffer(iev)
	if err != nil {
		return fmt.Errorf("writing abs event failed: %v", err)
	}

	_, err = deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write abs event to device file: %v", err)
	}

	return syncEvents(deviceFile)
}

func syncEvents(deviceFile *os.File) (err error) {
	buf, err := inputEventToBuffer(inputEvent{
		Time:  syscall.Timeval{0, 0},
		Type:  evSyn,
		Code:  0,
		Value: int32(synReport)})
	if err != nil {
		return fmt.Errorf("writing sync event failed: %v", err)
	}
	_, err = deviceFile.Write(buf)
	return err
}

func inputEventToBuffer(iev inputEvent) (buffer []byte, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, iev)
	if err != nil {
		return nil, fmt.Errorf("failed to write input event to buffer: %v", err)
	}
	return buf.Bytes(), nil
}

// original function taken from: https://github.com/tianon/debian-golang-pty/blob/master/ioctl.go
func ioctl(deviceFile *os.File, cmd, ptr uintptr) error {
	_, _, errorCode := syscall.Syscall(syscall.SYS_IOCTL, deviceFile.Fd(), cmd, ptr)
	if errorCode != 0 {
		return errorCode
	}
	return nil
}
