Uinput [![Build Status](https://travis-ci.org/bendahl/uinput.svg?branch=master)](https://travis-ci.org/bendahl/uinput) [![GoDoc](https://godoc.org/github.com/bendahl/uinput?status.png)](https://godoc.org/github.com/bendahl/uinput) [![Go Report Card](https://goreportcard.com/badge/github.com/bendahl/uinput)](https://goreportcard.com/report/github.com/bendahl/uinput)
====

This package provides pure go wrapper functions for the LINUX uinput device, which allows to create virtual input devices 
in userspace. At the moment this package offers a virtual keyboard implementation as well as a virtual mouse device and
a touch pad device. 

The keyboard can be used to either send single key presses or hold down a specified key and release it later 
(useful for building game controllers). The mouse device issues relative positional change events to the x and y axis 
of the mouse pointer and may also fire click events (left and right click). For implementing things like region selects
via a virtual mouse pointer, press and release functions for the mouse device are also included.

The touch pad, on the other hand can be used to move the mouse cursor to the specified position on the screen and to
issue left and right clicks. Note that you'll need to specify the region size of your screen first though (happens during
device creation).

Please note that you will need to make sure to have the necessary rights to write to uinput. You can either chmod your 
uinput device, or add a rule in /etc/udev/rules.d to allow your user's group or a dedicated group to write to the device.
You may use the following two commands to add the necessary rights for you current user to a file called 99-$USER.rules 
(where $USER is your current user's name):
<pre><code>
echo KERNEL==\"uinput\", GROUP=\"$USER\", MODE:=\"0660\" | sudo tee /etc/udev/rules.d/99-$USER.rules
sudo udevadm trigger
</code></pre>

Installation
-------------
Simply check out the repository and use the commands <pre><code>go build && go install</code></pre> 
The package will then be installed to your local respository, along with the package documentation. 
The documentation contains more details on the usage of this package. 

License
--------
The package falls under the MIT license. Please see the "LICENSE" file for details.

Current Status
--------------
2018-03-31: I am happy to announce that v1.0.0! Go ahead and use this library in your own projects! Feedback is always welcome.

TODO
----
The current API can be considered stable and the overall functionality (as originally envisioned) is complete. 
Testing on x86_64 and ARM platforms (specifically the RaspberryPi) has been successful. If you'd like to use this library
on a different platform that supports Linux, feel free to test it and share the results. This would be greatly appreciated.
One thing that I'd still like to improve, however, are the test cases. The basic functionality is covered, but more extensive
testing is something that needs to be worked on. 

- [x] Create Tests for the uinput package
- [x] Migrate code from C to GO
- [x] Implement relative input
- [x] Implement absolute input
- [x] Test on different platforms
    - [x] x86_64
    - [x] ARMv6 (RaspberryPi)
- [x] Implement functions to allow mouse button up and down events (for region selects)
- [ ] Extend test cases

