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

func createTouchPad(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create absolute axis input device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register key device: %v", err)
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

	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute axis input device: %v", err)
	}

	// register x and y axis events
	err = ioctl(deviceFile, uiSetAbsBit, uintptr(absX))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute x axis events: %v", err)
	}
	err = ioctl(deviceFile, uiSetAbsBit, uintptr(absY))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute y axis events: %v", err)
	}

	var absMin [absSize]int32
	absMin[absX] = minX
	absMin[absY] = minY

	var absMax [absSize]int32
	absMax[absX] = maxX
	absMax[absY] = maxY

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x4711,
				Product: 0x0817,
				Version: 1},
			Absmin: absMin,
			Absmax: absMax})
}

func createTouchScreen(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create absolute axis input device: %v", err)
	}

	BTN_TOUCH := 0x14a

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register key device: %v", err)
	}

	var uinp uinputUserDev

	uinp.Name = toUinputName([]byte(name))
	uinp.ID.Version = 4
	uinp.ID.Bustype = 0x6 //busUsb
	uinp.ID.Vendor = 0x4711
	uinp.ID.Product = 0x0817

	uinp.Absmin[AbsMtPositionX] = minX // screen dimension
	uinp.Absmax[AbsMtPositionX] = maxX // screen dimension
	uinp.Absmin[AbsMtPositionY] = minY // screen dimension
	uinp.Absmax[AbsMtPositionY] = maxY // screen dimension

	err = ioctl(deviceFile, uiSetEvBit, uintptr(evAbs))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute x axis events 2: %v", err)
	}
	err = ioctl(deviceFile, uiSetKeyBit, uintptr(BTN_TOUCH))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute x axis events 3: %v", err)
	}

	err = ioctl(deviceFile, uiSetAbsBit, uintptr(AbsMtPositionX))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute x axis events 4: %v", err)
	}

	err = ioctl(deviceFile, uiSetAbsBit, uintptr(AbsMtPositionY))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute x axis events 5: %v", err)
	}

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name:   toUinputName(name),
			ID:     uinp.ID,
			Absmin: uinp.Absmin,
			Absmax: uinp.Absmax})
}

func createMouse(path string, name []byte) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create relative axis input device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register key device: %v", err)
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
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
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

func sendAbsEvent(deviceFile *os.File, xPos int32, yPos int32) error {
	var ev [2]inputEvent
	ev[0].Type = evAbs
	ev[0].Code = absX
	ev[0].Value = xPos

	ev[1].Type = evAbs
	ev[1].Code = absY
	ev[1].Value = yPos

	for _, iev := range ev {
		buf, err := inputEventToBuffer(iev)
		if err != nil {
			return fmt.Errorf("writing abs event failed: %v", err)
		}

		_, err = deviceFile.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write abs event to device file: %v", err)
		}
	}

	return syncEvents(deviceFile)
}

func sendRelEvent(deviceFile *os.File, eventCode uint16, pixel int32) error {
	iev := inputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  evRel,
		Code:  eventCode,
		Value: pixel}

	buf, err := inputEventToBuffer(iev)
	if err != nil {
		return fmt.Errorf("writing abs event failed: %v", err)
	}

	_, err = deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write rel event to device file: %v", err)
	}

	return syncEvents(deviceFile)
}

func syncEvents(deviceFile *os.File) (err error) {
	buf, err := inputEventToBuffer(inputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  evSyn,
		Code:  0,
		Value: int32(synReport)})
	if err != nil {
		return fmt.Errorf("writing sync event failed: %v", err)
	}
	_, err = deviceFile.Write(buf)
	return err
}

func sendEvent(deviceFile *os.File, eventType uint16, eventCode uint16, eventValue int32) error {
	iev := inputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  eventType,
		Code:  eventCode,
		Value: eventValue}

	buf, err := inputEventToBuffer(iev)
	if err != nil {
		return fmt.Errorf("writing abs event failed: %v", err)
	}

	_, err = deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write rel event to device file: %v", err)
	}

	return nil
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
