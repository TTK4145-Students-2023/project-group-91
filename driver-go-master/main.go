package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

// SECTION - Orders
type Orders struct {
	HallUp   []bool
	HallDown []bool
	Cab      []bool
}

func (o *Orders) setOrder(floor int, button elevio.ButtonType) {

	if button == 0 {
		o.HallUp[floor] = true
		elevio.SetButtonLamp(button, floor, true)

	} else if button == 2 {
		o.Cab[floor] = true
		elevio.SetButtonLamp(button, floor, true)

	} else if button == 1 {
		o.HallDown[floor] = true
		elevio.SetButtonLamp(button, floor, true)

	}

}
func (o *Orders) clearAll() {
	for i := range o.Cab {
		o.HallDown[i] = false
		o.HallUp[i] = false
		o.Cab[i] = false
	}
	for i := range o.Cab {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, i, false)
		}
	}

}

func (o *Orders) completeOrder(floor int, dir elevio.MotorDirection) {

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	o.Cab[floor] = false
	o.HallUp[floor] = false
	o.HallDown[floor] = false

	for b := elevio.ButtonType(0); b < 3; b++ {
		elevio.SetButtonLamp(b, floor, false)
	}
	time.Sleep(2000 * time.Millisecond)
	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(dir)

}

//!SECTION

// SECTION - Elev
type Elev struct {
	dir       int
	curFloor  int
	doorOpen  bool
	numFloors int
	orders    Orders
}

func (e *Elev) setup(numFloors int) {
	e.dir = 0
	e.curFloor = 0
	e.doorOpen = false
	e.numFloors = numFloors
	e.orders = Orders{
		make([]bool, numFloors),
		make([]bool, numFloors),
		make([]bool, numFloors)}
	e.orders.clearAll()
}
func (e *Elev) updateFloor() int {
	e.curFloor = elevio.GetFloor()
	return e.curFloor
}
func (e *Elev) setFloor(f int) {
	e.curFloor = f
}
func (e Elev) getDirection() int {
	return e.dir
}
func (e Elev) getFloor() int {
	return e.curFloor
}
func (e *Elev) goUp() bool {
	if e.doorOpen {
		return false

	} else if e.curFloor == e.numFloors-1 {
		return false
	} else {
		e.dir = 1
		elevio.SetMotorDirection(elevio.MD_Up)
		return true
	}
}
func (e *Elev) goDown() bool {
	if e.doorOpen {
		return false

	} else if e.curFloor == 0 {
		return false

	} else {
		e.dir = 1
		elevio.SetMotorDirection(elevio.MD_Up)
		return true
	}
}
func (e *Elev) stop() {
	e.dir = 0
	elevio.SetMotorDirection(elevio.MD_Stop)
}
func (e *Elev) openDoors() {
	e.doorOpen = true
	elevio.SetDoorOpenLamp(true)
}
func (e *Elev) closeDoors() bool {

	if elevio.GetObstruction() {
		e.doorOpen = true
		return false
	} else {
		elevio.SetDoorOpenLamp(false)
		e.doorOpen = false
		return true
	}

}
func (e *Elev) completeOrder(floor int) {

}

//!SECTION

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	elevio.SetMotorDirection(d)

	// SETUP

	// creating orders list
	orders := Orders{
		make([]bool, numFloors),
		make([]bool, numFloors),
		make([]bool, numFloors)}
	orders.clearAll()

	elev := Elev{}
	elev.setup(numFloors)

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

		case a := <-drv_floors:
			fmt.Printf("--2:%+v\n", a)
			fmt.Println("GetFloor:", elevio.GetFloor())
			elev.setFloor(a)

			for i, v := range orders.Cab {
				if i == a && v == true {
					orders.completeOrder(a, d)

				} else if a == numFloors-1 {
					d = elevio.MD_Down
				} else if a == 0 {
					d = elevio.MD_Up
				}
				elevio.SetMotorDirection(d)
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
			// 	for f := 0; f < numFloors; f++ {
			// 		for b := elevio.ButtonType(0); b < 3; b++ {
			// 			elevio.SetButtonLamp(b, f, false)
			// 		}
			// 	}
		}
	}
}
