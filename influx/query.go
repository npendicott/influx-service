package influx

import (
	client "github.com/influxdata/influxdb1-client/v2" //"github.com/influxdata/influxdb/client/v2"

	"fmt"
	"log"
)

// TODO: Paramaterized queries, and maybe some query builder? All I know is this isn't so well structured

// QueryDateRangeWithMAC gets a response (Influx client official lib) from the current connection, based on a series (table) arugment, ID string, and date range.
// Two examples:
// SELECT * FROM daily_energy_readings WHERE time >= '2012-04-12 10:30:00.0000000' and time < '2012-05-12 10:30:00.0000000'
// SELECT * FROM daily_energy_readings WHERE time >= 1392249600000000000 and time < 1393545600000000000
// TODO: return pointer or contents?
// TODO: These start/ends need to be time types
func QueryDateRangeWithMAC(influxClient client.Client, series string, macID string, startDate string, endDate string) *client.Response {
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE time >= '%s' AND time < '%s' AND mac_id = '%s'", series, startDate, endDate, macID)

	fmt.Println(queryString)

	query := client.Query{
		Command:  queryString,
		Database: DatabaseName,
	}

	resp, err := influxClient.Query(query)
	if err != nil {
		log.Println(err)
		fmt.Println("Probably: could not connect to", InfluxAddress)
		// log.Fatal(err)
	}

	return resp
}

func QueryDateRange(influxClient client.Client, series string, startDate string, endDate string) *client.Response {
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE time >= '%s' AND time < '%s'", series, startDate, endDate)

	fmt.Println(queryString)

	query := client.Query{
		Command:  queryString,
		Database: DatabaseName,
	}

	resp, err := influxClient.Query(query)
	if err != nil {
		log.Println(err)
		fmt.Println("Probably: could not connect to", InfluxAddress)
		// log.Fatal(err)
	}

	return resp
}
