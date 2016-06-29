package contextLogs

import (
	"fmt"
	"github.com/mirrr/mgo-wrapper"
	"sync"
	//"time"
)

type (
	obj    map[string]interface{}
	Params struct {
		query map[string]interface{}
		limit int
	}
	contextLogs struct {

	}
)

var (
	streams map[string]contextLogs
	mu      sync.Mutex
	DBName string = "searchTDS"
	Collection string = "contextLogsTest"
)

func Get(id string) contextLogs {
	mu.Lock()
	defer mu.Unlock()
	return streams[id]
}

func All(p Params)[]contextLogs {
	arr := []contextLogs{}
	mu.Lock()

	qr := mongo.DB(DBName).C(Collection).Find(p.query)

	if p.limit > 0 {
		qr.Limit(p.limit)
	}

	err := qr.All(&arr)

	if err != nil {
		fmt.Println("Error: streams(1):", err)
	}

	defer mu.Unlock()

	return arr
}
