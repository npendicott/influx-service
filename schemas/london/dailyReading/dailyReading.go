package dailyReading

import (
	"influx-client-london/data/influx"
	"influx-client-london/data/schemas/csvParseUtil"

	"github.com/influxdata/influxdb1-client/v2"
	//"github.com/influxdata/influxdb/client/v2"

	"log"
	"time"

	_"strings"
	_"strconv"	

	"fmt"	
)

// Influx Consts
const (
	// TODO: SERIES
	TABLE_NAME = "daily_energy_readings"

	verbose = true
)

// Schema
// LCLid,day,energy_median,energy_mean,energy_max,energy_count,energy_std,energy_sum,energy_min
type DailyReading struct {
	Timestamp time.Time `json:"timestamp"`
	MacId string `json:"mac_id"`	
	// TODO: unit in name
	EnergyMedian float64 `json:"energy_median"`
	EnergyMean float64 `json:"energy_mean"`
	EnergyMax float64 `json:"energy_max"`
	EnergyCount float64 `json:"energy_count"`
	EnergyStd float64 `json:"energy_std"`
	EnergySum float64 `json:"energy_sum"`
	EnergyMin float64 `json:"energy_min"`
}


/// CSV
// Process string array
func ProcessCSVEnergyReading(record []string, timestampFormat string) DailyReading {
	// MAC ID
	macId := record[0]

	// Timestamp
	timestamp, err := time.Parse(timestampFormat, record[1])

	if err != nil {
 	   fmt.Println(err)
	}	

	// Energy
	energyMedian := csvParseUtil.StripAndParse(record[2])
	energyMean := csvParseUtil.StripAndParse(record[3])
	energyMax := csvParseUtil.StripAndParse(record[4])
	energyCount := csvParseUtil.StripAndParse(record[5])
	energyStd := csvParseUtil.StripAndParse(record[6])
	energySum := csvParseUtil.StripAndParse(record[7])
	energyMin := csvParseUtil.StripAndParse(record[8])

	energyReading := DailyReading {
		Timestamp: timestamp,
		MacId: macId,
		EnergyMedian: energyMedian,
		EnergyMean: energyMean,
		EnergyMax: energyMax,
		EnergyCount: energyCount,
		EnergyStd: energyStd,
		EnergySum: energySum,
		EnergyMin: energyMin,
	}

	return energyReading
}

/// InfluxDB
// Searialize
func createEnergyReadingPoint(reading DailyReading) *client.Point{	
	tags := map[string]string{
		"mac_id": reading.MacId,
	}

	fields := map[string]interface{}{
		"energy_median": reading.EnergyMedian,
		"energy_mean": reading.EnergyMean,
		"energy_max": reading.EnergyMax,
		"energy_count": reading.EnergyCount,
		"energy_std": reading.EnergyStd,
		"energy_sum": reading.EnergySum,
		"energy_min": reading.EnergyMin,
	}

	pt, err := client.NewPoint(
		TABLE_NAME,
		tags,
		fields,
		reading.Timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}

	return pt
}

// Writes
func WriteEnergyReading(clnt client.Client, reading DailyReading) {
	pt := createEnergyReadingPoint(reading)
	influx.WritePoint(clnt, pt)
}

func WriteEnergyReadingBatch(clnt client.Client, readings []DailyReading) {
	var pts []*client.Point

	for _, reading := range readings {
		pt := createEnergyReadingPoint(reading)

		pts = append(pts, pt)
	}
	fmt.Println(len(pts))
	influx.WritePointBatch(clnt, pts)
}

// TODO: Create an object from a point

// Reads
// TODO: Serialize to model from resp if we really need it (i.e. invert resp/args)
// func ReadEnergyReadingsBatch(startDate string, endDate string) *client.Response {
// 	resp := influx.QueryDateRange(TABLE_NAME, startDate, endDate)

// 	fmt.Println(resp.Results[0])

// 	return resp
// }




// Print
func PrintEnergyReading(reading DailyReading) {
	fmt.Println("Time: " + reading.Timestamp.String())
	fmt.Println("MAC: ", reading.MacId)
	fmt.Println("Energy: ", reading.EnergyMean)
	fmt.Println()
}
