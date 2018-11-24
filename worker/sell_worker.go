package worker

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/richardsric/tradebk/workers/worker/helpers"

	"time"
)

// SellWorker is func to update sell worker table
func SellWorker() {
	startProcess := time.Now() // get current time
	jobCount := 0
	//fmt.Println("Get Jobs From DB")
	btcPrice := getBTCPrice()
	fmt.Println("********************************************************")
	fmt.Println("")
	con, err := helpers.OpenConnection()
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer con.Close()

	row, err := con.Db.Query(`
	SELECT id_sell_worker,market,exchange_id,highest_bid_price,
	lowest_bid_price,highest_ask_price,
	lowest_ask_price,highest_volume,lowest_volume,actual_rate,
	actual_quantity,profit_keep,sell_trigger,last_volume,start_volume,
	ask_bid,order_type,quantity,profit_lock_start,work_status,account_id,
	order_date,order_id,threshold,sell_price,volume_diff,vol_percent,pl,
	percent_profit,profit_lock_start_price,cost,last_proceed,last_bid,
	exit_price,manual_exit,high_profit,high_profit_perc,
	actual_locked_profit,actual_locked_perc_profit,locked_proceed,
	percent_exit_profit,node_id,work_id,profit_locked,akey,
	stop_loss_active,stop_loss,txn_fee,tpfee,work_age,work_started_on,
	highest_proceed,btc_price_at_start,btc_last_price,btc_highest_ask,
	btc_lowest_ask,btc_price_perc, exit_quantity, parent_id
	 FROM sell_worker WHERE work_status = 0`)
	if err != nil {
		fmt.Println("Select Failed Due To: ", err)
	}
	defer row.Close()

	for row.Next() {
		start := time.Now() // get current time
		//	fmt.Println("Entered row dot Next")
		var market, exchangeID, orderType, aKey, workAge string
		var highBid, lowBid, highAsk, lowAsk, highVol, lowVol, lastVol, startVol, actualRate, actualQty, profitKeep, selTrigger, askBid, Qty float64
		var profitLockStart, thresHold, sellPrice, volumeDiff, volumePercent, PL, percentProfit, profitLockStartPrice, cost, lastProceed, lastBid float64
		var exitPrice, highestProceed, btcStartPrice, btcLastPrice, btcHighestAsk, btcLowestAsk, btcPricePerc, exitQuantity, cRemU float64
		var highProfit, highProfitPercent, actualLockProfit, actualLockPerProfit, lockProceed, percentExitProfit, stopLossP, txnFee, tpFee float64
		var workStatus, accountID, orderID, nodeID, workID, manualExit, profitLocked, stopLossActive, idSellWorker, parentID int
		var orderDate, workStartedOn time.Time

		//initialize cRemU for holding units remaining after sell command
		cRemU = 0
		err = row.Scan(&idSellWorker, &market, &exchangeID, &highBid, &lowBid, &highAsk, &lowAsk, &highVol, &lowVol, &actualRate, &actualQty, &profitKeep,
			&selTrigger, &lastVol, &startVol, &askBid, &orderType, &Qty, &profitLockStart, &workStatus, &accountID, &orderDate, &orderID, &thresHold, &sellPrice,
			&volumeDiff, &volumePercent, &PL, &percentProfit, &profitLockStartPrice, &cost, &lastProceed, &lastBid, &exitPrice, &manualExit, &highProfit, &highProfitPercent,
			&actualLockProfit, &actualLockPerProfit, &lockProceed, &percentExitProfit, &nodeID, &workID, &profitLocked, &aKey, &stopLossActive, &stopLossP, &txnFee, &tpFee,
			&workAge, &workStartedOn, &highestProceed, &btcStartPrice, &btcLastPrice, &btcHighestAsk, &btcLowestAsk, &btcPricePerc, &exitQuantity, &parentID)
		if err != nil {
			fmt.Println("Row Scan Failed Due To:\n", err)
		}

		//http://localhost:5000/pair/price?pair=btc-bcc&eid=1
		url := "http://localhost:5000/pair/price?pair=" + market + "&eid=" + exchangeID + ""
		// call the end point with the gotten values.
		//body, err := helpers.GetHTTPRequest("http://localhost:5000/pair/price?pair=" + market + "&eid=" + exchangeID + "")
		body, err := helpers.GetHTTPRequest(url)
		//		fmt.Println(string(body))
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
			//panic(err)
			fmt.Println(err)
			return
		}

		// Market data information
		ask := m["ask"]
		bid := m["bid"]
		vol := m["volume"]

		//Check to b sure it retrieves valid data
		if bid.(float64) > 0 && ask.(float64) > 0 {

			/// check Exit price for 1
			if manualExit == 1 {
				fmt.Println("..JOB ID: ", idSellWorker, ", Exit 1 Condition Detected. Command To Execute:: Sell at Current Bid Price.")
				if exitQuantity == 0 {
					exitQuantity = actualQty
				}
				aID := fmt.Sprintf("%v", accountID)
				aQty := fmt.Sprintf("%v", exitQuantity)
				Bid := fmt.Sprintf("%v", ReduceByPercent(bid.(float64), 0.001))
				//Reduce the bid price by 0.001 percent or any value so as to improve chances of fast sell.
				//check action to take. Live Trade submits order. Simulation sends IM

				//LIVE trade begins
				if orderType == "LIVE" {
					body, err := SellOrder(exchangeID, aID, Bid, aQty, aKey, market)

					// New Error check
					if err != nil {
						fmt.Println("Sell Order Failed ", err)
						return
					}

					if len(body) == 0 {
						fmt.Println("Nil Response Gotten From The Request For Sell Order")
						fmt.Println("Kindly Check Your Internet Connection")
						return
					}

					// unmarshal the json response.
					var m map[string]interface{}

					err = json.Unmarshal(body, &m)
					if err != nil {
						fmt.Println(err)
					}

					result := m["result"]
					message := m["message"]
					orderNo := m["order_number"]

					if result == "error" {
						fmt.Println("Order For Sell Encountred The Following error: ", message)
						return
					}
					if result == "success" {
						workStatus = 1
						exitPrice = ReduceByPercent(bid.(float64), 0.001)
						fmt.Println("Order For Sell Placed Order ID: ", orderNo)
					}
					//LIVE trade ends
				} else {
					//is it SIMULATION trade
					workStatus = 1
					exitPrice = ReduceByPercent(bid.(float64), 0.001)

				}
			}

			/// check Exit price for 2

			if manualExit == 2 {
				fmt.Println("..JOB ID: ", idSellWorker, ". Exit 2 Condition Detected. Command To Execute:: Sell at Specified Bid Price.")
				if exitQuantity == 0 {
					exitQuantity = actualQty
				}
				aID := fmt.Sprintf("%v", accountID)
				eID := fmt.Sprintf("%s", exchangeID)
				aQty := fmt.Sprintf("%v", exitQuantity)
				exP := fmt.Sprintf("%v", exitPrice)
				// Get api key
				//key := getKey(aID)
				fmt.Println("Quantity =", aQty)
				fmt.Println("Exit Price =", exP)

				if orderType == "LIVE" {
					//LIVE trade begins here
					body, err := SellOrder(eID, aID, exP, aQty, aKey, market)

					// New Error check
					if err != nil {
						fmt.Println("Sell Order Failed ", err)
						return
					}

					if len(body) == 0 {
						fmt.Println("Nil Response Gotten From The Request For Sell Order")
						fmt.Println("Kindly Check Your Internet Connection")
						return
					}

					// unmarshal the json response.
					var m map[string]interface{}

					err = json.Unmarshal(body, &m)
					if err != nil {
						fmt.Println(err)
					}

					result := m["result"]
					message := m["message"]
					orderNo := m["order_number"]

					if result == "error" {
						fmt.Println("Order For Sell Encountred The Following error: ", message)
						return
					}
					if result == "success" {
						workStatus = 1
						fmt.Println("Order For Sell Placed Order ID: ", orderNo)
					}

					//LIVE trade ends here
				} else {
					//SIMULATION trade
					workStatus = 1

				}
			}

			if lastBid == 0 {
				//		fmt.Println("Values first time:", lowAsk, lowBid, lowVol)

				lastBid = bid.(float64)
				//set work start timestamp here
				workStartedOn = time.Now()
				//workStartedOn = t.Format("20060102150405")
				//fmt.Println("WorkStartedOn: ", workStartedOn)
				//workStartedOn = time.Now()

			}
			if highAsk == 0 {
				highAsk = ask.(float64)
			}
			if lowAsk == 0 {
				lowAsk = ask.(float64)
			}

			if highBid == 0 {
				highBid = bid.(float64)
			}
			if lowBid == 0 {
				lowBid = bid.(float64)
			}
			if highVol == 0 {
				highVol = vol.(float64)
			}
			if lowVol == 0 {
				lowVol = vol.(float64)
			}
			if startVol == 0 {
				startVol = vol.(float64)
			}
			lastVol = vol.(float64)
			//Do BTC Price computations here
			//btcStartPrice, btcLastPrice, btcHighestAsk, btcLowestAsk, btcPricePerc
			if btcPrice > 0 {
				btcLastPrice = btcPrice

				if btcStartPrice == 0 {
					btcStartPrice = btcPrice
				}
				if btcHighestAsk == 0 {
					btcHighestAsk = btcPrice
				}
				if btcPrice > btcHighestAsk {
					btcHighestAsk = btcPrice
				}
				if btcLowestAsk == 0 {
					btcLowestAsk = btcPrice
				}
				if btcPrice < btcLowestAsk {
					btcLowestAsk = btcPrice
				}

				btcPricePerc = RoundDown((((btcPrice - btcStartPrice) / btcStartPrice) * 100), 4)

			}

			/// check if high ask price have changed
			if ask.(float64) > highAsk {
				highAsk = ask.(float64)
			}
			//		fmt.Println("HighestAsk: ", highAsk)

			/// check if low ask price have changed
			if ask.(float64) < lowAsk {
				lowAsk = ask.(float64)
			}
			//		fmt.Println("LowestAsk: ", lowAsk)

			/// check if high bid price have changed

			oldHighBid := highBid
			//		fmt.Println("Old Highest Bid: ", oldHighBid)
			if bid.(float64) > highBid {
				highBid = bid.(float64)
			}
			//		fmt.Println("New HighestBid: ", highBid)

			/// check if low bid price have changed
			//oldLowBid := lowBid
			if bid.(float64) < lowBid {
				lowBid = bid.(float64)
			}
			//		fmt.Println("LowestBid: ", lowBid)

			/// check if high vol have changed
			oldHighVol := highVol
			//		fmt.Println("Old Highest Volume: ", oldHighVol)
			if vol.(float64) > oldHighVol {
				highVol = vol.(float64)
			}
			//		fmt.Println("HighestVol: ", highVol)

			/// check if low bid price have changed
			if vol.(float64) < lowVol {
				lowVol = vol.(float64)
			}
			//		fmt.Println("Lowestvol: ", lowVol)

			// Set First time Value insert into the sell worker table.

			// Compute vol_diff
			// volume_diff (last_volume - start_volume)
			volumeDiff = RoundDown(vol.(float64)-startVol, 8)
			//		fmt.Println("Volume Difference: ", volumeDiff)

			// Compute vol percent
			// vol_percent (volume_diff/start_volume).
			volumePercent = RoundDown((volumeDiff/startVol)*100, 4)
			//		fmt.Println("Volume Percent: ", volumePercent)

			// profit lock is not disabled
			if profitLockStart != 0 {
				// Compute if we need to set d profit_lock_start_price.
				if profitLockStartPrice == 0 {
					profitLockStartPrice = RoundDown(actualRate+RoundDown(((profitLockStart/100)*actualRate), 8), 8)
				}
				// check if bid is upto price to start d locking.
				if bid.(float64) >= profitLockStartPrice && profitLocked == 0 {
					profitLocked = 1
				}

				if profitLocked == 1 && (thresHold == 0 || (bid.(float64) >= oldHighBid)) {
					//put highbid computations here####
					var highRange, threshval, sellval, highproceed float64
					highRange = highBid - actualRate

					/// compute threshold
					threshval = RoundDown((highRange * (profitKeep + selTrigger) / 100), 8)
					thresHold = RoundDown(actualRate+threshval, 8)
					//			fmt.Println("threshold: ", thresHold)

					// compute sell_price
					sellval = RoundDown((highRange * (profitKeep + (selTrigger / 2)) / 100), 8)
					sellPrice = RoundDown(actualRate+sellval, 8)
					//			fmt.Println("SellPrice: ", sellPrice)

					// compute high profit
					highproceed = RoundDown((actualQty * highBid), 8)
					highProfit = ReduceByPercent(highproceed, txnFee) - cost
					//			fmt.Println("Highest Profit: ", highProfit)
					highestProceed = ReduceByPercent(highproceed, txnFee)
					//			fmt.Println("Highest Proceed: ", highproceed)
					// compute high profit percentage
					highProfitPercent = RoundDown(((highBid-actualRate)/actualRate)*100, 4)
					//			fmt.Println("High Profit Percent: ", highProfitPercent)

					// compute locked_proceed
					//lockProceed = ReduceByPercent((actualQty * sellPrice), txnFee)
					lockProceed = ReduceByPercent((actualQty * sellPrice), txnFee)
					//			fmt.Println("Locked Proceed: ", lockProceed)

					// compute for actual_locked_profit
					actualLockProfit = RoundDown((lockProceed - cost), 8)
					//			fmt.Println("Actual Locked Profit: ", actualLockProfit)

					// check compute
					actualLockPerProfit = RoundDown(((sellPrice-actualRate)/actualRate)*100, 4)
					//			fmt.Println("Actual Locked Profit Percent: ", actualLockPerProfit)

				}

			} //profit lock computations done if enabled

			//Calculate normal values

			// Compute last_proceed
			fullproceed := RoundDown(actualQty*bid.(float64), 8)
			lastProceed = RoundDown(ReduceByPercent(fullproceed, txnFee), 8)
			//	fmt.Println("Last Proceed: ", lastProceed)

			// Compute PL
			PL = RoundDown((lastProceed - cost), 8)
			//	fmt.Println("PL: ", PL)

			// Compute per percent_profit
			percentProfit = RoundDown(((bid.(float64)-actualRate)/actualRate)*100, 4)
			//	fmt.Println("PL Percent: ", percentProfit)

			//end normal value computations

			if manualExit == 0 && bid.(float64) <= thresHold && ReduceByPercent(bid.(float64), txnFee) >= askBid && profitLockStartPrice > 0 && profitLocked == 1 {
				fmt.Println(market, "..JOBID: ", idSellWorker, ". Threshold Breached! Determining Best Sell Price.")
				var aID, eID, aQty, exP string
				//this means dat threshold has bn breached. We issue sell order here.
				oldExitQuantity := exitQuantity
				if exitQuantity == 0 {
					exitQuantity = actualQty
				}
				if ReduceByPercent(bid.(float64), 0.001) < sellPrice {
					//it doesnt make sense to submit at that sell price.
					//Determine to sell at the current bid price once it is higher than the bid_ask
					eID = fmt.Sprintf("%s", exchangeID)
					aQty = fmt.Sprintf("%v", exitQuantity)
					exP = fmt.Sprintf("%v", ReduceByPercent(bid.(float64), 0.001))
					//Reduce the exit price by 0.001 percent or any value so as to improve chances of fast sell.

				} else {
					aID = fmt.Sprintf("%v", accountID)
					eID = fmt.Sprintf("%s", exchangeID)
					aQty = fmt.Sprintf("%v", exitQuantity)
					exP = fmt.Sprintf("%v", ReduceByPercent(sellPrice, 0.001))
					//Reduce the exit price by 0.001 percent or any value so as to improve chances of fast sell.
				}
				if orderType == "LIVE" {
					//LIVE trade begins here
					body, err := SellOrder(eID, aID, exP, aQty, aKey, market)

					// New Error check
					if err != nil {
						fmt.Println("Sell Order Failed ", err)
						exitQuantity = oldExitQuantity
						return
					}

					if len(body) == 0 {
						fmt.Println("Nil Response Gotten From The Request For Sell Order")
						fmt.Println("Kindly Check Your Internet Connection")
						return
					}

					// unmarshal the json response.
					var m map[string]interface{}

					err = json.Unmarshal(body, &m)
					if err != nil {
						fmt.Println(err)
					}

					result := m["result"]
					message := m["message"]
					orderNo := m["order_number"]

					if result == "error" {
						fmt.Println("Order For Sell Encountred The Following error: ", message)
						return
					}
					if result == "success" {
						// compute for exit price
						exitPrice = ReduceByPercent(sellPrice, 0.001)
						workStatus = 1
						fmt.Println("Order For Sell Placed Order ID: ", orderNo)
					}

					//LIVE trade ends here
				} else {
					// SIMULATION trade
					exitPrice = ReduceByPercent(sellPrice, 0.001)
					workStatus = 1

				}
			}

			if manualExit == 0 && ReduceByPercent(bid.(float64), txnFee) < askBid && profitLockStartPrice > 0 && profitLocked == 1 {
				fmt.Println(market, "..JOBID: ", idSellWorker, "Market Crashed Too Fast. Disabling Profit Adjustment")
				fmt.Println(market, "..JOBID: ", idSellWorker, ". Auto Profit Adjust will switch ON when bid price reaches", profitLockStartPrice, "again")
				// disable profitlock to prevent false sell trigger during recovery when it reaches another trigger point
				profitLocked = 0
				//reset highest bid/ask to current bid so that when profit locking kicks in next time
				//false sell does not trigger due to higest profit computed from highest bid
				highBid = bid.(float64)
				highAsk = ask.(float64)
			}

			if exitPrice > 0 {
				percentExitProfit = RoundDown(((exitPrice-actualRate)/actualRate)*100, 4)
			}

			fmt.Println(market, orderType, "SELL WORKER ANALYSIS FOR JOB ID:", idSellWorker, "By Worker ID", workID)
			fmt.Println("TRADE TYPE:", orderType)
			fmt.Println("...................... ")
			fmt.Println("TICKER INFO")
			fmt.Println("...................... ")
			fmt.Println("Market:", market)
			fmt.Println("Last Bid:", bid)
			fmt.Println("Highest Bid:", highBid)
			fmt.Println("Lowest Bid:", lowBid)
			fmt.Println("Highest Ask:", highAsk)
			fmt.Println("Lowest Ask:", lowAsk)
			if btcPrice > 0 {
				fmt.Println("BTC Price @ Start:", btcStartPrice)
				fmt.Println("Last BTC Price:", btcPrice)
				fmt.Println("Highest BTC Price:", btcHighestAsk)
				fmt.Println("Lowest BTC Price:", btcLowestAsk)
				fmt.Println("BTC Pct Change:", btcPricePerc, "%")

			}
			fmt.Println("...................... ")
			fmt.Println("MARKET VOLUME INFO")
			fmt.Println("...................... ")
			fmt.Println("Vol @ Initially Work Start:", startVol)
			fmt.Println("Last Vol:", lastVol)
			fmt.Println("Highest Vol:", highVol)
			fmt.Println("Lowest Vol:", lowVol)
			fmt.Println("Volume Difference: ", volumeDiff)
			fmt.Println("Volume Pct:", volumePercent, "%")

			fmt.Println("......................... ")
			fmt.Println("COSTS/PROFIT/LOSS INFO ")
			fmt.Println("......................... ")
			fmt.Println("Actual Quanity:", actualQty)
			fmt.Println("Actual Entry Rate:", actualRate)
			fmt.Println("Cost:", cost)
			fmt.Println("Last Proceed:", lastProceed)
			if (PL >= 0) && (percentProfit >= 0) {
				fmt.Println("|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||")
				fmt.Println("|******>>>>>>---->> Profit/Loss:", PL, " <<-----<<<<<<<******|")
				fmt.Println("|******>>>>>>---->> Percent PL:", percentProfit, "%     <<-----<<<<<<<******|")
				fmt.Println("|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||")

			} else {
				fmt.Println("Profit/Loss:", PL)
				fmt.Println("Percent PL:", percentProfit, "%")
			}

			fmt.Println("................................. ")
			fmt.Println("AUTO TRADE PROFIT ADJUSTMENT INFO")
			fmt.Println("................................. ")
			if profitLocked > 0 {
				fmt.Println("Profit Lock/Adjustment: <ON>")
				//show all lock related figure
				fmt.Println("Profit Lock Engaged @:", profitLockStart, "% Profit")
				fmt.Println("Profit Lock Start Price:", profitLockStartPrice)
				fmt.Println("Minimum Profit To Keep:", profitKeep, "% of ", highProfit)
				fmt.Println("Threshold:", thresHold)
				fmt.Println("Auto Sell Price:", sellPrice)
				fmt.Println("Actual Locked Profit:", actualLockProfit, "out of", highProfit)
				fmt.Println("Actual Pct Locked Profit:", actualLockPerProfit, "%")
				fmt.Println("Locked Proceed:", lockProceed)
				fmt.Println("Highest Profit attained:", highProfit)
				fmt.Println("Highest Proceed reached:", highestProceed)
				fmt.Println("Highest Pct Profit attained:", highProfitPercent, "%")

			} else {
				//hide every lock related figure
				fmt.Println("Profit Lock/Adjustment: <Off>")

			}

			if stopLossActive > 0 || stopLossP > 0 {
				fmt.Println("Stop Loss: <ON>")
				fmt.Println("Stop Loss Set @ :", stopLossP)
			} else {
				fmt.Println("Stop Loss: <OFF>")
			}

			if exitPrice > 0 {
				fmt.Println("Pct Exit Price:", percentExitProfit, "%")
				fmt.Println("Exit Type:", manualExit)
			}
			//fmt.Println("Exit Type:", manualExit, "Current Price: ", bid, " Exit Price: ", exitPrice)
			//	fmt.Println(".....................................")
			elapsed := time.Since(start)
			//	fmt.Println(".....................................")
			fmt.Println("Job Processing Duration:", elapsed)
			//	fmt.Println("Time Started:", start)
			//	fmt.Println("Time End:", time.Now())
			//	fmt.Println("Time Difference/Elapsed:", elapsed)

			// Check if exitQuantity < actualQty
			if workStatus == 1 && (exitQuantity < actualQty) {
				//there is left over, therefore start a new worker job with existing information
				//recalculate cost based on reduced quantity
				//cost of 1 unit of sold coin, cS = cost/actualQty.
				//cost of remaining units, cRemU = cs * remUnits  = roundDown((cost/actualQty) * (actualQty - exitQuantity), 8)
				cRemU = RoundDown((cost/actualQty)*(actualQty-exitQuantity), 8)
			}

			if cRemU > 0 && workStatus == 1 {
				//Partial sale. insert new job with parentID
				//fmt.Println("Some Left overs detected. Creating new task with remaining units")
				parentID = idSellWorker

			} else {
				//Full Sale, Update Existing record
				//fmt.Println("Preparing Update")
				_, err = con.Db.Exec("UPDATE sell_worker SET highest_bid_price = $1, lowest_bid_price= $2,highest_ask_price =$3,lowest_ask_price=$4,last_bid=$5,highest_volume=$6,lowest_volume=$7,threshold=$8,sell_price=$9,high_profit=$10,high_profit_perc=$11,actual_locked_profit=$12,actual_locked_perc_profit=$13,locked_proceed=$14,last_proceed=$15,pl=$16,exit_price=$17,last_volume=$18,volume_diff=$19,vol_percent=$20,percent_profit=$21,percent_exit_profit=$22,profit_lock_start_price=$23,start_volume=$24, work_status=$25,stop_loss_active=$26,work_started_on=$27::timestamp,profit_locked=$28,work_age = age(current_timestamp, work_started_on::timestamp), work_count=work_count+1, highest_proceed=$29, last_work_duration=$30::int,btc_price_at_start=$31,btc_last_price=$32,btc_highest_ask=$33,btc_lowest_ask=$34,btc_price_perc=$35 WHERE id_sell_worker = $36",
					highBid, lowBid, highAsk, lowAsk, bid, highVol, lowVol, thresHold, sellPrice, highProfit, highProfitPercent, actualLockProfit,
					actualLockPerProfit, lockProceed, lastProceed, PL, exitPrice, vol, volumeDiff, volumePercent, percentProfit, percentExitProfit, profitLockStartPrice, startVol, workStatus, stopLossActive, workStartedOn, profitLocked, highestProceed, time.Since(start), btcStartPrice, btcLastPrice, btcHighestAsk, btcLowestAsk, btcPricePerc, idSellWorker)
				if err != nil {
					fmt.Println("Update Failed Due To: ", err)
				}
			}

			//Send Message
			SendMsg(idSellWorker)

			jobCount = jobCount + 1
			fmt.Println("Total Job Processing Duration:", time.Since(start))
			fmt.Println("*****************************************************************************************************")

		} else {
			fmt.Println("Invalid Market Data. ASK/BID Must be greater than 0. Please Check Network connection")

		}
	}
	fmt.Println("***********************************************")
	fmt.Println("Total Process Duration:", time.Since(startProcess))
	fmt.Println("Total Number Of Jobs Processed:", jobCount)
	fmt.Println("***********************************************")
	fmt.Println("*****************************************************************************************************")
	fmt.Println("")

}

// Round var f float64 = 514.89317306
// Round this rounds Output: 514.89317306 to 514
func Round(input float64) float64 {
	if input < 0 {
		return math.Ceil(input - 0.5)
	}
	return math.Floor(input + 0.5)
}

// RoundUp sample output in 4 decmial places
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

/*
"ReduceByPercent ... Reduce fv value by the percentage specified as perc
Returns value rounded to 8dp."
*/
func ReduceByPercent(fv float64, perc float64) (newfv float64) {

	newfv = RoundDown(fv-((perc/100)*fv), 8)
	return
}

/*
SendMsg ... Sends Message (IM and Support Email) Of Trade Summary of the Trade with JOB ID given to it.
*/
func SendMsg(jid int) (res bool) {
	//This retrieves JOB information and composes the messages to be sent to SendIM, SendEmail
	res = true
	return
}

// RoundDown Output in 4 decmial places
// var f float64 = 514.89317306
//RoundDown this round Output: 514.89317306 to 514.8931
func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}
