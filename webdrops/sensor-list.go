package webdrops

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cima-lexis/lexisdn/config"
)

// IDFromSensorsList ...
func (sess *Session) IDFromSensorsList(sensorList []byte, filter Domain) ([]string, error) {
	var sensors = []struct {
		ID  string
		Lng float64
		Lat float64
	}{}
	err := json.Unmarshal(sensorList, &sensors)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
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
	/*
		if len(result) > 1000 {
			result = result[:1000]
		}
	*/
	return result, nil
}

// SensorGroup ...
type SensorGroup int

const (
	// GroupWunderground ...
	GroupWunderground SensorGroup = iota
	// GroupDPC ...
	GroupDPC
)

func (g SensorGroup) String() string {
	if g == GroupDPC {
		return url.QueryEscape("Dewetra%Default")
	}

	if g == GroupWunderground {
		return url.QueryEscape("DewetraWorld%WunderEurope")
	}

	panic(fmt.Sprintf("Unknown group %d", g))
}

// SensorsList ...
func (sess *Session) SensorsList(class string, group SensorGroup) ([]byte, error) {
	url := fmt.Sprintf("%ssensors/list/%s?stationgroup=%s", config.Config.URL, class, group.String())
	return sess.DoGet(url, "application/json")
}
