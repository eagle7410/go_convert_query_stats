package main

import (
	"./db/contextLogs"
	"./utils/json"
	"runtime"
	"strconv"
	"time"
	"sync"
	"fmt"

	//"gopkg.in/mgo.v2/bson"
//	"github.com/davecgh/go-spew/spew"
)

const steep  = 20000
type (
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
	wg      sync.WaitGroup
	//mu      sync.Mutex
	//run     int
	//end     int
)

func main()  {
	start  = time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU()*2)
	max, err := contextLogs.Count()

	if err != nil {
		fmt.Println("Error get Count", err)
		return
	}

	max = 1000

	for i:=0; i < max;  {
		wg.Add(1)
		go queryCollected (i)
		i +=steep
	}

	wg.Wait()

	fmt.Println("the end", time.Since(start))
}

func queryCollected (skip int) {
	//p := contextLogs.Params{ Limit: steep, Query: bson.M{"geo.country" : bson.M{ "$exists" : true}}}
	Queris := make(TQueris)
	p := contextLogs.Params{ Limit: steep, Skip: skip }

	arContextLogs := contextLogs.All(p)

	//for i,log := range arContextLogs {
	for _,log := range arContextLogs {

//		contextLogs.Print(log);
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
//		if log.Uid !=  {
//
//		}

//		spew.Dump(statCountry)
//

		query.Count++

	}

	strSkip := strconv.Itoa(skip)
	_, err := json.ToFile(Queris, "./buffer/res_" + strSkip )

	if err != nil {
		fmt.Println("Error write to json ", err)
	}

	fmt.Println("Skip " + strSkip + " END")

	defer wg.Done()

//	Queris2 := make(TQueris)

//	json.FromFile("./buffer/res", &Queris2)


//	spew.Dump(Queris)
}

//	var
//	wg.Add(len(arContextLogs))
//
