package osql

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/elos/data"
	"github.com/jmoiron/sqlx"
)

const Type data.DBType = "sql"

type (
	Opts struct {
		Database   *sql.DB
		DriverName string
	}

	DB struct {
		*sqlx.DB
		tables map[data.Kind]string
	}
)

// DB Implementation {{{

func New(opts *Opts) (*DB, error) {
	if opts.DriverName == "" {
		return nil, errors.New("data/builtin/sql: must provide driver name")
	}

	if opts.Database == nil {
		return nil, errors.New("data/builtin/sql: must provide database")
	}

	db := sqlx.NewDb(opts.Database, opts.DriverName)

	return &DB{
		DB:     db,
		tables: make(map[data.Kind]string),
	}, nil
}

func (db *DB) RegisterKind(k data.Kind, tableName string) {
	db.tables[k] = tableName
}

func (db *DB) Type() data.DBType {
	return Type
}

func (db *DB) NewID() data.ID {
	return data.ID("temp")
}

func ID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func (db *DB) ParseID(s string) (data.ID, error) {
	_, err := ID(s)
	if err != nil {
		return data.ID(s), data.ErrInvalidID
	} else {
		return data.ID(s), nil
	}
}

// }}}
