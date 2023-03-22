// TODO - Fix door which stops the program
// TODO - serve orders only in one direction (dont complete each order on one floor)

// buglist:
// FIXME - sometimes when pressing buttons while the elevator is between two floors there is an error: "tried to set floor indicator to not existing floor"

package main

import (
	"Source/conf"
	"Source/elev"
	"Source/elevio"
	"Source/network/bcast"
	"Source/network/localip"
	"Source/network/peers"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const Open_Door_Time = 2
const Num_Of_Flors = 4

func IsMasterAlive(peers []string) bool {
	for _, v := range peers {
		if v[0] == 'M' {
			return true
		}
	}
	return false
}
func HowManyMasters(peers []string) int {
	i := 0
	for _, v := range peers {
		if v[0] == 'M' {
			i++
		}
	}
	return i
}
func MastersID(peers []string) string {

	for _, v := range peers {
		if v[0] == 'M' {
			return v[1:]
		}
	}
	return ""
}

func MaxIdAlive(peers []string) int {
	max := 0
	for _, v := range peers {
		x, _ := strconv.Atoi(v[1:])
		if x > max {
			max = x

		}
	}
	return max
}
func AmiAlone(peers []string) bool {
	if len(peers) == 1 {
		return true
	} else {
		return false
	}
}

type Msg struct {
	SenderID   int
	SenderRole string
	Message    string
}

func PrepareMsg(m string, e elev.Elev) Msg {
	// fmt.Println("FUNC:", e.GetID_I(), e.CurFloor, e.GetMode())
	msg := Msg{
		SenderID:   e.GetID_I(),
		SenderRole: e.GetMode(),
		Message:    m}

	return msg
}

func main() {

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
	elev := elev.Elev{}
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
			if IsMasterAlive(p.Peers) {
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

	// ----- Drivers channels -----------
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// ----- Drivers goroutines ---------
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// ----- Messaging goroutines ---------
	go bcast.Transmitter(Port_msgs, sendMsg)
	go bcast.Receiver(Port_msgs, rcvdMsg)

	for {
		select {

		case button := <-drv_buttons:

			elev.Orders.SetOrder(button.Floor, button.Button)

			if button.Floor == elevio.GetFloor() {
				if !elev.DoorOpen {

					go elev.CompleteOrder(elev.GetFloor())
				}

			} else if elev.Dir == 0 {

				go elev.NextOrder()
			}

		case floor := <-drv_floors:
			elev.UpdateFloor()
			elev.ShouldIstop(floor)
			fmt.Println("current floor:", elev.GetFloor())
			sendMsg <- PrepareMsg("Im on floor "+fmt.Sprint(elev.GetFloor()), elev)

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
					elev.NextOrder()
				}
			}

		// ------- Messages --------
		case m := <-rcvdMsg:
			if m.SenderID != elev.GetID_I() {
				if m.SenderRole == "M" {
					// ---- Messages from Master
					fmt.Println("m:	", m)

				} else {
					// ---- Messages from Others
					fmt.Println("SenderID:	", m.SenderID)
					fmt.Println("SenderRole:	", m.SenderRole)
					fmt.Println("Msg:		", m.Message)
				}

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
			if HowManyMasters(p.Peers) != 1 {
				if maxID := MaxIdAlive(p.Peers); maxID == elev.GetID_I() {
					elev.ChangeMode(conf.Master)
				} else {
					elev.ChangeMode(conf.Slave)

				}
				role_chan <- elev.GetMode()
			}

		}
	}
}
