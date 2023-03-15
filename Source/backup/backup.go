package backup

import (

    
    "fmt"
    
)




type Orders struct{
	
	HallUp []bool
	
	HallDown []bool
	
	Cab []bool
}

var orderList []Orders

var HallUpValueList []bool 
var HallDownValueList []bool
var CabValueList []bool

var uplist []bool
var downlist []bool
var cablist []bool


//Takes in an Order object and put it in a list
func structList(p Orders){


	orderList = append(orderList, p)	
}


//The get-functions appends the values in the different values of the Struct , to list to separate lists : "uplist, downlist, cablist". 
//Then it prints the values with their index
func getHallDown(){


	for _, j := range orderList{		

		downlist = append(j.HallDown)
	}

	fmt.Println("HallDown :")

	for i, k := range downlist{

		fmt.Println(i, k)
	}

	fmt.Println("------------------------------------")


}

func getHallUp(){


	for _, j := range orderList{

	

		uplist = append(j.HallUp)
	}

	fmt.Println("HallUp: ")

	for i, k := range uplist{

		fmt.Println( i, k)
	}

	fmt.Println("------------------------------------")
}

func getCab(){


	for _, j := range orderList{

		

		cablist = append(j.Cab)
	}

	fmt.Println("Cab :")


	for i, k := range cablist{

		fmt.Println(i, k)
	}

	fmt.Println("------------------------------------")
}



func printList(list []Orders)  {

	fmt.Println(orderList)
	getHallDown()
	getHallUp()
	getCab()

}



func clearlist(){

	orderList := []Orders{}

	fmt.Println(orderList)

}
func clearuplist(){

	uplist = []bool{}

}

func clearcablist(){

	cablist = []bool{}

}


func printsoloList(d []bool){


	fmt.Println(d)
}









/*func main(){
	

	HallUpValueList = append(HallUpValueList, true, false, true, false, true)
	
	HallDownValueList = append(HallDownValueList, false, false, true, true)
	
	CabValueList = append(CabValueList, true, true, false, false, false)
	
	
	
	x := Orders{ HallUpValueList, HallDownValueList, CabValueList }


	



	





	structList(x)

	
	
	printList(orderList)

	printsoloList(downlist)
	clearlist()
	//clearuplist()
	//clearcablist()

	printList(orderList)
	fmt.Println("Testing")*/
	/*---------------------------------------------------------------------------------------
	-----------------------------------------------------------------------------------------
	-----------------------------------------------------------------------------------------

	Sending list throug UDP*/


	// encode the list as a byte array
    /*var buffer bytes.Buffer
    encoder := gob.NewEncoder(&buffer)
    err := encoder.Encode(orderList)
    if err != nil {
        fmt.Println("Error encoding list:", err)
        return
    }

    // send the byte array through UDP
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Println("Error connecting to UDP:", err)
        return
    }
    defer conn.Close()

    _, err = conn.Write(buffer.Bytes())
    if err != nil {
        fmt.Println("Error sending data through UDP:", err)
        return
    }

	

	

	




	



	

	


	


}*/