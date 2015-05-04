package mongo_test

import (
	"os"
	"testing"

	"github.com/elos/data/builtin/mongo"
)

func TestMain(m *testing.M) {
	mongo.Testify(mongo.Runner)
	i := m.Run()
	mongo.Detestify(mongo.Runner)
	os.Exit(i)
}
