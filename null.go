package data

type NullID string

func (id NullID) String() string {
	return string(id)
}

func (id NullID) Hex() string {
	return string(id)
}

func (id NullID) Valid() bool {
	return true
}

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

func (db *NullDB) RegisterForUpdates(a Identifiable) *chan *Change {
	c := make(chan *Change)
	return &c
}

func (db *NullDB) NewObjectID() ID {
	return NullID("")
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

func (db *NullDB) NewQuery(k Kind) Query {
	return nil
}

type NullSchema struct {
	Schema
}

func NewNullSchema() *NullSchema {
	s, _ := NewSchema(new(RelationshipMap), 0)
	return &NullSchema{
		Schema: s,
	}
}

type NullStore struct {
	*NullDB
	*NullSchema
}

func NewNullStore() Store {
	return &NullStore{
		NullDB:     NewNullDB(),
		NullSchema: NewNullSchema(),
	}
}
