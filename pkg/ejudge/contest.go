package ejudge

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type EjContest struct {
	Problems []struct {
		ID        int    `json:"id"`
		ShortName string `json:"short_name"`
		LongName  string `json:"long_name"`
	} `json:"problems"`
}

func (ej *EjClient) GetContestStatus(id int) (*EjContest, error) {
	params := url.Values{
		"contest_id": {strconv.Itoa(id)},
	}
	answer, err := ej.shoot(context.TODO(), "ej/api/v1/client/contest-status-json", params)
	if err != nil {
		return nil, err
	}

	var contest EjContest
	if err := json.Unmarshal(answer.Result, &contest); err != nil {
		return nil, err
	}

	return &contest, nil
}
