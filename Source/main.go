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

func MaxIdAlive(peers []string) int {
	max := 0
	for _, p := range peers {
		x, _ := strconv.Atoi(p)
		if x > max {
			max = x

		}
	}
	return max
}

func main() {

	// Initialization

	const Port_msgs = 20000

	var id string
	var port string
	flag.StringVar(&id, "id", "", "id of this peer")         // custom id for localhost
	flag.StringVar(&port, "port", "15657", "simulator port") //custom port for localhost
	flag.Parse()

	//SECTION - Setting elevs ID

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

	//!SECTION

	fmt.Println("ID:", id)
	fmt.Println("PORT:", port)

	isMasterAlive := false

	elevio.Init("localhost:"+port, Num_Of_Flors)

	elev := elev.Elev{}
	elev.Init()
	elev.SetID(iid)
	// SETUP

	//SECTION - Setting elevs Role

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable_1 := make(chan bool)
	go peers.Receiver(15647, peerUpdateCh)
	go peers.Transmitter(15647, elev.GetID_S(), peerTxEnable_1)

	timer := 0

	for timer < conf.Wait_For_Master_Time && !isMasterAlive {

		select {
		case p := <-peerUpdateCh:
			if MaxIdAlive(p.Peers) != elev.GetID_I() {
				isMasterAlive = true
				elev.ChangeMode(conf.Slave)
			}
		default:
			time.Sleep(1 * time.Second)
			timer++
			println(timer, "second")

		}
	}

	if timer >= conf.Wait_For_Master_Time && !isMasterAlive {
		elev.ChangeMode(conf.Master)

	}
	peerTxEnable_1 <- false
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, elev.GetID_S(), peerTxEnable)

	//!SECTION

	//Messages channels
	sendMsg := make(chan string)
	rcvdMsg := make(chan string)
	//drivers channels
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// drivers goroutines
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// mesages
	go bcast.Transmitter(Port_msgs, sendMsg)
	go bcast.Receiver(Port_msgs, rcvdMsg)

	for {
		select {

		case button := <-drv_buttons:

			elev.Orders.SetOrder(button.Floor, button.Button)

			if button.Floor == elevio.GetFloor() {
				go elev.CompleteOrder(elev.GetFloor())

			} else if elev.Dir == 0 {

				go elev.NextOrder()
			}

		case floor := <-drv_floors:
			elev.UpdateFloor()
			elev.ShouldIstop(floor)
			fmt.Println("current floor:", elev.GetFloor())

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

		//messaging
		case p := <-peerUpdateCh:

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Elevs alive:    %q\n", p.Peers)
			// fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Disconnected Elev:     %q\n", p.Lost)

			if !IsMasterAlive(p.Peers) {
				maxID := MaxIdAlive(p.Peers)
				// fmt.Println("Max ID alive:", maxIDMode, maxID)
				if maxID == elev.GetID_I() {
					elev.ChangeMode(conf.Master)
				}
			}

		}
	}
}
