package fetcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
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
func WrfdaSensors(simulStartDate time.Time, domain webdrops.Domain) error {

	sensorClasses := []string{
		//"DIREZIONEVENTO",
		//"IGROMETRO",
		"TERMOMETRO",
		//"ANEMOMETRO",
		//"PLUVIOMETRO",
		//"BAROMETRO",
	}

	allDatesFetched := sync.WaitGroup{}
	errs := make(chan error, 3)

	fetchDate := func(date time.Time) {
		allDatesFetched.Add(1)
		go func() {
			defer allDatesFetched.Done()
			var sess webdrops.Session
			err := sess.Login()
			if err != nil {
				errs <- err
				return
			}

			fetcher := WrfdaSensorsSession{
				Sess:   sess,
				Domain: domain,
			}
			for _, class := range sensorClasses {
				fetcher.fetchSensor(class /*, ids*/, date, false)
			}
			if fetcher.sessError != nil {
				errs <- fetcher.sessError
			}

		}()
	}

	fetchDate(simulStartDate)
	fetchDate(simulStartDate.Add(-3 * time.Hour))
	fetchDate(simulStartDate.Add(-6 * time.Hour))

	allDatesFetched.Wait()
	var err error
	select {
	case err = <-errs:
	default:
	}

	close(errs)
	return err
}

// WrfdaSensorsSession ...
type WrfdaSensorsSession struct {
	sessError error
	Sess      webdrops.Session
	Domain    webdrops.Domain
}

// FetchSensorIDs ...
func (fetcher *WrfdaSensorsSession) FetchSensorIDs(class string, date time.Time, domain webdrops.Domain) []string {
	if fetcher.sessError != nil {
		return nil
	}

	fmt.Fprintf(os.Stderr, "Downloading sensors registry for %s\n", class)
	sensorAnag, err := fetcher.Sess.SensorsList(class, webdrops.GroupWunderground)
	if err != nil {
		fetcher.sessError = fmt.Errorf("error fetching sensors list: %w", err)
		return nil
	}

	jsonFilePath := filepath.Join(
		"WRFDA/SENSORS",
		//date.Format("2006010215"),
		fmt.Sprintf("%s-registry.json", class),
	)

	err = os.MkdirAll(filepath.Dir(jsonFilePath), os.FileMode(0755))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error creating directory `%s`: %w", filepath.Dir(jsonFilePath), err)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Saving observations to %s\n", jsonFilePath)
	err = ioutil.WriteFile(jsonFilePath, sensorAnag, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error saving sensors data to `%s`: %w", jsonFilePath, err)
	}

	ids, err := fetcher.Sess.IDFromSensorsList(sensorAnag, domain)
	fmt.Fprintf(os.Stderr, "Found %d sensors\n", len(ids))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error creating directory `%s`: %w", filepath.Dir(jsonFilePath), err)
		return nil
	}
	return ids
}

func (fetcher *WrfdaSensorsSession) fetchSensor(class string /*, ids []string*/, date time.Time, log bool) {
	if fetcher.sessError != nil {
		return
	}

	from := date.Add(-5 * time.Minute)
	to := date.Add(5 * time.Minute)

	fmt.Fprintf(os.Stderr, "Downloading observations for %s on %s\n", class, date.Format("02/01/2006 15"))
	observations, err := fetcher.Sess.SensorsData(class /*, ids*/, from, to, 60, webdrops.GroupWunderground)
	if err != nil {
		fetcher.sessError = fmt.Errorf("error fetching sensors data: %w", err)
		return
	}

	jsonFilePath := filepath.Join(
		"WRFDA/SENSORS",
		date.Format("2006010215"),
		fmt.Sprintf("%s.json", class),
	)

	err = os.MkdirAll(filepath.Dir(jsonFilePath), os.FileMode(0755))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error creating directory `%s`: %w", filepath.Dir(jsonFilePath), err)
		return
	}

	fmt.Fprintf(os.Stderr, "Saving observations to %s\n", jsonFilePath)
	err = ioutil.WriteFile(jsonFilePath, observations, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error saving sensors data to `%s`: %w", jsonFilePath, err)
	}

}
