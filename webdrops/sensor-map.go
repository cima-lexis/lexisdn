package webdrops

import (
	"fmt"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

//lonmin, latmin, lonmax, latmax
func (sess *Session) SensorsMap(class string, from, to time.Time, domain Domain) ([]byte, error) {
	sess.Refresh()
	fromS := from.Format("200601021504")
	toS := to.Format("200601021504")

	url := fmt.Sprintf(
		"%ssensors/map/%s/?from=%s&to=%s&geowin=%.0f,%.0f,%.0f,%.0f",
		config.Config.URL,
		class,
		fromS,
		toS,
		domain.MinLon, domain.MinLat, domain.MaxLon, domain.MaxLat,
	)
	fmt.Println(url)
	bodyResp, err := sess.DoGet(url)
	if err != nil {
		return nil, fmt.Errorf("Error performing GET: %w", err)

	}

	return bodyResp, nil
}
