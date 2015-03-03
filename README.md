Uinput
======

This go package provides wrapper functions for the LINUX uinput device. As it stands right now, only virtual keyboards are supported. Support for realative and absolute input devices will be added later on.

Please note that you will need to make sure to have the necessary rights to write to uinput. You can either chmod your uinput device, or add a rule in /etc/udev/rules.d to allow your user's group or a dedicated group to write to the device. An example file could be named "99-user.rules" and the line you would need to add for "user", belonging to the group "utest" would be <pre><code>KERNEL=="uinput", GROUP="utest", MODE:="0660"</code></pre> Also, make sure to restart in order for these settings to work. Which approach you'll take is up to you, although I would encourage the creation of a udev rule, as it is the clean approach.

Installation
-------------
Simply check out the repository and use the commands <pre><code>go build && go install</code></pre> The package will then be installed to your local respository, along with the package documentation. The documentation contains more details on the usage of this package. 

Don't worry about the C sources, as CGO will take care of compiling these for you as well. However, you will need to make sure to have the necessray header files for gcc installed on your system. They should be located underneath "/usr/include/linux".

License
--------
The package falls under the MIT license. Please see the "LICENSE" file for details.

ToDos/ Open Issues
------------------
The package is currently a work in progress and some more testing will need to be done. Also, as mentioned before, a few features will still need to be implemented as well. To get an idea of the things that are on the current todo list, check out the file "TODO.md". As always, helpful comments and ideas are always welcome. Feel free to do some testing on your own if you're up to it.

Version History
---------------
* 2014-05-11 - V0.1 
	
	This is an initial snapshot release, containing the basic functionality of a virtual keyboard.
