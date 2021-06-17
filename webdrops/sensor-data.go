package webdrops

import (
	"fmt"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

// SensorsData ...
func (sess *Session) SensorsData(class string, from, to time.Time, aggregation int, collection SensorGroup) ([]byte, error) {
	fromS := from.Format("200601021504")
	toS := to.Format("200601021504")

	url := fmt.Sprintf(
		"%ssensors/data/%s/%s?from=%s&to=%s&aggr=%d",
		config.Config.URL,
		class,
		collection.String(),
		fromS,
		toS,
		aggregation,
	)
	/*body := map[string][]string{
		"sensors": ids,
	}*/

	bodyResp, err := sess.DoGet(url /*, body*/)
	if err != nil {
		return nil, fmt.Errorf("error performing Post: %w", err)
	}

	return bodyResp, nil
}
