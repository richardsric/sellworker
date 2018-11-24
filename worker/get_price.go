package worker

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/richardsric/tradebk/workers/worker/helpers"
)

// BTCprice this holds the btc price in the memory. Other function and access it and get the btc price
var BTCprice float64

// SetBTCPrice This returns the price of BTC
func SetBTCPrice() {

	url := "https://api.coinbase.com/v2/prices/spot?currency=USD"
	body, err := helpers.GetHTTPRequest(url)
	//fmt.Println(string(body))
	if err != nil {
		fmt.Println("Error On Bittrex GetTicker Func", err)
		return
	}

	if len(body) == 0 {
		fmt.Println("Nil Response Gotten From The Request", url)
		fmt.Println("Kindly Check Your Internet Connection")
		return
	}
	// unmarshal the json response.
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Market data information
	for key, val := range m {
		//fmt.Println("key", key, "val", val)
		if key == "data" {
			price1 := val.(map[string]interface{})["amount"]
			BTCprice, _ = strconv.ParseFloat(price1.(string), 64)
		}

	}

}

func getBTCPrice() float64 {

	return BTCprice
}
