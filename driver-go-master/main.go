package main

import (
	"Driver-go/elevio"
	"fmt"
)

type Orders struct {
	HallUp   []bool
	HallDown []bool
	Cab      []bool
}

func (o *Orders) Restart() {
	for i, _ := range o.Cab {
		o.HallDown[i] = false
		o.HallUp[i] = false
		o.Cab[i] = false
	}

}

func main() {

	numFloors := 4

	elevio.Init("localhost:20017", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	elevio.SetMotorDirection(d)

	// SETUP

	// creating orders list
	orders := Orders{
		make([]bool, numFloors),
		make([]bool, numFloors),
		make([]bool, numFloors)}
	orders.Restart()

	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {

		case a := <-drv_buttons:
			switch a.Button {
			case 0:
				if a.Floor != numFloors {
					orders.HallUp[a.Floor] = true
					elevio.SetButtonLamp(a.Button, a.Floor, true)
				}
			case 1:
				if a.Floor != 0 {
					orders.HallDown[a.Floor] = true
					elevio.SetButtonLamp(a.Button, a.Floor, true)
				}
			case 2:
				orders.Cab[a.Floor] = true
				elevio.SetButtonLamp(a.Button, a.Floor, true)
			}
			fmt.Println(orders)

			// case a := <-drv_buttons:
			// 	fmt.Printf("--1:%+v\n", a)
			// 	if a.Button == 2 {

			// 		if a.Floor == 0 {
			// 			d = elevio.MD_Down
			// 		} else if a.Floor == 1 {
			// 			d = elevio.MD_Stop
			// 		} else if a.Floor == 2 {
			// 			d = elevio.MD_Up
			// 		}

			// 		elevio.SetMotorDirection(d)

			// 	}
			// case a := <-drv_buttons:
			// 	fmt.Printf("--1:%+v\n", a)
			// 	elevio.SetButtonLamp(a.Button, a.Floor, true)

			// case a := <-drv_floors:
			// 	fmt.Printf("--2:%+v\n", a)
			// 	if a == numFloors-1 {
			// 		d = elevio.MD_Down
			// 	} else if a == 0 {
			// 		d = elevio.MD_Up
			// 	}
			// 	elevio.SetMotorDirection(d)

			// case a := <-drv_obstr:
			// 	fmt.Printf("--3:%+v\n", a)
			// 	if a {
			// 		elevio.SetMotorDirection(elevio.MD_Stop)
			// 	} else {
			// 		elevio.SetMotorDirection(d)
			// 	}

			// case a := <-drv_stop:
			// 	fmt.Printf("--4:%+v\n", a)
			// 	for f := 0; f < numFloors; f++ {
			// 		for b := elevio.ButtonType(0); b < 3; b++ {
			// 			elevio.SetButtonLamp(b, f, false)
			// 		}
			// 	}
		}
	}
}
