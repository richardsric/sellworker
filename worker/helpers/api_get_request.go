package helpers

import (
	"io/ioutil"
	"net/http"
)

// GetHTTPRequest This is use to make http Get request. It returns byte an error
func GetHTTPRequest(url string) (bs []byte, err error) {
	//	fmt.Println("Getting HTTP Request From URL: " + url + "")
	//url := "https://poloniex.com/public?command=returnTicker"
	res, err := http.Get(url)
	if (err) != nil {
		//fmt.Println("ERROR: Failed To Connected to " + url + " For Market Data")
		return nil, err
	}
	defer res.Body.Close()
	bs, err = ioutil.ReadAll(res.Body)
	if (err) != nil {
		//panic(err)
		return nil, err
	}
	//fmt.Println(string(bs))
	res.Body.Close()
	return bs, nil
}
