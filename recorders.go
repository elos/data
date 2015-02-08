package data

const RecorderDBType DBType = "recorder"

type RecorderID struct {
	value string
	valid bool
}

func NewRecorderID(s string) *RecorderID {
	return &RecorderID{value: s}
}

func (id *RecorderID) Valid() bool {
	return id.valid
}

func (id *RecorderID) SetValidity(v bool) *RecorderID {
	id.valid = true
	return id
}

type RecorderDB struct {
	Connected        string
	ModelUpdates     chan *Change
	Saved            []Record
	Deleted          []Record
	PopulatedById    []Record
	PopulatedByField []Record
	Err              error
}

func NewRecorderDB() *RecorderDB {
	return (&RecorderDB{}).Reset()
}

func (db *RecorderDB) Reset() *RecorderDB {
	db.Connected = ""
	db.ModelUpdates = make(chan *Change)
	db.Saved = make([]Record, 0)
	db.Deleted = make([]Record, 0)
	db.PopulatedById = make([]Record, 0)
	db.PopulatedByField = make([]Record, 0)
	db.Err = nil

	return db
}

func (db *RecorderDB) Connect(addr string) error {
	if db.Err != nil {
		return db.Err
	}

	db.Connected = addr

	return nil
}

func (db *RecorderDB) NewID() ID {
	return NewRecorderID("RecorderID").SetValidity(true)
}

func (db *RecorderDB) CheckID(id ID) error {
	if db.Err != nil {
		return db.Err
	}

	if !id.Valid() {
		return ErrInvalidID
	}

	return nil
}

func (db *RecorderDB) Type() DBType {
	return RecorderDBType
}

func (db *RecorderDB) Save(r Record) error {
	if db.Err != nil {
		return db.Err
	}

	recordedSave(r)
	db.Saved = append(db.Saved, r)
	return nil
}

func (db *RecorderDB) Delete(r Record) error {
	if db.Err != nil {
		return db.Err
	}

	recordedDelete(r)
	db.Deleted = append(db.Deleted, r)
	return nil
}

func (db *RecorderDB) PopulateByID(r Record) error {
	if db.Err != nil {
		return db.Err
	}

	recordedPopulateID(r)
	db.PopulatedById = append(db.PopulatedById, r)
	return nil
}

func (db *RecorderDB) PopulateByField(field string, value interface{}, r Record) error {
	if db.Err != nil {
		return db.Err
	}

	recordedPopulateByField(field, value, r)
	db.PopulatedByField = append(db.PopulatedByField, r)
	return nil
}

func (db *RecorderDB) NewQuery(k Kind) Query {
	return nil
}

func (db *RecorderDB) RegisterForUpdates(a Identifiable) *chan *Change {
	return &db.ModelUpdates
}

var recordedSave = func(r Record) {}
var recordedDelete = func(r Record) {}
var recordedPopulateID = func(r Record) {}
var recordedPopulateByField = func(s string, v interface{}, r Record) {}

type RecorderSchema struct {
	Schema
	Err error
}

func NewRecorderSchema(rm *RelationshipMap, version int) (*RecorderSchema, error) {
	sch, err := NewSchema(rm, version)

	if err != nil {
		return nil, err
	}

	return &RecorderSchema{
		Schema: sch,
	}, nil
}

func (s *RecorderSchema) SetError(e error) *RecorderSchema {
	s.Err = e
	return s
}

func (s *RecorderSchema) Link(this Model, that Model, n LinkName) error {
	if s.Err != nil {
		return s.Err
	}

	recordedLink(this, that, n)
	return s.Schema.Link(this, that, n)
}

func (s *RecorderSchema) Unlink(this Model, that Model, n LinkName) error {
	if s.Err != nil {
		return s.Err
	}

	recordedUnlink(this, that, n)
	return s.Schema.Unlink(this, that, n)
}

var recordedLink = func(this Model, that Model, n LinkName) {}
var recordedUnlink = func(this Model, that Model, n LinkName) {}

type RecorderStore struct {
	Store
}

func NewRecorderStore(db DB, s Schema) *RecorderStore {
	return &RecorderStore{
		Store: NewStore(db, s),
	}
}

func (s *RecorderStore) Register(k Kind, c ModelConstructor) {
	recordedRegister(k, c)
	s.Store.Register(k, c)
}

func (s *RecorderStore) ModelFor(k Kind) (Model, error) {
	recordedModelFor(k)
	return s.Store.ModelFor(k)
}

func (s *RecorderStore) Unmarshal(k Kind, attrs AttrMap) (Model, error) {
	recordedUnmarshal(k, attrs)
	return s.Store.Unmarshal(k, attrs)
}

var recordedRegister = func(k Kind, c ModelConstructor) {}
var recordedModelFor = func(k Kind) {}
var recordedUnmarshal = func(k Kind, attrs AttrMap) {}
