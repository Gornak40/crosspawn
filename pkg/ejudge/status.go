package ejudge

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type RunStatus int

// ! Do not touch! It's from ejudge source.
const (
	RunStatusOK RunStatus = 0
	RunStatusRJ RunStatus = 17
)

type EjStatusChange struct {
	ContestID uint
	RunID     uint
	Status    RunStatus
}

func (ej *EjClient) ChangeRunStatus(status *EjStatusChange) error {
	params := url.Values{
		"contest_id": {strconv.Itoa(int(status.ContestID))},
		"run_id":     {strconv.Itoa(int(status.RunID))},
		"status":     {strconv.Itoa(int(status.Status))}, // TODO: status
	}

	data, err := ej.shootEjAPIRaw(context.TODO(), http.MethodPost, "ej/api/v1/master/change-status", params)
	if err != nil {
		return err
	}
	if string(data) == "" { // TODO: wait for bugfix in ejudge
		return nil
	}
	_, err = parseEjAPIAnswer(data)

	return err
}
