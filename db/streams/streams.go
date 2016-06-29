package streams

import (
	"fmt"
	"github.com/mirrr/mgo-wrapper"
	"sync"
	"time"
)

type (
	obj    map[string]interface{}
	Stream struct {
		ID     string `bson:"_id"`
		Domain uint64 `bson:"domain"`
		Plugin struct {
			RedirectResultsYandex bool   `bson:"redirectResultsYandex"`
			RedirectResultsGoogle bool   `bson:"redirectResultsGoogle"`
			ReplaceFormYandex     bool   `bson:"replaceFormYandex"`
			ReplaceFormGoogle     bool   `bson:"replaceFormGoogle"`
			Filtration            bool   `bson:"filtration"`
			FilterByStopDomains   bool   `bson:"filterByStopDomains"`
			FilterByIP            bool   `bson:"filterByIP"`
			FilterByAge           uint64 `bson:"filterByAge"`
		} `bson:"plugin"`
	}
)

var (
	streams map[string]Stream
	mu      sync.Mutex
)

func init() {
	update()
	go func() {
		for range time.Tick(time.Second) {
			update()
		}
	}()
}

func update() {
	arr := []Stream{}
	err := mongo.DB("searchTDS").C("streams").Find(obj{}).All(&arr)
	if err != nil {
		fmt.Println("Error: streams(1):", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	streams = map[string]Stream{}
	for _, current := range arr {
		streams[current.ID] = current
	}
}

func Get(id string) Stream {
	mu.Lock()
	defer mu.Unlock()
	return streams[id]
}

func All()map[string]Stream {
	mu.Lock()
	defer mu.Unlock()
	return streams
}
