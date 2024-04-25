package ejudge

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type EjContest struct {
	Contest struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"contest"`
	Problems []struct {
		ID        uint   `json:"id"`
		ShortName string `json:"short_name"`
		LongName  string `json:"long_name"`
	} `json:"problems"`
}

func (ej *EjClient) GetContestStatus(id uint) (*EjContest, error) {
	params := url.Values{
		"contest_id": {strconv.Itoa(int(id))},
	}
	answer, err := ej.shootAPI(context.TODO(), "ej/api/v1/client/contest-status-json", params)
	if err != nil {
		return nil, err
	}

	var contest EjContest
	if err := json.Unmarshal(answer.Result, &contest); err != nil {
		return nil, err
	}

	return &contest, nil
}
