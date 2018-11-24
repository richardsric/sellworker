package helper

import (
	"math"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type psCoinInfo struct {
	P string
	S string
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomNumber This use to generate random number in a given length
func RandomNumber(length int) string {
	const chars = "0123456789"
	result := ""
	for i := 0; i < length; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

// RandomString This use to generate random string in a given length
func RandomString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for i := 0; i < length; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

// RandomNumberString This use to generate random number string in a given length
func RandomNumberString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < length; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

// ReduceByPercent is use to reduce a percentage of a given input
func ReduceByPercent(fv float64, perc float64, dp int, roundMode string) (newfv float64) {

	newfv = Round(fv-((perc/100)*fv), dp, roundMode)
	return
}

// IncreaseByPercent is use to increae a percentage of a given input
func IncreaseByPercent(fv float64, perc float64, dp int, roundMode string) (newfv float64) {

	newfv = Round(fv+((perc/100)*fv), dp, roundMode)
	return
}

// Round is use to round up or down a given number. Use c - for ceiling to a whole number
// Use u - for rounding up a number. Use d - for rounding down a number.
func Round(input float64, places int, roundMode string) float64 {
	var newVal float64
	roundMode = strings.ToLower(roundMode)
	if roundMode == "c" {
		if input < 0 {
			newVal = math.Ceil(input - 0.5)
		}
		newVal = math.Floor(input + 0.5)
	}

	if roundMode == "u" {
		var round float64
		pow := math.Pow(10, float64(places))
		digit := pow * input
		round = math.Ceil(digit)
		newVal = round / pow
	}

	if roundMode == "d" {
		var round float64
		pow := math.Pow(10, float64(places))
		digit := pow * input
		round = math.Floor(digit)
		newVal = round / pow

	}
	return newVal
}

// sample Output in 4 decmial places
// var f float64 = 514.89317306
//RoundUp this rounds Output: 514.89317306 to 514.8932
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

// sample Output in 4 decmial places
// var f float64 = 514.89317306
//RoundUp this round Output: 514.89317306 to 514.8931
func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}

//PsCoin returns the Primary or Secondary coin based on market(in string) passed
func PsCoin(market string) psCoinInfo {
	marketPair := strings.SplitN(strings.Trim(market, " "), "-", 2)
	return psCoinInfo{P: strings.ToUpper(strings.Trim(marketPair[0], " ")), S: strings.ToUpper(strings.Trim(marketPair[1], " "))}
}

//Checks for Nil Interface and recovers
func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
