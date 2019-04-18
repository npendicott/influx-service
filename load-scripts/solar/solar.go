package main

import (
	"github.com/npendicott/influx-service/csvReader"
	"github.com/npendicott/influx-service/schemas/solar/controlReading"

	"github.com/influxdata/influxdb1-client/v2"

	"log"
	_"strings"
	_"strconv"
	"fmt"
	
)

const (

	HH_TIMESTAMP_FORMAT = "2006-01-02 15:04:05.0000000"
)


func main() {



	// Influx
	// TODO: new client ugh
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://127.0.0.1:8086",
		// Addr: "http://influxdb:8086",
		// Addr: "http://192.168.0.24:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer influxClient.Close()

	
	// Control
	controlData := csvReader.GetDataArray("../../data/solar/control.csv")

	// processedReading := controlReading.ProcessCSVEnergyReading(controlData[5])
	// controlReading.Print(processedReading)

	
	for i, reading := range controlData {
		fmt.Println("Time", reading[4])
		fmt.Println("Reading:", i)

		processedReading := controlReading.ProcessCSVEnergyReading(reading)
		controlReading.Print(processedReading)
		// _ = processedReading
		// _ = i		

		// readingBatch = append(readingBatch, processedReading)
		// batch++

		// // Block write logic
		// if batch == 4000 || i == len(blockData) - 1 {
		// 	halfhourlyReading.WriteEnergyReadingBatch(influxClient, readingBatch)
		// 	batch = 0
		// 	// TODO: Better way to clear this slice? Maybe preserve the space
		// 	readingBatch = nil
		// }

		fmt.Println()
	}
	
	
	
	// for blockIndex := 0; blockIndex < 3; blockIndex++ {
	// 	blockData := csvReader.GetDataArray(londonPathRoot + frequencyLevel + "block_" + strconv.Itoa(blockIndex) +".csv")

	// 	batch := 0
	// 	var readingBatch []halfhourlyReading.HalfhourlyReading
	// 	for i, reading := range blockData {
	// 		fmt.Println("Block:", blockIndex)
	// 		fmt.Println("Reading:", i)

	// 		processedReading := halfhourlyReading.ProcessCSVEnergyReading(reading, HH_TIMESTAMP_FORMAT)

	// 		readingBatch = append(readingBatch, processedReading)
	// 		batch++

	// 		// Block write logic
	// 		if batch == 4000 || i == len(blockData) - 1 {
	// 			halfhourlyReading.WriteEnergyReadingBatch(influxClient, readingBatch)
	// 			batch = 0
	// 			// TODO: Better way to clear this slice? Maybe preserve the space
	// 			readingBatch = nil
	// 		}

	// 		fmt.Println()
	// 	}
	// }

}
