package sql

import (
	"database/sql"
	"fmt"

	"github.com/elos/data"
)

func (db *DB) Save(r data.Record) error {
	return nil
}

func (db *DB) Delete(r data.Record) error {
	return nil
}

func (db *DB) PopulateByID(r data.Record) error {
	id, err := ID(r.ID().String())
	if err != nil {
		return data.ErrInvalidID
	}

	kind := r.Kind()
	table, ok := db.tables[kind]
	if !ok {
		panic(fmt.Sprintf("data/builtin/sql: undefined table for kind: %s", kind))
	}

	stmt := fmt.Sprintf("SELECT * FROM %s WHERE id=?", table)
	err = db.database.QueryRowx(stmt, id).StructScan(r)

	switch {
	case err == sql.ErrNoRows:
		return data.ErrNotFound
	case err != nil:
		return err
	default:
		return nil
	}
}

func (db *DB) PopulateByField(r data.Record) error {

	return nil
}
