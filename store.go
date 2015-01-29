package data

type Store interface {
	DB
	Schema
}

type store struct {
	DB
	Schema
}

func NewStore(db DB, s Schema) *store {
	return &store{
		DB:     db,
		Schema: s,
	}
}
