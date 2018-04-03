package main

import (
	"fmt"

	"github.com/iwittkau/auto-frontend/button"
	"github.com/iwittkau/auto-frontend/frontend"
)

var (
	number  = 0
	getFunc = func() map[string]interface{} {
		return map[string]interface{}{
			"Number": number,
		}
	}

	setFunc = func(data map[string]interface{}) {
		fmt.Sscanf(data["Number"].(string), "%d", &number)
	}
)

func main() {

	// a button for getting the number variable
	btnGet1 := button.NewGet("Get Number", button.DefaultClassGet, button.Keys{"Number"}, getFunc)

	btnSet1 := button.NewSet("Set Number", button.DefaultClassSet, button.Keys{"Number"}, setFunc)

	// Set a name, address with port, and a HTML template path (leave empty to use the default template)
	f := frontend.New("Auto Frontend", "localhost:8080", "")

	// register get button
	f.RegisterGetButtons(btnGet1)

	// register set button
	f.RegisterSetButtons(btnSet1)

	// start the frontend
	if err := f.Start(); err != nil {
		panic(err)
	}

}
