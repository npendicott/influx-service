package hhourly

import (

	//	"github.com/influxdata/influxdb/client/v2"

	"log"
	"time"

	"strconv"
	"strings"

	"fmt"
)

// Influx Consts
const (
	TimestampFormat = "2006-01-02 15:04:05.0000000"

	verbose = true
)

// Schema
// LCLid,tstp,energy(kWh/hh)
type HalfhourlyReading struct {
	Timestamp time.Time `json:"timestamp" influx:"timestamp"`
	MacID     string    `json:"mac_id" influx:"tag"`
	Energy    float64   `json:"energy_median" influx:"field"`
}

/// CSV
// Process string array
func ProcessCSVEnergyReading(record []string) HalfhourlyReading {
	// MAC ID
	macId := record[0]

	// Timestamp
	timestamp, err := time.Parse(TimestampFormat, record[1])

	if err != nil {
		fmt.Println(err)
	}

	// Energy
	energy := stripAndParse(record[2])

	energyReading := HalfhourlyReading{
		Timestamp: timestamp,
		MacID:     macId,
		Energy:    energy,
	}

	return energyReading
}

func stripAndParse(str string) float64 {
	// float, err := strconv.ParseFloat(str, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	strpString := strings.Replace(str, " ", "", -1)

	fmt.Println(strpString)

	if strpString == "Null" {
		return -1
	}
	float, err := strconv.ParseFloat(strpString, 64)
	if err != nil {
		log.Fatal(err)
	}

	return float
}

// /// InfluxDB
// // Searialize
// func createEnergyReadingPoint(reading HalfhourlyReading) *client.Point {
// 	tags := map[string]string{
// 		"mac_id": reading.MacId,
// 	}

// 	fields := map[string]interface{}{
// 		"energy_kwh": reading.Energy,
// 	}

// 	pt, err := client.NewPoint(
// 		TABLE_NAME,
// 		tags,
// 		fields,
// 		reading.Timestamp,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return pt
// }

// // Writes
// func WriteEnergyReading(clnt client.Client, reading HalfhourlyReading) {
// 	pt := createEnergyReadingPoint(reading)
// 	influx.WritePoint(clnt, pt)
// }

// func WriteEnergyReadingBatch(clnt client.Client, readings []HalfhourlyReading) {
// 	var pts []*client.Point

// 	for _, reading := range readings {
// 		pt := createEnergyReadingPoint(reading)

// 		pts = append(pts, pt)
// 	}

// 	influx.WritePointBatch(clnt, pts)
// }

// // TODO: Create an object from a point
// // TODO: Reads

// // Print
// func PrintEnergyReading(reading HalfhourlyReading) {
// 	fmt.Println("Time: " + reading.Timestamp.String())
// 	fmt.Println("MAC: ", reading.MacId)
// 	fmt.Println("Energy: ", reading.Energy)
// 	fmt.Println()
// }
