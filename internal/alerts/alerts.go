package alerts

import "github.com/gin-contrib/sessions"

const flashesGroup = "alerts"

type AlertType int

const (
	AlertSuccess AlertType = iota
	AlertWarning
	AlertDanger
	AlertInfo
)

type Alert struct {
	Message string
	Type    AlertType
}

func Add(session sessions.Session, a Alert) {
	session.AddFlash(a, flashesGroup)
	_ = session.Save()
}

func Get(session sessions.Session) []Alert {
	flashes := session.Flashes(flashesGroup)
	res := make([]Alert, 0, len(flashes))
	for _, f := range flashes {
		if flash, ok := f.(Alert); ok {
			res = append(res, flash)
		}
	}
	_ = session.Save()

	return res
}
