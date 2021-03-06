package data

import "golang.org/x/net/context"

type (
	// A changeKind indicates the nature of a Chage
	ChangeKind int

	// A change represents a changeification to the data state
	// backed by a store. Any succesful changeification to the underlying
	// DB should trigger a Change
	Change struct {
		Record     `json:"record"`
		ChangeKind `json:"kind"`
	}

	ChangeHub struct {
		subs     []chan *Change
		register chan chan *Change
		Inbound  chan *Change
	}

	FilterFunc func(c *Change) bool
)

const (
	// Update is the ChangeKind triggered on a save
	// As such, it covers both creation and changeification
	Update ChangeKind = iota + 1

	// Delete is the ChangeKind triggered on Delete
	Delete

	Create
)

// Change Implementation {{{

func NewChange(k ChangeKind, r Record) *Change {
	return &Change{r, k}
}

func NewUpdate(r Record) *Change {
	return NewChange(Update, r)
}

func NewDelete(r Record) *Change {
	return NewChange(Delete, r)
}

func NewCreate(r Record) *Change {
	return NewChange(Create, r)
}

// }}}

func NewChangeHub(ctx context.Context) *ChangeHub {
	hub := &ChangeHub{
		subs:     make([]chan *Change, 0),
		register: make(chan chan *Change),
		Inbound:  make(chan *Change),
	}
	go hub.start(ctx)
	return hub
}

func (h *ChangeHub) start(ctx context.Context) {
Run:
	for {
		select {
		// add registstrations to subs
		case sub := <-h.register:
			h.subs = append(h.subs, sub)
		// fan out changes
		case change := <-h.Inbound:
			for _, sub := range h.subs {
				go func(s chan<- *Change, c *Change) {
					s <- c // send the change to the subscribe
				}(sub, change)
			}
		// end
		case <-ctx.Done():
			break Run
		}
	}
}

func (h *ChangeHub) Changes() *chan *Change {
	// make the channel
	c := make(chan *Change)
	// register it
	h.register <- c
	// return it
	return &c
}

func (h *ChangeHub) Notify(c *Change) {
	go func() {
		h.Inbound <- c
	}()
}

// Filtering {{{

// TODO make name clearer
func Filter(ch *chan *Change, fn FilterFunc) *chan *Change {
	nc := make(chan *Change)

	go func() {
		for change := range *ch {
			if fn(change) {
				nc <- change
			}
		}
	}()

	return &nc
}

// TODO make name clearer
func FilterKind(ch *chan *Change, k Kind) *chan *Change {
	return Filter(ch, func(change *Change) bool {
		return change.Record.Kind() == k
	})
}

// }}}
