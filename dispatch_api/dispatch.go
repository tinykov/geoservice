package dispatch_api

import(
	//"golang.org/x/net/context"
	"fmt"
	//"github.com/uber/tchannel-go"
	//"github.com/uber/tchannel-go/raw"
	"github.com/uber/tchannel-go/thrift"
)

const port = 4040

type VehiclesNearRider struct{
	vehicle_id string
	s2_position string
	latitude float32
	longitude float32
}

type tripService struct{
	// getVehiclesNearRider
}

//TChannel thrift method handler for geo-fencing vehicles near rider
func (h *VehiclesNearRider) getVehiclesNearRider(ctx thrift.Context, lat string,lon string,rider_id string) (v *[]VehiclesNearRider){

	return nil
}

//TChannel thrift method handler for tracking driver vehicles
func (h *tripService) updateDriverLocation(lat string,lon string,vehicle_id string){
}

func Run2(){
	fmt.Println("dispatch running at port:",port)
}