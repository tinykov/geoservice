package main

import(
	"fmt"
	geo "geoservice/dispatch_api"
	//redis "geoservice/redisprovider"
	//s2prov "geoservice/s2provider"
	//"github.com/go-redis/redis"
)

func main(){
	fmt.Println("hello main.go")
	geo.Run()
	//redis.Run()
}