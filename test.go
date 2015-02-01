package data

// TestID {{{

type TestID struct {
	value string
	valid bool
}

func (id *TestID) Valid() bool {
	return id.valid
}

// TestID }}}

// TestDB {{{

type TestDB struct {
	ModelUpdates     chan *Change
	Saved            []Record
	Deleted          []Record
	PopulatedById    []Record
	PopulatedByField []Record
	Err              error
}

func NewTestDB() *TestDB {
	db := &TestDB{}
	db.Reset()
	return db
}

const TestDBType DBType = "test"

func (db *TestDB) Type() DBType {
	return TestDBType
}

func (db *TestDB) Reset() {
	db.ModelUpdates = make(chan *Change)
	db.Saved = make([]Record, 0)
	db.Deleted = make([]Record, 0)
	db.PopulatedById = make([]Record, 0)
	db.PopulatedByField = make([]Record, 0)
	db.Err = nil
}

func (db *TestDB) Connect(addr string) error {
	if db.shouldError() {
		return db.Err
	}

	return nil
}

func (db *TestDB) RegisterForUpdates(a Identifiable) *chan *Change {
	return &db.ModelUpdates
}

func (db *TestDB) NewID() ID {
	return &TestID{valid: true}
}

func (db *TestDB) CheckID(id ID) error {
	if db.shouldError() {
		return db.Err
	}

	if !id.Valid() {
		return ErrInvalidID
	}

	return nil
}

func (db *TestDB) Save(r Record) error {
	if db.shouldError() {
		return db.Err
	}

	db.Saved = append(db.Saved, r)
	return nil
}

func (db *TestDB) Delete(r Record) error {
	if db.shouldError() {
		return db.Err
	}

	db.Deleted = append(db.Deleted, r)
	return nil
}

func (db *TestDB) PopulateByID(r Record) error {
	if db.shouldError() {
		return db.Err
	}

	db.PopulatedById = append(db.PopulatedById, r)
	return nil
}

func (db *TestDB) PopulateByField(field string, value interface{}, r Record) error {
	if db.shouldError() {
		return db.Err
	}

	db.PopulatedByField = append(db.PopulatedByField, r)
	return nil
}

func (db *TestDB) NewQuery(k Kind) Query {
	return nil
}

func (db *TestDB) shouldError() bool {
	if db.Err != nil {
		return true
	} else {
		return false
	}
}

// TestDB }}}

// NullSchema {{{

type TestSchema struct {
	Schema
}

func NewTestSchema() *TestSchema {
	return &TestSchema{
		Schema: NewNullSchema(),
	}
}

// NullSchema }}}

// TestStore {{{

type TestStore struct {
	*store
}

func NewTestStore(db DB, s Schema) *TestStore {
	return &TestStore{
		store: NewStore(NewTestDB(), NewTestSchema()),
	}
}

func (s *TestStore) RegisteredModels() map[Kind]ModelConstructor {
	return s.store.registered
}

// TestStore }}}
