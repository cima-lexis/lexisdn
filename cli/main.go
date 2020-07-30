package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cima-lexis/lexisdn/config"
	"github.com/cima-lexis/lexisdn/fetcher"
	"github.com/cima-lexis/lexisdn/webdrops"
)

func usage(errmsg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, errmsg, args...)
	fmt.Fprint(os.Stderr, "\n\n")
	fmt.Fprintln(os.Stderr, `Usage: lexisdn STARTDATE [DOWNLOAD_TYPE ...]
	STARTDATE - Satrt date/time of the simulation, in format YYYYMMDDHH
	DOWNLOAD_TYPE - types of data to download. One or more of "RISICO" | "CONTINUUM" | "WRFDA", separated by space
	`)
	os.Exit(1)
}

func checkArguments() {
	if len(os.Args) < 2 {
		usage("Missing STARTDATE argument.")
	}

	if len(os.Args) < 3 {
		usage("At least one DOWNLOAD_TYPE argument required.")
	}

	_, err := time.Parse("2006010215", os.Args[1])
	if err != nil {
		usage("Invalid STARTDATE argument `%s`.", os.Args[1])
	}

	for _, downloadType := range os.Args[2:] {
		switch downloadType {
		case "RISICO":
		case "CONTINUUM":
		case "WRFDA":
			continue
		default:
			usage("Invalid DOWNLOAD_TYPE argument `%s`.", downloadType)
		}
	}
}

func fatal(msgerr string) {

}

func main() {
	config.Init()

	checkArguments()

	sess := webdrops.Session{}
	err := sess.Login()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Error during login: %w", err))
		os.Exit(1)
	}

	startDateWRF, err := time.Parse("2006010215", os.Args[1])

	for _, downloadType := range os.Args[2:] {
		switch downloadType {
		case "RISICO":
			err = fetcher.RisicoSensorsMaps(sess, startDateWRF)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("Error fetching wunderground observations maps for RISICO: %w", err))
				os.Exit(1)
			}
		case "CONTINUUM":
			err = fetcher.ContinuumSensors(sess, startDateWRF, webdrops.ItalyDomain)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("Error fetching wunderground observations for CONTINUUM: %w", err))
				os.Exit(1)
			}
		case "WRFDA":
			err = fetcher.WrfdaRadars(sess, startDateWRF)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("Error fetching radar data for WRFDA: %w", err))
				os.Exit(1)
			}

			err = fetcher.WrfdaSensors(sess, startDateWRF, webdrops.ItalyDomain)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("Error fetching wunderground observations for WRFDA: %w", err))
				os.Exit(1)
			}

		}
	}

}
