package main

import (
	"fmt"
	"time"

	"github.com/studiolambda/immersim"
	"github.com/studiolambda/immersim/event"
	"github.com/studiolambda/immersim/resource"
	"github.com/studiolambda/immersim/storage"
)

func isAbove(name string, storage *storage.Storage) bool {
	setpoint, _ := storage.Read("setpoint")
	tank_temperature, _ := storage.Read("tank_temperature")

	return tank_temperature.(float32) > float32(setpoint.(int32))
}

func main() {
	events := event.NewEvents()
	storage := storage.NewStorage(map[string]storage.Resource{
		"tank_temperature": resource.NewSineWave(
			0.15,                // Frequency
			50,                  // Amplitude
			50,                  // Offset
			50*time.Millisecond, // Interval update
		),
		"setpoint": resource.NewStaticReadWriter[int32](25),
		"is_above": resource.NewComputed(isAbove, []string{"tank_temperature", "setpoint"}),
		"rand":     resource.NewRandom[int32](0, 20, 100*time.Millisecond),
	})

	app := immersim.NewApplication(storage, events)

	app.Start()
	defer app.Stop()

	for {
		tank_temperature, _ := storage.Read("tank_temperature")
		setpoint, _ := storage.Read("setpoint")
		is_above, _ := storage.Read("is_above")
		random, _ := storage.Read("rand")
		fmt.Printf("\rtank_temperature: %f\tsetpoint: %d\tis_above: %t\t rand: %d          ", tank_temperature, setpoint, is_above, random)
		time.Sleep(100 * time.Millisecond)
	}
}
