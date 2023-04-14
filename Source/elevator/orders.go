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
func (o Orders) CheckOrder(floor int) (bool, int) {
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
func (o *Orders) CompleteOrder(floor int, dir int) {
	fmt.Println("dir:", dir)
	if dir < 0 {
		// fmt.Println(" down compleated")
		o.HallDown[floor] = false
		elevio.SetButtonLamp(1, floor, false)

	} else if dir > 0 {
		o.HallUp[floor] = false
		elevio.SetButtonLamp(0, floor, false)
		// fmt.Println(" up compleated")

	}
	if floor == conf.Num_Of_Flors-1 {
		o.HallDown[floor] = false
		elevio.SetButtonLamp(1, floor, false)
	}
	if floor == 0 {
		o.HallUp[floor] = false
		elevio.SetButtonLamp(0, floor, false)
	}

	// fmt.Println(" and cabin compleated")
	o.Cab[floor] = false
	// o.HallUp[floor] = false
	// o.HallDown[floor] = false
	o.NumOfOrders--
	// elevio.SetButtonLamp(0, floor, false)
	// elevio.SetButtonLamp(1, floor, false)
	elevio.SetButtonLamp(2, floor, false)
	o.Print()

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

func (o Orders) HowManyOrders() int {
	return o.NumOfOrders
}

func (o *Orders) UpdateLights() {
	for i, v := range o.HallUp {
		elevio.SetButtonLamp(elevio.BT_HallUp, i, v)
	}
	for i, v := range o.HallDown {
		elevio.SetButtonLamp(elevio.BT_HallDown, i, v)
	}
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
