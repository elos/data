package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

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
