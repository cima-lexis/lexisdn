package webdrops

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) RadarTimeline(date time.Time, log bool) ([]time.Time, error) {
	from := date.Add(-30 * time.Minute)
	to := date.Add(30 * time.Minute)

	url := fmt.Sprintf(
		"%scoverages/RADAR_DPC_HDF5_CAPPI/?from=%s&to=%s",
		config.Config.URL,
		from.Format("200601021504"),
		to.Format("200601021504"),
	)

	body, err := sess.DoGet(url)
	if err != nil {
		return nil, fmt.Errorf("Error performing Post: %w", err)
	}
	if log {
		//fmt.Println(string(body))
	}
	var timeline []string
	err = json.Unmarshal(body, &timeline)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}

	sort.Strings(timeline)

	results := []time.Time{}
	lastTime := ""
	for _, singleTime := range timeline {
		if singleTime == lastTime {
			continue
		}
		lastTime = singleTime
		dt, err := time.Parse("200601021504", singleTime)
		if err != nil {
			return nil, fmt.Errorf("Error parsing timeline: %w", err)
		}
		results = append(results, dt)
	}

	return results, nil
}
