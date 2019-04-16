package influx

import (
	"log"
	"github.com/influxdata/influxdb1-client/v2"
	//"github.com/influxdata/influxdb/client/v2"

	"os"
	"fmt"
)

const (
	DATABASE_NAME = "londondb"  // TODO: Move to schema

	verbose = false
)

var (
	InfluxClient client.Client
	
	// ENVs
	INFLUX_ADDRESS string
)

func Connect() {
	// ENVs
	INFLUX_ADDRESS = os.Getenv("INFLUX_ADDRESS")


	// Influx
	//https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
	// TODO: This works but seems bad/ugly. Need to work out maybe a pointer structure or something
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: INFLUX_ADDRESS,
	})
	if err != nil {
		log.Fatal(err)
	}

	InfluxClient = influxClient
	defer InfluxClient.Close()
	// TODO: Test connection here?
}

// Writes
func WritePoint(clnt client.Client, pt *client.Point) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DATABASE_NAME,
		// Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}
	if verbose {fmt.Println("Created Batchpoint")}
	
	// Add to batch
	bp.AddPoint(pt)
	if verbose {fmt.Println("Added point to table")	
	fmt.Println()}

	// Write
	if verbose {fmt.Println("Write Batch")}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

// TODO: Batch
func WritePointBatch(clnt client.Client, pts []*client.Point) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DATABASE_NAME,
		// Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}
	if verbose {fmt.Println("Created Batchpoint")}
	
	for _, pt := range pts {
	
		// Add to batch
		bp.AddPoint(pt)

		if verbose {fmt.Println("Added point to table")	
		fmt.Println()}
	}

	// Write
	if verbose {fmt.Println("Write Batch")}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}	

// Reads
// Two examples:
// SELECT * FROM daily_energy_readings WHERE time >= '2012-04-12 10:30:00.0000000' and time < '2012-05-12 10:30:00.0000000'
// SELECT * FROM daily_energy_readings WHERE time >= 1392249600000000000 and time < 1393545600000000000
// TODO: return pointer or contents?
// TODO: These start/ends need to be time types
func QueryDateRange(series string, macId string, startDate string, endDate string) *client.Response {
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE time >= '%s' AND time < '%s' AND mac_id = '%s'", series, startDate, endDate, macId)
	
	query := client.Query{
		Command: queryString,
		Database: DATABASE_NAME,
	}
	
	resp, err := InfluxClient.Query(query)
	if err != nil {
		log.Println(err)
		fmt.Println("Probably: could not connect to", INFLUX_ADDRESS)
		// log.Fatal(err)
	}

	return resp

	// res, err := resp.Results[0].Series[0].Values[0][1].(json.Number).Float64()

	// fmt.Println(fmt.Sprintf("SELECT * FROM %s WHERE time >= '%s' and time < '%s'", series, startDate, endDate))
		
	// q := client.Query{
	// 	Command: fmt.Sprintf("SELECT * FROM '%s' WHERE time >= '%s' and time < '%s'", series, startDate, endDate)
	// 	//Command:  fmt.Sprintf("select mean(cpu_usage) from node_status where cluster = '%s'", cluster),
	// 	Database: DATABASE_NAME,
	// }
	
	// resp, err := c.Query(q)
	
	// // ...
	
	// res, err := resp.Results[0].Series[0].Values[0][1].(json.Number).Float64()

}

