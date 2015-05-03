package data

type Store struct {
	DB
	*ModHub
}

func NewStore(db DB) *Store {
	return &Store{
		DB:     db,
		ModHub: NewModHub(),
	}
}

func (s *Store) Save(m Model) error {
	err := s.DB.Save(m)
	if err != nil {
		s.ModHub.notify(newUpdate(m))
	}
	return err
}

func (s *Store) Delete(m Model) error {
	err := s.DB.Delete(m)
	if err != nil {
		s.ModHub.notify(newDelete(m))
	}
	return s.DB.Delete(m)
}
