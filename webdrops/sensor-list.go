package webdrops

import (
	"encoding/json"
	"fmt"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) SensorsList(class string, filter Domain, log bool) ([]string, error) {
	sess.Refresh()
	if log {
		fmt.Println(config.Config.URL + "sensors/list/" + class)
	}
	body, err := sess.DoGet(config.Config.URL + "sensors/list/" + class)
	if err != nil {
		return nil, fmt.Errorf("Error performing GET: %w", err)
	}

	var sensors = []struct {
		ID  string
		Lng float64
		Lat float64
	}{}
	err = json.Unmarshal(body, &sensors)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}
	if log {
		fmt.Println(sensors)
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
