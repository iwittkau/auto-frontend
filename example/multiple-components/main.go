package main

import (
	"fmt"
	"time"

	"github.com/iwittkau/auto-frontend/button"
	"github.com/iwittkau/auto-frontend/frontend"
)

func main() {

	// some values to
	number1 := 0
	number2 := 0
	name := "Ingmar"

	// a channel to be closed for stopping
	stop := make(chan struct{})

	// a button for getting the numbers
	btnGet1 := button.NewGet("Get Numbers", button.DefaultClassGet, button.Keys{"Number1", "Number2"}, func() map[string]interface{} {
		return map[string]interface{}{
			"Number1": number1,
			"Number2": number2,
		}
	})

	// two inputs with keys to set number1 and number2
	btnSet1 := button.NewSet("Set Numbers", button.DefaultClassSet, button.Keys{"Number1", "Number2"}, func(data map[string]interface{}) {
		fmt.Sscanf(data["Number1"].(string), "%d", &number1)
		fmt.Sscanf(data["Number2"].(string), "%d", &number2)
	})

	// an input to set number1
	btnSet1a := button.NewSet("Set Number1", button.DefaultClassSet, button.Keys{"Number1"}, func(data map[string]interface{}) {
		fmt.Sscanf(data["Number1"].(string), "%d", &number1)
	})

	// an input to set number2
	btnSet1b := button.NewSet("Set Number2", button.DefaultClassSet, button.Keys{"Number2"}, func(data map[string]interface{}) {
		fmt.Sscanf(data["Number2"].(string), "%d", &number2)
	})

	// a button to stop, closes the channel
	btnDo1 := button.NewDo("Stop", button.ClassDanger, func() {
		close(stop)
	})

	// a button to get the time
	btnGet2 := button.NewGet("Get Time", button.DefaultClassGet, button.Keys{"Time"}, func() map[string]interface{} {
		t := time.Now()
		return map[string]interface{}{
			"Time": t,
		}
	})

	// a button to display the variable 'name'
	btnGet3 := button.NewGet("Get Name", button.DefaultClassGet, button.Keys{"Name"}, func() map[string]interface{} {

		return map[string]interface{}{
			"Name": name,
		}
	})

	// an input to set the variable 'name'
	btnSet2 := button.NewSet("Set Name", button.ClassDanger, button.Keys{"Name"}, func(data map[string]interface{}) {
		name = data["Name"].(string)
	})

	// create a new frontend with a name, address, and absolute path to your template. You can leave the template path blank to use the default template
	f := frontend.New("Auto Frontend", "localhost:8080", "")
	// f := frontend.New("Auto Frontend", "localhost:8080", "/Users/Ingmar/go/src/github.com/iwittkau/auto-frontend/frontend/tmpl/frontend.html")

	// register get buttons
	f.RegisterGetButtons(btnGet1, btnGet2, btnGet3)

	// register set buttons
	f.RegisterSetButtons(btnSet1, btnSet1a, btnSet1b, btnSet2)

	// register do button
	f.RegisterDoButton(btnDo1)

	// start the frontend in a goroutine
	go f.Start()

	// block until stop channel is closed
	select {
	case <-stop:
	}
}
