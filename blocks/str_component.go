package blocks

import (
	"fmt"
	"strings"
)

func (c Component) String() string {
	return "@" + c.Name
}

func (c PropertySet) String() string {
	return fmt.Sprintf("%v.%v = %v", c.Component, c.Property, c.Value)
}

func (c PropertyGet) String() string {
	return fmt.Sprintf("%v.%v", c.Component, c.Property)
}

func (c GenericPropertySet) String() string {
	// Button1:Text = "Hello, World!"
	pFormat := "%v->%v = %v"
	if !c.Component.Continuous() {
		pFormat = "(%v)->%v = %v"
	}
	return fmt.Sprintf(pFormat, c.Component, c.Property, c.Value)
}

func (c GenericPropertyGet) String() string {
	// Button:Text
	pFormat := "%v->%v"
	if !c.Component.Continuous() {
		pFormat = "(%v)->%v"
	}
	return fmt.Sprintf(pFormat, c.Component, c.Property)
}

func (c Event) String() string {
	pFormat := "when %v.%v(%v) {\n%v}"
	if c.IsGeneric {
		pFormat = "when any %v.%v(%v) {\n%v}"
	}
	return fmt.Sprintf(pFormat, c.Component, c.Event, strings.Join(c.Parameters, ", "), PadBody(c.Body))
}

func (c MethodCall) String() string {
	return fmt.Sprintf("%v.%v(%v)", c.Component, c.Method, JoinBlocks(c.Args, ", "))
}

func (c GenericMethodCall) String() string {
	pFormat := "%v->%v(%v)"
	if !c.Component.Continuous() {
		pFormat = "(%v)->%v(%v)"
	}
	return fmt.Sprintf(pFormat, c.Component, c.Method, JoinBlocks(c.Args, ", "))
}

func (c AllComponent) String() string {
	return "Every." + c.Type
}
