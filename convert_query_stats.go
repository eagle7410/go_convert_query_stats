package main

import (
	"./db/contextLogs"
	"labix.org/v2/mgo/bson"
	"runtime"
	"sync"
	"fmt"

)

var (
	queris  []string
	mu      sync.Mutex
)

func main()  {
//	p := contextLogs.Params{Limit: 5}
//	p := contextLogs.Params{ Query: bson.M{"price" : bson.M{"$exists" : true}}, Limit: 5}
	queryCollected ()

}

func AddQuery(query string, wgPoint *sync.WaitGroup) {

	l := len(queris)

	fmt.Println("Query ", query)

	for i:=0; i<l; i++ {
		if queris[i] == query {
			return
		}
	}

	queris = append(queris, query)

	wgPoint.Done()
}

func queryCollected () {
	runtime.GOMAXPROCS(runtime.NumCPU())

	p := contextLogs.Params{Select: bson.M{"query" : 1}, Limit: 5}

	arContextLogs := contextLogs.All(p)

	var wg sync.WaitGroup

	wg.Add(len(arContextLogs))

	for _,log := range arContextLogs {

		go AddQuery(log.Query, &wg)

	}

	wg.Wait()

	fmt.Println("the end", len(queris), queris)
}

