package data_test

import (
	. "github.com/elos/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("null.go", func() {

	// NullID {{{

	Describe("NullID", func() {
		It("Subtypes a string", func() {
			var id NullID = NullID("null id")
			Expect(id.Valid()).To(BeTrue())
		})

		It("Is always valid", func() {
			Expect(NullID("").Valid()).To(BeTrue())
		})

		It("Satisfies the ID interface", func() {
			// Won't compile if NullID fails implementation
			_ = func() ID { return NullID("") }
		})
	})

	// NullID }}}

	Describe("NullDB", func() {
		It("Satisfies the DB interface", func() {
			// Won't compile if NullDB fails implementation
			_ = func() DB { return NewNullDB() }
		})

	})

	Describe("NullSchema", func() {
		It("Satisfies the Schema interface", func() {
			// Won't compile if NullSchema fails implementation
			_ = func() Schema { return NewNullSchema() }
		})
	})

	Describe("NullStore", func() {
		It("Satisfies the Store interface", func() {
			// Won't compile if NullStore fails implementation
			_ = func() Store { return NewNullStore() }
		})
	})

})