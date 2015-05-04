package mongo_test

import (
	"testing"

	"github.com/elos/data/builtin/mongo"
)

func TestNewDB(t *testing.T) {
	_, err := mongo.New(&mongo.Opts{})
	if err != nil {
		t.Fatal(err)
	}
}
