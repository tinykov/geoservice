package dispatch_api

import (
	"time"
	"fmt"
	"math"
	"math/rand"
	"bufio"
	"strings"
	"log"
	"net"
	"os"
	"runtime"
	
	tchannel "github.com/uber/tchannel-go"
	thrift "github.com/uber/tchannel-go/thrift"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"

	"github.com/opentracing/opentracing-go"
	"github.com/go-redis/redis"
	gen "geoservice/gen-go/tripservice"
	s2service "geoservice/s2provider"
)

func kmToAngle(km float64) s1.Angle {
	// The Earth's mean radius in kilometers (according to NASA).
	const earthRadiusKm = 6371.01
	//fmt.Printf(km / earthRadiusKm)
	return s1.Angle(km / earthRadiusKm)
}

func generateLatLng(y0,x0 float64, radius int) s2.LatLng {
	radiusInDegrees := float64(radius) / 111000.0

	u := rand.Float64()
	v := rand.Float64()
	w := radiusInDegrees * math.Sqrt(u)
	t := 2 * math.Pi * v
	x := w * math.Cos(t)
	y := w * math.Sin(t)

	newX := x / math.Cos(y0)

	foundLng := newX + x0
	foundLat := y + y0

	var latlng = s2.LatLngFromDegrees(foundLat, foundLng)
	//var cellId = s2.CellIDFromLatLng(latlng)
	var cellId = s2service.S2CellIDfromLatLng(foundLat,foundLng)
	var s2_token = cellId.ToToken()
	fmt.Println(s2_token)
	var s2_pos = cellId.Pos()
	fmt.Println(s2_pos)

	var s2_level = cellId.Level()
	fmt.Println(s2_level)

	return latlng
}

//entry point
func Run() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		})
	//fmt.Println(client)
	response,err2 := client.Get("city").Result()
	if(err2 != nil){
		log.Fatal("Error getting city cells ",err2)
	}
	fmt.Println(response)

	//client.NewScript()
	//-26.029246,28.033959 - wroxham
	//-26.087295,28.048183 - sandton 
	var ll = s2.LatLngFromDegrees(-26.087295,28.048183)
	var answer = fmt.Sprintf("[%v, %v]", ll.Lat, ll.Lng)
	fmt.Println(answer)
    
    starttime := time.Now()
    var count = 1;
    for i := 0; i < count; i++ {
    	latlng := generateLatLng(-26.087295,28.048183,26000)
    	s2cap := s2service.GetS2Cap(latlng,2680)
    	fmt.Println(i)
    	s2service.GetCovering(latlng,s2cap)
	}
	elapsed := time.Since(starttime)
	fmt.Println("----------%s",(elapsed))
    //---------------------------------------------------------
    var (
		listener net.Listener
		err      error
	)

	if listener, err = setupServer(); err != nil {
		fmt.Println("setupServer failed: %v...%v", err,listener)
	}
	fmt.Println("Server setup ... ")
	go listenConsole()
	select {} //lets a goroutine wait on multiple communication operations
	
	// Run for 10 seconds, then stop
	//time.Sleep(time.Second * 10)
}

func setupServer() (net.Listener, error) {
	tchan, err := tchannel.NewChannel("t-server", optsFor("TChannel-server"))
	if err != nil {
		fmt.Println("setup error %v",err)
		return nil, err
	}

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		fmt.Println("Listen error %v",err)
		return nil, err
	}

	server := thrift.NewServer(tchan)
	//&tripServiceHandler{}
	server.Register(gen.NewTChanTripServiceServer(&tripServiceHandler{}))

	// Serve will set the local peer info, and start accepting sockets in a separate goroutine.
	tchan.Serve(listener)
	//err2 := tchan.ListenAndServe("127.0.0.1:2001")
	fmt.Println("Server listening %v",tchan)
	return listener, nil
}

/*
  TChannel Thrift service handler (must implement all methods defined in IDL)
*/
type tripServiceHandler struct{
	//ringpop *ringpop.Ringpop
	//channel *tchannel.Channel
	//logger  *log.Logger
}

var (Tracer opentracing.Tracer)

func (c *tripServiceHandler) UpdateDriverLocation(ctx thrift.Context, lat float64, lon float64, vehicleID string) error {
	//fmt.Println("TChannel-thrift call...coming for vehicle_id %v",vehicleID)
	fmt.Println("addDriverLocation called over = %v,%v,%v",lat,lon,vehicleID)
	return nil
}

//
func (c *tripServiceHandler) GetVehiclesNearRider(ctx thrift.Context, lat float64, lon float64) (r gen.VehicleList, err error) {
	fmt.Println("Get VehiclesNearRider called over TChannel-thrift = %v,%v",lat,lon)
	var ll = s2.LatLngFromDegrees(lat,lon)
	s2cap := s2service.GetS2Cap(ll,2680)
    region := s2service.GetCovering(ll,s2cap)
    v := gen.VehiclesNearRider{"7721","2203841155517382656",-26.087295,28.048183}
    var size = len(*region)
    list:= make(gen.VehicleList,0,size)
  
    for i := 0; i < size; i ++ {
    	list = append(list,&v)
	}
    
    fmt.Println("vehicle list...%v",len(list))
	return list,nil
}

//listen on standard console for input
func listenConsole() {
	rdr := bufio.NewReader(os.Stdin)
	for {
		line, _ := rdr.ReadString('\n')
		switch strings.TrimSpace(line) {
		case "s":
			printStack()
		default:
			fmt.Println("Unrecognized command:", line)
		}
	}
}

func printStack() {
	buf := make([]byte, 1000)
	runtime.Stack(buf, true /* all */)
	fmt.Println("Stack:\n", string(buf))
}

func optsFor(processName string) *tchannel.ChannelOptions {
	return &tchannel.ChannelOptions{
		ProcessName: processName,
		//Tracer: tracer,
		Logger:      tchannel.NewLevelLogger(tchannel.SimpleLogger, tchannel.LogLevelWarn),
	}
}