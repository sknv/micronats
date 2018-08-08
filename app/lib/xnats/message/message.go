package message

type statusCode string

const (
	statusCodeKey = "meta.statusCode"

	statusOK    statusCode = "ok"
	statusError statusCode = "error"
)

func (m *Message) HasError() bool {
	if m.GetMeta() == nil {
		return false
	}
	status, _ := m.Meta[statusCodeKey]
	return statusError == statusCode(status)
}
