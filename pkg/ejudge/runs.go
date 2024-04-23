package ejudge

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type ejMaskBit uint

// Do not touch it! It's from ejudge source.
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

type EjRuns struct {
	Runs []struct {
		RunID     uint   `json:"run_id"`
		UserLogin string `json:"user_login"`
		UserName  string `json:"user_name"`
		StatusStr string `json:"status_str"`
	} `json:"runs"`
}

func (ej *EjClient) GetContestRuns(id uint, filter string) (*EjRuns, error) {
	fieldMask := getFieldMask(ejRunID, ejUserLogin, ejUserName, ejResult)

	params := url.Values{
		"contest_id":  {strconv.Itoa(int(id))},
		"filter_expr": {filter},
		"field_mask":  {fieldMask},
	}
	answer, err := ej.shoot(context.TODO(), "ej/api/v1/master/list-runs-json", params)
	if err != nil {
		return nil, err
	}

	var contest EjRuns
	if err := json.Unmarshal(answer.Result, &contest); err != nil {
		return nil, err
	}

	return &contest, nil
}
