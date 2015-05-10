package transfer

import "github.com/elos/data"

type (
	Action string

	Conn interface {
		WriteJSON(interface{}) error
	}

	// Inbound
	Envelope struct {
		Conn   `json:"-"`
		Action `json:"action"`
		Data   map[data.Kind]data.AttrMap `json:"data"`
	}

	// Outbound
	Package struct {
		Action Action       `json:"action"`
		Data   data.KindMap `json:"data"`
	}

	Dispatcher interface {
		Dispatch(e *Envelope) error
	}
)

func NewEnvelope(c Conn, a Action, data map[data.Kind]data.AttrMap) *Envelope {
	return &Envelope{
		Conn:   c,
		Action: a,
		Data:   data,
	}
}

func NewPackage(a Action, data map[data.Kind]interface{}) *Package {
	return &Package{
		Action: a,
		Data:   data,
	}
}
