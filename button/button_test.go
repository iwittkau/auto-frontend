package button_test

import (
	"testing"
	"time"

	"github.com/iwittkau/auto-frontend/button"
)

var (
	result = map[string]interface{}{"test": time.Now()}

	getFunc = func() map[string]interface{} {
		return result
	}
	setFunc = func(_ map[string]interface{}) {

	}

	doFunc = func() {

	}
)

func TestNewGet(t *testing.T) {

	keys := button.Keys{"test"}

	btnGet1 := button.NewGet("test", button.DefaultClassGet, keys, getFunc)
	btnGet2 := button.NewGet("test", "", keys, nil)

	if btnGet1.ID() == "" {
		t.Error("ID is not set")
	}

	if btnGet1.Class() != button.DefaultClassGet {
		t.Error("Class is not set to button.DefaultClassGet")
	}

	if btnGet1.Label() != "test" {
		t.Error("Label is not set to button.DefaultClassGet")
	}

	if btnGet1.Keys()[0] != keys[0] {
		t.Error("Keys are not set")
	}

	if btnGet1.CallbackGet()["test"] != result["test"] {
		t.Error("Result incorrect")
	}

	// Should not panic
	btnGet1.CallbackDo()
	btnGet1.CallbackSet(result)

	if btnGet2.ID() == "" {
		t.Error("ID is not set")
	}

	if btnGet2.Class() != button.DefaultClassGet {
		t.Error("Class is not set to button.DefaultClassGet")
	}

	if btnGet2.CallbackGet() == nil {
		t.Error("Callback result is nil")
	}
}

func TestNewSet(t *testing.T) {

	keys := button.Keys{"test"}

	btnSet1 := button.NewSet("test", button.DefaultClassSet, keys, setFunc)
	btnSet2 := button.NewSet("test", "", keys, nil)

	if btnSet1.ID() == "" {
		t.Error("ID is not set")
	}

	if btnSet1.Class() != button.DefaultClassSet {
		t.Error("Class is not set to button.DefaultClassGet")
	}

	if btnSet1.Label() != "test" {
		t.Error("Label is not set to button.DefaultClassGet")
	}

	if btnSet1.Keys()[0] != keys[0] {
		t.Error("Keys are not set")
	}

	// Should not panic
	btnSet1.CallbackDo()
	btnSet1.CallbackGet()
	btnSet1.CallbackSet(result)

	if btnSet2.ID() == "" {
		t.Error("ID is not set")
	}

	if btnSet2.Class() != button.DefaultClassSet {
		t.Error("Class is not set to button.DefaultClassGet")
	}

	if btnSet2.CallbackGet() == nil {
		t.Error("Callback result is nil")
	}

}

func TestNewDo(t *testing.T) {

	btnDo1 := button.NewDo("test", button.DefaultClassSet, doFunc)
	btnDo2 := button.NewDo("test", "", nil)

	if btnDo1.ID() == "" {
		t.Error("ID is not set")
	}

	if btnDo1.Class() != button.DefaultClassSet {
		t.Error("Class is not set to button.DefaultClassGet")
	}

	if btnDo1.Label() != "test" {
		t.Error("Label is not set to button.DefaultClassGet")
	}

	if btnDo1.Keys() != nil {
		t.Error("Keys are set")
	}

	// Should not panic
	btnDo1.CallbackDo()
	btnDo1.CallbackGet()
	btnDo1.CallbackSet(result)

	if btnDo2.ID() == "" {
		t.Error("ID is not set")
	}

	if btnDo2.Class() != button.DefaultClassDo {
		t.Error("Class is not set to button.DefaultClassGet")
	}

	if btnDo2.CallbackGet() == nil {
		t.Error("Callback result is nil")
	}

}
