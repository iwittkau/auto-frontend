package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/rs/xid"

	"github.com/gorilla/websocket"
	"github.com/iwittkau/auto-frontend"

	"html/template"

	"net/http"
)

// Frontend implements the auto_frontend Frontend interface
type Frontend struct {
	name         string
	address      string
	mux          *sync.Mutex
	templatePath string
	buttons      []auto_frontend.Button
	setButtons   []auto_frontend.SetButton
	getButtons   []auto_frontend.GetButton
	doButtons    []auto_frontend.DoButton
}

// Data is given to the template and holds all data needed to render
type Data struct {
	Version     string
	BackendName string
	BackendUrl  string
	Buttons     []auto_frontend.Button
}

// New retruns a new Frontend
func New(name string, address string, templatePath string) *Frontend {
	return &Frontend{name, address, &sync.Mutex{}, templatePath, []auto_frontend.Button{}, []auto_frontend.SetButton{}, []auto_frontend.GetButton{}, []auto_frontend.DoButton{}}
}

// Render
func (f *Frontend) Render(w http.ResponseWriter) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	tmpl := template.New("frontend")
	var err error
	if f.templatePath == "" {
		tmpl, err = tmpl.Parse(frontendTmpl)
	} else {
		tmpl, err = template.ParseFiles(f.templatePath)
	}

	if err != nil {
		return err
	}

	var buttons []auto_frontend.Button

	for _, b := range f.getButtons {
		btn := auto_frontend.Button{b.ID(), b.Label(), auto_frontend.ButtonTypeGet, http.MethodGet, b.Class(), b.Keys(), "/" + b.ID()}
		buttons = append(buttons, btn)
	}

	for _, b := range f.setButtons {
		btn := auto_frontend.Button{b.ID(), b.Label(), auto_frontend.ButtonTypeSet, http.MethodPost, b.Class(), b.Keys(), "/" + b.ID()}
		buttons = append(buttons, btn)
	}

	for _, b := range f.doButtons {
		btn := auto_frontend.Button{b.ID(), b.Label(), auto_frontend.ButtonTypeDo, http.MethodGet, b.Class(), nil, "/" + b.ID()}
		buttons = append(buttons, btn)
	}

	data := Data{
		Version:     xid.New().String(),
		BackendName: f.name,
		BackendUrl:  "http://" + f.address,
		Buttons:     buttons,
	}

	return tmpl.Execute(w, data)

}

// RegisterSetButton
func (f *Frontend) RegisterSetButton(button auto_frontend.SetButton) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, b := range f.setButtons {
		if b.ID() == button.ID() {
			return auto_frontend.ErrButtonAlreadyRegistered
		}
	}
	f.setButtons = append(f.setButtons, button)
	return nil
}

// RegisterSetButtons
func (f *Frontend) RegisterSetButtons(buttons ...auto_frontend.SetButton) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, b := range f.setButtons {
		for _, btn := range buttons {
			if b.ID() == btn.ID() {
				return auto_frontend.ErrButtonAlreadyRegistered
			}
		}

	}
	f.setButtons = append(f.setButtons, buttons...)
	return nil
}

// RegisterGetButton
func (f *Frontend) RegisterGetButton(button auto_frontend.GetButton) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, b := range f.getButtons {
		if b.ID() == button.ID() {
			return auto_frontend.ErrButtonAlreadyRegistered
		}
	}
	f.getButtons = append(f.getButtons, button)
	return nil
}

// RegisterGetButtons
func (f *Frontend) RegisterGetButtons(buttons ...auto_frontend.GetButton) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, b := range f.getButtons {
		for _, btn := range buttons {
			if b.ID() == btn.ID() {
				return auto_frontend.ErrButtonAlreadyRegistered
			}
		}

	}
	f.getButtons = append(f.getButtons, buttons...)
	return nil
}

// RegisterDoButton
func (f *Frontend) RegisterDoButton(button auto_frontend.DoButton) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, b := range f.doButtons {
		if b.ID() == button.ID() {
			return auto_frontend.ErrButtonAlreadyRegistered
		}
	}
	f.doButtons = append(f.doButtons, button)
	return nil
}

// RegisterDoButtons
func (f *Frontend) RegisterDoButtons(buttons ...auto_frontend.DoButton) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, b := range f.doButtons {
		for _, btn := range buttons {
			if b.ID() == btn.ID() {
				return auto_frontend.ErrButtonAlreadyRegistered
			}
		}

	}
	f.doButtons = append(f.doButtons, buttons...)
	return nil
}

// GetButtons
func (f *Frontend) GetButtons() []auto_frontend.GetButton {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.getButtons
}

// SetButtons
func (f *Frontend) SetButtons() []auto_frontend.SetButton {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.setButtons
}

// DoButtons
func (f *Frontend) DoButtons() []auto_frontend.DoButton {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.doButtons
}

// Start
func (f *Frontend) Start() error {

	setButtons := f.SetButtons()

	getButtons := f.GetButtons()

	doButtons := f.DoButtons()

	for _, setButton := range setButtons {
		b := setButton
		// fmt.Println("'" + b.Label() + "' -> " + http.MethodGet + " /" + b.Function())
		http.HandleFunc("/"+b.ID(), func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var data map[string]interface{}
			err := decoder.Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				defer r.Body.Close()
				b.CallbackSet(data)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}
		})
	}

	for _, getButton := range getButtons {
		b := getButton
		// fmt.Println("'" + b.Label() + "' -> " + http.MethodGet + " /" + b.Function())
		http.HandleFunc("/"+b.ID(), func(w http.ResponseWriter, r *http.Request) {
			result := b.CallbackGet()
			jsn, err := json.Marshal(result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(jsn)
			}

		})
	}

	for _, doButton := range doButtons {
		b := doButton
		// fmt.Println("'" + b.Label() + "' -> " + http.MethodGet + " /" + b.Function())
		http.HandleFunc("/"+b.ID(), func(w http.ResponseWriter, r *http.Request) {
			b.CallbackDo()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("done"))
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := f.Render(w)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	})

	fmt.Println(fmt.Sprintf("Auto Frontend '%s' listening on http://%s", f.name, f.address))

	return http.ListenAndServe(f.address, nil)
}
