package elevio

import (
	"fmt"
	"net"
)

func DirectMsg(msg string, ip string, port int) {

}
func Broadcast(msg string, port int) {

}
func ListenAll(port int, msg chan<- string, sndr_addr chan<- string, stop_chan chan bool) {

	addr := net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: port}
	conn, err := net.ListenUDP("udp", &addr)
	buf := make([]byte, 1024)
	stop := false

	if err != nil {
		fmt.Println("Listen-Conn error:", err)
		return
	}

	for !stop {
		select {
		case stop = <-stop_chan:

		default:
			bytes_num, sender_addr, err := conn.ReadFromUDP(buf)

			if err != nil {
				fmt.Println("Listen-Read error:", err)
				return
			}
			msg <- fmt.Sprint(buf[0:bytes_num])
			sndr_addr <- sender_addr.IP.String()

		}

	}

}
