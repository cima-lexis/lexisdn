package webdrops

import (
	"encoding/json"
	"fmt"

	"github.com/cima-lexis/lexisdn/config"
)

type interpolatorReq struct {
	Lats        []float64 `json:"lats"`
	Lons        []float64 `json:"lons"`
	Values      []float64 `json:"values"`
	GridLatN    float64   `json:"gridLatN"`    //max latitude of the interpolation grid
	GridLatS    float64   `json:"gridLatS"`    //min latitude of the interpolation grid
	GridLatStep float64   `json:"gridLatStep"` //latitude step of the interpolation grid
	GridLonE    float64   `json:"gridLonE"`    //max longitude of the interpolation grid
	GridLonStep float64   `json:"gridLonStep"` //longitude step of the interpolation grid
	GridLonW    float64   `json:"gridLonW"`    //min longitude of the interpolation grid
	Radius      float64   `json:"radius"`      //interpolator radius (negative for GRISO)

}

type ValuePoint struct {
	Lat   float64
	Lon   float64
	Value float64
}

func (sess *Session) Interpolator(values []ValuePoint) ([]byte, error) {
	sess.Refresh()

	url := fmt.Sprintf(
		"%ssensors/interpolator",
		config.Config.URL,
	)

	req := interpolatorReq{
		Lats:        make([]float64, len(values)),
		Lons:        make([]float64, len(values)),
		Values:      make([]float64, len(values)),
		GridLatN:    47.5,
		GridLatS:    36,
		GridLatStep: 0.02,
		GridLonE:    18.6,
		GridLonStep: 0.02,
		GridLonW:    6,
		Radius:      -0.5,
	}

	for i, vp := range values {
		req.Lats[i] = vp.Lat
		req.Lons[i] = vp.Lon
		req.Values[i] = vp.Value
	}

	bodyJ, err := json.MarshalIndent(req, " ", " ")
	if err != nil {
		return nil, fmt.Errorf("Error converting body to JSON: %w", err)
	}
	fmt.Println(string(bodyJ))

	bodyResp, err := sess.DoPost(url, req)
	bodyRespS := string(bodyResp)
	fmt.Println(bodyRespS)
	if err != nil {
		return nil, fmt.Errorf("Error performing Interpolator: %w", err)

	}

	return bodyResp, nil
}
