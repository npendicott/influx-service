package halfhourlyReading

import (
	"github.com/npendicott/influx-service/influx"

	"github.com/influxdata/influxdb1-client/v2"
//	"github.com/influxdata/influxdb/client/v2"

	"log"
	"time"

	"strings"
	"strconv"	

	"fmt"	
)

// Influx Consts
const (
	TABLE_NAME = "halfhourly_energy_readings"

	verbose = true
)

// Schema
// LCLid,tstp,energy(kWh/hh)
type HalfhourlyReading struct {
	Timestamp time.Time `json:"timestamp"`
	MacId string `json:"mac_id"`	
	// TODO: unit in name
	Energy float64 `json:"energy"`
}


/// CSV
// Process string array
func ProcessCSVEnergyReading(record []string, timestampFormat string) HalfhourlyReading {
	// MAC ID
	macId := record[0]

	// Timestamp
	timestamp, err := time.Parse(timestampFormat, record[1])

	if err != nil {
 	   fmt.Println(err)
	}	

	// Energy
	energy := StripAndParse(record[2])

	energyReading := HalfhourlyReading {
		Timestamp: timestamp,
		MacId: macId,
		Energy: energy,
	}

	return energyReading
}

func StripAndParse(str string) float64 {
	// float, err := strconv.ParseFloat(str, 64)		
	// if err != nil {
	// 	log.Fatal(err)
	// }

	strpString := strings.Replace(str, " ", "", -1)

	fmt.Println(strpString)

	if strpString == "Null"{
		return -1
	}
	float, err := strconv.ParseFloat(strpString, 64)		
	if err != nil {
		log.Fatal(err)
	}

	return float
}


/// InfluxDB
// Searialize
func createEnergyReadingPoint(reading HalfhourlyReading) *client.Point{	
	tags := map[string]string{
		"mac_id": reading.MacId,
	}

	fields := map[string]interface{}{
		"energy_kwh": reading.Energy,
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
func WriteEnergyReading(clnt client.Client, reading HalfhourlyReading) {
	pt := createEnergyReadingPoint(reading)
	influx.WritePoint(clnt, pt)
}

func WriteEnergyReadingBatch(clnt client.Client, readings []HalfhourlyReading) {
	var pts []*client.Point

	for _, reading := range readings {
		pt := createEnergyReadingPoint(reading)

		pts = append(pts, pt)
	}

	influx.WritePointBatch(clnt, pts)
}

// TODO: Create an object from a point
// TODO: Reads


// Print
func PrintEnergyReading(reading HalfhourlyReading) {
	fmt.Println("Time: " + reading.Timestamp.String())
	fmt.Println("MAC: ", reading.MacId)
	fmt.Println("Energy: ", reading.Energy)
	fmt.Println()
}
