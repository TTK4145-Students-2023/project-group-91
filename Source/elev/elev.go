package elev

import (
	"Source/conf"
	"Source/elevio"
	"time"
)

type ElevMode int

const (
	Master ElevMode = 0
	Backup ElevMode = 1 // do we need that? or each slave will be a backup?
	Slave  ElevMode = 2
)

type Directions int

const (
	Up   Directions = 1
	Cab  Directions = 0
	Down Directions = -1
)

type Elev struct {
	Dir      int
	PrevDir  int
	CurFloor int
	DoorOpen bool
	Orders   Orders
	Mode     ElevMode
	ID       int
}

func (e *Elev) ChangeMode(mode ElevMode) {
	e.Mode = mode
	// "reboot" into different mode
}
func (e *Elev) Whoami() ElevMode {
	return e.Mode
}

func (e *Elev) Init() {
	// setup default values for elevator
	e.Dir = 0
	e.PrevDir = 0
	e.CurFloor = elevio.GetFloor()
	e.DoorOpen = false
	e.Mode = Slave
	e.Orders = Orders{
		HallUp:   make([]bool, conf.Num_Of_Flors),
		HallDown: make([]bool, conf.Num_Of_Flors),
		Cab:      make([]bool, conf.Num_Of_Flors)}
	e.Orders.ClearAll()

	if e.CurFloor == -1 {
		e.GoDown()
	}
}
func (e *Elev) UpdateFloor() int {
	e.CurFloor = elevio.GetFloor()
	elevio.SetFloorIndicator(e.CurFloor)

	if e.CurFloor == conf.Num_Of_Flors-1 || e.CurFloor == 0 {
		e.Stop()
	}

	return e.CurFloor
}
func (e Elev) GetDirection() int {
	return e.Dir
}
func (e Elev) GetFloor() int {
	return e.CurFloor
}
func (e *Elev) GoUp() bool {
	if e.DoorOpen {
		return false

	} else if e.CurFloor == conf.Num_Of_Flors-1 {
		return false
	} else {
		e.Dir = 1
		elevio.SetMotorDirection(elevio.MD_Up)
		return true
	}
}
func (e *Elev) GoDown() bool {
	if e.DoorOpen {
		return false

	} else if e.CurFloor == 0 {
		return false

	} else {
		e.Dir = -1
		elevio.SetMotorDirection(elevio.MD_Down)
		return true
	}
}
func (e *Elev) Stop() {
	e.PrevDir = e.Dir
	e.Dir = 0
	elevio.SetMotorDirection(elevio.MD_Stop)
}
func (e *Elev) OpenDoors() {
	e.DoorOpen = true
	elevio.SetDoorOpenLamp(true)
}
func (e *Elev) CloseDoors() bool {

	if elevio.GetObstruction() {
		e.DoorOpen = true
		return false
	} else {
		elevio.SetDoorOpenLamp(false)
		e.DoorOpen = false
		return true
	}

}
func (e *Elev) NextOrder() {
	// going for the next order
	if e.PrevDir > 0 {
		for i := e.CurFloor; i < conf.Num_Of_Flors; i++ {

			if e.Orders.HallUp[i] || e.Orders.Cab[i] {
				e.GoUp()
				return
			}

		}
	}

	if e.PrevDir < 0 {
		for i := e.CurFloor; i > 0; i-- {

			if e.Orders.HallDown[i] || e.Orders.Cab[i] {
				e.GoDown()
				return
			}

		}

	}

	for i := 0; i < conf.Num_Of_Flors; i++ {

		if e.Orders.Cab[i] || e.Orders.HallUp[i] || e.Orders.HallDown[i] {
			if i < e.CurFloor {
				e.GoDown()
				return
			}
			if i > e.CurFloor {
				e.GoUp()
				return
			}
			// case when someone press the button of the floor where they currently are
			e.ShouldIstop(e.CurFloor)
			return
		}

	}

}

func (e *Elev) CompleteOrder(floor int) bool {
	// stop, open doors, wait, close, go for next order
	e.Stop()
	e.OpenDoors()
	e.Orders.CompleteOrder(floor)
	time.Sleep(conf.Open_Door_Time * time.Second)
	for !e.CloseDoors() {

	}
	e.NextOrder()

	return true
}

func (e *Elev) ShouldIstop(floor int) bool {
	// return true if there is an order on current floor "for us" - in same direction or order was from cabin
	e.UpdateFloor()
	if floor == e.CurFloor {

		tf, d := e.Orders.CheckOrder(floor)

		if tf {

			//  cab order || not moving || same dir	  || there is no orders in
			//                          || as order   || curr dir
			if d == 0 || e.Dir == 0 || d == e.Dir {
				// return e.CompleteOrder(floor)
				go e.CompleteOrder(floor)
			}

			if e.Orders.FirstDown() > e.CurFloor {
				return false
			}

			//
			if e.Dir < 0 && e.Orders.FirstUp() == -1 {
				// return e.CompleteOrder(floor)
				go e.CompleteOrder(floor)

			}
			if e.Dir > 0 && e.Orders.FirstDown() == -1 {
				// return e.CompleteOrder(floor)
				go e.CompleteOrder(floor)

			}
			if e.Orders.FirstUp() == e.CurFloor || e.Orders.FirstDown() == e.CurFloor {

				// return e.CompleteOrder(floor)
				go e.CompleteOrder(floor)
			}

		}
	}

	return false
}
