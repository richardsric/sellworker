package main

import (
	"fmt"

	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/richardsric/tradebk/workers/worker"
	"github.com/richardsric/tradebk/workers/worker/helpers"
)

func main() {

	//worker.SellWorker()

	workers()

	http.HandleFunc("/", index)
	http.ListenAndServe(":6000", nil)
}

func workers() {
	count := 1
	for {
		//fmt.Println("Service Have Run This Number Of times ", count)
		timeInterval := helpers.GetTimerInterval("BuyOrderUpdateWorker")
		time.Sleep(timeInterval * time.Second)
		go worker.SetBTCPrice()
		//go worker.BuyOrderUpdateWorker()
		//go worker.SellOrderUpdateWorker()
		go worker.SellWorker()
		count = count + 1

	}
}

func index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "iTradeCoin Sell Worker Service Is Running On Port 6000")
}

func init() {
	//Load sell_worker settings from database
	var name = "iTradeCoin Sell Worker Service"
	var version = "0.001 DEVEL"
	var developer = "iYOCHU Nig LTD"

	fmt.Println("App Name: ", name)
	fmt.Println("App Version: ", version)
	fmt.Println("Developer Name: ", developer)

	//go microservices.TruncateMarketData()
}
