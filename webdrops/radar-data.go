package webdrops

import (
	"fmt"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

func (sess *Session) RadarData(date time.Time, varName string) ([]byte, error) {

	url := fmt.Sprintf(
		"%scoverages/RADAR_DPC_HDF5_%s/%s/%s/-/all",
		config.Config.URL,
		varName,
		date.Format("200601021504"),
		varName,
	)

	bodyResp, err := sess.DoGet(url, "application/octet-stream")
	if err != nil {
		return nil, fmt.Errorf("error performing Post: %w", err)
	}

	return bodyResp, nil
}
