package ejudge

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

func (ej *EjClient) GetRunSource(contestID uint, runID uint) (string, error) {
	params := url.Values{
		"contest_id": {strconv.Itoa(int(contestID))},
		"run_id":     {strconv.Itoa(int(runID))},
	}
	answer, err := ej.shootEjAPIRaw(context.TODO(), http.MethodGet, "ej/api/v1/master/download-run", params)
	if err != nil {
		return "", err
	}

	return string(answer), err
}
