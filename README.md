Uinput [![Build Status](https://travis-ci.org/bendahl/uinput.svg?branch=master)](https://travis-ci.org/bendahl/uinput) [![GoDoc](https://godoc.org/github.com/bendahl/uinput?status.png)](https://godoc.org/github.com/bendahl/uinput) [![Go Report Card](https://goreportcard.com/badge/github.com/bendahl/uinput)](https://goreportcard.com/report/github.com/bendahl/uinput)
====

This package provides pure go wrapper functions for the LINUX uinput device, which allows to create virtual input devices in userspace. At the moment this package offers a virtual keyboard implementation as well as a virtual mouse device. The keyboard can be used to either send single key presses or hold down a specified key and release it later (useful for building game controllers). The mouse device issues relative positional change events to the x and y asix of the mouse pointer and may also fire click events (left and right click). More functionality will be added in future version. 

Please note that you will need to make sure to have the necessary rights to write to uinput. You can either chmod your uinput device, or add a rule in /etc/udev/rules.d to allow your user's group or a dedicated group to write to the device.
You may use the following two commands to add the necessary rights for you current user to a file called 99-$USER.rules (where $USER is your current user's name):
<pre><code>
echo KERNEL==\"uinput\", GROUP=\"$USER\", MODE:=\"0660\" | sudo tee /etc/udev/rules.d/99-$USER.rules
sudo udevadm trigger
</code></pre>

Installation
-------------
Simply check out the repository and use the commands <pre><code>go build && go install</code></pre> The package will then be installed to your local respository, along with the package documentation. The documentation contains more details on the usage of this package. 

License
--------
The package falls under the MIT license. Please see the "LICENSE" file for details.

ToDos
------------------
Besides mouse and keyboard events it would be great to introduce an input device that allows absolute mouse pointer movements (similar to a touchpad). This still needs to be done. Also, all testing has been done on Ubunu 14.04 and 16.04 x86\_64. Testing for other platforms will need to be done. To get an idea of the things that are on the current todo list, check out the file "TODO.md". As always, helpful comments and ideas are always welcome. Feel free to do some testing on your own if you're up to it.

