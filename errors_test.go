package data_test

import (
	. "github.com/elos/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	// Not redundant, this keeps the data package honest about which
	// error types it provides
	Context("Exported, necessarily supported ErrTypes", func() {
		It("Defines ErrNotFound", func() {
			Expect(ErrNotFound).NotTo(BeNil())
		})

		It("Defines ErrInvalidID", func() {
			Expect(ErrInvalidID).NotTo(BeNil())
		})

		It("Defines ErrInvalidDBType", func() {
			Expect(ErrInvalidDBType).NotTo(BeNil())
		})

		It("Defines ErrInvalidSchema", func() {
			Expect(ErrInvalidSchema).NotTo(BeNil())
		})

		It("Defines ErrUndefinedKind", func() {
			Expect(ErrUndefinedKind).NotTo(BeNil())
		})

		It("Defines ErrUndefinedLink", func() {
			Expect(ErrUndefinedLink).NotTo(BeNil())
		})

		It("Defines ErrUndefinedLinkKind", func() {
			Expect(ErrUndefinedLinkKind).NotTo(BeNil())
		})

		It("Defines ErrIncompatibleModels", func() {
			Expect(ErrIncompatibleModels).NotTo(BeNil())
		})

		It("Defines AttError", func() {
			Expect(AttrError{}).NotTo(BeNil())
		})
	})

	Describe("NewAttrError", func() {
		It("Prints correctly", func() {
			e := NewAttrError("first", "second")
			Expect(e.Error()).To(Equal("attribute first must second"))
		})
	})

})
