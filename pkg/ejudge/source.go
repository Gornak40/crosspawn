package ejudge

import (
	"context"
	"net/url"
	"strconv"
)

func (ej *EjClient) GetRunSource(contestID uint, runID uint) (string, error) {
	params := url.Values{
		"contest_id": {strconv.Itoa(int(contestID))},
		"run_id":     {strconv.Itoa(int(runID))},
	}
	answer, err := ej.shootRaw(context.TODO(), "ej/api/v1/master/download-run", params)
	if err != nil {
		return "", err
	}

	return string(answer), err
}
