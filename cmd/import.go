package cmd

import (
	"cli-go-test/internal"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	isegosdk "github.com/CiscoISE/ciscoise-go-sdk/sdk"
	"github.com/joho/godotenv"
)

type Device struct {
	Mac           string
	Hostname      string
	EndpointGroup string
}

type ImportCmd struct {
	Filename string `help:"Filename of csv to import from"`
}

func importCSV(filename string) []Device {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Cannot open %s: %s", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Cannot find any records: %s", err)
	}

	var deviceRecords []Device
	for _, record := range records[1:] {
		row := Device{
			Mac:           record[0],
			Hostname:      record[1],
			EndpointGroup: record[2],
		}
		deviceRecords = append(deviceRecords, row)
	}
	return deviceRecords
}

func connectToISE(deviceRecords []Device) *isegosdk.Client {

	iseUri := os.Getenv("ISE_URL")
	iseUser := os.Getenv("ISE_USER")
	isePasswd := os.Getenv("ISE_PASSWD")

	client, err := isegosdk.NewClientWithOptions("https://"+iseUri,
		iseUser, isePasswd,
		"false", "false",
		"false", "false",
	)
	if err != nil {
		log.Fatalf("Connection error: %s", err)
	}

	return client

}

func (i *ImportCmd) Run(ctx *internal.Context) error {
	//
	// TODO:
	// connect to ISE
	// get Groupid from ISE
	// Push new devices to ISE

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env: %s", err)
	}
	deviceRecords := importCSV(i.Filename)

	for _, v := range deviceRecords {
		fmt.Printf("mac: %s, hostname: %s, endpointgroup: %s\n", v.Mac, v.Hostname, v.EndpointGroup)
	}

	client := connectToISE(deviceRecords)

	res, _, err := client.EndpointIDentityGroup.GetEndpointGroupByName("Cisco-AP")

	endpointGroupID := res.EndPointGroup.ID

	truePtr := true

	for _, device := range deviceRecords {
		endpoint := &isegosdk.RequestEndpointCreateEndpoint{
			ERSEndPoint: &isegosdk.RequestEndpointCreateEndpointERSEndPoint{
				Name:                  device.Hostname,
				Description:           device.Hostname,
				GroupID:               endpointGroupID,
				StaticGroupAssignment: &truePtr,
				Mac:                   device.Mac,
			},
		}
		res, err := client.Endpoint.CreateEndpoint(endpoint)
		if err != nil {
			log.Printf("Failed to create endpoint for device %s: %s", device.Hostname, err)
			continue
		}

		log.Printf("Successfully created endpoint %s: %v", device.Hostname, res)
	}

	return nil
}
