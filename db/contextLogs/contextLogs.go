package contextLogs

import (
	"fmt"
	"github.com/mirrr/mgo-wrapper"
	"reflect"
	"sync"
	"time"
)

type (
	obj    map[string]interface{}
	Geo    map[string] string
	Params struct {
		Select map[string]interface{}
		Query  map[string]interface{}
		Skip   int
		Limit  int
	}
	contextLogs struct {
		Timestamp time.Time
		Price float32
		Query string
		Uid   string
		Geo Geo
	}
)

var (
	streams map[string]contextLogs
	mu      sync.Mutex
	DBName string = "searchTDS"
	Collection string = "contextLogs"
)

func Get(id string) contextLogs {
	mu.Lock()
	defer mu.Unlock()
	return streams[id]
}

func Count() (n int, err error){
	return mongo.DB(DBName).C(Collection).Count()
}

func All(p Params)[]contextLogs {
	arr := []contextLogs{}
	mu.Lock()

	qr := mongo.DB(DBName).C(Collection).Find(p.Query).Skip(p.Skip)

	if p.Limit > 0 {
		qr.Limit(p.Limit)
	}

	if len(p.Select) > 0 {
		qr.Select(p.Select)
	}

	err := qr.All(&arr)

	if err != nil {
		fmt.Println("Error: streams(1):", err)
	}

	defer mu.Unlock()

	return arr
}

func Print (c contextLogs) {
	fmt.Println("contextLogs map")

	val := reflect.ValueOf(&c).Elem()

	for i := 0; i < val.NumField(); i++ {

		fieldType := val.Type().Field(i)
		fieldValue := val.Field(i)

		fmt.Printf("--- %s, %v \n", fieldType.Name, fieldValue.Interface())
	}

	fmt.Println("---------------")
}
