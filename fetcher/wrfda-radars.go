package fetcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cima-lexis/lexisdn/webdrops"
)

// WrfdaSensors retrieves a set of sensors datasets and save them to files.
// Datasets downloaded are all wunderground sensors data to assimilate in a
// WRFDA simulation.
//
// Needed dewetra sensor classes are:
//  * IGROMETRO
//  * TERMOMETRO
//  * DIREZIONEVENTO
//  * ANEMOMETRO
//  * PLUVIOMETRO
//  * BAROMETRO
//
// Being D the start date and time of WRF simulation, needed dates
// of observations are at time D-6H, D-3H and D. For each hour this function downloads
// all observations from 30 minutes before to 30 minutes after. The
// observation that will be assimilated for each sensors is the one
// near the exact hour.
//
// Observations are saved, under cwd, on directory WRFDA/SENSORS/<DATE>
// with name <SENSORCLASS>.json
func WrfdaRadars(sess webdrops.Session, simulStartDate time.Time) error {
	fetcher := wrfdaRadarsSession{
		sess: sess,
	}

	fetchDate := func(instant time.Time) {
		bestInstant /*timeline*/, err := fetcher.sess.RadarTimeline(instant, false)
		if err != nil {
			fetcher.sessError = fmt.Errorf("Error downloading radars timeline: %w", err)
			return
		}
		fetcher.fetchRadar(bestInstant, "CAPPI2", true)
		fetcher.fetchRadar(bestInstant, "CAPPI3", false)
		fetcher.fetchRadar(bestInstant, "CAPPI5", false)

	}

	fetchDate(simulStartDate)
	fetchDate(simulStartDate.Add(-3 * time.Hour))
	fetchDate(simulStartDate.Add(-6 * time.Hour))

	return fetcher.sessError
}

type wrfdaRadarsSession struct {
	sessError error
	sess      webdrops.Session
	//domain    webdrops.Domain
}

func (fetcher *wrfdaRadarsSession) fetchRadar(date time.Time, varName string, log bool) {
	if fetcher.sessError != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Downloading radars for %s\n", date.Format("02/01/2006 15"))
	fileContent, err := fetcher.sess.RadarData(date, varName)
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error downloading radars: %w", err)
		return
	}

	radarFilePath := filepath.Join(
		"WRFDA/RADARS",
		date.Format("2006010215"),
		varName+".nc",
	)

	err = os.MkdirAll(filepath.Dir(radarFilePath), os.FileMode(0755))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error creating directory `%s`: %w", filepath.Dir(radarFilePath), err)
		return
	}

	fmt.Fprintf(os.Stderr, "Saving radars to %s\n", radarFilePath)
	err = ioutil.WriteFile(radarFilePath, fileContent, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error saving radars to `%s`: %w", radarFilePath, err)
	}

}
