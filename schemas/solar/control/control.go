package control

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Influx Consts
const (
	TimestampFormat = "1/02 (15:4:5)" // "2006-01-02 15:04:05.0000000"

	verbose = true
)

// Reading is the base schema for all measurements collected by the PV control system.
type Reading struct {
	Timestamp time.Time `json:"timestamp" influx:"timestamp"` // 4/15 (13:32:20)
	VoltAL    float64   `json:"volt_al" influx:"field"`       // Val = 15.22
	VoltCE    float64   `json:"volt_ce" influx:"field"`       // Vce = 12.46
	Current   float64   `json:"current" influx:"field"`       // Ic = 0.73
	Power     float64   `json:"power" influx:"field"`         // P = 60.96
}

// Print all the struct vals to the console.
func Print(reading Reading) {
	fmt.Println("Time: " + reading.Timestamp.String())
	fmt.Println("VoltAL: ", reading.VoltAL)
	fmt.Println("VoltCE: ", reading.VoltCE)
	fmt.Println("Current: ", reading.Current)
	fmt.Println("Power: ", reading.Power)
}

// CSV

// ProcessCSVReading takes a string array (from the CSV output of the microcontroller) and builds a Reading.
func ProcessCSVReading(record []string) Reading {
	// Timestamp
	timestamp, err := time.Parse(TimestampFormat, strings.TrimSpace(record[4]))
	timestamp = timestamp.AddDate(2015, 0, 0)
	// fmt.Println("\nAdd 1 Year:", after)

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
