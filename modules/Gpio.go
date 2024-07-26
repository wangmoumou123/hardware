package modules

import (
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"time"
)

type Gpio interface {
	Read() int
}
type gpio struct {
	Pin int
}

func (g *gpio) Read() int {
	return 0
}

func Tt() {
	err := rpio.Open()
	if err != nil {
		log.Println(err)
	}
	pin := rpio.Pin(4)

	pin.Output() // Output mode
	//pin.Mode(rpio.Output) // Alternative syntax
	for {
		time.Sleep(time.Second)
		pin.High() // Set pin High
		time.Sleep(time.Second)
		pin.Low() // Set pin Low
	}

	//pin.Toggle() // Toggle pin (Low -> High -> Low)

	//pin.Input()       // Input mode
	//res := pin.Read() // Read state from pin (High / Low)
	//log.Println(res)
	//pin.Mode(rpio.Output) // Alternative syntax
	//pin.Write(rpio.High)  // Alternative syntax
}
