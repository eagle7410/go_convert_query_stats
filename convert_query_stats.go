package main

import (
	"./db/contextLogs"
	"./utils/json"
	//"labix.org/v2/mgo/bson"

	"runtime"
	//"strconv"
	"time"
	//"sync"
	"fmt"

	//"gopkg.in/mgo.v2/bson"
	//"github.com/davecgh/go-spew/spew"
)

const steep  = 100000
type (
	Count struct {
		val int
	}
	StatsDay  map[string]int
	StatsDays map[string]StatsDay
	Stats     struct {
		Data      StatsDays
		Count     int
	}
	TQueris map[string] *Stats
)

var (
	start   time.Time
	Queris  TQueris
	//mu      sync.Mutex
	//run     int
	//end     int
)

func main()  {
	start  = time.Now()
	Queris = make(TQueris)
	runtime.GOMAXPROCS(runtime.NumCPU())
	queryCollected ()

}

//func AddQuery(query string, wgPoint *sync.WaitGroup) {
//func AddQuery(query string) {
//	//Queris = append(queris, query)
//}
//func UnlockQueris (wgPoint *sync.WaitGroup) {
//
//	fmt.Println("End tread", run, end)
//	wgPoint.Done()
//}


func queryCollected () {
	//p := contextLogs.Params{ Limit: steep, Query: bson.M{"geo.country" : bson.M{ "$exists" : true}}}
	p := contextLogs.Params{ Limit: steep}

	arContextLogs := contextLogs.All(p)

	//for i,log := range arContextLogs {
	for _,log := range arContextLogs {

		//contextLogs.Print(log);
		country := "ALL";

		if len(log.Geo) != 0 {
			for k, val:= range log.Geo {
				if k == "country" && val != "" {
					country = val
				}
			}
		}

		stDate := log.Timestamp.Format("2006-01-02")

		query := Queris[log.Query]

		if query == nil {
			Queris[log.Query] = &Stats{}
			query = Queris[log.Query]
		}


		statDays := query.Data

		if statDays == nil {
			query.Data  = make(StatsDays)
			statDays = query.Data
		}

		statDay := statDays[stDate]

		if statDay == nil {
			query.Data[stDate] = make(StatsDay)
			statDay = query.Data[stDate]
		}

		query.Data[stDate]["ALL"]++

		if country != "ALL" {
			query.Data[stDate][country]++
		}

		query.Count++

		//if i % 1000 == 0 {
		//	fmt.Println("__ITER ", i)
		//}
		//fmt.Println("__ITER ", i)
	}

	_, err := json.ToFile(Queris, "./buffer/res")

	if err != nil {
		fmt.Println("Error write to json ", err)
	}

	Queris2 := make(TQueris)

	json.FromFile("./buffer/res", &Queris2)

	fmt.Println("the end", time.Since(start), len(Queris), len(Queris2))
	//spew.Dump(Queris2)
}

//	var wg sync.WaitGroup
//	wg.Add(len(arContextLogs))
//	wg.Wait()
