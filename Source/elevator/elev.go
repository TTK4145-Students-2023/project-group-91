package elevator

import (
	"Source/conf"
	"Source/elevio"
	"fmt"
	"time"
)

type Directions int

const (
	Up   Directions = 1
	Cab  Directions = 0
	Down Directions = -1
)

type SemiElev struct {
	Dir       int
	CurFloor  int
	Orders    Orders
	OrdersNum int
	Mode      string
	Dist      int
	ID        int
}

func (e *SemiElev) CountOrders() {
	// count number of the orders without the cabin orders

	i := 0
	// for _, v := range e.Orders.Cab {
	// 	if v {
	// 		i++
	// 	}
	// }
	for _, v := range e.Orders.HallDown {
		if v {
			i++
		}
	}
	for _, v := range e.Orders.HallUp {
		if v {
			i++
		}
	}
}

type Elev struct {
	Dir      int
	PrevDir  int
	CurFloor int
	DoorOpen bool
	Orders   Orders
	Elevs    []SemiElev
	Mode     conf.ElevMode
	ID       int
}

func (e Elev) IsMoving() bool {
	if e.Dir == 0 {
		return false
	} else {
		return true
	}
}

func (e Elev) ImTheMaster() bool {
	if e.Mode == conf.Master {
		return true
	} else {
		return false
	}
}

func (e *Elev) ChangeMode(mode conf.ElevMode) {
	e.Mode = mode

	if e.Mode == conf.Master {
		fmt.Println("I'm The Master now!")
	} else if e.Mode == conf.Slave {
		fmt.Println("I'm Slave now")
	}
	// "reboot" into different mode
}
func (e Elev) Whoami() conf.ElevMode {
	return e.Mode
}
func (e Elev) GetMode() string {
	if e.Mode == 0 {
		return "M"
	} else if e.Mode == 1 {
		return "S"
	} else {
		return ""
	}
}
func (e Elev) GetID_I() int {
	// returns ID as int
	return e.ID
}
func (e Elev) GetID_S() string {
	// returns ID as string
	return fmt.Sprint(e.ID)
}
func (e *Elev) SetID(id int) {
	e.ID = id
}
func (e *Elev) Init() {
	// setup default values for elevator
	e.Dir = 0
	e.PrevDir = 0
	e.CurFloor = elevio.GetFloor()
	e.DoorOpen = false
	e.Mode = conf.Slave
	e.ID = -1
	e.Orders = Orders{
		HallUp:   make([]bool, conf.Num_Of_Flors),
		HallDown: make([]bool, conf.Num_Of_Flors),
		Cab:      make([]bool, conf.Num_Of_Flors)}
	e.Orders.ClearAll()

	if e.CurFloor == -1 {
		e.GoDown()
	}

	e.CloseDoors()

	e.Elevs = append(e.Elevs, SemiElev{
		ID:        e.ID,
		Dir:       e.Dir,
		CurFloor:  e.CurFloor,
		Orders:    e.Orders,
		OrdersNum: 0,
	})
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
				// fmt.Print("__1_GOUP__")
				e.GoUp()
				return
			}

		}
	}

	if e.PrevDir < 0 {
		for i := e.CurFloor; i > 0; i-- {

			if e.Orders.HallDown[i] || e.Orders.Cab[i] {
				// fmt.Print("__1_GODOWN__")
				e.GoDown()
				return
			}

		}

	}

	for i := 0; i < conf.Num_Of_Flors; i++ {

		if e.Orders.Cab[i] || e.Orders.HallUp[i] || e.Orders.HallDown[i] {
			if i < e.CurFloor {
				// fmt.Print("__2_GODOWN__")
				e.GoDown()
				return
			}
			if i > e.CurFloor {
				// fmt.Print("__2_GOUP__")
				e.GoUp()
				return
			}
			// case when someone press the button of the floor where they currently are

			e.Orders.CompleteOrder(e.CurFloor, 1)
			e.Orders.CompleteOrder(e.CurFloor, -1)
			// fmt.Print("__4__")

			return
		}

	}

}

func (e *Elev) CompleteOrder(floor int) bool {
	// stop, open doors, wait, close, go for next order
	// e.Orders.Print()
	e.Stop()

	if !e.DoorOpen {
		e.OpenDoors()
	}

	e.Orders.CompleteOrder(floor, e.PrevDir)
	time.Sleep(conf.Open_Door_Time * time.Second)

	if e.DoorOpen {
		e.CloseDoors()
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

			// if e.Orders.NumOfOrders == 1 {
			// 	return true
			// }

			if d != 0 && e.Dir != d && e.Dir != 0 && e.Orders.NumOfOrders > 1 {
				return false
			}
			//  cab order || not moving || same dir	  || there is no orders in
			//                          || as order   || curr dir
			if d == 0 || e.Dir == 0 || d == e.Dir {
				// return e.CompleteOrder(floor)
				// go e.CompleteOrder(floor)
				return true
			}
			//FIXME
			if e.Orders.FirstDown() > e.CurFloor {
				return false
			}

			//
			if e.Dir < 0 && e.Orders.FirstUp() == -1 {
				// return e.CompleteOrder(floor)
				// go e.CompleteOrder(floor)
				return true

			}
			if e.Dir > 0 && e.Orders.FirstDown() == -1 {
				// return e.CompleteOrder(floor)
				// go e.CompleteOrder(floor)
				return true

			}
			if e.Orders.FirstUp() == e.CurFloor || e.Orders.FirstDown() == e.CurFloor {

				// return e.CompleteOrder(floor)
				// go e.CompleteOrder(floor)
				return true
			}

		}
	}

	return false
}

func (e Elev) ShouldIstopGPT(floor int) bool {
	// If the elevator has an active order for the current floor, it should stop
	if e.Orders.Cab[floor] || e.Orders.HallUp[floor] || e.Orders.HallDown[floor] {
		return true
	}

	// Check if there are any orders in the direction that the elevator is going
	var ordersInDirection []int
	for f := floor + 1; f < conf.Num_Of_Flors; f++ {
		if e.Orders.HallUp[f] || e.Orders.Cab[f] {
			ordersInDirection = append(ordersInDirection, f)
		}
	}
	for f := floor - 1; f >= 0; f-- {
		if e.Orders.HallDown[f] || e.Orders.Cab[f] {
			ordersInDirection = append(ordersInDirection, f)
		}
	}

	if len(ordersInDirection) == 0 {
		// If there are no orders in the direction that the elevator is going, it should stop
		return true
	} else if len(ordersInDirection) == 1 {
		// If there is only one order in the direction that the elevator is going, it should stop
		return true
	} else {
		// If there are multiple orders in the direction that the elevator is going,
		// we need to check if we should stop at the current floor or continue to the next one

		// Check if there is an order on the current floor that is in the direction that the elevator is going
		orderExists, orderType := e.Orders.CheckOrder(floor)
		if orderExists && ((e.Dir == int(Up) && orderType == int(Up)) || (e.Dir == int(Down) && orderType == int(Down))) {
			return true
		}

		// Find the nearest order in the direction that the elevator is going
		var nearestOrder int
		if e.Dir == int(Up) {
			nearestOrder = ordersInDirection[0]
			for _, f := range ordersInDirection {
				if f < nearestOrder {
					nearestOrder = f
				}
			}
		} else if e.Dir == int(Down) {
			nearestOrder = ordersInDirection[0]
			for _, f := range ordersInDirection {
				if f > nearestOrder {
					nearestOrder = f
				}
			}
		}

		// Check if we should skip the current floor to go to the nearest order first
		if e.Dir == int(Up) && nearestOrder < floor {
			return false
		} else if e.Dir == int(Down) && nearestOrder > floor {
			return false
		} else {
			return true
		}
	}
}

func (e Elev) NoOrders() bool {
	if e.Orders.HowManyOrders() == 0 {
		return true
	} else {
		return false
	}
}
func (e *Elev) AddElev(a SemiElev) {
	for i, el := range e.Elevs {
		if el.ID < 0 {
			e.Elevs[i] = a
			e.Elevs[i].Orders = a.Orders
			e.Elevs[i].CountOrders()
			return
		}
	}
	e.Elevs = append(e.Elevs, a)
	e.Elevs[len(e.Elevs)-1].CountOrders()
}
func (e *Elev) UpdateElev(id int, a SemiElev) {
	for i := 0; i < len(e.Elevs); i++ {
		if e.Elevs[i].ID == id {
			e.Elevs[i] = a
			e.Elevs[i].Orders = a.Orders
			e.Elevs[i].CountOrders()
		}
	}
}

func (e *Elev) RemElev(id int) {
	for i := 0; i < len(e.Elevs); i++ {
		if e.Elevs[i].ID == id {
			e.Elevs[i].ID = -1
			e.Elevs[i].CountOrders()
		}
	}
}

func (e *Elev) giveOrderTo(elevID int, floor int, dir Directions, order bool) {

	if dir == Up {

		for i := range e.Elevs {
			if e.Elevs[i].ID == elevID {
				e.Elevs[i].Orders.HallUp[floor] = order
				e.Orders.HallUp[floor] = !order
			}
		}
	} else if dir == Down {

		for i := range e.Elevs {
			if e.Elevs[i].ID == elevID {
				e.Elevs[i].Orders.HallDown[floor] = order
				e.Orders.HallDown[floor] = !order
			}
		}
	}
}

// TODO - Distribute order sytem:
func (e *Elev) DistributeOrders() {

	n := len(e.Elevs)

	for i, o := range e.Orders.HallDown {
		if o {
			minDist := conf.Num_Of_Flors + 1
			maxOrdNum := 0
			for el := 0; el < n; el++ {

				if e.Elevs[el].ID >= 0 {
					dist := i - e.Elevs[el].CurFloor
					if dist < 0 {
						dist *= -1
					}
					if dist < minDist {
						minDist = dist
					}
					e.Elevs[el].Dist = dist

					if maxOrdNum < e.Elevs[el].OrdersNum {
						maxOrdNum = e.Elevs[el].OrdersNum
					}

				}
			}
			for el := 0; el < n; el++ {
				// if maxOrdNum != 0 && el.OrdersNum < maxOrdNum {

				if e.Elevs[el].Dist == minDist {
					e.giveOrderTo(e.Elevs[el].ID, i, Down, o)
					// e.Elevs[el].Orders.HallDown[i] = o
					// e.Orders.HallDown[i] = !o
					break
				}
				if e.Elevs[el].Dir == int(Down) {
					e.giveOrderTo(e.Elevs[el].ID, i, Down, o)
					// e.Elevs[el].Orders.HallDown[i] = o
					// e.Orders.HallDown[i] = !o
					break
				}
				// }

			}
		}
	}
	// ---------------
	for i, o := range e.Orders.HallUp {
		if o {
			minDist := conf.Num_Of_Flors + 1
			maxOrdNum := 0
			for el := 0; el < n; el++ {

				if e.Elevs[el].ID >= 0 {
					dist := i - e.Elevs[el].CurFloor
					if dist < 0 {
						dist *= -1
					}
					if dist < minDist {
						minDist = dist
					}
					e.Elevs[el].Dist = dist

					if maxOrdNum < e.Elevs[el].OrdersNum {
						maxOrdNum = e.Elevs[el].OrdersNum
					}

				}
			}
			for el := 0; el < n; el++ {
				// if maxOrdNum != 0 && el.OrdersNum < maxOrdNum {

				if e.Elevs[el].Dist == minDist {
					e.giveOrderTo(e.Elevs[el].ID, i, Up, o)
					// e.Elevs[el].Orders.HallUp[i] = o
					// e.Orders.HallUp[i] = !o
					break
				}
				if e.Elevs[el].Dir == int(Up) {
					e.giveOrderTo(e.Elevs[el].ID, i, Up, o)
					// e.Elevs[el].Orders.HallUp[i] = o
					// e.Orders.HallUp[i] = !o
					break
				}
				// }

			}
		}
	}

	// for el := 0; el < n; el++ {
	// 	if e.Elevs[el].ID == e.GetID_I() {
	// 		e.Orders.HallDown = e.Elevs[el].Orders.HallDown
	// 		e.Orders.HallUp = e.Elevs[el].Orders.HallUp
	// 	}
	// }
}

// func (elev *Elev) DistributeOrdersGPT() {
// 	// Create a map to keep track of which orders have already been assigned to an elevator
// 	assignedOrders := make(map[int]bool)

// 	// Iterate over all elevators and their orders
// 	for i := 0; i < len(elev.Elevs); i++ {
// 		for j := 0; j < len(elev.Elevs[i].Orders.HallUp); j++ {
// 			// Check if the current order has already been assigned to an elevator
// 			if !assignedOrders[j] && elev.Elevs[i].Orders.HallUp[j] {
// 				// Calculate the distance between the elevator and the order's floor
// 				dist := int(math.Abs(float64(elev.Elevs[i].CurFloor - j)))

// 				// Update the elevator's order list and mode
// 				elev.Elevs[i].OrdersNum++
// 				elev.Elevs[i].Orders.HallUp[j] = true
// 				elev.Elevs[i].Dist += dist
// 				// elev.Elevs[i].Mode = "Idle"

// 				// Mark the order as assigned
// 				assignedOrders[j] = true
// 			}
// 		}
// 		for j := 0; j < len(elev.Elevs[i].Orders.HallDown); j++ {
// 			// Check if the current order has already been assigned to an elevator
// 			if !assignedOrders[j] && elev.Elevs[i].Orders.HallDown[j] {
// 				// Calculate the distance between the elevator and the order's floor
// 				dist := int(math.Abs(float64(elev.Elevs[i].CurFloor - j)))

// 				// Update the elevator's order list and mode
// 				elev.Elevs[i].OrdersNum++
// 				elev.Elevs[i].Orders.HallDown[j] = true
// 				elev.Elevs[i].Dist += dist
// 				elev.Elevs[i].Mode = "Idle"

// 				// Mark the order as assigned
// 				assignedOrders[j] = true
// 			}
// 		}
// 	}

// 	// Iterate over all unassigned orders and assign them to the nearest elevator
// 	for i := 0; i < len(elev.Orders.HallUp); i++ {
// 		if !assignedOrders[i] && elev.Orders.HallUp[i] {
// 			// Find the elevator with the shortest distance to the order's floor
// 			shortestDist := math.MaxInt32
// 			var closestElev *SemiElev

// 			for j := 0; j < len(elev.Elevs); j++ {
// 				dist := int(math.Abs(float64(elev.Elevs[j].CurFloor - i)))
// 				if dist < shortestDist {
// 					shortestDist = dist
// 					closestElev = &elev.Elevs[j]
// 				}
// 			}

// 			// Update the elevator's order list and mode
// 			closestElev.OrdersNum++
// 			closestElev.Orders.HallUp[i] = true
// 			closestElev.Dist += shortestDist
// 			closestElev.Mode = "Idle"

// 			// Mark the order as assigned
// 			assignedOrders[i] = true
// 		}
// 	}

// 	for i := 0; i < len(elev.Orders.HallDown); i++ {
// 		if !assignedOrders[i] && elev.Orders.HallDown[i] {
// 			// Find the elevator
// 			// Find the elevator with the shortest distance to the order's floor
// 			shortestDist := math.MaxInt32
// 			var closestElev *SemiElev

// 			for j := 0; j < len(elev.Elevs); j++ {
// 				dist := int(math.Abs(float64(elev.Elevs[j].CurFloor - i)))
// 				if dist < shortestDist {
// 					shortestDist = dist
// 					closestElev = &elev.Elevs[j]
// 				}
// 			}

// 			// Update the elevator's order list and mode
// 			closestElev.OrdersNum++
// 			closestElev.Orders.HallDown[i] = true
// 			closestElev.Dist += shortestDist
// 			closestElev.Mode = "Idle"

// 			// Mark the order as assigned
// 			assignedOrders[i] = true
// 		}
// 	}
// }
