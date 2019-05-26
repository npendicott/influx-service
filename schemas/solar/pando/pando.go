package pando

// Vsp = 22.37 ,Vce = 22.37 ,Ic = 0.00 ,P = 0.00, Vsp = 22.50 ,Vce = 22.53 ,Ic = -0.01 ,P = -0.20, 4/2 (15:48:31)

import (
	"fmt"
	"strconv"
	"strings"

	"log"
	"time"
)

// Influx Consts
const (
	// TableName       = "pando_reading" // TODO: Common readings schema on Influx?
	TimestampFormat = "1/2 (15:4:5)" // "2006-01-02 15:04:05.0000000"

	verbose = true
)

// Reading is the base schema for all measurements collected by the PV control system.
type Reading struct {
	Timestamp time.Time `json:"timestamp" influx:"timestamp"` // 4/15 (13:32:20)
	// Pre-perturb
	VoltSP  float64 `json:"volt_sp" influx:"field"` // Vsp = 22.37
	VoltCE  float64 `json:"volt_ce" influx:"field"` // Vce = 12.46
	Current float64 `json:"current" influx:"field"` // Ic = 0.73
	Power   float64 `json:"power" influx:"field"`   // P = 60.96
	// Post-perturb
	VoltSPAdj  float64 `json:"volt_sp_adj" influx:"field"` // Vsp = 22.37
	VoltCEAdj  float64 `json:"volt_ce_adj" influx:"field"` // Vce = 12.46
	CurrentAdj float64 `json:"current_adj" influx:"field"` // Ic = 0.73
	PowerAdj   float64 `json:"power_adj" influx:"field"`   // P = 60.96
}

// ProcessCSVReading takes a string array (from the CSV output of the microcontroller) and builds a Reading.
func ProcessCSVReading(record []string) Reading {
	// Timestamp
	timestamp, err := time.Parse(TimestampFormat, strings.TrimSpace(record[8]))
	timestamp = timestamp.AddDate(2015, 0, 0)
	// fmt.Println("\nAdd 1 Year:", after)

	if err != nil {
		fmt.Println(err)
	}

	// Pre
	// Voltage
	voltSP := stripAndParse(record[0], "Vsp =")
	voltCE := stripAndParse(record[1], "Vce =")

	// Current
	current := stripAndParse(record[2], "Ic =")

	// Power
	power := stripAndParse(record[3], "P =")

	// Post
	// Voltage
	voltSPAdj := stripAndParse(record[4], "Vsp =")
	voltCEAdj := stripAndParse(record[5], "Vce =")

	// Current
	currentAdj := stripAndParse(record[6], "Ic =")

	// Power
	powerAdj := stripAndParse(record[7], "P =")

	energyReading := Reading{
		Timestamp:  timestamp,
		VoltSP:     voltSP,
		VoltCE:     voltCE,
		Current:    current,
		Power:      power,
		VoltSPAdj:  voltSPAdj,
		VoltCEAdj:  voltCEAdj,
		CurrentAdj: currentAdj,
		PowerAdj:   powerAdj,
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

// Print all the struct vals to the console.
func Print(reading Reading) {
	fmt.Println("Time: " + reading.Timestamp.String())
	fmt.Println("VoltSP: ", reading.VoltSP)
	fmt.Println("VoltCE: ", reading.VoltCE)
	fmt.Println("Current: ", reading.Current)
	fmt.Println("Power: ", reading.Power)
}
