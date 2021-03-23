package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/cima-lexis/lexisdn/config"
	"github.com/cima-lexis/lexisdn/fetcher"
	"github.com/cima-lexis/lexisdn/webdrops"
	"github.com/meteocima/dewetra2wrf"
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
//			getConvertStationsSync(startDateWRF)
			getConvertRadarSync(startDateWRF)

//			os.RemoveAll("WRFDA/SENSORS")
//			os.RemoveAll("WRFDA/RADARS")
		case "WRFDAFR":
			// TODO: use france domain here
			getConvertStationsSync(startDateWRF)
			// will be provided via DDI
			//getRadars(err, sess, startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
		}
	}

}

func getConvertRadarSync(dt time.Time) {
	err := fetcher.WrfdaRadars(dt)
	fatalIfError(err, "Error convertRadar for WRFDA: %w")
	//time.Sleep(time.Second)

	// TODO: move all this stuff to a conversion module
	instants := []time.Time{
		dt,
		dt.Add(-3 * time.Hour),
		dt.Add(-6 * time.Hour),
	}

	//	allDatesConverted := sync.WaitGroup{}
	for _, dt := range instants {
		//		allDatesConverted.Add(1)
		//		go func(dt time.Time) {
		convertRadar(dt, &err)
		//			allDatesConverted.Done()
		//		}(dt)
	}
	fatalIfError(err, "Error convertRadar for WRFDA: %w")

	//	allDatesConverted.Wait()
}

func getConvertStationsSync(dt time.Time) {
	err := fetcher.WrfdaSensors(dt, webdrops.ItalyDomain)

	// TODO: move all this stuff to a conversion module
	fatalIfError(err, "Error fetching wunderground observations for WRFDA: %w")
	instants := []time.Time{
		dt,
		dt.Add(-3 * time.Hour),
		dt.Add(-6 * time.Hour),
	}

	allDatesConverted := sync.WaitGroup{}
	for _, dt := range instants {
		allDatesConverted.Add(1)
		go func(dt time.Time) {
			var err error
			convertStations(dt, &err)
			if err != nil {
				fatalIfError(err, "Error convertStations wunderground observations for WRFDA: %w")
			}
			allDatesConverted.Done()
		}(dt)
	}

	allDatesConverted.Wait()

}

/*
func getConvertObs(startDateWRF time.Time, getConvertFn getConvertFnT) {
	//var sess webdrops.Session
	//err := sess.Login()
	//fatalIfError(err, "Error during login: %w")

	getConvertFn(startDateWRF)
	//fatalIfError(err, "cannot convert radar data: %w")
}
*/

// TODO: move all this stuff to a conversion module
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
	outfile, e := os.OpenFile("WRFDA/ob.radar."+dtS, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
	if e != nil {
		*err = e
		return
	}
	defer outfile.Close()
	outfileBuff := bufio.NewWriter(outfile)

	_, *err = io.Copy(outfileBuff, reader)
	if *err == nil {
		outfileBuff.Flush()
	}

}

// TODO: move all this stuff to a conversion module
func convertStations(date time.Time, err *error) {
	if *err != nil {
		return
	}

	dtS := date.Format("2006010215")
	fmt.Printf("Converting stations %s\n", dtS)

	*err = dewetra2wrf.Convert(
		dewetra2wrf.DewetraFormat,
		"WRFDA/SENSORS/"+dtS,
		"24,64,-19,48",
		date,
		"WRFDA/ob.ascii."+dtS,
	)

}

/*
leftlon = -19.0
rightlon = 48.0
toplat = 64.0
bottomlat = 24.0
*/
