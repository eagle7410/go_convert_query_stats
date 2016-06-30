package main

import (
	"./db/contextLogs"
	"labix.org/v2/mgo/bson"
	"runtime"
	"time"
	"sync"
	"fmt"

)

var (
	queris  []string
	mu      sync.Mutex
	run     int
	end     int
)

func main()  {
//	p := contextLogs.Params{Limit: 5}
//	p := contextLogs.Params{ Query: bson.M{"price" : bson.M{"$exists" : true}}, Limit: 5}
	queryCollected ()

}

//func AddQuery(query string, wgPoint *sync.WaitGroup) {
func AddQuery(query string) {

//	defer UnlockQueris(wgPoint)
	run++
	for i:=0; i<len(queris); i++ {
		if queris[i] == query {
			return
		}
	}

//	mu.Lock()
	queris = append(queris, query)
//	mu.Unlock()

}

func UnlockQueris (wgPoint *sync.WaitGroup) {
//
	end++
	fmt.Println("End tread", run, end)
	wgPoint.Done()
}

func queryCollected () {
	start := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())

	p := contextLogs.Params{Select: bson.M{"query" : 1}, Limit: 1000}

	arContextLogs := contextLogs.All(p)

//	var wg sync.WaitGroup

//	wg.Add(len(arContextLogs))

	for _,log := range arContextLogs {

//		go AddQuery(log.Query, &wg)
		AddQuery(log.Query)

	}

//	wg.Wait()

	fmt.Println("the end", time.Since(start),len(queris))
}

