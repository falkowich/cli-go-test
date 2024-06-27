package cmd

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"cli-go-test/internal"
	"github.com/labstack/gommon/log"
	"github.com/schollz/progressbar/v3"
)

type Response struct {
	Results []struct {
		Location struct {
			City string `json:"city"`
		} `json:"location"`
	} `json:"results"`
}

type GenerateCmd struct {
	Records  int    `help:"Number of records to generate." default:"10"`
	Group    string `help:"Groupname for the device." default:"Cisco-AP"`
	Filename string `help:"Path to the saved csv output."`
}

func (g *GenerateCmd) Run(ctx *internal.Context) error {

	if g.Filename != "" {

		file, err := os.Create(g.Filename)
		if err != nil {
			log.Fatalf("Error creating file: %s", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"mac", "hostname", "endpointgroup", "ipv4"}
		if err := writer.Write(header); err != nil {
			log.Fatalf("Error writing header: %s", err)
		}

		fmt.Printf("Generating %d records and saving to %s\n", g.Records, g.Filename)
		bar := progressbar.Default(int64(g.Records))
		for i := 0; i < g.Records; i++ {
			mac, err := GenerateMacAddr()
			if err != nil {
				log.Fatalf("Mac address generation failed:  %s", err)
			}
			hostname, err := GenerateHostname()
			if err != nil {
				log.Fatalf("Hostname generation failed: %s", err)
			}

			ipAddr, err := GenerateIPv4()
			if err != nil {
				log.Fatalf("IPv4 generation failed: %s", err)
			}
			record := []string{mac.String(), hostname, g.Group, ipAddr.String()}
			if err := writer.Write(record); err != nil {
				log.Fatalf("Error writing to csv: %s")
			}
			bar.Add(1)
		}
	} else {
		fmt.Printf("Generating %d records\n", g.Records)

		for i := 0; i < g.Records; i++ {
			mac, err := GenerateMacAddr()
			if err != nil {
				log.Fatalf("Mac address generation failed:  %s", err)
			}
			hostname, err := GenerateHostname()
			if err != nil {
				log.Fatalf("Hostname generation failed: %s", err)
			}

			ipAddr, err := GenerateIPv4()
			if err != nil {
				log.Fatalf("IPv4 generation failed: %s", err)
			}

			fmt.Printf("Mac: %s, Hostname: %s, EndpointGroup: %s, IPv4: %s\n", mac, hostname, g.Group, ipAddr)
		}
	}
	return nil
}

func cleanAndConcatenate(input string) string {
	cleanedString := strings.TrimSpace(input)
	cleanedString = strings.ToLower(cleanedString)
	cleanedString = strings.ReplaceAll(cleanedString, " ", "")

	result := "xy-ap-" + cleanedString + "-01.net.example.com"

	return result
}

func GenerateMacAddr() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	buf[0] |= 2

	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac, nil
}

func GenerateHostname() (string, error) {
	res, err := http.Get("https://randomuser.me/api/?nat=gb,us")
	if err != nil {
		return "", err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resObject Response
	err = json.Unmarshal(resData, &resObject)
	if err != nil {
		return "", err
	}

	if len(resObject.Results) > 0 && resObject.Results[0].Location.City != "" {
		hostname := resObject.Results[0].Location.City
		cleanedHostname := cleanAndConcatenate(hostname)

		return cleanedHostname, err

	}

	return "", fmt.Errorf("city not found in response")
}

func GenerateIPv4() (net.IP, error) {
	ip := make([]byte, 4)
	_, err := rand.Read(ip)
	if err != nil {
		return nil, err
	}

	return net.IPv4(ip[0], ip[1], ip[2], ip[3]), nil
}
