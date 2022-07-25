package exp

// Experimental gamepad package with a focus on usability.

type Button struct {
	code uint16
}

func (b Button) Tap(vgp vGamepad) error {
	vgp.ButtonPress(b.code)
}

func (b Button) Down(vgp vGamepad) error {
	vgp.ButtonDown(b.code)
}

func (b Button) Up(vgp vGamepad) error {
	vgp.ButtonUp(b.code)
}

type ButtonCluster struct {
	// Dpad
	DPad Dpad

	// Face buttons
	North Button
	East  Button
	South Button
	West  Button

	// Bumpers
	LB Button
	RB Button

	// Triggers
	LT Button
	RT Button

	// Thumbsticks
	L3 Button
	R3 Button

	// Center cluster
	Select Button
	Start  Button
	Mode   Button // Center button usually bearing the manufacturer's logo
}

type Dpad struct {
	Up    Button
	Right Button
	Down  Button
	Left  Button
}

func newDpad() Dpad {
	return Dpad{
		Up:    Button{code: ButtonDpadUp},
		Right: Button{code: ButtonDpadRight},
		Down:  Button{code: ButtonDpadDown},
		Left:  Button{code: ButtonDpadLeft},
	}
}

type Axis struct {
	code uint16
}

func (a Axis) MoveTo(value float32, vgp vGamepad) error {
	vgp.sendStickAxisEvent(a.code, value)
}

type Stick struct {
	x, y Axis
}

func (s Stick) MoveTo(x, y float32, vgp vGamepad) error {
	err := vgp.sendStickAxisEvent(x.code, x)
	err = vgp.sendStickAxisEvent(y.code, y)
	return err
}

func (s Stick) MoveX(value float32, vgp vGamepad) error {
	vgp.sendStickAxisEvent(x.code, x)
}

func (s Stick) MoveY(value float32, vgp vGamepad) error {
	vgp.sendStickAxisEvent(y.code, y)
}

func (s Stick) Center(vgp vGamepad) error {
	err := vgp.sendStickAxisEvent(x.code, 1)
	err = vgp.sendStickAxisEvent(y.code, 1)
	return err
}

type Thumbsticks struct {
	Left  Stick
	Right Stick
}

type Gamepad struct {
	vGP vGamepad

	Buttons ButtonCluster

	Thumbsticks Thumbsticks
}

func NewGamepad(vgp vGamepad) Gamepad {
	buttons := ButtonCluster{
		DPad: newDpad(),
		// Face buttons
		North: Button{code: ButtonNorth},
		East:  Button{code: ButtonEast},
		South: Button{code: ButtonSouth},
		West:  Button{code: ButtonWest},

		// Bumpers
		LB: Button{code: ButtonBumperLeft},
		RB: Button{code: ButtonBumperRight},

		// Triggers
		LT: Button{code: ButtonTriggerLeft},
		RT: Button{code: ButtonTriggerRight},

		// Thumbsticks
		L3: Button{code: ButtonThumbLeft},
		R3: Button{code: ButtonThumbRight},

		// Center cluster
		Select: Button{code: ButtonSelect},
		Start:  Button{code: ButtonStart},
		Mode:   Button{code: ButtonMode},
	}

	axes := Thumbsticks{
		Left:  Stick{x: Axis{code: absX}, y: Axis{code: absY}},
		Rigth: Stick{x: Axis{code: absRX}, y: Axis{code: absRY}},
	}

	return Gamepad{
		vGP:     vgp,
		Buttons: buttons,
		Axes:    axes,
	}
}
