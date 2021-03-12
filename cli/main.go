package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cima-lexis/lexisdn/config"
	"github.com/cima-lexis/lexisdn/fetcher"
	"github.com/cima-lexis/lexisdn/webdrops"
	"github.com/meteocima/dewetra2wrf/trusted"
	"github.com/meteocima/radar2wrf/radar"
)

func usage(errmsg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, errmsg, args...)
	fmt.Fprint(os.Stderr, "\n\n")
	fmt.Fprintln(os.Stderr, `Usage: lexisdn STARTDATE [DOWNLOAD_TYPE ...]
	STARTDATE - Satrt date/time of the simulation, in format YYYYMMDDHH
	DOWNLOAD_TYPE - types of data to download. One or more of "RISICO" | "CONTINUUM" | "WRFDAIT" | "WRFDAFR", separated by space
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
			continue
		case "CONTINUUM":
			continue
		case "WRFDAIT":
			continue
		case "WRFDAFR":
			continue
		default:
			usage("Invalid DOWNLOAD_TYPE argument `%s`.", downloadType)
		}
	}
}

func fatalIfError(err error, msgerr string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf(msgerr, err))
		os.Exit(1)
	}
}

func main() {

	config.Init()

	checkArguments()

	startDateWRF, err := time.Parse("2006010215", os.Args[1])
	fatalIfError(err, "date not valid: %w")

	fmt.Println(startDateWRF.Format("2006010215"))

	for _, downloadType := range os.Args[2:] {
		switch downloadType {
		case "RISICO":
			err = fetcher.RisicoSensorsMaps(startDateWRF)
			fatalIfError(err, "Error fetching wunderground observations maps for RISICO: %w")
		case "CONTINUUM":
			err = fetcher.ContinuumSensors(startDateWRF, webdrops.ItalyDomain)
			fatalIfError(err, "Error fetching wunderground observations for CONTINUUM: %w")
		case "WRFDAIT":
			getConvertObs(startDateWRF, getConvertStationsSync)
			getConvertObs(startDateWRF, getConvertRadarSync)

			os.RemoveAll("WRFDA")
		case "WRFDAFR":
			// TODO: use france domain here
			getConvertObs(startDateWRF, getConvertStationsSync)
			//getRadars(err, sess, startDateWRF)

			os.RemoveAll("WRFDA")
		}
	}

}

type getConvertFnT func(sess webdrops.Session, dt time.Time)

func getConvertRadarSync(sess webdrops.Session, dt time.Time) {
	err := fetcher.WrfdaRadars(sess, dt)
	fatalIfError(err, "Error WrfdaRadars wunderground observations for WRFDA: %w")

	convertRadar(dt, &err)
	convertRadar(dt.Add(-3*time.Hour), &err)
	convertRadar(dt.Add(-6*time.Hour), &err)
	fatalIfError(err, "Error convertRadar wunderground observations for WRFDA: %w")

}

func getConvertStationsSync(sess webdrops.Session, dt time.Time) {
	err := fetcher.WrfdaSensors(sess, dt, webdrops.ItalyDomain)
	fatalIfError(err, "Error fetching wunderground observations for WRFDA: %w")

	convertStations(dt, &err)
	convertStations(dt.Add(-3*time.Hour), &err)
	convertStations(dt.Add(-6*time.Hour), &err)
	fatalIfError(err, "Error convertStations wunderground observations for WRFDA: %w")

}

func getConvertObs(startDateWRF time.Time, getConvertFn getConvertFnT) {
	var sess webdrops.Session
	err := sess.Login()
	fatalIfError(err, "Error during login: %w")

	getConvertFn(sess, startDateWRF)
	fatalIfError(err, "cannot convert radar data: %w")
}

func convertRadar(date time.Time, err *error) {
	if *err != nil {
		return
	}

	dtS := date.Format("2006010215")
	fmt.Printf("Converting radar %s\n", dtS)
	dir := "WRFDA/RADARS/" + dtS

	reader, e := radar.Convert(dir, dtS)
	if e != nil {
		*err = e
		return
	}
	outfile, e := os.OpenFile("ob.radar."+dtS, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
	if e != nil {
		*err = e
		return
	}

	outfileBuff := bufio.NewWriter(outfile)

	defer outfile.Close()
	_, *err = io.Copy(outfileBuff, reader)

}

func convertStations(date time.Time, err *error) {
	if *err != nil {
		return
	}

	dtS := date.Format("2006010215")
	fmt.Printf("Converting stations %s\n", dtS)

	*err = trusted.Get(
		trusted.DewetraFormat,
		"WRFDA/SENSORS/"+dtS,
		"ob.ascii."+dtS,
		"24,64,-19,48",
		date,
	)
}

/*
leftlon = -19.0
rightlon = 48.0
toplat = 64.0
bottomlat = 24.0
*/
