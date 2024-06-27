package cmd

import (
	"cli-go-test/internal"
	"encoding/csv"
	"fmt"
	"log"
	"os"

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

func (i *ImportCmd) Run(ctx *internal.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env: %s", err)
	}

	iseUri := os.Getenv("ISE_URL")
	iseUser := os.Getenv("ISE_USER")
	isePasswd := os.Getenv("ISE_PASSWD")

	fmt.Printf(iseUri, iseUser, isePasswd)
	fmt.Println("")

	deviceRecords := importCSV(i.Filename)

	for _, v := range deviceRecords {
		fmt.Printf("mac: %s, hostname: %s, endpointgroup: %s\n", v.Mac, v.Hostname, v.EndpointGroup)
	}
	return nil
}
