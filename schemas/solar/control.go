package control

import (
	client "github.com/influxdata/influxdb1-client/v2"

	"fmt"
	"strconv"
	"strings"

	"log"
	"time"
)

// Influx Consts
const (
	TableName       = "control_readings" // TODO: Common readings schema on Influx?
	TimestampFormat = "1/02 (15:4:5)"    // "2006-01-02 15:04:05.0000000"

	verbose = true
)

// Reading is the base schema for all measurements collected by the PV control system.
type Reading struct {
	Timestamp time.Time `json:"timestamp"` // 4/15 (13:32:20)
	VoltAL    float64   `json:"volt_al"`   // Val = 15.22
	VoltCE    float64   `json:"volt_ce"`   // Vce = 12.46
	Current   float64   `json:"current"`   // Ic = 0.73
	Power     float64   `json:"power"`     // P = 60.96
}

// ProcessCSVReading takes a string array (from the CSV output of the microcontroller) and builds a Reading.
func ProcessCSVReading(record []string) Reading {
	// Timestamp
	timestamp, err := time.Parse(TimestampFormat, strings.TrimSpace(record[4]))

	if err != nil {
		fmt.Println(err)
	}

	// Voltage
	voltAL := stripAndParse(record[0], "Val =")
	voltCE := stripAndParse(record[1], "Vce =")

	// Current
	current := stripAndParse(record[2], "Ic =")

	// Power
	power := stripAndParse(record[3], "P =")

	energyReading := Reading{
		Timestamp: timestamp,
		VoltAL:    voltAL,
		VoltCE:    voltCE,
		Current:   current,
		Power:     power,
	}

	return energyReading
}

func stripAndParse(str string, label string) float64 {
	strpString := strings.Replace(str, label, "", -1)
	strpString = strings.Replace(strpString, " ", "", -1)

	if strpString == "Null" {
		return -1
	}
	float, err := strconv.ParseFloat(strpString, 64)
	if err != nil {
		log.Fatal(err)
	}

	return float
}

// Influx

// CreatePoint takes a Reading and generates an InfluxDB Point.
func CreatePoint(reading Reading) *client.Point {
	tags := map[string]string{
		//"mac_id": reading.MacId,
	}

	fields := map[string]interface{}{
		"volt_al": reading.VoltAL,
		"volt_ce": reading.VoltCE,
		"current": reading.Current,
		"power":   reading.Power,
	}

	pt, err := client.NewPoint(
		TableName,
		tags,
		fields,
		reading.Timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}

	return pt
}

// Print all the struct vals to the console.
func Print(reading Reading) {
	fmt.Println("Time: " + reading.Timestamp.String())
	fmt.Println("VoltAL: ", reading.VoltAL)
	fmt.Println("VoltCE: ", reading.VoltCE)
	fmt.Println("Current: ", reading.Current)
	fmt.Println("Power: ", reading.Power)
}
