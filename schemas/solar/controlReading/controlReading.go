package controlReading

import (
	_"github.com/npendicott/influx-service/influx"

	_"github.com/influxdata/influxdb1-client/v2"

	"strings"
	"strconv"	
	"fmt"	

	"log"
	"time"
)

// Influx Consts
const (
	//TABLE_NAME = "halfhourly_energy_readings"
	TABLE_NAME = "control_readings"
	// TIMESTAMP_FORMAT = "2006-01-02"
	// TIMESTAMP_FORMAT = "1/02 (15:04:05)"  // "2006-01-02 15:04:05.0000000"
	TIMESTAMP_FORMAT = "1/02 (15:4:5)"  // "2006-01-02 15:04:05.0000000"

	verbose = true
)

// Schema
// Val = 15.22 ,Vce = 12.46 ,Ic = 0.73 ,P = 60.96, 4/15 (13:32:20)
type ControlReading struct {
	Timestamp time.Time `json:"timestamp"`  // 4/15 (13:32:20)
	VoltAL float64 `json:"volt_al"`  // Val = 15.22
	VoltCE float64 `json:"volt_ce"`  // Vce = 12.46
	Current float64 `json:"current"`  // Ic = 0.73
	Power float64 `json:"power"`  // P = 60.96
}


/// CSV
// Process string array
func ProcessCSVEnergyReading(record []string) ControlReading {
	// Timestamp
	// t := strings.TrimSpace(s)
	timestamp, err := time.Parse(TIMESTAMP_FORMAT, strings.TrimSpace(record[4]))
	// timestamp, err := time.Parse(TIMESTAMP_FORMAT, record[4])
	// timestamp, err := time.Parse(TIMESTAMP_FORMAT, "3/13 (13:32:48)")


	if err != nil {
 	   fmt.Println(err)
	}	

	// Voltage
	voltAL := StripAndParse(record[0], "Val =")
	voltCE := StripAndParse(record[1], "Vce =")
	
	// Current
	current := StripAndParse(record[2], "Ic =")

	// Power
	power := StripAndParse(record[3], "P =")

	energyReading := ControlReading {
		Timestamp: timestamp,
		VoltAL: voltAL,
		VoltCE: voltCE,
		Current: current,
		Power: power,
	}

	return energyReading
}

func StripAndParse(str string, label string) float64 {
	strpString := strings.Replace(str, label, "", -1)
	strpString = strings.Replace(strpString, " ", "", -1)

	if strpString == "Null"{
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
// func createEnergyReadingPoint(reading HalfhourlyReading) *client.Point{	
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


// Print
func Print(reading ControlReading) {
	fmt.Println("Time: " + reading.Timestamp.String())
	fmt.Println("VoltAL: ", reading.VoltAL)
	fmt.Println("VoltCE: ", reading.VoltCE)
	fmt.Println("Current: ", reading.Current)
	fmt.Println("Power: ", reading.Power)
}
