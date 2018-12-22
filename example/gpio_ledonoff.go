// +build netbsd,arm

// Example for Raspberry PI with GPIO Switch / LED
// Hardware:
//   GPIO 5, 6:    output: LED 
//   GPIO 22, 23: input:   Switch (Close to GND)
// Configuration at /etc/gpio.conf:
//   gpio0 5 set out
//   gpio0 6 set out
//   gpio0 22 set in
//   gpio0 22 set pu
//   gpio0 23 set in
//   gpio0 23 set pu

package main

import (
	"github.com/oshimaya/dev"
	"time"
)

func main() {

	gpio := dev.NewGpio("/dev/gpio0")

	err := gpio.Open()
	if err != nil {
		return
	}
	defer gpio.Close()

	for {
		s1, _ := gpio.ReadPin(22)
		s2, _ := gpio.ReadPin(23)
		if s1 == 0 {
			gpio.WritePin(5, 1)
		} else {
			gpio.WritePin(5, 0)
		}
		if s2 == 0 {
			gpio.WritePin(6, 1)
		} else {
			gpio.WritePin(6, 0)
		}
		time.Sleep(50 * time.Millisecond)
	}

}
