package data

type Store struct {
	DB
	Schema
}

func NewStore(db DB, s Schema) *Store {
	return &Store{
		DB:     db,
		Schema: s,
	}
}
