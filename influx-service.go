package main

import (
	// Client
	"influx-client-london/data/influx"
	// Schemas
	// TODO: Fix this: https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc
	"github.com/npendicott/influx-service/schemas/london/dailyReading"
	_"github.com/npendicott/influx-service/schemas/london/haflhourlyReading"
	// Net stuff
	"net/http"
	"encoding/json"
	"github.com/joho/godotenv"
	// OS stuff
	"log"
	"os"
	"fmt"
	// TODO: RM
	_"github.com/influxdata/influxdb1-client/v2"
	//_"github.com/influxdata/influxdb/client/v2"
	_"strings"
	_"strconv"	
)

const (
	PORT = ":9090" // TODO: Remove
)

// Utility
func addCrossSiteOriginHeader(resWrt http.ResponseWriter) http.ResponseWriter{
	resWrt.Header().Set("Access-Control-Allow-Origin", "*")
	// These all have to be lowercase, this might be bad HTTP
	// TODO: maybe an "expectedTokens" slice?
	resWrt.Header().Set("Access-Control-Allow-Headers", "mac-id")

	return resWrt
}

func unpackTime(req *http.Request, key string) (string, []byte) {
	time, ok := req.URL.Query()[key]
	if !ok || len(time[0]) < 1 {
		errString := "Url Param {0} is missing"
		errOut := []byte(errString)
		log.Println(errString, key)
		
		return "", errOut
	}

	return time[0], nil  // I guess I will leave the array accessor until I need to take it out?
}

// Routs
func root(resWrt http.ResponseWriter, req *http.Request) {
	out := []byte("Root!")
	resWrt.Write(out)
	
}

func readingsHHourly(resWrt http.ResponseWriter, req *http.Request) {
		// Parse Request 
		macId := req.Header.Get("mac-id")
		fmt.Println(macId)
	
		startString, out := unpackTime(req, "start")
		if out != nil {
			resWrt.Write(out)
			return
		}
		//fmt.Println(startString)
	
		endString, out := unpackTime(req, "end")
		if out != nil {
			resWrt.Write(out)
			return
		}
		
		// Call Influx
		resp := influx.QueryDateRange(dailyReading.TABLE_NAME, macId, startString, endString)

		// Write Response
		resWrt = addCrossSiteOriginHeader(resWrt)
		out, err := json.MarshalIndent(resp, "", "     ")
		if err != nil {
			log.Fatal(err)
		}
	
		resWrt.Write(out)
		return
}

// TODO: Err handleing
func readingsDaily(resWrt http.ResponseWriter, req *http.Request) {
	// Headers 
	macId := req.Header.Get("mac-id")
	fmt.Println(macId)

	// Parse start time
	startString, errOut := unpackTime(req, "start")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(startString)

	// Parse end time
	endString, errOut := unpackTime(req, "end")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(endString)

	// Get response
	resp := influx.QueryDateRange(dailyReading.TABLE_NAME, macId, startString, endString)

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

// Server
func main() {
	// ENVs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to InfluxDB
	INFLUX_ADDRESS := os.Getenv("INFLUX_ADDRESS")
	fmt.Println(INFLUX_ADDRESS)

	influx.Connect()

	// Routes
	http.HandleFunc("/", root)
	http.HandleFunc("/readings/daily", readingsDaily)
	http.HandleFunc("/readings/hhourly", readingsHHourly)

	// Start server
	fmt.Println("Listeninggggggg")
	serveErr := http.ListenAndServe(PORT, nil) 
    if serveErr != nil {
        log.Fatal("ListenAndServe: ", serveErr)
	}
}
