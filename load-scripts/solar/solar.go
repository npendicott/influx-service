package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2" //"github.com/influxdata/influxdb/client/v2"

	"github.com/npendicott/influx-service/csvReader"
	"github.com/npendicott/influx-service/influx"
	"github.com/npendicott/influx-service/schemas/solar/control"

	_ "strconv"
	_ "strings"
)

func main() {
	// // ENVs
	// err := godotenv.Load("../..")  // TODO: ENV in the folder dangus
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// Influx
	influxClient := influx.GetConnection()
	defer influxClient.Close()

	// Control
	controlData := csvReader.GetDataArray("../../data/solar/control.csv")

	// var readingBatch []*client.Point
	readingBatch := make([]*client.Point, 4000)
	for i, reading := range controlData {
		// fmt.Println("Time", reading[4])
		// fmt.Println("Reading:", i)

		controlReading := control.ProcessCSVReading(reading)
		// controlPoint := control.CreatePoint(controlReading) // Old Create
		controlPoint, err := influx.Marshall(controlReading, "control")
		if err != nil {
			panic(err)
		}
		readingBatch = append(readingBatch, controlPoint)
		// control.Print(controlReading)

		if len(readingBatch)%4000 == 0 || i == len(readingBatch)-1 {
			fmt.Println(len(readingBatch))
			influx.WritePointBatch(influxClient, readingBatch)
			readingBatch = make([]*client.Point, 4000)
		}

		// fmt.Println()
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
