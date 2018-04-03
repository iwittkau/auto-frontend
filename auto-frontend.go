package auto_frontend

import "net/http"

// Button Types
const (
	ButtonTypeGet = "get_btn"
	ButtonTypeSet = "set_btn"
	ButtonTypeDo  = "do_btn"
)

type Frontend interface {
	Render(w http.ResponseWriter) error
	RegisterSetButton(button SetButton) error
	RegisterSetButtons(buttons ...SetButton) error
	SetButtons() []SetButton
	RegisterGetButton(button GetButton) error
	RegisterGetButtons(buttons ...GetButton) error
	GetButtons() []GetButton
	RegisterDoButton(button GetButton) error
	RegisterDoButtons(buttons ...GetButton) error
	DoButtons() []DoButton
	Start() error
}

type GetCallback func() map[string]interface{}
type SetCallback func(map[string]interface{})
type DoCallback func()

type ButtonType string

type Button struct {
	ID     string
	Label  string
	Type   ButtonType
	Method string
	Class  string
	Keys   []string
	Path   string
}

type GetButton interface {
	ID() string
	Label() string
	Class() string
	Keys() []string
	CallbackGet() map[string]interface{}
}

type SetButton interface {
	ID() string
	Label() string
	Class() string
	Keys() []string
	CallbackSet(map[string]interface{})
}

type DoButton interface {
	ID() string
	Label() string
	Class() string
	Keys() []string
	CallbackDo()
}
