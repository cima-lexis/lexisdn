package webdrops

import (
	"encoding/json"
	"fmt"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) IdFromSensorsList(sensorList []byte, filter Domain) ([]string, error) {
	var sensors = []struct {
		ID  string
		Lng float64
		Lat float64
	}{}
	err := json.Unmarshal(sensorList, &sensors)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}

	result := []string{}
	for _, sensor := range sensors {
		if sensor.Lat < filter.MinLat || sensor.Lat > filter.MaxLat {
			continue
		}
		if sensor.Lng < filter.MinLon || sensor.Lng > filter.MaxLon {
			continue
		}

		result = append(result, sensor.ID)
	}

	return result, nil
}

func (sess *Session) SensorsList(class string) ([]byte, error) {
	sess.Refresh()

	return sess.DoGet(config.Config.URL + "sensors/list/" + class)

}
