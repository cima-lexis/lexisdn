package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	DOWNLOAD_TYPE - types of data to download. One of "RISICO" | "CONTINUUM" | "ADMS" | "LIMAGRAIN" | "WRFIT" | "WRFFR"
	`)
	os.Exit(1)
}

func checkArguments() {
	if len(os.Args) < 2 {
		usage("Missing STARTDATE argument.")
	}

	if len(os.Args) < 3 {
		usage("DOWNLOAD_TYPE argument required.")
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
		case "ADMS":
			continue
		case "LIMAGRAIN":
			continue
		case "WRFIT":
			continue
		case "WRFFR":
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

type domainDef struct{ coord string }

var franceDomain = domainDef{"38,55,-10,12"}
var italyDomain = domainDef{"24,64,-19,48"}
var emptyDomain webdrops.Domain

// ToStruct returns a new Domain pointer
// accordingly to the given string, that must
// contains  MinLat,MaxLat,MinLon,MaxLon values,
// in that sequence, separated by commas and
// represented as floats.
func (d domainDef) ToStruct() (webdrops.Domain, error) {
	if d.coord == "" {
		return webdrops.Domain{
			MinLat: -180,
			MinLon: -90,
			MaxLat: 90,
			MaxLon: 180,
		}, nil
	}
	coords := strings.Split(d.coord, ",")

	MinLat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return emptyDomain, err
	}

	MaxLat, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return emptyDomain, err
	}

	MinLon, err := strconv.ParseFloat(coords[2], 64)
	if err != nil {
		return emptyDomain, err
	}

	MaxLon, err := strconv.ParseFloat(coords[3], 64)
	if err != nil {
		return emptyDomain, err
	}

	return webdrops.Domain{
		MinLat: MinLat,
		MinLon: MinLon,
		MaxLat: MaxLat,
		MaxLon: MaxLon,
	}, nil
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

			getConvertStationsSync(startDateWRF, italyDomain)
			getConvertRadarSync(startDateWRF)
			getConvertStationsSync(startDateWRF.Add(-24*time.Hour), italyDomain)
			getConvertRadarSync(startDateWRF.Add(-24 * time.Hour))
			getConvertStationsSync(startDateWRF.Add(-48*time.Hour), italyDomain)
			getConvertRadarSync(startDateWRF.Add(-48 * time.Hour))

			os.RemoveAll("WRFDA/SENSORS")
			os.RemoveAll("WRFDA/RADARS")

		case "CONTINUUM":
			d, err := italyDomain.ToStruct()
			fatalIfError(err, "Error parsing domain: %w")

			err = fetcher.ContinuumSensors(startDateWRF, d)
			fatalIfError(err, "Error fetching wunderground observations for CONTINUUM: %w")

			getConvertStationsSync(startDateWRF, italyDomain)
			getConvertRadarSync(startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
			os.RemoveAll("WRFDA/RADARS")

		case "WRFIT":
			getConvertStationsSync(startDateWRF, italyDomain)
			getConvertRadarSync(startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
			os.RemoveAll("WRFDA/RADARS")
		case "ADMS", "LIMAGRAIN", "WRFFR":
			// TODO: use france domain here
			getConvertStationsSync(startDateWRF, franceDomain)
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

func copyFile(src, target string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer w.Close()

	bufSrc := bufio.NewReader(r)
	bufTarget := bufio.NewWriter(w)

	_, err = io.Copy(bufTarget, bufSrc)
	return err
}

func getConvertStationsSync(dt time.Time, domain domainDef) {
	d, err := domain.ToStruct()
	fatalIfError(err, "Error parsing domain: %w")

	var sess webdrops.Session
	err = sess.Login()
	if err != nil {
		return
	}

	d, err = domain.ToStruct()
	fatalIfError(err, "Error parsing domain")

	f := fetcher.WrfdaSensorsSession{
		Sess:   sess,
		Domain: d,
	}
	f.FetchSensorIDs("TERMOMETRO", dt, d)

	err = fetcher.WrfdaSensors(dt, d)
	fatalIfError(err, "Error fetching wunderground observations for WRFDA: %w")

	// qui, ricopiare il file del registry su tutte le altre date
	// scaricate

	registrySrc := "WRFDA/SENSORS/TERMOMETRO-registry.json"

	dtCycle3 := dt
	dtCycle2 := dt.Add(-3 * time.Hour)
	dtCycle1 := dt.Add(-6 * time.Hour)

	registry1 := filepath.Join(
		"WRFDA/SENSORS",
		dtCycle1.Format("2006010215"),
		fmt.Sprintf("%s-registry.json", "TERMOMETRO"),
	)

	registry2 := filepath.Join(
		"WRFDA/SENSORS",
		dtCycle2.Format("2006010215"),
		fmt.Sprintf("%s-registry.json", "TERMOMETRO"),
	)

	registry3 := filepath.Join(
		"WRFDA/SENSORS",
		dtCycle3.Format("2006010215"),
		fmt.Sprintf("%s-registry.json", "TERMOMETRO"),
	)

	fatalIfError(copyFile(registrySrc, registry1), "unable to copy registry for cycle 1: %w")
	fatalIfError(copyFile(registrySrc, registry2), "unable to copy registry for cycle 2: %w")
	fatalIfError(copyFile(registrySrc, registry3), "unable to copy registry for cycle 3: %w")

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
			convertStations(dt, domain, &err)
			if err != nil {
				msg := fmt.Sprintf("Error converting wunderground observations of date %s: %%w", dt.Format("200601021504"))
				fatalIfError(err, msg)
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
func convertStations(date time.Time, domain domainDef, err *error) {
	if *err != nil {
		return
	}

	dtS := date.Format("2006010215")
	fmt.Printf("Converting stations %s\n", dtS)

	*err = dewetra2wrf.Convert(
		dewetra2wrf.DewetraFormat,
		"WRFDA/SENSORS/"+dtS,
		domain.coord,
		date,
		"WRFDA/ob.ascii."+dtS,
	)

}
