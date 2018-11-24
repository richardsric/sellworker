package worker

import (
	"fmt"

	"github.com/richardsric/tradebk/workers/worker/helpers"
)

// BuyOrder this places buy order to the gateway
func BuyOrder(eID, aID, rate, Qty, key, market string) (result []byte, err error) {
	fmt.Println(".....Issuing BUY Order......")
	//if eID == nil || aID == nil || rate == nil || Qty == nil || key == "" || market =="" {

	//}
	url := "http://localhost:5000/buyOrder?market=" + market + "&quantity=" + Qty + "&rate=" + rate + "&eid=" + eID + "&apiKey=" + key + "&aid=" + aID + ""

	result, err = helpers.GetHTTPRequest(url)

	return result, err
}

// SellOrder this places buy order to the gateway
func SellOrder(eID, aID, rate, Qty, key, market string) (result []byte, err error) {
	fmt.Println(".....Issuing SELL Order......")

	//if eID == nil || aID == nil || rate == nil || Qty == nil || key == "" || market =="" {

	//}

	url := "http://localhost:5000/sellOrder?market=" + market + "&quantity=" + Qty + "&rate=" + rate + "&eid=" + eID + "&apiKey=" + key + "&aid=" + aID + ""
	//fmt.Println(url)

	result, err = helpers.GetHTTPRequest(url)
	if err != nil {
		fmt.Println("Request Failed with: ", err)
		//time.Sleep(5 * time.Second)
	}
	return result, err
}
