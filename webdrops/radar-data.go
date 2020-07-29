package webdrops

import (
	"fmt"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) RadarData(date time.Time, varName string) ([]byte, error) {
	sess.Refresh()

	url := fmt.Sprintf(
		"%scoverages/RADAR_DPC_HDF5_CAPPI/%s/%s/-/all",
		config.Config.URL,
		date.Format("200601021504"),
		varName,
	)

	bodyResp, err := sess.DoGet(url)
	if err != nil {
		return nil, fmt.Errorf("Error performing Post: %w", err)
	}

	return bodyResp, nil
}
