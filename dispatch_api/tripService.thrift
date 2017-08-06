struct VehiclesNearRider{
	1: required string vehicle_id;
	2: required string s2_position;
	3: required double latitude;
	4: required double longitude;
}

typedef list<VehiclesNearRider> vehicleList

service tripService {
	vehicleList getVehiclesNearRider(
	1: double lat,
	2: double lon);

	void updateDriverLocation(
	1: double lat,
	2: double lon,
	3: string vehicle_id);
}

