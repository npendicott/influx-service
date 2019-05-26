package main

import (
	// TODO: Remove
	"errors"

	"github.com/npendicott/influx-service/influx"
	_ "github.com/npendicott/influx-service/schemas/london/daily"
	_ "github.com/npendicott/influx-service/schemas/london/haflhourlyReading"

	client "github.com/influxdata/influxdb1-client/v2" //"github.com/influxdata/influxdb/client/v2"
	"github.com/joho/godotenv"

	// Net stuff
	"encoding/json"
	"net/http"

	// OS stuff
	"fmt"
	"log"
	"os"
	"regexp"

	// TODO: RM
	_ "github.com/influxdata/influxdb1-client/v2"
	//_"github.com/influxdata/influxdb/client/v2"
	_ "strconv"
	_ "strings"
)

const (
	PORT = ":9090" // TODO: Remove
)

var validPath = regexp.MustCompile("^/(readings|reading)/([a-zA-Z0-9]+)$")
var influxClient client.Client

// Utility
func addCrossSiteOriginHeader(resWrt http.ResponseWriter) http.ResponseWriter {
	resWrt.Header().Set("Access-Control-Allow-Origin", "*")
	// These all have to be lowercase, this might be bad HTTP
	// TODO: maybe an "expectedTokens" slice?
	resWrt.Header().Set("Access-Control-Allow-Headers", "mac-id")

	return resWrt
}

func getSeries(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Series")
	}
	return m[2], nil // The title is the second subexpression.
}

// So, not usre if this is the best way, but seems better than copying for each.
// I guess these should just be headers?? I don't think so though
func getRequiredParam(req *http.Request, key string) (string, []byte) {
	result, ok := req.URL.Query()[key]
	if !ok || len(result[0]) < 1 {
		errString := fmt.Sprintf("One %s param is required.\n", key)
		errOut := []byte(errString)
		log.Println(errString)

		return "", errOut
	}

	return result[0], nil
}

// Routs
func root(resWrt http.ResponseWriter, req *http.Request) {
	out := []byte("Root!")
	resWrt.Write(out)

}

func mac_readings(resWrt http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/mac_readings/"):] // TODO: Is this really the only way to do this?

	fmt.Println(req.URL.Path[2:])
	// Headers
	macID := req.Header.Get("mac-id")
	fmt.Println(macID)

	// Get date range
	start, errOut := getRequiredParam(req, "start")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(start)

	end, errOut := getRequiredParam(req, "end")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(end)

	// Get response
	resp := influx.QueryDateRangeWithMAC(influxClient, title, macID, start, end)

	// Writing Response Body
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Fatal(err)
	}

	// Header
	resWrt = addCrossSiteOriginHeader(resWrt)

	resWrt.Write(out)
	return
}

func readings(resWrt http.ResponseWriter, req *http.Request) {
	series := req.URL.Path[len("/readings/"):] // TODO: Is this really the only way to do this?

	// Get date range
	start, errOut := getRequiredParam(req, "start")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(start)

	end, errOut := getRequiredParam(req, "end")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(end)

	// Get response
	resp := influx.QueryDateRange(influxClient, series, start, end)

	// Writing Response Body
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Fatal(err)
	}

	// Header
	resWrt = addCrossSiteOriginHeader(resWrt)

	resWrt.Write(out)
	return
}

// Specific

// func readingsHHourly(resWrt http.ResponseWriter, req *http.Request) {
// 	// Parse Request
// 	macID := req.Header.Get("mac-id")
// 	fmt.Println(macID)

// 	startString, out := unpackTime(req, "start")
// 	if out != nil {
// 		resWrt.Write(out)
// 		return
// 	}
// 	//fmt.Println(startString)

// 	endString, out := unpackTime(req, "end")
// 	if out != nil {
// 		resWrt.Write(out)
// 		return
// 	}

// 	// Call Influx
// 	resp := influx.QueryDateRange(dailyReading.TABLE_NAME, macID, startString, endString)

// 	// Write Response
// 	resWrt = addCrossSiteOriginHeader(resWrt)
// 	out, err := json.MarshalIndent(resp, "", "     ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resWrt.Write(out)
// 	return
// }

// func readingsDaily(resWrt http.ResponseWriter, req *http.Request) {
// 	fmt.Println(req.URL.Path[2:])
// 	// Headers
// 	macID := req.Header.Get("mac-id")
// 	fmt.Println(macID)

// 	// Get date range
// 	start, errOut := getRequiredParam(req, "start")
// 	if errOut != nil {
// 		resWrt.Write(errOut)
// 		return
// 	}
// 	fmt.Println(start)

// 	end, errOut := getRequiredParam(req, "end")
// 	if errOut != nil {
// 		resWrt.Write(errOut)
// 		return
// 	}
// 	fmt.Println(end)

// 	// Get response
// 	resp := influx.QueryDateRange(influxClient, "daily", macID, start, end)

// 	// Writing Response Body
// 	out, err := json.MarshalIndent(resp, "", "     ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Header
// 	resWrt = addCrossSiteOriginHeader(resWrt)

// 	resWrt.Write(out)
// 	return
// }

// Server
func main() {
	// ENVs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to InfluxDB
	InfluxAddress := os.Getenv("INFLUX_ADDRESS")
	fmt.Println(InfluxAddress)

	influxClient = influx.GetConnection()
	defer influxClient.Close()

	// Routes
	http.HandleFunc("/", root)
	http.HandleFunc("/readings/", readings)
	http.HandleFunc("/mac_readings/", mac_readings)
	// http.HandleFunc("/readings/daily", readingsDaily)
	// http.HandleFunc("/readings/hhourly", readingsHHourly)

	// Start server
	fmt.Println("Listeninggggggg")
	serveErr := http.ListenAndServe(PORT, nil)
	if serveErr != nil {
		log.Fatal("ListenAndServe: ", serveErr)
	}
}
