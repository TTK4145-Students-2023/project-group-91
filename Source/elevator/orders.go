package elevator

import (
	"Source/conf"
	"Source/elevio"
	"fmt"
)

type Orders struct {
	HallUp      []bool
	HallDown    []bool
	Cab         []bool
	NumOfOrders int
}

func (o *Orders) SetOrder(floor int, button elevio.ButtonType) {
	if button == 0 {
		o.HallUp[floor] = true
		elevio.SetButtonLamp(button, floor, true)
		o.NumOfOrders++

	} else if button == 2 {
		o.Cab[floor] = true
		elevio.SetButtonLamp(button, floor, true)
		o.NumOfOrders++

	} else if button == 1 {
		o.HallDown[floor] = true
		elevio.SetButtonLamp(button, floor, true)
		o.NumOfOrders++

	}
}

func (o *Orders) SetOrderTMP(floor int, button elevio.ButtonType) {
	if button == 0 {
		o.HallUp[floor] = true
		o.NumOfOrders++

	} else if button == 2 {
		o.Cab[floor] = true
		o.NumOfOrders++

	} else if button == 1 {
		o.HallDown[floor] = true
		o.NumOfOrders++

	}
}
func (o Orders) CheckOrder(floor int) (bool, conf.Directions) {
	// check of there is any order on current floor and return true/false and order type
	if floor == -1 {
		return false, 0
	}
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
func (o *Orders) ClearAll() {
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
	o.NumOfOrders = 0

}
func (o *Orders) AddOrders(orders Orders, ordType ...string) {
	up := false
	down := false
	cab := false

	for _, v := range ordType {
		if v[0] == 'U' || v[0] == 'u' {
			up = true
		}
		if v[0] == 'D' || v[0] == 'd' {
			down = true
		}
		if v[0] == 'C' || v[0] == 'c' {
			cab = true
		}

	}

	if up {

		for i := range o.HallUp {
			if orders.HallUp[i] {
				o.HallUp[i] = true
				o.NumOfOrders++
			}
		}

	}
	if down {
		for i := range o.HallDown {
			if orders.HallDown[i] {
				o.HallDown[i] = true
				o.NumOfOrders++
			}
		}

	}
	if cab {
		for i := range o.Cab {
			if orders.Cab[i] {
				o.Cab[i] = true
				o.NumOfOrders++
			}
		}

	}
}
func (o *Orders) CompleteOrder(floor int, dir conf.Directions, Elevs []SemiElev) conf.Directions {
	// fmt.Println("dir:", dir)
	nextDir := conf.None
	if dir < 0 {
		o.HallDown[floor] = false
		//FIXME e.Orders.HallDown[floor] = false
		elevio.SetButtonLamp(1, floor, false)
		nextDir = -1

	} else if dir > 0 {
		o.HallUp[floor] = false
		//FIXME e.Orders.HallUp[floor] = false
		elevio.SetButtonLamp(0, floor, false)
		nextDir = 1

	}
	if floor == conf.Num_Of_Flors-1 {
		o.HallDown[floor] = false
		//FIXME e.Orders.HallDown[floor] = false
		elevio.SetButtonLamp(1, floor, false)
		nextDir = -1
	}
	if floor == 0 {
		o.HallUp[floor] = false
		//FIXME e.Orders.HallUp[floor] = false
		elevio.SetButtonLamp(0, floor, false)
		nextDir = 1
	}

	o.Cab[floor] = false
	//FIXME e.Orders.Cab[floor] = false
	o.NumOfOrders--
	//FIXME e.Orders.NumOfOrders--
	elevio.SetButtonLamp(2, floor, false)

	// NOTE show the list of orders
	o.Print()
	return nextDir
}

func (o *Orders) IsAny(butType int) bool {
	if butType < 0 {
		for _, v := range o.HallDown {
			if v {
				return true
			}
		}
		return false
	}
	if butType > 0 {
		for _, v := range o.HallUp {
			if v {
				return true
			}
		}
		return false
	}
	if butType == 0 {
		for _, v := range o.Cab {
			if v {
				return true
			}
		}
		return false
	}
	return false
}

func (o Orders) FirstUp() int {
	for i, v := range o.HallUp {
		if v {
			return i
		}
	}
	return -1
}

func (o Orders) FirstDown() int {
	for i := len(o.HallDown) - 1; i > 0; i-- {
		if o.HallDown[i] {
			return i
		}
	}
	return -1
}

func (o Orders) CountOrders(ordType ...string) int {
	up := false
	down := false
	cab := false
	i := 0
	for _, v := range ordType {
		if v[0] == 'U' || v[0] == 'u' {
			up = true
		}
		if v[0] == 'D' || v[0] == 'd' {
			down = true
		}
		if v[0] == 'C' || v[0] == 'c' {
			cab = true
		}

	}
	if up {
		for _, v := range o.HallUp {
			if v {
				i++
			}
		}
	}
	if down {
		for _, v := range o.HallDown {
			if v {
				i++
			}
		}
	}
	if cab {
		for _, v := range o.Cab {
			if v {
				i++
			}
		}
	}

	return i

}

func (o Orders) Print() {
	fmt.Print("\n\t\t======Orders======\n")

	fmt.Print("Floor:\t\t")
	for i := 0; i < conf.Num_Of_Flors; i++ {
		fmt.Printf("[%v] ", i)
	}
	fmt.Print("\nHall up:\t")
	for _, v := range o.HallUp {
		if v {
			fmt.Print("[*] ")
		} else {
			fmt.Print("[ ] ")

		}
	}
	fmt.Print("\nHall Down:\t")
	for _, v := range o.HallDown {
		if v {
			fmt.Print("[*] ")
		} else {
			fmt.Print("[ ] ")

		}
	}
	fmt.Print("\nCab:\t\t")
	for _, v := range o.Cab {
		if v {
			fmt.Print("[*] ")
		} else {
			fmt.Print("[ ] ")

		}
	}
	fmt.Print("\n\t\t==================\n")
	fmt.Println()
}
