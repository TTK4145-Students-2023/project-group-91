package dist

import (
	"Source/elev"
)


type Manage struct{
	Elevators []elev.Elev
	NumOfElevators int
	
}

func (m *Manage) SetNumOfElevators(n int) {
	m.NumOfElevators = n
}

func (m *Manage) InitializeElevators(){
	for i := 0; i < m.NumOfElevators; i++{
		var elv elev.Elev
		m.Elevators[i] = elv
		elv.Init()
	}
}

func (m Manage) HowManyElevators() int{
	return m.NumOfElevators
}

func (m Manage) GetMode(index int) int{
	return int(m.Elevators[index].Mode)
}



func (m Manage) Slice(class elev.Orders) {
	// splits from orders to each elevator
	var o elev.Orders
	numButtons := 8
	size := numButtons / m.NumOfElevators 

	up := o.HallUp
	down := o.HallDown

	sliceUp1 := up[:size]
	sliceUp2 := up[size : size*2]
	sliceUp3 := up[size*2:]

	sliceDown1 := down[:size]
	sliceDown2 := down[size : size*2]
	sliceDown3 := down[size*2:]

	m.Elevators[0].Orders.HallUp = sliceUp1
	m.Elevators[1].Orders.HallUp = sliceUp2
	m.Elevators[2].Orders.HallUp = sliceUp3

	m.Elevators[0].Orders.HallDown = sliceDown1
	m.Elevators[1].Orders.HallDown = sliceDown2
	m.Elevators[2].Orders.HallDown = sliceDown3
}


func (m Manage) OnEqualNumOfOrders() bool{
	numOrd := m.Elevators[0].Orders.NumOfOrders
	isIt := true
	for i := range m.Elevators{
		if m.Elevators[i].Orders.NumOfOrders != numOrd{
			isIt = false
			break
		}
	}
	return isIt 
}


func (m Manage) OnZeroNumOfOrders() bool{
	temp := true
	for i := 0; i < m.NumOfElevators; i++ {
		if m.Elevators[i].Orders.NumOfOrders != 0{
			temp = false
			break
		}
	}
	return temp
}

func (m *Manage) WhoGetsOrder(class elev.Orders) {
	// decides by num of orders

}



func (m Manage) Choose(){
	// all on same floor and same amount of orders
}


// e1  e2  e3
// 4 _ 4 _ 4 
// 3 _ 3 _ 3 
// 2 _ 2 _ 2 
// 1 _ 1 _ 1 