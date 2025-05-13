package main

import (
	"fmt"

	"github.com/robfig/cron"
)

func printActualCurrency() {

	fmt.Println("1$ = 41.54 ")

}

func main() {

	c := cron.New()

	c.AddFunc("@every 5m", printActualCurrency)

	c.Run()

}
