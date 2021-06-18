package webdrops

import (
	"fmt"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

// SensorsMap ...
func (sess *Session) SensorsMap(class string, from, to time.Time, group SensorGroup) ([]byte, error) {
	fromS := from.Format("200601021504")
	toS := to.Format("200601021504")

	url := fmt.Sprintf(
		"%ssensors/map/%s/?from=%s&to=%s&stationgroup=%s",
		config.Config.URL,
		class,
		fromS,
		toS,
		group.String(),
	)
	//fmt.Println(url)
	bodyResp, err := sess.DoGet(url, "application/octet-stream")
	if err != nil {
		return nil, fmt.Errorf("error performing GET: %w", err)

	}

	return bodyResp, nil
}
