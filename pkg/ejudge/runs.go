package ejudge

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type ejMaskBit uint

// ! Do not touch! It's from ejudge source.
const (
	ejRunID ejMaskBit = iota
	ejSize
	ejTime
	ejAbsoluteTime
	ejRelativeTime
	ejNsec
	ejUserID
	ejUserLogin
	ejUserName
	ejProbID
	ejProbName
	ejLangID
	ejLangName
	ejIP
	ejSHA1
	ejScore
	ejTest
	ejScoreAdj
	ejResult
	ejVariant
	ejMimeType
	ejSavedScore
	ejSavedTest
	ejSavedResult
	ejUUID
	ejEOLNType
	ejStorageFlags
	ejTokens
	ejVerdictBits
	ejLastChangeTime
	ejExternalUser
	ejNotificationInfo
)

func getFieldMask(need ...ejMaskBit) string {
	var mask uint64
	for _, f := range need {
		mask |= 1 << f
	}

	return strconv.FormatUint(mask, 10)
}

type EjRun struct {
	RunID     uint   `json:"run_id"`
	ContestID uint   `json:"-"` // not working, but it's not a problem
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	StatusStr string `json:"status_str"`
	ProbName  string `json:"prob_name"`
}

type EjRuns struct {
	Runs      []EjRun `json:"runs"`
	TotalRuns uint    `json:"total_runs"`
}

func (ej *EjClient) GetContestRuns(id uint, filter string, count int) (*EjRuns, error) {
	fieldMask := getFieldMask(ejRunID, ejUserLogin, ejUserName, ejResult, ejProbName)

	params := url.Values{
		"contest_id":  {strconv.Itoa(int(id))},
		"filter_expr": {filter},
		"last_run":    {strconv.Itoa(-count)}, // some ejudge stuff
		"field_mask":  {fieldMask},
	}
	answer, err := ej.shootEjAPIGet(context.TODO(), "ej/api/v1/master/list-runs-json", params)
	if err != nil {
		return nil, err
	}

	var contest EjRuns
	if err := json.Unmarshal(answer.Result, &contest); err != nil {
		return nil, err
	}

	for i := range contest.Runs {
		contest.Runs[i].ContestID = id
	}

	return &contest, nil
}
