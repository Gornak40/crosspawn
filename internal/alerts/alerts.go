package alerts

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/sirupsen/logrus"
)

const alertsFlashesGroup = "alerts"

type AlertType string

const (
	TypeSuccess AlertType = "success"
	TypeWarning AlertType = "warning"
	TypeDanger  AlertType = "danger"
	TypeInfo    AlertType = "info"
)

type Alert struct {
	Message string    `json:"message"`
	Type    AlertType `json:"type"`
}

func Add(session sessions.Session, a Alert) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	session.AddFlash(string(data), alertsFlashesGroup)

	return session.Save()
}

func Get(session sessions.Session) []Alert {
	flashes := session.Flashes(alertsFlashesGroup)
	if len(flashes) > 0 {
		_ = session.Save()
	}

	result := make([]Alert, 0, len(flashes))
	for _, f := range flashes {
		s, ok := f.(string)
		if !ok {
			logrus.Errorf("bad flash: %v", f)

			continue
		}

		var a Alert
		if err := json.Unmarshal([]byte(s), &a); err != nil {
			logrus.WithError(err).Errorf("failed to unmarshal flash: %s", s)

			continue
		}
		result = append(result, a)
	}

	return result
}
