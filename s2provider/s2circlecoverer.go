package s2provider

import (
"time"
"fmt"
"math"
"math/rand"
"github.com/golang/geo/s1"
"github.com/golang/geo/s2"

)

// The Earth's mean radius in kilometers (according to NASA).
const earthRadiusKm = 6371.01

func kmToAngle(km float64) s1.Angle {
	return s1.Angle(km / earthRadiusKm)
}

//create random LatLng given a center-point and radius
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
	var cellId = s2.CellIDFromLatLng(latlng)
	var s2_token = cellId.ToToken()
	fmt.Println(s2_token)
	var s2_pos = cellId.Pos()
	fmt.Println(s2_pos)

	var s2_level = cellId.Level()
	fmt.Println(s2_level)

	return latlng
}

//get covering region given a LatLng and disk-shaped cap
/*func GetCovering(latlng s2.LatLng,cap *s2.Cap) *s2.CellUnion {
	starttime := time.Now()
	rc := &s2.RegionCoverer{MinLevel: 12, MaxLevel: 16, LevelMod: 1, MaxCells: 100}
	var region = rc.Covering(cap)
	//fmt.Println(len(region))
	elapsed := time.Since(starttime)
	fmt.Println("-------------------------------",elapsed)
	return &region
}*/

func GetCovering(latlng s2.LatLng,cap s2.Cap) *s2.CellUnion {
	starttime := time.Now()
	rc := &s2.RegionCoverer{MinLevel: 12, MaxLevel: 16, LevelMod: 1, MaxCells: 100}
	var region = rc.Covering(cap)
	//fmt.Println(len(region))
	for _, ci := range region {
		ci++
		//var s2_pos = ci.Pos()
		//fmt.Printf("%d,\n",s2_pos);
	}

	elapsed := time.Since(starttime)
	fmt.Println("----------%s",elapsed)
	return &region
}

//create s2.Cap from LatLng and radius meters
func GetS2Cap(ll s2.LatLng,radius_meters float64) s2.Cap{
	var radians = (2 * math.Pi) * (radius_meters /(1000 * 40075.017));
	var point = s2.PointFromLatLng(ll);
	var axis_height = (radians * radians)/2
	cap := s2.CapFromCenterHeight(point,axis_height)
	fmt.Println(cap.String())
	return cap
}

func Run1() {

	//-26.029246,28.033959 - wroxham
	//-26.087295,28.048183 - sandton
	var ll = s2.LatLngFromDegrees(-26.087295,28.048183)
	var answer = fmt.Sprintf("[%v, %v]", ll.Lat, ll.Lng)
	fmt.Println(answer)
	s2cap := GetS2Cap(ll,2680)

	starttime := time.Now()
	var count = 3;
	for i := 0; i < count; i++ {
		region := GetCovering(ll,s2cap)
		for _, ci := range *region {
			var s2_pos = ci.Pos()
			fmt.Printf("%d,\n",s2_pos);
		}
	}
	elapsed := time.Since(starttime)
	fmt.Println("----------%s",(elapsed))

}
