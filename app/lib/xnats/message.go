package xnats

type Message struct {
	Body   []byte            `json:"body,omitempty"`
	Meta   map[string]string `json:"meta,omitempty"`
	Status *Status           `json:"status,omitempty"`
}
