package webdrops

import (
	"fmt"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) SensorsData(class string, ids []string, from, to time.Time, aggregation int, log bool) ([]byte, error) {
	sess.Refresh()
	fromS := from.Format("200601021504")
	toS := to.Format("200601021504")

	url := fmt.Sprintf(
		"%ssensors/data/%s/?from=%s&to=%s&aggr=%d",
		config.Config.URL,
		class,
		fromS,
		toS,
		aggregation,
	)
	body := map[string][]string{
		"sensors": ids,
	}

	bodyResp, err := sess.DoPost(url, body)
	if err != nil {
		return nil, fmt.Errorf("Error performing Post: %w", err)
	}

	return bodyResp, nil
}
