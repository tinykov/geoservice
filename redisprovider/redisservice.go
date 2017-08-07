package redisprovider

import(
	"fmt"
	"log"
	"os"
	"time"
	"reflect"
	"bytes"
    "io/ioutil"
    "geoservice/s2provider"
	"github.com/go-redis/redis"
	"encoding/json"
	//"crypto/sha1"
	//"encoding/hex"
)

var client *redis.Client
var source = readFiletoString("geo-radius.lua")
var script1 = redis.NewScript(source)

type vehicle struct{
	vehicle_id string
	s2_pos string
	busy bool
}

func readFiletoString(file string) (string) {
	byteArray, err := ioutil.ReadFile(file)
	if(err != nil){
		fmt.Printf("Error reading file %v",err)
	}
	//s := string(byteArray[:])
	s := bytes.NewBuffer(byteArray).String()
	//fmt.Printf("lua script = %v",s)
	return s
}

func initialize(){
	fmt.Println("------------------------------------")
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		})
}

type RedisCache struct {
    client *redis.Client
}

func (c *RedisCache) GetClient() *redis.Client {
	//os.Getenv("REDIS_HOST")
    if c.client == nil {
    	fmt.Printf("Host %v\n",os.Getenv("REDIS_HOST"))
        c.client = redis.NewClient(&redis.Options{
            Addr:       "localhost:6379",
            Password:    "", // no password set
            DB:          0,  // use default DB
            DialTimeout: 4 * time.Second,
            PoolSize:    1,
        })
    }
    return c.client
}

func (c *RedisCache) Get(key string) (string, error) {
    return c.GetClient().Get(key).Result()
}

//get vehicles near rider in redis
func RedisVehiclesInCellArray(cellarray []string) *map[string]interface{}{
	redis_c := RedisCache{}
	resp, err := redis_c.GetClient().Get("city").Result()
	fmt.Printf("city = %v\n",resp)
	
	key_array := []string{"1","12"}
	fmt.Printf("RedisVehiclesInCellArray() ->args array %v\n",cellarray)

	resp2,err := script1.Run(redis_c.GetClient(),key_array,cellarray[0],cellarray[1],cellarray[2],cellarray[3],
		cellarray[4],cellarray[5],cellarray[6],cellarray[7],
		cellarray[8],cellarray[9],cellarray[10],cellarray[11]).Result()

	//resp2,err := script1.Run(redis_c.getClient(),key_array,make([]interface{"2203794242663350272"}, 6)).Result()
	if(err != nil){
		log.Fatal("error occured %v\n",err)
	}
	fmt.Println("Type of resp2 : %v",reflect.TypeOf(resp2))
	fmt.Println("---------CJSON payload from redis ------")
	fmt.Printf("response %v\n", resp2)
	
	/*const lit = 
		'{"7448":["2203681420263724327"],' +
		'"6975":["2203793846057184235"],' +
		'"6320":["2203794206341922299"]}' +
	bytes := []byte(lit)*/
	var v_data map[string]interface{}
	value, ok := resp2.(string)
	if ok == true{
		bytes := []byte(value)

		json.Unmarshal(bytes,&v_data)
		fmt.Println("V_DATA = %v",len(v_data))

	}
	//vehiclelist := vehicle{}p
	/*return vehicle{
		vehicle_id:"8888",
		s2_pos:"2203788195349397504",
	}*/
	return &v_data
}

func Run(){
	initialize()
	cellArray := []string{"2203794242663350272","2203793830346489856","2203681817599410176","2203681405282549760","2203681817599410176"}
	args_array := make([]string,0,12)

	for i := 0; i <12; i++{
		if i < len(cellArray){
			args_array = append(args_array,cellArray[i])
		}
		args_array = append(args_array,"")
	}

	v := RedisVehiclesInCellArray(args_array)
	s2cell := s2provider.S2CellIDfromLatLng(-26.087295,28.048183)
	fmt.Println("hello redis %v ...%v",v,s2cell)
}