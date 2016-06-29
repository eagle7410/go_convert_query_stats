package main

import (
	"fmt"
	"./db/contextLogs"

)
type (
	obj    map[string]interface{}

)

func main()  {
	p := contextLogs.Params{}

	arContextLogs := contextLogs.All(p)
	fmt.Println("arContextLogs", arContextLogs)
}
