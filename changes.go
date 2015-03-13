package data

import "sync"

// A ChangeKind represents the nature of a Change.
type ChangeKind int

const (
	// Update is the ChangeKind triggered on a save;
	// As such it covers boths creation and modification.
	Update ChangeKind = 1

	// Delete is the ChangeKind triggered on a delete.
	Delete ChangeKind = 2
)

// A Change represents a modification to the data state
// a DB represents. Any succesful modification to the underlying
// should trigger a Change to be sent of a channel.
// Implementations of a the DB interface should implement all
// defined ChangeKinds
type Change struct {
	ChangeKind
	Record
}

// NewChange is a simple constructor for a Change object
func NewChange(kind ChangeKind, r Record) *Change {
	return &Change{
		ChangeKind: kind,
		Record:     r,
	}
}

func NewUpdate(r Record) *Change {
	return NewChange(Update, r)
}

func NewDelete(r Record) *Change {
	return NewChange(Delete, r)
}

func NewChangeHub() *ChangeHub {
	return &ChangeHub{
		subscribers: make(map[ID][]*chan *Change),
	}
}

type ChangeHub struct {
	sync.Mutex
	subscribers map[ID][]*chan *Change
}

func (ch *ChangeHub) RegisterForChanges(client Client) *chan *Change {
	ch.Lock()
	defer ch.Unlock()

	id := client.ID()
	c := make(chan *Change)
	alreadySubscribed, ok := ch.subscribers[id]

	if !ok {
		alreadySubscribed = make([]*chan *Change, 0)
	}

	ch.subscribers[id] = append(alreadySubscribed, &c)

	return &c
}

func (ch *ChangeHub) Notify(c *Change) {
	ch.Lock()
	defer ch.Unlock()

	for _, concernedID := range c.Concerned() {
		channels, ok := ch.subscribers[concernedID]
		if ok {
			for _, channel := range channels {
				go send(*channel, c)
			}
		}
	}
}

func send(channel chan *Change, change *Change) {
	channel <- change
}
