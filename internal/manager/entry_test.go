package manager_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

var _ = Expect

var _ = Describe("Entry", func() {
	existing := []*MockEntryObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
	var client *MockEntryClient[*MockEntryObject]
	var service manager.SDKEntryService[*MockEntryObject, MockLocation]
	var sdk *manager.EntryObjectManager[*MockEntryObject, MockLocation, manager.SDKEntryService[*MockEntryObject, MockLocation]]

	location := MockLocation{}

	ctx := context.Background()

	BeforeEach(func() {
		client = NewMockEntryClient(existing)
		service = NewMockEntryService[*MockEntryObject, MockLocation](client)
		sdk = manager.NewEntryObjectManager(client, service, MockEntrySpecifier, MockEntryMatcher)
	})

	Context("Read()", func() {
		When("reading entry that does not exist", func() {
			It("should return nil object and ErrObjectNotFound error", func() {
				object, err := sdk.Read(ctx, location, "4")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(MatchRegexp("Object not found")))
				Expect(object).To(BeNil())
			})
		})
		When("reading entry that exists", func() {
			It("should return nil error and the existing entry", func() {
				object, err := sdk.Read(ctx, location, "1")
				Expect(err).ToNot(HaveOccurred())
				Expect(object.Name).To(Equal("1"))
			})
		})
	})

	Context("ReadMany()", func() {
		When("no entries are in the state", func() {
			It("should return an empty list of entries", func() {
				entries := []*MockEntryObject{}
				processed, err := sdk.ReadMany(ctx, location, entries)
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(0))
			})
		})
		When("there are entries in the state", func() {
			It("should return a list of entries from the server that match state entries", func() {
				entries := []*MockEntryObject{{Name: "1"}, {Name: "2"}}
				processed, err := sdk.ReadMany(ctx, location, entries)
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(2))

				Expect(processed[0].EntryName()).To(Equal("1"))
				Expect(processed[1].EntryName()).To(Equal("2"))
			})
		})
	})

	Context("Create()", func() {
		Context("creating a single entry on the server", func() {
			Context("when the entry already exists on the server", func() {
				It("should return an error back to the caller", func() {
					entry := &MockEntryObject{Name: "1", Value: "A"}
					processed, err := sdk.Create(ctx, location, entry)
					Expect(err).To(MatchError(MatchRegexp("already exists")))
					Expect(processed).To(BeNil())

				})
			})
			Context("when there is no conflict between plan and remote state", func() {
				It("should return a pointer to the created object", func() {
					entry := &MockEntryObject{Name: "4", Value: "D"}
					processed, err := sdk.Create(ctx, location, entry)
					Expect(err).ToNot(HaveOccurred())
					Expect(processed).ToNot(BeNil())
					Expect(processed.EntryName()).To(Equal(entry.Name))

				})
			})
		})
	})

	Context("CreateMany()", func() {
		Context("when creating new entries on the server", func() {
			Context("and some entries already exist on the server", func() {
				It("should return an error about conflict", func() {
					entries := []*MockEntryObject{{Name: "1", Value: "A"}, {Name: "4", Value: "D"}}
					processed, err := sdk.CreateMany(ctx, location, entries)

					Expect(err).To(MatchError(manager.ErrConflict))
					Expect(processed).To(BeNil())
					Expect(client.list()).To(HaveExactElements(existing))
				})
			})
		})

		Context("when creating new entries on the server", func() {
			Context("and the list of entries is empty", func() {
				It("should not make any changes", func() {
					entries := []*MockEntryObject{}
					processed, err := sdk.CreateMany(ctx, location, entries)

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(BeEmpty())
					Expect(client.list()).To(HaveExactElements(existing))
				})
			})
		})

		Context("when creating new entries on the server", func() {
			Context("and there are new entries on the list", func() {
				It("should add those entries to the server and return only those new entries", func() {
					entries := []*MockEntryObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}}
					processed, err := sdk.CreateMany(ctx, location, entries)

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveExactElements(entries))

					expected := append(existing, entries...)
					Expect(client.list()).To(HaveExactElements(expected))
				})
			})
		})
	})

	Context("Update()", func() {

	})

	Context("UpdateMany()", func() {
		Context("when entries from the plan are missing from the server", func() {
			It("should recreate them, and return back list of all managed entries", func() {
				entries := []*MockEntryObject{{Name: "4", Value: "D"}}
				processed, err := sdk.UpdateMany(ctx, location, entries, entries)

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveExactElements(entries))

				expected := append(existing, entries...)
				Expect(client.list()).To(HaveExactElements(expected))
			})

		})
	})

	Context("Delete()", func() {
		Context("when entries from the plan are missing from the server", func() {
			It("should not delete anything from the server", func() {
				entries := []string{"4"}
				err := sdk.Delete(ctx, location, entries)

				Expect(err).ToNot(HaveOccurred())
				Expect(client.list()).To(HaveExactElements(existing))
			})
		})
	})
})
