package s2provider

import(
	"fmt"
	"errors"
	"github.com/golang/geo/s2"
)

//get s2 cellid from a lat,lon 
func S2CellIDfromLatLng(lat float64,lon float64) s2.CellID{
	var latlng = s2.LatLngFromDegrees(lat,lon)
	var s2cell = s2.CellIDFromLatLng(latlng)
	return s2cell
}

func S2CellKey(lat float64,lon float64) uint64{
	return S2CellIDfromLatLng(lat,lon).Pos()
}

//get parent id at level
func GetParentIdAtLevel(level int,leaf s2.CellID) (s2.CellID,error){
	var cell_struct s2.CellID
	if level > 30 || level < 1 {
		fmt.Printf("level of cell out of bounds for level :%v\n",level)
		return cell_struct, errors.New("invalid")
	}
	var cell = leaf.Parent(level)
	return cell,nil
}


//entry point
func Run(){
	var s2key = S2CellIDfromLatLng(-26.087295,28.048183)
	fmt.Printf("s2cell = %v\n",s2key)
	var latlng = s2.LatLngFromDegrees(-26.087295,28.048183)
	var cell,err = GetParentIdAtLevel(14,s2.CellIDFromLatLng(latlng))

	fmt.Printf("cell at level %v\n",cell.Pos(),err)
}

