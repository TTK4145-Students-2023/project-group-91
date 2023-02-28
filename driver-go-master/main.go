package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

const Open_Door_Time = 2

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
func (o Orders) checkOrder(floor int) (bool, int) {

	if o.Cab[floor] {
		return true, 0

	} else if o.HallUp[floor] {

		if o.HallUp[floor] == o.HallDown[floor] {
			return true, 0
		} else {
			return true, 1
		}

	} else if o.HallDown[floor] {
		return true, -1

	} else {
		return false, 0

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

func (o *Orders) completeOrder(floor int) {

	o.Cab[floor] = false
	o.HallUp[floor] = false
	o.HallDown[floor] = false

	for b := elevio.ButtonType(0); b < 3; b++ {
		elevio.SetButtonLamp(b, floor, false)
	}
}

//!SECTION

// SECTION - Elev
type Elev struct {
	dir       int
	prevDir   int
	curFloor  int
	doorOpen  bool
	numFloors int
	orders    Orders
}

func (e *Elev) setup(numFloors int) {
	e.dir = 0
	e.prevDir = 0
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
	e.prevDir = e.dir
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
func (e *Elev) nextOrder() {

	if e.prevDir > 0 {
		for i := e.curFloor; i < e.numFloors; i++ {

			if e.orders.HallUp[i] || e.orders.Cab[i] {
				e.goUp()
				return
			}

		}
	}

	if e.prevDir < 0 {
		for i := e.curFloor; i > 0; i-- {

			if e.orders.HallDown[i] || e.orders.Cab[i] {
				e.goDown()
				return
			}

		}

	}

	for i := 0; i < e.numFloors; i++ {

		if e.orders.Cab[i] || e.orders.HallUp[i] || e.orders.HallDown[i] {
			if i < e.curFloor {
				e.goDown()
				return
			}
			if i > e.curFloor {
				e.goUp()
				return
			}
			// case when someone press the button of the floor where they currently are
			e.completeOrder(e.curFloor)
			return
		}

	}

}
func (e *Elev) completeOrder(floor int) bool {
	if floor == e.curFloor {

		tf, d := e.orders.checkOrder(floor)

		if tf {

			if d == 0 || d == e.dir {
				e.stop()
				e.openDoors()
				e.orders.completeOrder(floor)
				time.Sleep(Open_Door_Time * time.Second)
				e.nextOrder()
				e.closeDoors()

				return true

			}
		}
	}

	return false
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
			elev.orders.setOrder(a.Floor, a.Button)
			fmt.Println(elev.orders)

		case a := <-drv_floors:
			// fmt.Printf("--2:%+v\n", a)
			elev.updateFloor()
			elev.completeOrder(a)
			fmt.Println("current floor:", elev.getFloor())

			// for i, v := range orders.Cab {
			// 	if i == a && v == true {
			// 		orders.completeOrder(a, d)

			// 	} else if a == numFloors-1 {
			// 		d = elevio.MD_Down
			// 	} else if a == 0 {
			// 		d = elevio.MD_Up
			// 	}
			// 	elevio.SetMotorDirection(d)
			// }

		case a := <-drv_stop:
			if a {
				elev.stop()
			}

		case a := <-drv_obstr:
			if a {
				elev.stop()
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
