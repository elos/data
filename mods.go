package data

import "sync"

type (
	// A ModKind indicates the nature of a Chage
	ModKind int

	// A Mod represents a modification to the data state
	// backed by a store. Any succesful modification to the underlying
	// DB should trigger a Change
	Mod struct {
		Model
		ModKind
	}

	ModHub struct {
		m    sync.Mutex
		subs map[ID][]*chan *Mod
	}
)

const (
	// Update is the ChangeKind triggered on a save
	// As such, it covers both creation and modification
	Update ModKind = iota

	// Delete is the ChangeKind triggered on Delete
	Delete
)

func NewMod(k ModKind, m Model) *Mod {
	return &Mod{m, k}
}

func newUpdate(m Model) *Mod {
	return NewMod(Update, m)
}

func newDelete(m Model) *Mod {
	return NewMod(Delete, m)
}

func NewModHub() *ModHub {
	return &ModHub{
		subs: make(map[ID][]*chan *Mod),
	}
}

func (h *ModHub) GetMods(client Client) *chan *Mod {
	h.m.Lock()
	defer h.m.Unlock()

	id := client.ID()
	c := make(chan *Mod)

	if currentSubs, someExist := h.subs[id]; someExist {
		h.subs[id] = append(currentSubs, &c)
	} else {
		h.subs[id] = []*chan *Mod{&c}
	}

	return &c
}

func (h *ModHub) notify(m *Mod) {
	h.m.Lock()
	defer h.m.Unlock()

	for _, id := range m.Concerned() {
		if chans, ok := h.subs[id]; ok {
			for _, c := range chans {
				go func() { *c <- m }()
			}
		}
	}
}
