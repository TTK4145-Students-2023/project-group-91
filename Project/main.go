package main

import (
	"Driver-go/elevio"
	"fmt"
)

const Open_Door_Time = 2
const Num_Of_Flors = 4

func main() {

	elevio.Init("localhost:15657", Num_Of_Flors)

	// SETUP
	elev := elevio.Elev{}
	elev.Init()

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

			elev.Orders.SetOrder(a.Floor, a.Button)
			// fmt.Println(elev.orders)
			if elev.Dir == 0 {
				elev.NextOrder()
			}

		case a := <-drv_floors:
			elev.UpdateFloor()
			elev.CheckOrder(a)
			fmt.Println("current floor:", elev.GetFloor())

		case a := <-drv_stop:
			if a {
				elev.Stop()
			}

		case a := <-drv_obstr:
			if a && elev.DoorOpen {

				elev.Stop()

			} else if elev.DoorOpen {

				elev.CloseDoors()

				if !elev.CheckOrder(elev.CurFloor) {
					elev.NextOrder()
				}
			}

			//MANUAL CONTROL
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

			// case a := <-drv_obstr:
			// 	fmt.Printf("--3:%+v\n", a)
			// 	if a {
			// 		elevio.SetMotorDirection(elevio.MD_Stop)
			// 	} else {
			// 		elevio.SetMotorDirection(d)
			// 	}

			// case a := <-drv_stop:
			// 	fmt.Printf("--4:%+v\n", a)
			// 	for f := 0; f < Num_Of_Flors; f++ {
			// 		for b := elevio.ButtonType(0); b < 3; b++ {
			// 			elevio.SetButtonLamp(b, f, false)
			// 		}
			// 	}
		}
	}
}
