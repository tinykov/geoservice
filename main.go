package main

import(
	"fmt"
	//geo "geoservice/dispatch_api"
	redis "geoservice/redisprovider"
	//s2prov "geoservice/s2provider"
)

func main(){
	fmt.Println("hello main.go")
	//geo.Run()
	redis.Run()
}