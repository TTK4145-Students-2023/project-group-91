package elevio

import (
	"time"
)

type Elev struct {
	Dir      int
	PrevDir  int
	CurFloor int
	DoorOpen bool
	Orders   Orders
}

func (e *Elev) Init() {
	e.Dir = 0
	e.PrevDir = 0
	e.CurFloor = GetFloor()
	e.DoorOpen = false
	e.Orders = Orders{
		make([]bool, Num_Of_Flors),
		make([]bool, Num_Of_Flors),
		make([]bool, Num_Of_Flors)}
	e.Orders.ClearAll()

	if e.CurFloor == -1 {
		e.GoDown()
	}
}
func (e *Elev) UpdateFloor() int {
	e.CurFloor = GetFloor()
	SetFloorIndicator(e.CurFloor)

	if e.CurFloor == Num_Of_Flors-1 || e.CurFloor == 0 {
		e.Stop()
	}

	return e.CurFloor
}
func (e *Elev) SetFloor(f int) {
	e.CurFloor = f
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

	} else if e.CurFloor == Num_Of_Flors-1 {
		return false
	} else {
		e.Dir = 1
		SetMotorDirection(MD_Up)
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
		SetMotorDirection(MD_Down)
		return true
	}
}
func (e *Elev) Stop() {
	e.PrevDir = e.Dir
	e.Dir = 0
	SetMotorDirection(MD_Stop)
}
func (e *Elev) OpenDoors() {
	e.DoorOpen = true
	SetDoorOpenLamp(true)
}
func (e *Elev) CloseDoors() bool {

	if GetObstruction() {
		e.DoorOpen = true
		return false
	} else {
		SetDoorOpenLamp(false)
		e.DoorOpen = false
		return true
	}

}
func (e *Elev) NextOrder() {

	if e.PrevDir > 0 {
		for i := e.CurFloor; i < Num_Of_Flors; i++ {

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

	for i := 0; i < Num_Of_Flors; i++ {

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
			e.CheckOrder(e.CurFloor)
			return
		}

	}

}

func (e *Elev) CompleteOrder(floor int) bool {
	e.Stop()
	e.OpenDoors()
	e.Orders.CompleteOrder(floor)
	time.Sleep(Open_Door_Time * time.Second)
	e.CloseDoors()
	e.NextOrder()

	return true
}

func (e *Elev) CheckOrder(floor int) bool {
	e.UpdateFloor()
	if floor == e.CurFloor {

		tf, d := e.Orders.CheckOrder(floor)

		if tf {

			//  cab order || not moving || same dir	  || there is no orders in
			//                             as order    	 curr dir
			if d == 0 || e.Dir == 0 || d == e.Dir {
				return e.CompleteOrder(floor)
			}

			if e.Orders.FirstDown() > e.CurFloor {
				return false
			}

			if e.Dir < 0 && e.Orders.FirstUp() == -1 {
				return e.CompleteOrder(floor)

			}
			if e.Dir > 0 && e.Orders.FirstDown() == -1 {
				return e.CompleteOrder(floor)

			}
			if e.Orders.FirstUp() == e.CurFloor || e.Orders.FirstDown() == e.CurFloor {

				return e.CompleteOrder(floor)
			}

		}
	}

	return false
}
