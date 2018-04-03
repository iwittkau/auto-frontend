package button

import (
	"github.com/iwittkau/auto-frontend"
	"github.com/rs/xid"
)

// Button sytling class constants
const (
	ClassPrimary    = "primary"
	ClassSecondary  = "secondary"
	ClassSuccess    = "success"
	ClassDanger     = "danger"
	ClassWarning    = "warning"
	ClassInfo       = "info"
	ClassLight      = "light"
	ClassDark       = "dark"
	DefaultClassGet = ClassPrimary
	DefaultClassDo  = ClassPrimary
	DefaultClassSet = ClassWarning
)

// Button implements all auto_frontend button interfaces
type Button struct {
	id          string
	label       string
	class       string
	keys        []string
	getCallback auto_frontend.GetCallback
	setCallback auto_frontend.SetCallback
	doCallback  auto_frontend.DoCallback
}

// Keys is a string slice
type Keys []string

// NewGet returns a new button with a auto_frontend.GetCallback
func NewGet(label string, class string, keys Keys, cb auto_frontend.GetCallback) *Button {
	guid := xid.New()
	if class == "" {
		class = DefaultClassGet
	}
	return &Button{guid.String(), label, class, keys, cb, nil, nil}
}

// NewSet returns a new button with a auto_frontend.SetCallback
func NewSet(label string, class string, keys Keys, cb auto_frontend.SetCallback) *Button {
	guid := xid.New()
	if class == "" {
		class = DefaultClassSet
	}
	return &Button{guid.String(), label, class, keys, nil, cb, nil}
}

// NewDo returns a new button with a auto_frontend.DoCallback
func NewDo(label string, class string, cb auto_frontend.DoCallback) *Button {
	guid := xid.New()
	if class == "" {
		class = DefaultClassDo
	}
	return &Button{guid.String(), label, class, nil, nil, nil, cb}
}

// ID returns Button.id
func (b *Button) ID() string {
	return b.id
}

// Label returns Button.label
func (b *Button) Label() string {
	return b.label
}

// Class returns Button.class
func (b *Button) Class() string {
	return b.class
}

// Keys returns Button.keys
func (b *Button) Keys() []string {
	return b.keys
}

// CallbackGet returns the call of the getCallback function
func (b *Button) CallbackGet() map[string]interface{} {
	if b.getCallback != nil {
		return b.getCallback()
	}
	return map[string]interface{}{}
}

// CallbackSet calls the setCallback function with data
func (b *Button) CallbackSet(data map[string]interface{}) {
	if b.setCallback != nil {
		b.setCallback(data)
	}
}

// CallbackDo calls the doCallback function
func (b *Button) CallbackDo() {
	if b.doCallback != nil {
		b.doCallback()
	}
}
