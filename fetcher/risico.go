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
// Downloaded maps contains an interpolation of all wunderground sensors observations
// needed for a Risico simulation.
//
// Needed dewetra sensor classes are:
//  * IGROMETRO
//  * TERMOMETRO
//  * PLUVIOMETRO
//
// Being D the start date and time of Risico simulation, needed
// observations are all that from time D-72H to D.
// Maps are generated in step of 12 hours each, so to produce 72 hours of map
// 6 sets of maps must be created.
//
// Observations are saved, under cwd, on directory RISICO/SENSORS/<STEP START DATE>
// with name <SENSORCLASS>.nc
func RisicoSensorsMaps(sess webdrops.Session, simulStartDate time.Time) error {
	fetcher := risicoSession{
		sess: sess,
	}

	for step := 6; step >= 1; step-- {
		from := simulStartDate.Add(-12 * time.Duration(step) * time.Hour)
		to := from.Add(12 * time.Hour)

		fetcher.fetchSensorMap("PLUVIOMETRO", from, to)
		fetcher.fetchSensorMap("IGROMETRO", from, to)
		fetcher.fetchSensorMap("TERMOMETRO", from, to)

		if fetcher.sessError != nil {
			break
		}
	}

	return fetcher.sessError
}

type risicoSession struct {
	sessError error
	sess      webdrops.Session
}

func (fetcher *risicoSession) fetchSensorMap(class string, from, to time.Time) {
	if fetcher.sessError != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Downloading observations map for %s from %s to %s\n", class, from.Format("02/01/2006 15"), to.Format("02/01/2006 15"))
	sensorsMap, err := fetcher.sess.SensorsMap(class, from, to, webdrops.GroupDPC)
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error fetching observations map: %w", err)
		return
	}

	mapFilePath := filepath.Join(
		"RISICO/SENSORS",
		from.Format("2006010215"),
		fmt.Sprintf("%s.nc", class),
	)

	err = os.MkdirAll(filepath.Dir(mapFilePath), os.FileMode(0755))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error creating directory `%s`: %w", filepath.Dir(mapFilePath), err)
		return
	}

	fmt.Fprintf(os.Stderr, "Saving observations map to %s\n", mapFilePath)
	err = ioutil.WriteFile(mapFilePath, sensorsMap, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error saving observations map to `%s`: %w", mapFilePath, err)
	}

}
