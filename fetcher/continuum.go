package fetcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cima-lexis/lexisdn/webdrops"
)

// ContinuumSensors retrieves a set of sensors datasets and save them to files.
// Datasets downloaded are all wunderground sensors data need for a Continuum simulation
//
// Needed dewetra sensor classes are:
//  * IGROMETRO
//  * TERMOMETRO
//  * ANEMOMETRO
//  * PLUVIOMETRO
//  * RADIOMETRO
//
// Being D the start date and time of Continuum simulation, needed
// observations are all that from time D-60H to D. Observations are
// aggregated on a hourly basis.
//
// Observations are saved, under cwd, on directory CONTINUUM/SENSORS/
// with name <SENSORCLASS>.json
func ContinuumSensors(sess webdrops.Session, simulStartDate time.Time, domain webdrops.Domain) error {
	fetcher := continuumSession{
		sess:   sess,
		domain: domain,
	}

	from := simulStartDate.Add(-60 * time.Hour)
	to := simulStartDate

	fetcher.fetchSensor("RADIOMETRO", from, to, false)
	fetcher.fetchSensor("IGROMETRO", from, to, false)
	fetcher.fetchSensor("TERMOMETRO", from, to, false)
	fetcher.fetchSensor("ANEMOMETRO", from, to, false)
	fetcher.fetchSensor("PLUVIOMETRO", from, to, false)

	return fetcher.sessError
}

type continuumSession struct {
	sessError error
	sess      webdrops.Session
	domain    webdrops.Domain
}

func (fetcher *continuumSession) fetchSensor(class string, from, to time.Time, log bool) {
	if fetcher.sessError != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Downloading sensors registry for %s\n", class)
	sensorRegistry, err := fetcher.sess.SensorsList(class, webdrops.GroupDPC)
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error fetching sensors list: %w", err)
		return
	}

	ids, err := fetcher.sess.IdFromSensorsList(sensorRegistry, fetcher.domain)
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error readings ids: %w", err)
		return
	}
	fmt.Fprintf(os.Stderr, "Found %d sensors\n", len(ids))

	if len(ids) > 0 {
		fmt.Fprintf(os.Stderr, "Downloading observations for %s from %s to %s\n", class, from.Format("02/01/2006 15"), to.Format("02/01/2006 15"))
		observations, err := fetcher.sess.SensorsData(class, ids, from, to, 3600, false)
		if err != nil {
			fetcher.sessError = fmt.Errorf("Error fetching sensors data: %w", err)
			return
		}

		jsonFilePath := filepath.Join(
			"CONTINUUM/SENSORS",
			fmt.Sprintf("%s.json", class),
		)

		err = os.MkdirAll(filepath.Dir(jsonFilePath), os.FileMode(0755))
		if err != nil {
			fetcher.sessError = fmt.Errorf("Error creating directory `%s`: %w", filepath.Dir(jsonFilePath), err)
			return
		}

		fmt.Fprintf(os.Stderr, "Saving observations to %s\n", jsonFilePath)
		err = ioutil.WriteFile(jsonFilePath, observations, os.FileMode(0644))
		if err != nil {
			fetcher.sessError = fmt.Errorf("Error saving sensors data to `%s`: %w", jsonFilePath, err)
		}
	}
	jsonAnagFilePath := filepath.Join(
		"CONTINUUM/SENSORS",
		fmt.Sprintf("anag-%s.json", class),
	)
	err = ioutil.WriteFile(jsonAnagFilePath, sensorRegistry, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error saving sensors registry data to `%s`: %w", jsonAnagFilePath, err)
	}

}
