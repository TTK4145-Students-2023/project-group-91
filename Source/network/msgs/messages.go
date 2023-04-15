package msgs

import (
	"Source/conf"
	"Source/elevator"
	"Source/elevio"
)

type Msg struct {
	SenderID   int
	SenderRole string
	Dir        conf.Directions
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
type MsgBackup struct {
	SenderID   int
	SenderRole string
	Elevs      []elevator.SemiElev
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

func PrepareBackupMsg(e elevator.Elev) MsgBackup {
	msg := MsgBackup{
		SenderID:   e.ID,
		SenderRole: e.GetMode(),
		Elevs:      e.Elevs,
	}
	return msg
}
