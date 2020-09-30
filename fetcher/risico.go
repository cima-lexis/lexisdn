package fetcher

import (
	"encoding/json"
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
		domain: webdrops.Domain{
			MinLat: 36.0,
			MaxLat: 47.5,
			MinLon: 6,
			MaxLon: 18.6,
		},
	}

	for step := 6; step >= 1; step-- {
		from := simulStartDate.Add(-12 * time.Duration(step) * time.Hour)
		to := from.Add(12 * time.Hour)

		fetcher.fetchSensorData("PLUVIOMETRO", from, to)
		fetcher.fetchSensorData("IGROMETRO", from, to)
		fetcher.fetchSensorData("TERMOMETRO", from, to)

		if fetcher.sessError != nil {
			break
		}
	}

	return fetcher.sessError
}

type risicoSession struct {
	sessError error
	sess      webdrops.Session
	domain    webdrops.Domain
}

func (fetcher *risicoSession) fetchSensorData(class string, from, to time.Time) {
	if fetcher.sessError != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Downloading sensors registry for %s\n", class)
	sensorRegistry, err := fetcher.sess.SensorsList(class)
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

	fmt.Fprintf(os.Stderr, "Downloading observations for %s from %s to %s\n", class, from.Format("02/01/2006 15"), to.Format("02/01/2006 15"))
	observations, err := fetcher.sess.SensorsData(class, ids, from, to, 3600, false)
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error fetching sensors data: %w", err)
		return
	}

	jsonFilePath := filepath.Join(
		"RISICO/SENSORS",
		fmt.Sprintf("%s.json", class),
	)
	jsonAnagFilePath := filepath.Join(
		"RISICO/SENSORS",
		fmt.Sprintf("anag-%s.json", class),
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
		return
	}
	err = ioutil.WriteFile(jsonAnagFilePath, sensorRegistry, os.FileMode(0644))
	if err != nil {
		fetcher.sessError = fmt.Errorf("Error saving sensors reg data to `%s`: %w", jsonFilePath, err)
		return
	}

	instants, err := aggregateData(sensorRegistry, observations)
	for _, inst := range instants {
		mapPath := filepath.Join(
			"RISICO/SENSORS",
			fmt.Sprintf("%s.nc", inst.Time.Format("2006010215")),
		)
		mapData, err := fetcher.sess.Interpolator(inst.Values)
		if err != nil {
			fetcher.sessError = fmt.Errorf("Error interpolating data: %w", err)
			return
		}
		err = ioutil.WriteFile(mapPath, mapData, os.FileMode(0644))
		if err != nil {
			fetcher.sessError = fmt.Errorf("Error saving map data to `%s`: %w", mapPath, err)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Found %d instants\n", len(instants))
}

type LatLon struct {
	Lat, Lon float64
}

type Instant struct {
	Time   time.Time
	Values []webdrops.ValuePoint
}

func aggregateData(sensorRegistry, sensorData []byte) ([]*Instant, error) {
	var sensors = []struct {
		ID  string
		Lng float64
		Lat float64
	}{}
	err := json.Unmarshal(sensorRegistry, &sensors)
	if err != nil {
		return nil, fmt.Errorf("Error parsing registry JSON: %w", err)
	}

	var values = []struct {
		SensorID string
		Timeline []string
		Values   []float64
	}{}
	err = json.Unmarshal(sensorData, &values)
	if err != nil {
		return nil, fmt.Errorf("Error parsing values JSON: %w", err)
	}

	sensorMap := map[string]LatLon{}

	for _, sensor := range sensors {
		sensorMap[sensor.ID] = LatLon{sensor.Lat, sensor.Lng}
	}

	resultByTime := map[time.Time]*Instant{}

	for _, val := range values {
		pos := sensorMap[val.SensorID]
		for idx, tmS := range val.Timeline {
			v := val.Values[idx]
			if v == -9998 {
				continue
			}
			tm, err := time.Parse(time.RFC3339, tmS)
			if err != nil {
				return nil, fmt.Errorf("Error parsing date '%s': %w", tmS, err)
			}
			instant, ok := resultByTime[tm]
			if !ok {
				instant = &Instant{Time: tm, Values: []webdrops.ValuePoint{}}
				resultByTime[tm] = instant
			}
			/*if len(instant.Values) > 30 {
				continue
			}*/
			instant.Values = append(instant.Values, webdrops.ValuePoint{
				Lat:   pos.Lat,
				Lon:   pos.Lon,
				Value: v,
			})
		}
	}

	result := make([]*Instant, len(resultByTime))
	i := 0
	for _, inst := range resultByTime {

		result[i] = inst
		i++
	}
	return result, nil
}

func (fetcher *risicoSession) fetchSensorMap(class string, from, to time.Time) {
	if fetcher.sessError != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Downloading observations map for %s from %s to %s\n", class, from.Format("02/01/2006 15"), to.Format("02/01/2006 15"))
	sensorsMap, err := fetcher.sess.SensorsMap(class, from, to)
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
