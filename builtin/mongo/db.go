package mongo

import (
	"fmt"
	"sync"

	"github.com/elos/data"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

const defaultName = "test"
const defaultAddr = "localhost"

type (
	Opts struct {
		Addr string
		Name string
	}

	Conn struct {
		s *mgo.Session
	}

	DB struct {
		name        string
		conn        *Conn
		collections map[data.Kind]string
		m           sync.Mutex
		hub         *data.ChangeHub
	}
)

// Conn Implementation {{{

func Connect(addr string) (*Conn, error) {
	if sesh, err := mgo.Dial(addr); err != nil {
		return nil, err
	} else {
		return &Conn{sesh}, nil
	}
}

func (c *Conn) Close() {
	c.s.Close()
}

// }}}

// DB Implementation {{{

func New(o *Opts) (*DB, error) {
	name := o.Name
	if name == "" {
		name = defaultName
	}

	addr := o.Addr
	if addr == "" {
		addr = defaultAddr
	}

	c, err := Connect(addr)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn:        c,
		name:        name,
		collections: make(map[data.Kind]string),
		hub:         data.NewChangeHub(context.TODO()),
	}, nil
}

func (db *DB) Name() string {
	return db.name
}

func (db *DB) SetName(n string) {
	db.m.Lock()
	defer db.m.Unlock()
	db.name = n
}

func (db *DB) RegisterKind(k data.Kind, collectionName string) {
	db.collections[k] = collectionName
}

func (db *DB) Fork() (*mgo.Session, error) {
	if db.conn == nil {
		return nil, data.ErrNoConnection
	}

	return db.conn.s.Copy(), nil
}

func (db *DB) Collection(s *mgo.Session, k data.Kind) (*mgo.Collection, error) {
	c, ok := db.collections[k]
	if !ok {
		panic(fmt.Sprintf("data/builtin/mongo: undefined collection for kind: %s", k))
	}
	return s.DB(db.name).C(c), nil
}

// data.DB implementation
func (db *DB) NewID() data.ID {
	return data.ID(NewObjectID().Hex())
}

// data.DB implementation
func (db *DB) ParseID(s string) (data.ID, error) {
	if bid, err := ParseObjectID(s); err != nil {
		return "", err
	} else {
		return data.ID(bid.Hex()), nil
	}
}

// data.DB implementation
func (db *DB) Changes() *chan *data.Change {
	return db.hub.Changes()
}

// }}}
