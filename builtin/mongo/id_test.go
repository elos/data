package mongo_test

import (
	"math/rand"
	"testing"

	"github.com/elos/data/builtin/mongo"
	"github.com/elos/testing/expect"
)

func TestNewObjectID(t *testing.T) {
	id := mongo.NewObjectID()

	if !id.Valid() {
		t.Errorf("NewObjectID should return a valid id")
	}
}

func TestParseID(t *testing.T) {
	id := mongo.NewObjectID()

	s := id.String()

	id, err := mongo.ParseObjectID(s)
	expect.NoError("parsing object id", err, t)
}

func TestIDSet(t *testing.T) {
	ids := make(mongo.IDSet, 0)

	ids = mongo.AddID(mongo.AddID(ids, mongo.NewObjectID()), mongo.NewObjectID())

	for i := 0; i < rand.Intn(10); i++ {
		ids = mongo.AddID(ids, mongo.NewObjectID())
	}

	id := mongo.NewObjectID()

	ids = mongo.AddID(ids, id)
	_, ok := ids.IndexID(id)
	if !ok {
		t.Errorf("Id be in set")
	}

	ids = mongo.AddID(ids, mongo.NewObjectID())
	ids = mongo.DropID(ids, id)

	_, ok = ids.IndexID(id)
	if ok != false {
		t.Errorf("Id should no longer be in set")
	}

	ids = make(mongo.IDSet, 0)
	i, ok := ids.IndexID(id)
	if i != -1 || ok != false {
		t.Errorf("IndexID return -1, false for a non-member id")
	}

	ids = make(mongo.IDSet, 0)
	ids = mongo.AddID(ids, id)
	ids = mongo.AddID(ids, id)
	ids = mongo.AddID(ids, id)

	if len(ids) != 1 {
		t.Errorf("IDSet should be a set, no duplicates")
	}

	ids = make(mongo.IDSet, 0)
	ids = mongo.DropID(ids, id)
	if len(ids) != 0 {
		t.Errorf("Should be able to drop a non-member id")
	}
}
