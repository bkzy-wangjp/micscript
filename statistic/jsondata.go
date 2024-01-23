package statistic

import (
	"encoding/json"
	"os"
)

func GetDataFromJSON(fpath string) ([]TimeSeriesDataS, error) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	var tsds []TimeSeriesDataS
	if err := json.Unmarshal(data, &tsds); err != nil {
		return nil, err
	}
	return tsds, nil
}

func GetItDataFromJSON(fpath string) ([]TimeSeriesDataI, error) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	var tsds []TimeSeriesDataI
	if err := json.Unmarshal(data, &tsds); err != nil {
		return nil, err
	}
	return tsds, nil
}
