package data

// nullID {{{

type nullID string

func NewNullID(s string) nullID {
	return nullID(s)
}

func (id nullID) Valid() bool {
	return true
}

// nullID }}}

// nullDB {{{

type nullDB struct{}

func NewNullDB() *nullDB {
	return &nullDB{}
}

var nullDBType DBType = "dev/null"

func (db *nullDB) Type() DBType {
	return nullDBType
}

func (db *nullDB) Connect(addr string) error {
	return nil
}

func (db *nullDB) RegisterForUpdates(a Identifiable) *chan *Change {
	c := make(chan *Change)
	return &c
}

func (db *nullDB) NewID() ID {
	return nullID("")
}

func (db *nullDB) CheckID(id ID) error {
	return nil
}

func (db *nullDB) Save(m Record) error {
	return nil
}

func (db *nullDB) Delete(m Record) error {
	return nil
}

func (db *nullDB) PopulateByID(m Record) error {
	return nil
}

func (db *nullDB) PopulateByField(field string, value interface{}, m Record) error {
	return nil
}

func (db *nullDB) NewQuery(k Kind) Query {
	return nil
}

// nullDB }}}

// nullSchema {{{

type nullSchema struct {
	Schema
}

func NewNullSchema() *nullSchema {
	s, _ := NewSchema(new(RelationshipMap), 0)
	return &nullSchema{
		Schema: s,
	}
}

// nullSchema }}}

// nullStore {{{

type nullStore struct {
	*nullDB
	*nullSchema
	dbType DBType
}

func (s *nullStore) Register(k Kind, c ModelConstructor) {
	return
}

func (s *nullStore) ModelFor(kind Kind) (Model, error) {
	return nil, nil
}

func (s *nullStore) Unmarshal(k Kind, attrs AttrMap) (Model, error) {
	return nil, nil
}

func (s *nullStore) Type() DBType {
	if s.dbType == "" {
		return s.nullDB.Type()
	}

	return s.dbType
}

func NewNullStore() Store {
	return &nullStore{
		nullDB:     NewNullDB(),
		nullSchema: NewNullSchema(),
	}
}

func NewNullStoreWithType(t DBType) Store {
	return &nullStore{
		nullDB:     NewNullDB(),
		nullSchema: NewNullSchema(),
		dbType:     t,
	}
}

// nullStore }}}
