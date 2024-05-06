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
	tank_temperature, _ := storage.Read("tmp")

	return tank_temperature.(float32) > float32(setpoint.(int32))
}

func main() {
	events := event.NewEvents()
	storage := storage.NewStorage(map[string]storage.Resource{
		"tmp": resource.NewSineWave(
			0.15,                // Frequency
			50,                  // Amplitude
			50,                  // Offset
			50*time.Millisecond, // Interval update
		),
		"setpoint": resource.NewStaticReadWriter[int32](25),
		"above":    resource.NewComputed(isAbove, []string{"tmp", "setpoint"}),
		"rand":     resource.NewRandom[int32](0, 20, 10*time.Second),
		"feedback": resource.NewLinearFeedback[int32](1, 500*time.Millisecond, "rand"),
	})

	app := immersim.NewApplication(storage, events)

	app.Start()
	defer app.Stop()

	for {
		tmp, _ := storage.Read("tmp")
		setpoint, _ := storage.Read("setpoint")
		above, _ := storage.Read("above")
		random, _ := storage.Read("rand")
		feedback, _ := storage.Read("feedback")
		fmt.Printf("tmp: %f\tsetpoint: %d\tabove: %t\t rand: %d\t feedback: %d\n", tmp, setpoint, above, random, feedback)
		time.Sleep(100 * time.Millisecond)
	}
}
