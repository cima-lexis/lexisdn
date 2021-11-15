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

// WrfdaRadars retrieves
func WrfdaRadars(simulStartDate time.Time) error {

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
			fetcher := wrfdaRadarsSession{
				sess: sess,
			}
			bestInstant /*timeline*/, err := fetcher.sess.RadarTimeline(date, false)
			if err != nil {
				errs <- fmt.Errorf("error downloading radars timeline: %w", err)
				return
			}
			fetcher.fetchRadar(bestInstant, "CAPPI2", date)
			fetcher.fetchRadar(bestInstant, "CAPPI3", date)
			fetcher.fetchRadar(bestInstant, "CAPPI4", date)
			fetcher.fetchRadar(bestInstant, "CAPPI5", date)
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

type wrfdaRadarsSession struct {
	sessError error
	sess      webdrops.Session
	//domain    webdrops.Domain
}

func (fetcher *wrfdaRadarsSession) fetchRadar(date time.Time, varName string, dateRequested time.Time) {
	if fetcher.sessError != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Downloading radars for %s\n", date.Format("02/01/2006 15"))
	fileContent, err := fetcher.sess.RadarData(date, varName)
	if err != nil {
		fetcher.sessError = fmt.Errorf("error downloading radars: %w", err)
		return
	}

	dtReq := dateRequested.Format("2006010215")
	radarFilePath := fmt.Sprintf("WRFDA/RADARS/%s/%s-%s.nc", dtReq, dtReq, varName)

	err = os.MkdirAll(filepath.Dir(radarFilePath), os.FileMode(0755))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error creating directory `%s`: %w", filepath.Dir(radarFilePath), err)
		return
	}

	fmt.Fprintf(os.Stderr, "Saving radars to %s\n", radarFilePath)
	err = ioutil.WriteFile(radarFilePath, fileContent, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("error saving radars to `%s`: %w", radarFilePath, err)
	}

}
