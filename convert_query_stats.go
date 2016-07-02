package main

import (
	"./db/contextLogs"
	"./utils/json"
	"runtime"
	"strconv"
	"time"
	"sync"
	"fmt"
	"os"

	//"gopkg.in/mgo.v2/bson"
//	"github.com/davecgh/go-spew/spew"
)

const (
	steep          = 20000
	filePathBuffer = "./buffer/res_"
	fileMerged     = "./buffer/mergeFiles/merge_"
)

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
	fromSteep byte
	start     time.Time
	wg        sync.WaitGroup
	max       int
	err       error
	//mu      sync.Mutex
)

func ( i TQueris ) Add (a TQueris)  {
	fmt.Println("Adding instance")
}

func SteepGet () {
	fromSteep := strconv.Itoa(os.Getenv("fromSteep"))
	if fromSteep == 0 {
		fromSteep = 1
	}
}

func init () {
	start  = time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU()*3)
	SteepGet()
	max, err = contextLogs.Count()
}

func main()  {

	if fromSteep == 1 {

		if err != nil {
			fmt.Println("Error get Count", err)
			return
		}

		for i:=0; i < max;  {
			wg.Add(1)
			go queryCollected (i)
			i +=steep
		}

		wg.Wait()
		fromSteep++
	}

	if fromSteep == 2 {
		Queris := make(TQueris)

		err = json.FromFile(filePathBuffer + "0", &Queris)

		if err != nil {
			fmt.Println("[Merge part]Error get data from first file", err)
			return
		}

		for i:=steep; i < max;  {
			queryMerge(i, &Queris)
			i +=steep
		}

		fromSteep++
	}

	fmt.Println("the end", time.Since(start))
}

func queryMerge(skip int, Queris *TQueris) {

	strSkip := strconv.Itoa(skip)
	querisAdding := make(TQueris)

	err = json.FromFile(filePathBuffer + strSkip, &querisAdding)

	if err != nil {
		fmt.Println("[queryMerge]Error get data from file ", err)
		return
	}

	Queris.Add(querisAdding)

	if skip % 100000 == 0 {
		_, err = json.ToFile(Queris, fileMerged + strSkip )
		if err != nil {
			fmt.Println("[queryMerge]Error save merged ", err)
		}
	}

}

func queryCollected (skip int) {
	Queris := make(TQueris)
	p := contextLogs.Params{ Limit: steep, Skip: skip }

	arContextLogs := contextLogs.All(p)

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

		query.Count++

	}

	strSkip := strconv.Itoa(skip)
	_, err := json.ToFile(Queris, filePathBuffer + strSkip )

	if err != nil {
		fmt.Println("Error write to json ", err)
	}

	fmt.Println("Skip " + strSkip + " END")

	defer wg.Done()

}

//	Queris2 := make(TQueris)
//
//	spew.Dump(Queris)
//	wg.Add(len(arContextLogs))
//
