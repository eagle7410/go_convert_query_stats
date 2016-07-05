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
	//"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
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
	Queris, querisAdding TQueris
	fromSteep, max, i    int
	strSkip, labelLast   string
	start                time.Time
	wg                   sync.WaitGroup
	err                  error
	//mu      sync.Mutex
)

func ( inst TQueris ) Add (a *TQueris)  {
	for query, stat := range *a {
		qr := inst[query]

		if qr == nil {
			inst[query] = stat
			continue
		}

		inst[query].Count += stat.Count

		for day, statCountries := range stat.Data {
			statDay := qr.Data[day]

			if statDay == nil {
				inst[query].Data[day] = statCountries
				continue
			}

			for country, count := range statCountries {
				statCount := statDay[country]

				if statCount == 0 {
					inst[query].Data[day][country] = count
					continue
				}

				inst[query].Data[day][country] += count
			}

		}

	}
}

func SteepGet () {

	strEnvFromSteep := os.Getenv("fromSteep")
	envFromSteep, _ := strconv.Atoi(strEnvFromSteep)
	fromSteep = envFromSteep

	if fromSteep == 0 {
		fromSteep = 1
	}

}

func init () {
	runtime.GOMAXPROCS(runtime.NumCPU()*3)
	start  = time.Now()
	SteepGet()
	max, err = contextLogs.Count()
	labelLast = "Last"
}

func main()  {
	fmt.Println("fromSteep", fromSteep)

	if err != nil {
		fmt.Println("Error get Count", err)
		return
	}

	if fromSteep == 1 {

		for i=0; i < max;  {
			wg.Add(1)
			go queryCollected (i)
			i +=steep
		}

		wg.Wait()
		fromSteep++
	}

	if fromSteep == 2 {
		max = 3300000
		Queris = make(TQueris)


		err = json.FromFile(filePathBuffer + "0", &Queris )

		if err != nil {
			fmt.Println("[Merge part]Error get data from first file", err)
			return
		}

		for i=steep; i < max;  {
			strSkip = strconv.Itoa(i)
			querisAdding = make(TQueris)
			queryMerge(&i, &strSkip, &Queris, &querisAdding)
			i += steep
			fmt.Println( "Next steep" + strconv.Itoa(i) )
		}

		saveMerge(&labelLast)

		fromSteep++

	}

	fmt.Println("the end", time.Since(start))
}

func queryMerge(skip *int, strSkip *string, Queris *TQueris, querisAdding *TQueris) {

	if json.FromFile(filePathBuffer + *strSkip, querisAdding) != nil {
		fmt.Println("[queryMerge]Error get data from file ", err)
		return
	}


	Queris.Add(querisAdding)

	if *skip % 10000000 == 0 {
		saveMerge(strSkip)
	}

}
func saveMerge(strSkip *string) {
	if json.ToFile(Queris, fileMerged + *strSkip ) != nil {
		fmt.Println("[queryMerge]Error save merged ", err)
	} else {
		fmt.Println("[queryMerge]Save ok " + fileMerged + *strSkip)
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
	err := json.ToFile(Queris, filePathBuffer + strSkip )

	if err != nil {
		fmt.Println("Error write to json ", err)
	}

	fmt.Println("Skip " + strSkip + " END")

	defer wg.Done()

}
