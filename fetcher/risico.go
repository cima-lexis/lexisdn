package fetcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cima-lexis/lexisdn/webdrops"
)

// RisicoSensorsMaps retrieves a set of sensors maps and save them to files.
// Maps downloaded are all wunderground sensors data need for a Risico simulation
//
// Needed dewetra sensor classes are:
//  * IGROMETRO
//  * TERMOMETRO
//  * PLUVIOMETRO
//
// Being D the start date and time of Continuum simulation, needed
// observations are all that from time D-72H to D.
//
// Observations are saved, under cwd, on directory RISICO/SENSORS/
// with name <SENSORCLASS>.nc
func RisicoSensorsMaps(sess webdrops.Session, simulStartDate time.Time, domain webdrops.Domain) error {
	fetcher := risicoSession{
		sess:   sess,
		domain: domain,
	}

	from := simulStartDate.Add(-24 * time.Hour)
	to := simulStartDate

	fetcher.fetchSensorMap("IGROMETRO", from, to)
	fetcher.fetchSensorMap("TERMOMETRO", from, to)
	fetcher.fetchSensorMap("PLUVIOMETRO", from, to)

	return fetcher.sessError
}

type risicoSession struct {
	sessError error
	sess      webdrops.Session
	domain    webdrops.Domain
}

func (fetcher *risicoSession) fetchSensorMap(class string, from, to time.Time) {
	if fetcher.sessError != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Downloading observations for %s from %s to %s\n", class, from.Format("02/01/2006 15"), to.Format("02/01/2006 15"))
	sensorsMap, err := fetcher.sess.SensorsMap(class, from, to, fetcher.domain)
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error fetching sensors map: %w", err)
		return
	}

	mapFilePath := filepath.Join(
		"RISICO/SENSORS",
		fmt.Sprintf("%s.nc", class),
	)

	err = os.MkdirAll(filepath.Dir(mapFilePath), os.FileMode(0755))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error creating directory `%s`: %w", filepath.Dir(mapFilePath), err)
		return
	}

	fmt.Fprintf(os.Stderr, "Saving observations to %s\n", mapFilePath)
	err = ioutil.WriteFile(mapFilePath, sensorsMap, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error saving sensors data to `%s`: %w", mapFilePath, err)
	}

}
