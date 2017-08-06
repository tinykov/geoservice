package redisprovider

import "fmt"

type vehicle_data struct{
	vehicle_id string
	s2_pos string
	busy bool
}


func Run4(){
	var v1 = vehicle_data{"3456", "2203842213153079296",false}
	fmt.Println(v1)
}

