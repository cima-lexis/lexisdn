package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
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
	DOWNLOAD_TYPE - types of data to download. One of "RISICO" | "CONTINUUM" | "ADMS" | "LIMAGRAIN" | "WRFIT" | "WRFITDPC" | "WRFFR"
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
		case "RISICO", "CONTINUUM", "ADMS", "LIMAGRAIN", "WRFIT", "WRFITDPC", "WRFFR":
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

			getConvertStationsSync(startDateWRF, italyDomain, webdrops.GroupWunderground)
			getConvertRadarSync(startDateWRF)
			getConvertStationsSync(startDateWRF.Add(-24*time.Hour), italyDomain, webdrops.GroupWunderground)
			getConvertRadarSync(startDateWRF.Add(-24 * time.Hour))
			getConvertStationsSync(startDateWRF.Add(-48*time.Hour), italyDomain, webdrops.GroupWunderground)
			getConvertRadarSync(startDateWRF.Add(-48 * time.Hour))

			os.RemoveAll("WRFDA/SENSORS")
			os.RemoveAll("WRFDA/RADARS")

		case "CONTINUUM":
			d, err := italyDomain.ToStruct()
			fatalIfError(err, "Error parsing domain: %w")

			err = fetcher.ContinuumSensors(startDateWRF, d)
			fatalIfError(err, "Error fetching wunderground observations for CONTINUUM: %w")

			getConvertStationsSync(startDateWRF, italyDomain, webdrops.GroupWunderground)
			getConvertRadarSync(startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
			os.RemoveAll("WRFDA/RADARS")

		case "WRFIT":
			getConvertStationsSync(startDateWRF, italyDomain, webdrops.GroupWunderground)
			getConvertRadarSync(startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
			//os.RemoveAll("WRFDA/RADARS")
		case "WRFITDPC":
			getConvertStationsSync(startDateWRF, italyDomain, webdrops.GroupDPC)
			getConvertRadarSync(startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
			//os.RemoveAll("WRFDA/RADARS")
		case "ADMS", "LIMAGRAIN", "WRFFR":
			// TODO: use france domain here
			getConvertStationsSync(startDateWRF, franceDomain, webdrops.GroupWunderground)
			// will be provided via DDI
			//getRadars(err, sess, startDateWRF)

			os.RemoveAll("WRFDA/SENSORS")
		}
	}
}

func getConvertRadarSync(dt time.Time) {
	var err error
	err = fetcher.WrfdaRadars(dt)
	fatalIfError(err, "Error convertRadar for WRFDA: %w")

	// TODO: move all this stuff to a conversion module
	instants := []time.Time{
		dt,
		dt.Add(-3 * time.Hour),
		dt.Add(-6 * time.Hour),
	}

	//	allDatesConverted := sync.WaitGroup{}
	for _, dt := range instants {
		convertRadar(dt, 1, &err)
		convertRadar(dt, 2, &err)
		convertRadar(dt, 3, &err)

	}
	//fatalIfError(os.RemoveAll("./dom_01"), "Error removing temp directory for domain 1")
	//fatalIfError(os.RemoveAll("./dom_02"), "Error removing temp directory for domain 2")
	//fatalIfError(os.RemoveAll("./dom_03"), "Error removing temp directory for domain 3")

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

func getConvertStationsSync(dt time.Time, domain domainDef, group webdrops.SensorGroup) {
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
	f.FetchSensorIDs("TERMOMETRO", dt, d, group)

	err = fetcher.WrfdaSensors(dt, d, group)
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

func filenameForVar(dirname, varname, dt string) string {

	pt := fmt.Sprintf("%s/%s-%s.nc", dirname, dt, varname)
	return pt
}

const regridTmplDir = "~/regrid-tmpl"

func remapBilinear(dir string, radarTime time.Time, varname string, domain int) error {
	// regrid radar netcdf file
	sourceFile := filenameForVar(dir, varname, radarTime.Format("2006010215"))
	targetFile := fmt.Sprintf("%s_dom%02d.remapped", sourceFile, domain)
	operator := fmt.Sprintf("remapbil,%s/wrfinput_d%02d.template", regridTmplDir, domain)

	cmd := exec.Command("cdo", operator, sourceFile, targetFile)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"Cannot apply bilinear remapping for variable %s of radar %s:\n"+
				"CMD: cdo %s %s %s\n"+
				"ERR: %w\n",
			varname, radarTime,
			operator, sourceFile, targetFile, err,
		)
	}
	return nil
}

func filterOutLowValues(dir string, radarTime time.Time, varname string, domain int) error {
	operator := fmt.Sprintf("where(%s < 10) %s=-9999", varname, varname)

	file := filenameForVar(dir, varname, radarTime.Format("2006010215"))
	sourceFile := fmt.Sprintf("%s_dom%02d.remapped", file, domain)
	targetFile := fmt.Sprintf("%s_dom%02d.filtered", file, domain)

	cmd := exec.Command("ncap2", "-s", operator, sourceFile, targetFile)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"Cannot filter low values for variable %s of radar %s:\n"+
				"CMD: ncap2 -s %s %s %s\n"+
				"ERR: %w\n",
			varname, radarTime,
			operator, sourceFile, targetFile, err,
		)
	}

	operator = "time=int(time)"
	sourceFile = fmt.Sprintf("%s_dom%02d.filtered", file, domain)
	targetFile = fmt.Sprintf("%s_dom%02d.timefixed", file, domain)

	cmd = exec.Command("ncap2", "-s", operator, sourceFile, targetFile)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"Cannot change type for variable %s of radar %s to `int`:\n"+
				"CMD: ncap2 -s %s %s %s\n"+
				"ERR: %w\n",
			varname, radarTime,
			operator, sourceFile, targetFile, err,
		)
	}

	//if err := os.Remove(fmt.Sprintf("%s_dom%02d.remapped", file, domain)); err != nil {
	//	return err
	//}
	//
	//if err := os.Remove(fmt.Sprintf("%s_dom%02d.filtered", file, domain)); err != nil {
	//	return err
	//}

	domainDir := fmt.Sprintf("./dom_%02d", domain)
	if err := os.MkdirAll(path.Dir(path.Join(domainDir, file)), 0755); err != nil {
		return err
	}

	targetFile = path.Join(domainDir, file)
	if err := os.Rename(fmt.Sprintf("%s_dom%02d.timefixed", file, domain), targetFile); err != nil {
		return err
	}

	return nil
}

var varnames = []string{"CAPPI2", "CAPPI3" /*, "CAPPI4"*/, "CAPPI5"}

// TODO: move all this stuff to a conversion module
func convertRadar(date time.Time, domain int, err *error) {
	if *err != nil {
		return
	}

	dtS := date.Format("2006010215")
	fmt.Printf("Converting radar %s domain %d\n", dtS, domain)
	dir := "WRFDA/RADARS/" + dtS

	for _, varname := range varnames {
		if e := remapBilinear(dir, date, varname, domain); e != nil {
			*err = e
			return
		}
		if e := filterOutLowValues(dir, date, varname, domain); e != nil {
			*err = e
			return
		}
	}

	reader, e := radar.Convert(fmt.Sprintf("./dom_%02d/%s", domain, dir), "", dtS)
	if e != nil {
		*err = e
		return
	}
	radarOutFilePath := fmt.Sprintf("WRFDA/ob.radar.%s_dom%02d", dtS, domain)
	outfile, e := os.OpenFile(radarOutFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
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
