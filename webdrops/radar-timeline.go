package webdrops

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) timelineForVar(date time.Time, cappivar int) ([]string, error) {
	from := date.Add(-30 * time.Minute)
	to := date.Add(30 * time.Minute)

	fromS := from.Format("200601021504")
	toS := to.Format("200601021504")
	urlFormat := "%scoverages/RADAR_DPC_HDF5_CAPPI%d/?from=%s&to=%s"

	url := fmt.Sprintf(urlFormat, config.Config.URL, cappivar, fromS, toS)

	body, err := sess.DoGet(url, "application/json")
	if err != nil {
		return nil, fmt.Errorf("error performing get: %w", err)
	}

	var timeline []string
	err = json.Unmarshal(body, &timeline)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	sort.Strings(timeline)
	fmt.Printf("Radar availability for date %s, variable CAPPI%d: %v\n", date.Format("2006-01-02-15-04"), cappivar, timeline)

	return timeline, nil
}

// RadarTimeline ...
func (sess *Session) RadarTimeline(date time.Time, log bool) (time.Time, error) {

	var timelines [4][]string
	var err error

	timelines[0], err = sess.timelineForVar(date, 2)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting timeline: %w", err)
	}
	timelines[1], err = sess.timelineForVar(date, 3)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting timeline: %w", err)
	}
	timelines[2], err = sess.timelineForVar(date, 4)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting timeline: %w", err)
	}
	timelines[3], err = sess.timelineForVar(date, 5)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting timeline: %w", err)
	}

	commonInstants := intersect(timelines[0], timelines[1])
	commonInstants = intersect(timelines[2], commonInstants)
	commonInstants = intersect(timelines[3], commonInstants)

	if len(commonInstants) == 0 {
		return time.Time{}, fmt.Errorf("no radar found for %s", date.Format("200601021504"))
	}

	bestFound := time.Time{}
	for _, instantS := range commonInstants {
		instant, err := time.Parse("200601021504", instantS)
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing timeline: %w", err)
		}

		if instant == date {
			return instant, nil
		}

		if math.Abs(instant.Sub(date).Seconds()) <
			math.Abs(bestFound.Sub(date).Seconds()) {
			bestFound = instant

		}

	}

	return bestFound, nil
}

func intersect(a []string, b []string) []string {
	set := make([]string, 0)

	for _, item := range a {
		idx := sort.SearchStrings(b, item)
		if idx < len(b) && b[idx] == item {
			set = append(set, item)
		}
	}

	return set
}
