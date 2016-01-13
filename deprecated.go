package data

// a (largely) unsuccesful attempt at defining Access interfaces
// it has become clear the best way to do this is wrapping a db interface
// in another structure which implements the db interface and therefore
// manages access. Client is severely to restricting, you must know more
// about your access control schema and user structure.
type (
	Client interface {
		Record
	}

	Access interface {
		DB
		Client() Client
	}
)
