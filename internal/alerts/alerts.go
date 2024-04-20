package alerts

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
