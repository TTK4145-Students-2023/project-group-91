/* TODO - 	distribute orders system, for now it gives order only by distance to the order (not good),
go to */ /*DistributeOrders()*/ /* to see how does it work.

 */
/* TODO - 	connect the main with backup
Master elevator have a array field called Elevs with all nessesary data of others (and himself also) elevators
it stores their orders, roles, id, direction etc Elevs is an array of SemiElev struct type to have less data than
normal Elev struct
*/

// buglist:
// FIXME[epic=bugs] - 	sometimes when pressing buttons while the elevator is between two floors there is an error: "core.exception.AssertError@src/sim_server.d(536): Tried to set floor indicator to non-existent floor 255
// 						std.concurrency.OwnerTerminated@std/concurrency.d(236): Owner terminated
// 						std.concurrency.OwnerTerminated@std/concurrency.d(236): Owner terminated"

// FIXME[epic=bugs] - sometimes the orders (espesially on the last floor) is served but when elev will move somewhere else the light is lighting up again on its own
package main

import (
	"Source/conf"
	"Source/elevator"
	"Source/elevio"
	"Source/network/bcast"
	"Source/network/localip"
	"Source/network/peers"
	"Source/roles"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const Open_Door_Time = 2
const Num_Of_Flors = 4

type Msg struct {
	SenderID   int
	SenderRole string
	Dir        int
	Floor      int
	Orders     elevator.Orders
	Message    string
}
type MsgOrder struct {
	SenderID   int
	SenderRole string
	BType      elevio.ButtonType
	BFloor     int
	ReciverID  int
}
type MsgOrders struct {
	SenderID   int
	SenderRole string
	Orders     elevator.Orders
	ReciverID  int
}

func PrepareMsg(m string, e elevator.Elev) Msg {
	// preparing message as a special struct
	msg := Msg{
		SenderID:   e.GetID_I(),
		SenderRole: e.GetMode(),
		Floor:      e.GetFloor(),
		Dir:        e.GetDirection(),
		Orders:     e.Orders,
		Message:    m}

	return msg
}

func PrepareMsgOrder(floor int, button elevio.ButtonType, e elevator.Elev, recivID int) MsgOrder {
	msg := MsgOrder{
		SenderID:   e.GetID_I(),
		SenderRole: e.GetMode(),
		BType:      button,
		BFloor:     floor,
		ReciverID:  recivID}

	return msg
}
func PrepareMsgOrders(e elevator.Elev, o elevator.Orders, recivID int) MsgOrders {
	msg := MsgOrders{
		SenderID:   e.GetID_I(),
		SenderRole: e.GetMode(),
		Orders:     o,
		ReciverID:  recivID}

	return msg
}

func main() {
	// runtime.GOMAXPROCS(10)
	// Initialization

	const Port_msgs = 20000

	var id string
	var port string
	flag.StringVar(&id, "id", "", "id of this peer")         // custom id for localhost
	flag.StringVar(&port, "port", "15657", "simulator port") //custom port for localhost
	flag.Parse()

	//SECTION ---- Setting elevs ID -----
	// our ID is based on the last IP's octet
	// e.g. 192.168.0.35
	// 				   ^----[35]---- this is our ID

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		idt := strings.Split(localIP, ".")
		id = idt[3]

	}
	iid, _ := strconv.Atoi(id)

	//!SECTION -----------------------

	fmt.Println("ID:", id)
	fmt.Println("PORT:", port)

	elevio.Init("localhost:"+port, Num_Of_Flors)

	// ----- creating elev struct and initialization -----
	elev := elevator.Elev{}
	elev.Init()
	elev.SetID(iid)

	//SECTION ----- Setting elevs Role -----

	isMasterAlive := false

	role_chan := make(chan string, 1)
	peerUpdateCh := make(chan peers.PeerUpdate)

	go peers.Receiver(15647, peerUpdateCh)

	timer := 0
	// waiting for any signal from Master ...
	for timer < conf.Wait_For_Master_Time && !isMasterAlive {

		select {
		case p := <-peerUpdateCh:
			if roles.IsMasterAlive(p.Peers) {
				isMasterAlive = true
				elev.ChangeMode(conf.Slave)
			}
		default:
			time.Sleep(1 * time.Second)
			timer++
			println(timer, "second")

		}
	}

	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, elev.GetID_S(), role_chan, peerTxEnable)

	// ... if there is no signal after specific time it means that there is no master
	// so I can become a master

	if timer >= conf.Wait_For_Master_Time && !isMasterAlive {
		elev.ChangeMode(conf.Master)
	}

	role_chan <- elev.GetMode()

	//!SECTION -----------------------

	// ----- Messages channels ---------
	sendMsg := make(chan Msg)
	rcvdMsg := make(chan Msg)

	sendOrderChan := make(chan MsgOrder)
	rcvdOrderChan := make(chan MsgOrder)
	sendOrdersChan := make(chan MsgOrders)
	rcvdOrdersChan := make(chan MsgOrders)

	// ----- Drivers channels -----------
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// ------ Other channels ------------
	sleeperDetected := make(chan bool)
	// ----- Drivers goroutines ---------
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// ----- Messaging goroutines ---------
	go bcast.Transmitter(Port_msgs, sendMsg, sendOrderChan, sendOrdersChan)
	go bcast.Receiver(Port_msgs, rcvdMsg, rcvdOrderChan, rcvdOrdersChan)

	// ----- Other goroutines ---------
	go elev.SleeperDetection(sleeperDetected)

	for {
		select {

		// SECTION ---- Button is pressed ----
		case button := <-drv_buttons:

			if elev.NoOrders() && button.Floor == elev.CurFloor { // if the button pressed is our only order

				elev.Orders.CompleteOrder(button.Floor, 1, elev.GetMeFromSemiElevs())
				elev.Orders.CompleteOrder(button.Floor, -1, elev.GetMeFromSemiElevs())

			} else if button.Button == elevio.BT_Cab { // if the order is cabin order add it to OUR orders

				elev.Orders.SetOrder(button.Floor, button.Button)

				if elev.Dir == 0 {
					go elev.MoveOn()
				}

			} else /*if !elev.ImTheMaster()*/ { // else send order to master (if it is the master it will send it to itself)

				sendOrderChan <- PrepareMsgOrder(button.Floor, button.Button, elev, -1)

			} /*else {
			elev.Orders.SetOrder(button.Floor, button.Button)

			} */
			sendMsg <- PrepareMsg("U", elev) // sending updating msg about our state
			//!SECTION ----------

			//SECTION ---- When arrive on the floor ----
		case floor := <-drv_floors:

			//NOTE just printing some debugging stuff
			fmt.Println("elevs:")
			for _, e := range elev.Elevs {
				fmt.Println(e.ID)
			}
			elev.UpdateFloor()
			if elev.ShouldIstop(floor) { // check if it should stop to serve the order, and serve it if yes
				go elev.CompleteOrder(floor)
			}

			//NOTE just printing some debugging stuff
			fmt.Println("current floor:", elev.GetFloor())
			elev.Orders.Print()

			sendMsg <- PrepareMsg("U", elev) // sending updating msg about our state
		//!SECTION ---------

		// SECTION ------- Messaging -------
		// SECTION ---- Recived single order msg ----
		case o := <-rcvdOrderChan:

			if elev.ImTheMaster() { // master got an order from other elev
				elev.Orders.SetOrderTMP(o.BFloor, o.BType) // add it to its orders (without activation)
				elev.DistributeOrders()                    //TODO distribute orders among elevs
				// elev.DistributeOrdersGPT()
				for _, e := range elev.Elevs { // send distributed orders to all elevs
					sendOrdersChan <- PrepareMsgOrders(elev, e.Orders, e.ID)
				}

			}

			// if o.ReciverID == elev.ID { // we got a message for us
			// 	elev.Orders.SetOrder(o.BFloor, o.BType)
			// 	go elev.NextOrder()

			// }
			sendMsg <- PrepareMsg("U", elev) // sending updating msg about our state
		// !SECTION

		// SECTION ---- Recived orders obj msg ---
		case ors := <-rcvdOrdersChan:

			fmt.Println("ReciverID:", ors.ReciverID)
			fmt.Println("ElevID:", elev.GetID_I())

			if ors.ReciverID == elev.GetID_I() { // check if the message is for us (based on id)
				// fmt.Println("HALOOO")
				// fmt.Println("ors")
				// ors.Orders.Print()
				elev.Orders.HallUp = ors.Orders.HallUp
				elev.Orders.HallDown = ors.Orders.HallDown
				elev.Orders.UpdateLights()

				// checking if we got some orders to compleate
				if !elev.IsMoving() {

					if elev.Orders.HallUp[elev.GetFloor()] /*&& elev.Orders.NumOfOrders == 1*/ {

						go elev.CompleteOrder(elev.GetFloor())

					} else if elev.Orders.HallDown[elev.GetFloor()] /*&& elev.Orders.NumOfOrders == 1*/ {

						go elev.CompleteOrder(elev.GetFloor())

					} else if elev.GetDirection() == 0 {

						go elev.MoveOn()

					}
				}
				//NOTE printing debbuging list of orders
				elev.Orders.Print()
			}
		// !SECTION

		// SECTION ---- Recived msg ----
		case m := <-rcvdMsg:

			// adding elevator from network to local database of elevators
			if elev.ImTheMaster() {
				alreadyInNet := false
				for i := 0; i < len(elev.Elevs); i++ {
					if elev.Elevs[i].ID == m.SenderID {
						alreadyInNet = true
					}
				}
				if !alreadyInNet { // adding new elevator to the network

					elev.AddElev(elevator.SemiElev{
						ID:       m.SenderID,
						Mode:     m.SenderRole,
						Dir:      m.Dir,
						CurFloor: m.Floor,
						Orders:   m.Orders})

				} else if m.Message[0] == 'U' { // updating current elevator

					elev.UpdateElev(m.SenderID, elevator.SemiElev{
						ID:       m.SenderID,
						Mode:     m.SenderRole,
						Dir:      m.Dir,
						CurFloor: m.Floor,
						Orders:   m.Orders})
				}
			}

		// if m.SenderID != elev.GetID_I() {
		// 	if m.SenderRole == "M" {
		// 		// ---- Messages from Master
		// 		// fmt.Println("m:	", m)

		// 	} else {
		// 		// ---- Messages from Others
		// 		fmt.Println("SenderID:	", m.SenderID)
		// 		fmt.Println("SenderRole:	", m.SenderRole)
		// 		fmt.Println("Msg:		", m.Message)
		// 	}

		// }
		// !SECTION
		// !SECTION

		case s := <-sleeperDetected:
			if s {
				fmt.Println("HALo")
				elev.MoveOn()

			}
		// ----------- peer system, M/S control -------------
		case p := <-peerUpdateCh:

			// ---- network info ----

			// fmt.Printf("Peer update:\n")
			fmt.Printf("  Elevs alive:    %q\n", p.Peers)
			// fmt.Printf("  New:      %q\n", p.New)
			// fmt.Printf("  Disconnected Elev:     %q\n", p.Lost)

			// ----------------------

			// if there is 0 or more than one masters it means that we have to choose the new one
			// new master is choosen based on the highest ID in the nwtwork
			// if there is more masters rest will be degraded
			if roles.HowManyMasters(p.Peers) != 1 {
				if maxID := roles.MaxIdAlive(p.Peers); maxID == elev.GetID_I() {
					elev.ChangeMode(conf.Master)
				} else {
					elev.ChangeMode(conf.Slave)

				}
				role_chan <- elev.GetMode()
			}

			for _, id := range p.Lost {
				for _, e := range elev.Elevs {
					if len(id) > 1 {

						if idtmp, _ := strconv.Atoi(id[0:]); e.ID == idtmp {
							if elev.ImTheMaster() {
								elev.RemElev(e.ID)
								elev.Orders.AddOrders(e.Orders)
								fmt.Println("HHHHHHHHHHHHHHHHHHH")
								e.Orders.Print()
								elev.DistributeOrders()
							}
						}
					}
				}
			}

		case stop := <-drv_stop:
			if stop {
				elev.Stop()
			}

		case obstr := <-drv_obstr:
			if obstr && elev.DoorOpen {

				elev.Stop()

			} else if elev.DoorOpen {

				elev.CloseDoors()

				if !elev.ShouldIstop(elev.CurFloor) {
					elev.MoveOn()
				}
			}

		}
	}
}
