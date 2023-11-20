package uinput

import (
	"fmt"
	"io"
	"os"
)

// MultiTouch is an input device that uses absolute axis events.
// Unlike the TouchPad, MultiTouch supports the simulation of multiple inputs (contacts)
// allowing for different gestures, for exmaple pinch to zoom.
// Each contact point is assigned a slot, making it necessary to define the maxmimum
// expected amount of contact points .
// Since MultiTouch uses absolute axis events, it is necessary to define the size
// of the rectangle in which the contacs may move upon creation of the device.
type MultiTouch interface {
	//Gets all contacts which can then be manipulated
	GetContacts() []multiTouchContact

	// FetchSyspath will return the syspath to the device file.
	FetchSyspath() (string, error)

	io.Closer
}

type vMultiTouch struct {
	name       []byte
	deviceFile *os.File
	contacts   []multiTouchContact
}

// The contact can be described as a finger contacting the surface of the MultiTouch device.
type multiTouchContact struct {
	multitouch  *vMultiTouch
	slot        int32
	tracking_id int32
}

// CreateMultiTouch will create a new multitouch device. Note that you will need to define the x and y-axis boundaries
// (min and max) within which the contacs maybe moved around, as well as the maximum amount of contacts allowed.
func CreateMultiTouch(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32, maxContacts int32) (MultiTouch, error) {
	err := validateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = validateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createMultiTouch(path, name, minX, maxX, minY, maxY, maxContacts)
	if err != nil {
		return nil, err
	}

	var multitouch vMultiTouch = vMultiTouch{name: name, deviceFile: fd}

	for i := int32(0); i < maxContacts; i++ {
		multitouch.contacts = append(multitouch.contacts, multiTouchContact{slot: i, multitouch: &multitouch})
	}

	return multitouch, nil
}

func (vMulti vMultiTouch) GetContacts() []multiTouchContact {
	return vMulti.contacts
}

func (vMulti vMultiTouch) FetchSyspath() (string, error) {
	return fetchSyspath(vMulti.deviceFile)
}

func (vMulti vMultiTouch) Close() error {
	return closeDevice(vMulti.deviceFile)
}

func createMultiTouch(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32, maxContacts int32) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create absolute axis input device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to register key device: %v", err)
	}

	for _, event := range []int{evBtnTouch} {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(event))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register button event %v: %v", event, err)
		}
	}

	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute axis input device: %v", err)
	}

	for _, event := range []int{
		absMtSlot,
		absMtTrackingId,
		absMtPositionX,
		absMtPositionY,
	} {
		err = ioctl(deviceFile, uiSetAbsBit, uintptr(event))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register absolute axis event %v: %v", event, err)
		}
	}

	var absMin [absSize]int32
	absMin[absMtPositionX] = minX
	absMin[absMtPositionY] = minY
	absMin[absMtTrackingId] = 0x00
	absMin[absMtSlot] = 0x00

	var absMax [absSize]int32
	absMax[absMtPositionX] = maxX
	absMax[absMtPositionY] = maxY
	absMax[absMtTrackingId] = maxContacts
	absMax[absMtSlot] = maxContacts

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x0,
				Product: 0x0,
				Version: 0},
			Absmin: absMin,
			Absmax: absMax})
}

// The contact will be held down at the coordinates specified
func (c multiTouchContact) TouchDownAt(x int32, y int32) error {
	var events []inputEvent

	events = append(events, inputEvent{
		Type:  evAbs,
		Code:  absMtPositionX,
		Value: x,
	})

	if x == 0 && y == 0 {
		y--
	}

	events = append(events, inputEvent{
		Type:  evAbs,
		Code:  absMtPositionY,
		Value: y,
	})

	c.tracking_id = c.slot

	return c.sendAbsEvent(events)
}

// The contact will be raised off of the surface
func (c multiTouchContact) TouchUp() error {
	c.tracking_id = -1
	return c.sendAbsEvent(nil)
}

func (c multiTouchContact) sendAbsEvent(events []inputEvent) error {
	var ev []inputEvent

	ev = append(ev, inputEvent{
		Type:  evAbs,
		Code:  absMtSlot,
		Value: c.slot,
	})

	ev = append(ev, inputEvent{
		Type:  evAbs,
		Code:  absMtTrackingId,
		Value: c.tracking_id,
	})

	if events != nil {
		ev = append(ev, events...)
	}

	for _, iev := range ev {
		buf, err := inputEventToBuffer(iev)
		if err != nil {
			return fmt.Errorf("writing abs event failed: %v", err)
		}

		_, err = c.multitouch.deviceFile.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write abs event to device file: %v", err)
		}
	}

	return syncEvents(c.multitouch.deviceFile)
}
