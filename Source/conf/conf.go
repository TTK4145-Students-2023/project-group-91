package conf

const Open_Door_Time = 2
const Num_Of_Flors = 4
const Wait_For_Master_Time = 2
const Detect_Sleeper_Time = 5
const Wait_For_Order_Time = 5
const Port_msgs = 2000
const Update_Time_Interval = 500 // [in miliseconds]

type ElevMode int

var (
	Master ElevMode = 0
	Slave  ElevMode = 1
)

type Directions int

const (
	Up   Directions = 1
	Cab  Directions = 0
	Down Directions = -1
	None Directions = 0
)
