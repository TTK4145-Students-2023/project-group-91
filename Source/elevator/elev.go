package elevator

import (
	"Source/conf"
	"Source/elevio"
	"fmt"
	"math"
	"time"
)

type SemiElev struct {
	// mini elevator structure to store other elevators data
	Dir       conf.Directions
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
	e.Orders.NumOfOrders = i
}

type Elev struct {
	// Out main elevator structure responible for everything
	Dir      conf.Directions
	PrevDir  conf.Directions
	NextDir  conf.Directions
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
func (e Elev) GetMeFromSemiElevs() SemiElev {
	// rturn a semiElev obj from database with our ID
	for i := range e.Elevs {
		if e.Elevs[i].ID == e.GetID_I() {
			return e.Elevs[i]
		}
	}
	return SemiElev{ID: -1}
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
func (e Elev) GetDirection() conf.Directions {
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
func (e *Elev) compleateElevsOrders(floor int, dir conf.Directions) {
	// remove orders from elevs database
	for i, el := range e.Elevs {
		if dir >= 1 {
			for f := range el.Orders.HallUp {
				if f == floor {
					e.Elevs[i].Orders.HallUp[f] = false
				}
			}
		}
		if dir <= -1 {
			for f := range el.Orders.HallDown {
				if f == floor {
					e.Elevs[i].Orders.HallDown[f] = false
				}
			}
		}
	}
}
func (e Elev) UpdateLightsSum() {
	// update ordder lighs so they will be syncronized among all elevators

	var lampsUp [conf.Num_Of_Flors]int
	var lampsDown [conf.Num_Of_Flors]int

	for _, el := range e.Elevs {
		for i, v := range el.Orders.HallUp {
			if v {
				lampsUp[i]++
			}
		}
		for i, v := range el.Orders.HallDown {
			if v {
				lampsDown[i]++
			}

		}
	}
	for i, v := range lampsUp {
		if v > 0 {
			elevio.SetButtonLamp(elevio.BT_HallUp, i, true)
		} else {

			elevio.SetButtonLamp(elevio.BT_HallUp, i, false)
		}
	}
	for i, v := range lampsDown {
		if v > 0 {
			elevio.SetButtonLamp(elevio.BT_HallDown, i, true)
		} else {
			elevio.SetButtonLamp(elevio.BT_HallDown, i, false)
		}
	}
}
func (e *Elev) MoveOn() {
	// Choosing next direction to move

	if e.NextDir == conf.Up && e.Orders.CountOrders("Up") > 0 && e.CurFloor != conf.Num_Of_Flors-1 {
		e.GoUp()
		return
	}

	if e.PrevDir == conf.Up {
		for i := e.CurFloor; i < conf.Num_Of_Flors; i++ {

			if e.Orders.HallUp[i] || e.Orders.Cab[i] {
				// fmt.Print("__1_GOUP__")
				e.GoUp()
				return
			}

		}
	}

	if e.NextDir == conf.Down && e.Orders.CountOrders("Down") > 0 && e.CurFloor != 0 {
		e.GoDown()
		return
	}

	if e.PrevDir == conf.Down {
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

			e.Orders.CompleteOrder(e.CurFloor, 1, e.Elevs)
			e.Orders.CompleteOrder(e.CurFloor, -1, e.Elevs)
			// fmt.Print("__4__")
			return
		}
	}
}

func (e *Elev) SleeperDetection(sleeperDetected chan bool) {
	// Goroutine to detect elevators that stuck and doing nothing bcs of uknown reason
	timer := 0
	for {
		time.Sleep(time.Millisecond * 1000)
		timer++
		if e.IsMoving() || e.NoOrders() {
			timer = 0
		}
		if timer >= conf.Detect_Sleeper_Time {
			sleeperDetected <- true
			timer = 0
		}

	}
}

func (e *Elev) CompleteOrder(floor int) bool {
	// stop, open doors, wait, close, go for next order
	// e.Orders.Print()
	e.UpdateLightsSum()

	e.Stop()

	if !e.DoorOpen {
		e.OpenDoors()
	}

	e.NextDir = e.Orders.CompleteOrder(floor, e.PrevDir, e.Elevs)
	e.compleateElevsOrders(floor, e.NextDir)
	time.Sleep(conf.Open_Door_Time * time.Second)

	if e.DoorOpen {
		e.CloseDoors()
	}

	e.MoveOn()
	e.UpdateLightsSum()
	return true
}
func (e *Elev) ShouldIstop(floor int) bool {
	// return true if there is an order on current floor "for us" - in same direction or order was from cabin
	e.UpdateFloor()
	if floor == e.CurFloor {

		orderExists, d := e.Orders.CheckOrder(floor)

		if orderExists {

			// if e.Orders.NumOfOrders == 1 {
			// 	return true
			// }
			if d != 0 && e.Dir != d && e.Orders.CountOrders() > 1 {
				return false
			}
			//  cab order || not moving || same dir	  || there is no orders in
			//                          || as order   || curr dir
			if d == 0 || e.Dir == 0 || d == e.Dir {
				return true
			}
			if e.Orders.FirstDown() > e.CurFloor {
				return false
			}

			//
			if e.Dir < 0 && e.Orders.FirstUp() == -1 {
				return true

			}
			if e.Dir > 0 && e.Orders.FirstDown() == -1 {
				return true

			}
			if e.Orders.FirstUp() == e.CurFloor || e.Orders.FirstDown() == e.CurFloor {

				return true
			}

		}
	}

	return false
}

func (e Elev) NoOrders() bool {
	count := 0
	for _, v := range e.Orders.Cab {
		if v {
			count++
		}
	}
	for _, v := range e.Orders.HallDown {
		if v {
			count++
		}
	}
	for _, v := range e.Orders.HallUp {
		if v {
			count++
		}
	}
	if count == 0 {
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
	for i, el := range e.Elevs {
		if el.ID == id {
			e.Elevs = append(e.Elevs[:i], e.Elevs[i+1:]...)
			return
		}
	}
}
func (e Elev) GiveElev(id int) SemiElev {

	for _, el := range e.Elevs {
		if el.ID == id {
			return el
		}
	}
	return SemiElev{ID: -1}
}

func (e *Elev) giveOrderTo(elevID int, floor int, dir conf.Directions, order bool) {

	if dir == conf.Up {

		for i := range e.Elevs {
			if e.Elevs[i].ID == elevID {
				e.Elevs[i].Orders.HallUp[floor] = order
				e.Orders.HallUp[floor] = !order
			}
		}
	} else if dir == conf.Down {

		for i := range e.Elevs {
			if e.Elevs[i].ID == elevID {
				e.Elevs[i].Orders.HallDown[floor] = order
				e.Orders.HallDown[floor] = !order
			}
		}
	}
}

func (e *Elev) DistributeOrders() {
	// not so great

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
					e.giveOrderTo(e.Elevs[el].ID, i, conf.Down, o)
					// e.Elevs[el].Orders.HallDown[i] = o
					// e.Orders.HallDown[i] = !o
					break
				}
				if e.Elevs[el].Dir == conf.Down {
					e.giveOrderTo(e.Elevs[el].ID, i, conf.Down, o)
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
					e.giveOrderTo(e.Elevs[el].ID, i, conf.Up, o)
					// e.Elevs[el].Orders.HallUp[i] = o
					// e.Orders.HallUp[i] = !o
					break
				}
				if e.Elevs[el].Dir == conf.Up {
					e.giveOrderTo(e.Elevs[el].ID, i, conf.Up, o)
					// e.Elevs[el].Orders.HallUp[i] = o
					// e.Orders.HallUp[i] = !o
					break
				}
				// }

			}
		}
	}

}

func (e *Elev) findClosestElevator(floor int) *SemiElev {
	minDist := conf.Num_Of_Flors + 1
	minOrders := conf.Num_Of_Flors + 1
	var closestElev *SemiElev

	for i := range e.Elevs {
		if e.Elevs[i].OrdersNum < conf.Max_Orders_Per_Elevator {
			dist := int(math.Abs(float64(e.Elevs[i].CurFloor - floor)))
			if dist < minDist || (dist == minDist && e.Elevs[i].OrdersNum < minOrders) {
				minDist = dist
				minOrders = e.Elevs[i].OrdersNum
				closestElev = &e.Elevs[i]
			}
		}
	}

	return closestElev
}

func (e *Elev) DistributeOrdersV2() {
	for floor := 0; floor < conf.Num_Of_Flors; floor++ {
		if e.Orders.HallUp[floor] || e.Orders.HallDown[floor] {
			// Find the closest elevator to the floor
			closestElev := e.findClosestElevator(floor)

			// Add the order to the closest elevator's orders
			if e.Orders.HallUp[floor] {
				closestElev.Orders.HallUp[floor] = true
			} else {
				closestElev.Orders.HallDown[floor] = true
			}
			closestElev.Orders.NumOfOrders++

			// Remove the order from this elevator's orders
			e.Orders.HallUp[floor] = false
			e.Orders.HallDown[floor] = false
			e.Orders.NumOfOrders--
		}
	}
}

func (e *Elev) DistributeOrdersV3() {
	// Loop over all hall up and down buttons and distribute the orders
	for i, btn := range e.Orders.HallUp {
		if btn {
			// Find the closest elevator moving in the upward direction
			elev := e.findClosestElevator(i)
			if elev != nil && (elev.Dir == conf.Up || elev.Dir == conf.None) {
				elev.Orders.HallUp[i] = true
				elev.Orders.NumOfOrders++
				elev.OrdersNum++
				e.Orders.HallUp[i] = false
			} else {
				// If there are no available elevators moving in the upward direction,
				e.Orders.HallUp[i] = true
				e.Orders.NumOfOrders++
			}
		}
	}
	for i, btn := range e.Orders.HallDown {
		if btn {
			// Find the closest elevator moving in the downward direction
			elev := e.findClosestElevator(i)
			if elev != nil && (elev.Dir == conf.Down || elev.Dir == conf.None) {
				elev.Orders.HallDown[i] = true
				elev.Orders.NumOfOrders++
				elev.OrdersNum++
				e.Orders.HallDown[i] = false
			} else {
				// If there are no available elevators moving in the downward direction,
				e.Orders.HallDown[i] = true
				e.Orders.NumOfOrders++
			}
		}
	}

	// Handle scenario when both hall up and hall down buttons are pressed on the same floor
	for i := range e.Orders.HallUp {
		if e.Orders.HallUp[i] && e.Orders.HallDown[i] {
			// Find the closest elevator moving in the upward direction
			elevUp := e.findClosestElevator(i)
			// Find the closest elevator moving in the downward direction
			elevDown := e.findClosestElevator(i)
			if elevUp != nil && (elevUp.Dir == conf.Up || elevUp.Dir == conf.None) {
				elevUp.Orders.HallUp[i] = true
				elevUp.Orders.NumOfOrders++
				elevUp.OrdersNum++
				e.Orders.HallUp[i] = false
			} else if elevDown != nil && (elevDown.Dir == conf.Down || elevDown.Dir == conf.None) {
				elevDown.Orders.HallDown[i] = true
				elevDown.Orders.NumOfOrders++
				elevDown.OrdersNum++
				e.Orders.HallDown[i] = false
			} else {
				// If there are no available elevators moving in the required directions,
				e.Orders.HallUp[i] = true
				e.Orders.HallDown[i] = true
				e.Orders.NumOfOrders += 2
			}
		}
	}
}
