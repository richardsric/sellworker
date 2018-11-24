package helpers

import (
	"fmt"
	"time"
)

// GetTimerInterval This is use to get function time interval from database
func GetTimerInterval(functionName string) time.Duration {
	var timeInterval time.Duration
	con, err := OpenConnection()
	if err != nil {
		//return err
		fmt.Println(err)
	}
	defer con.Close()

	row, err := con.Db.Query("SELECT time_interval FROM time_interval_settings WHERE process_name = $1", functionName)
	if err != nil {
		fmt.Println("Select Failed Due To: ", err)
		return 5
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&timeInterval)
		if err != nil {
			fmt.Println("Row Scan Failed Due To: ", err)
			return 5
		}

	}

	return timeInterval
}
