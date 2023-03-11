package main

import (
	"Source/elev"
	"Source/elevio"
	"fmt"
)

const Open_Door_Time = 2
const Num_Of_Flors = 4

func main() {

	elevio.Init("localhost:15657", Num_Of_Flors)

	// SETUP
	elev := elev.Elev{}
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

		case button := <-drv_buttons:

			elev.Orders.SetOrder(button.Floor, button.Button)

			if button.Floor == elevio.GetFloor() {
				go elev.CompleteOrder(elev.GetFloor())

			} else if elev.Dir == 0 {

				go elev.NextOrder()
			}

		case floor := <-drv_floors:
			elev.UpdateFloor()
			elev.ShouldIstop(floor)
			fmt.Println("current floor:", elev.GetFloor())

		case stop := <-drv_stop:
			if stop {
				elev.Stop()
			}

		case obstr := <-drv_obstr:
			if obstr && elev.DoorOpen {

				elev.Stop()

			} else if elev.DoorOpen {

				elev.CloseDoors()

				if !elev.ShouldIstop(elev.CurFloor) {
					elev.NextOrder()
				}
			}

		}
	}
}
