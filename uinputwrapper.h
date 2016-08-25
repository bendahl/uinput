/*
 * File:   uinputwrapper.h
 * Author: Benjamin Dahlmanns
 *
 * Created on  May 1, 2014
 */

#ifndef UINPUTWRAPPER_H
#define	UINPUTWRAPPER_H

#include <linux/input.h>
#include <linux/uinput.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sys/fcntl.h>
#include <unistd.h>

/*
 * This Method has to be called first in order to create and initialize the
 * virtual keyboard. It will return 0 if successful, otherwise
 * -1 will be returned.
 */
int initVKeyboardDevice(char* uinputPath, char* virtDeviceName);

/*
 * Send a button event to the virutal keyboard. Possible values for the key
 * variable can be found in input.h
 */
int sendBtnEvent(int deviceHandle, int key, int btnState);

/*
 * Releases the previously created device. Returns 0 if successful
 * and -1 if not successful.
 */
int releaseDevice(int deviceHanlde);



#endif	/* UINPUTWRAPPER_H */

