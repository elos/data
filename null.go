package data

import "time"

// NullID {{{

type NullID string

func NewNullID(s string) NullID {
	return NullID(s)
}

func (id NullID) Valid() bool {
	return true
}

func (id NullID) String() string {
	return string(id)
}

// NullID }}}

// NullDB {{{

type NullDB struct{}

func NewNullDB() *NullDB {
	return &NullDB{}
}

var NullDBType DBType = "dev/null"

func (db *NullDB) Type() DBType {
	return NullDBType
}

func (db *NullDB) Connect(addr string) error {
	return nil
}

func (db *NullDB) RegisterForChanges(client Client) *chan *Change {
	c := make(chan *Change)
	return &c
}

func (db *NullDB) NewID() ID {
	return NullID("")
}

func (db *NullDB) ParseID(id string) (ID, error) {
	return NullID(id), nil
}

func (db *NullDB) CheckID(id ID) error {
	return nil
}

func (db *NullDB) Save(m Record) error {
	return nil
}

func (db *NullDB) Delete(m Record) error {
	return nil
}

func (db *NullDB) PopulateByID(m Record) error {
	return nil
}

func (db *NullDB) PopulateByField(field string, value interface{}, m Record) error {
	return nil
}

func (db *NullDB) NewQuery(k Kind) RecordQuery {
	return nil
}

// NullDB }}}

// NullSchnma {{{

type NullSchema struct {
	Schema
}

func NewNullSchema() *NullSchema {
	s, _ := NewSchema(new(RelationshipMap), 0)
	return &NullSchema{
		Schema: s,
	}
}

// NullSchnma }}}

// NullStore {{{

type NullStore struct {
	*NullDB
	*NullSchema
	dbType DBType
}

func (s *NullStore) Register(k Kind, c ModelConstructor) {
	return
}

func (s *NullStore) ModelFor(kind Kind) (Model, error) {
	return nil, nil
}

func (s *NullStore) Unmarshal(k Kind, attrs AttrMap) (Model, error) {
	return nil, nil
}

func (s *NullStore) Registered() []Kind {
	return make([]Kind, 0)
}

func (s *NullStore) Save(m Model) error {
	return s.NullDB.Save(m)
}

func (s *NullStore) Delete(m Model) error {
	return s.NullDB.Delete(m)
}

func (s *NullStore) PopulateByID(m Model) error {
	return s.NullDB.PopulateByID(m)
}

func (s *NullStore) PopulateByField(str string, v interface{}, m Model) error {
	return s.NullDB.PopulateByField(str, v, m)
}

func (s *NullStore) Type() DBType {
	if s.dbType == "" {
		return s.NullDB.Type()
	}

	return s.dbType
}

func (s *NullStore) Query(k Kind) ModelQuery {
	return nil
}

func (s *NullStore) Compatible(p Persistable) bool {
	return p.DBType() == s.Type()
}

func NewNullStore() *NullStore {
	return &NullStore{
		NullDB:     NewNullDB(),
		NullSchema: NewNullSchema(),
	}
}

func NewNullStoreWithType(t DBType) *NullStore {
	return &NullStore{
		NullDB:     NewNullDB(),
		NullSchema: NewNullSchema(),
		dbType:     t,
	}
}

// NullStore }}}

// NullModel {{{
var NullKind Kind = "null"

type NM struct {
	String string
	Int    int
	kind   Kind   `json:"-" bson"-"`
	dbtype DBType `json:"-" bson"-"`
	id     ID
}

func NewNullModel() *NM {
	return &NM{
		kind:   NullKind,
		dbtype: NullDBType,
		id:     NewNullID("example"),
	}
}

// Model Constructor
func NewNM(s Store) Model {
	return NewNullModel()
}

func (nm *NM) DBType() DBType {
	return nm.dbtype
}

func (nm *NM) Kind() Kind {
	return nm.kind
}

func (nm *NM) SetKind(k Kind) {
	nm.kind = k
}

func (nm *NM) SetDBType(t DBType) {
	nm.dbtype = t
}

func (nm *NM) ID() ID {
	return nm.id
}

func (nm *NM) Version() int {
	return 0
}

func (nm *NM) Valid() bool {
	return true
}

func (nm *NM) Concerned() []ID {
	return make([]ID, 0)
}

func (nm *NM) SetID(id ID) error {
	nm.id = id
	return nil
}

var exampleCanRead = func() bool { return true }
var exampleCanWrite = func() bool { return true }

func (nm *NM) CanRead(c Client) bool {
	return exampleCanRead()
}

func (nm *NM) CanWrite(c Client) bool {
	return exampleCanWrite()
}

var exampleLink = func(m Model, l Link) error { return nil }
var exampleUnlink = func(m Model, l Link) error { return nil }

func (nm *NM) Link(m Model, l Link) error {
	return exampleLink(m, l)
}

func (nm *NM) Unlink(m Model, l Link) error {
	return exampleUnlink(m, l)
}

func (nm *NM) SetCreatedAt(t time.Time) {
}

func (nm *NM) SetUpdatedAt(t time.Time) {
}

func (nm *NM) UpdatedAt() time.Time {
	return time.Now()
}

func (nm *NM) CreatedAt() time.Time {
	return time.Now()
}

func (nm *NM) Schema() Schema {
	return NewNullSchema()
}

var NullLink = &Link{
	Name:    "example",
	Kind:    MulLink,
	Other:   NullKind,
	Inverse: "null's_inverse",
}
