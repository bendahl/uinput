/*
 * File:   uinputwrapper.c
 * Author: Benjamin Dahlmanns
 *
 * Created on  May 1, 2014
 *
 * The handling of the virtual keyboard is done here. The uinput device will need to be writeable in
 * order to use this code.
 *
 * Special thanks goes to Gregory Thiemonge for the excellent introduction to uinput on this site:
 * http://thiemonge.org/getting-started-with-uinput
 *
 */

#include "uinputwrapper.h"

int initVKeyboardDevice(char* uinputPath) {
    int i;
    int deviceHandle = -1;
    struct uinput_user_dev uidev;

    deviceHandle = open(uinputPath, O_WRONLY | O_NONBLOCK | O_NDELAY);

    // if a valid handle could be determined, try to enable key events
    if(deviceHandle > 0) {
        if(ioctl(deviceHandle, UI_SET_EVBIT, EV_KEY) < 0) {
            if(releaseDevice(deviceHandle) < 0) {
                exit(EXIT_FAILURE);
            } else {
                deviceHandle = -1;
            }
        } else {
            ioctl (deviceHandle, UI_SET_EVBIT, EV_REP);

            for(i=0; i<255; i++) {
                ioctl(deviceHandle, UI_SET_KEYBIT, i);
            }

            memset(&uidev, 0, sizeof (uidev));
            snprintf(uidev.name, UINPUT_MAX_NAME_SIZE, "uinput_vkeyboard");
            uidev.id.bustype = BUS_USB;
            uidev.id.vendor  = 0x4711;
            uidev.id.product = 0x0815;
            uidev.id.version = 1;

            if (write(deviceHandle, &uidev, sizeof (uidev)) < 0) {
                exit(EXIT_FAILURE);
            }

            if (ioctl(deviceHandle, UI_DEV_CREATE) < 0) {
                exit(EXIT_FAILURE);
            }

            sleep(2);

        }
    }

    return deviceHandle;
}

int sendBtnEvent(int deviceHandle, int key, int btnState) {
    struct input_event ev;
    memset(&ev, 0, sizeof(ev));
    
    ev.type = EV_KEY;
    ev.code = key;  
    ev.value= btnState;
    int ret = write(deviceHandle, &ev, sizeof(ev));


	memset(&ev, 0, sizeof(ev));
	ev.type  = EV_SYN;
	ev.code  = 0;
	ev.value = SYN_REPORT;
	ret = write(deviceHandle, &ev, sizeof(ev));

	return ret;
}

int releaseDevice(int deviceHandle) {
	return ioctl(deviceHandle, UI_DEV_DESTROY);
}
