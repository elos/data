package mongo

import (
	"errors"
	"strings"
	"sync"

	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

var emptyID = *new(bson.ObjectId)

func EmptyID(id bson.ObjectId) bool {
	return id.String() == emptyID.String()
}

func NewEmptyID() bson.ObjectId {
	return *new(bson.ObjectId)
}

func ParseObjectID(idS string) (id bson.ObjectId, err error) {
	defer func() { // cause the bson pkg isn't idiomatic
		if r := recover(); r != nil {
			id = emptyID
			err = errors.New("Could not parse id")
		}
	}()

	if strings.Contains(idS, "ObjectIdHex(") {
		ss := strings.Split(idS, "\"")
		idS = ss[1]
	}

	id = bson.ObjectIdHex(idS) // can panic

	return id, nil
}

func NewObjectID() bson.ObjectId {
	return bson.NewObjectId()
}

func NewObjectIDFromHex(hex string) data.ID {
	return data.ID(bson.ObjectIdHex(hex).Hex())
}

func IsObjectIDHex(hex string) bool {
	return bson.IsObjectIdHex(hex)
}

func NewIDSetFromStrings(strings []string) IDSet {
	set := make(IDSet, len(strings))
	for _, s := range strings {
		pid, _ := ParseObjectID(s)
		set = AddID(set, pid)
	}
	return set
}

type IDSet []bson.ObjectId

func AddID(s IDSet, id bson.ObjectId) IDSet {
	_, ok := s.IndexID(id)

	if !ok {
		ns := append(s, id)
		return ns
	}

	return s
}

func DropID(s IDSet, id bson.ObjectId) IDSet {
	i, ok := s.IndexID(id)

	if !ok {
		return s
	}

	s = s[:i+copy(s[i:], s[i+1:])]

	return s
}

func (s IDSet) IndexID(id bson.ObjectId) (int, bool) {
	for i, iid := range s {
		if id == iid {
			return i, true
		}
	}

	return -1, false
}

func (s IDSet) NotIn(i IDSet) IDSet {
	ids := make(IDSet, 0)
	for _, id := range s {
		_, in := i.IndexID(id)
		if !in {
			ids = append(ids, id)
		}
	}
	return ids
}

func (s IDSet) IDs() []data.ID {
	ids := make([]data.ID, len(s))
	for i, id := range s {
		ids[i] = data.ID(id.Hex())
	}
	return ids
}

type IDIter struct {
	data.DB
	ids   IDSet
	place int
	err   error

	*sync.Mutex
}

func NewIDIter(set IDSet, s data.DB) *IDIter {
	return &IDIter{
		place: 0,
		DB:    s,
		ids:   set,
		Mutex: new(sync.Mutex),
	}
}

func (i *IDIter) Next(r data.Record) bool {
	if i.place >= len(i.ids) {
		return false
	}

	r.SetID(data.ID(i.ids[i.place].Hex()))

	if err := i.DB.PopulateByID(r); err != nil {
		i.err = err
		return false
	} else {
		i.place += 1
		return true
	}
}

func (i *IDIter) Close() error {
	return i.err
}
