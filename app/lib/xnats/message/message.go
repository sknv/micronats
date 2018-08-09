package message

type statusCode string

const (
	statusCodeKey = "meta.statusCode"

	statusError statusCode = "error"
)

func (m *Message) HasError() bool {
	if m.Meta == nil {
		return false
	}
	status, _ := m.Meta[statusCodeKey]
	return statusError == statusCode(status)
}

func (m *Message) WithError() *Message {
	if m.Meta == nil {
		m.Meta = make(map[string]string)
	}
	m.Meta[statusCodeKey] = string(statusError) // upsert the value
	return m
}
