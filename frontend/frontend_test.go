package frontend_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/emicklei/forest"

	"github.com/iwittkau/auto-frontend/button"
	"github.com/iwittkau/auto-frontend/frontend"

	. "github.com/emicklei/forest"
)

var (
	keys   = button.Keys{"test"}
	result = map[string]interface{}{"test": time.Now()}

	getFunc = func() map[string]interface{} {
		return result
	}
	setFunc = func(_ map[string]interface{}) {

	}

	doFunc = func() {

	}

	localhost = NewClient("http://localhost:8080", new(http.Client))
)

func TestNew(t *testing.T) {
	f := frontend.New("Test", "localhost:8080", "")
	btnGet1 := button.NewGet("test", button.DefaultClassGet, keys, getFunc)
	btnGet2 := button.NewGet("test", button.DefaultClassGet, keys, getFunc)
	if err := f.RegisterGetButton(btnGet1); err != nil {
		t.Error(err)
	}
	if err := f.RegisterGetButton(btnGet1); err == nil {
		t.Error("Should be an error")
	}
	if err := f.RegisterGetButtons(btnGet1); err == nil {
		t.Error("Should be an error")
	}
	if err := f.RegisterGetButtons(btnGet2); err != nil {
		t.Error(err)
	}

	btnSet1 := button.NewSet("test", button.DefaultClassSet, keys, setFunc)
	btnSet2 := button.NewSet("test", "", keys, nil)
	if err := f.RegisterSetButton(btnSet1); err != nil {
		t.Error(err)
	}
	if err := f.RegisterSetButton(btnSet1); err == nil {
		t.Error("Should be an error")
	}
	if err := f.RegisterSetButtons(btnSet1); err == nil {
		t.Error("Should be an error")
	}
	if err := f.RegisterSetButtons(btnSet2); err != nil {
		t.Error(err)
	}

	btnDo1 := button.NewDo("test", button.DefaultClassSet, doFunc)
	btnDo2 := button.NewDo("test", "", nil)

	if err := f.RegisterDoButton(btnDo1); err != nil {
		t.Error(err)
	}
	if err := f.RegisterDoButton(btnDo1); err == nil {
		t.Error("Should be an error")
	}
	if err := f.RegisterDoButtons(btnDo1); err == nil {
		t.Error("Should be an error")
	}
	if err := f.RegisterDoButtons(btnDo2); err != nil {
		t.Error(err)
	}

	go f.Start()
	time.Sleep(time.Second)

	root := localhost.GET(t, forest.Path("/"))
	ExpectStatus(t, root, http.StatusOK)

	w := httptest.NewRecorder()

	f.Render(w)

	btnGet1Req := localhost.GET(t, forest.Path("/"+btnGet1.ID()))
	ExpectStatus(t, btnGet1Req, http.StatusOK)

	wsReq := localhost.GET(t, forest.Path("/ws"))
	ExpectStatus(t, wsReq, http.StatusBadRequest)

}
