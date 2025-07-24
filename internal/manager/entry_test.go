package manager_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

var _ = Expect

var _ = Describe("Entry", func() {
	var existing []*MockEntryObject
	var client *MockEntryClient[*MockEntryObject]
	var service *MockEntryService[*MockEntryObject, MockLocation]
	var sdk *manager.EntryObjectManager[*MockEntryObject, MockLocation, *MockEntryService[*MockEntryObject, MockLocation]]
	var batchSize int

	var location MockLocation

	ctx := context.Background()

	JustBeforeEach(func() {
		batchSize = 500
		client = NewMockEntryClient(existing)
		service = NewMockEntryService[*MockEntryObject, MockLocation](client)
		sdk = manager.NewEntryObjectManager[*MockEntryObject, MockLocation, *MockEntryService[*MockEntryObject, MockLocation]](client, service, batchSize, MockEntrySpecifier, MockEntryMatcher)
	})

	BeforeEach(func() {
		existing = []*MockEntryObject{
			{Location: "parent", Name: "1", Value: "A"},
			{Location: "child", Name: "2", Value: "B"},
			{Location: "child", Name: "3", Value: "C"},
		}
	})

	Context("Read()", func() {
		When("reading entry that does not exist", func() {
			It("should return nil object and ErrObjectNotFound error", func() {
				object, err := sdk.Read(ctx, location, []string{}, "4")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(MatchRegexp("Object not found")))
				Expect(object).To(BeNil())
			})
		})
		When("reading entry that exists without any location filter", func() {
			It("should return nil error and the existing entry", func() {
				object, err := sdk.Read(ctx, location, []string{}, "1")
				Expect(err).ToNot(HaveOccurred())
				Expect(object.Name).To(Equal("1"))
			})
		})
		When("reading entry that exists with matching location filter", func() {
			BeforeEach(func() {
				location = MockLocation{
					Filter: "child",
				}
			})
			It("should return nil error and the existing entry", func() {
				object, err := sdk.Read(ctx, location, []string{}, "3")
				Expect(err).ToNot(HaveOccurred())
				Expect(object.Name).To(Equal("3"))
			})
		})
		When("reading entry that exists with non-matching location filter", func() {
			BeforeEach(func() {
				location = MockLocation{
					Filter: "child",
				}
			})
			It("should return nil error and the existing entry", func() {
				object, err := sdk.Read(ctx, location, []string{}, "1")
				Expect(err).Should(MatchError(manager.ErrObjectNotFound))
				Expect(object).To(BeNil())
			})
		})
	})

	Context("ReadMany()", func() {
		When("location has a filter", func() {
			BeforeEach(func() {
				location = MockLocation{
					Filter: "child",
				}
			})
			It("should return a list entries filtered by location", func() {
				processed, err := sdk.ReadMany(ctx, location, []string{})
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(2))

				Expect(processed[0].EntryName()).To(Equal("2"))
				Expect(processed[1].EntryName()).To(Equal("3"))
			})
		})
		When("location has no filter", func() {
			BeforeEach(func() {
				location = MockLocation{}
			})
			It("should return a list of entries from the server that match state entries", func() {
				processed, err := sdk.ReadMany(ctx, location, []string{})
				Expect(err).ToNot(HaveOccurred())

				Expect(processed).To(MatchEntries(existing))
			})
		})
	})

	Context("Create()", func() {
		Context("creating a single entry on the server", func() {
			Context("when the entry already exists on the server", func() {
				It("should return an error back to the caller", func() {
					entry := &MockEntryObject{Name: "1", Value: "A"}
					processed, err := sdk.Create(ctx, location, []string{entry.EntryName()}, entry)
					Expect(err).To(MatchError(MatchRegexp("already exists")))
					Expect(processed).To(BeNil())

				})
			})
			Context("when there is no conflict between plan and remote state", func() {
				It("should return a pointer to the created object", func() {
					entry := &MockEntryObject{Name: "4", Value: "D"}
					processed, err := sdk.Create(ctx, location, []string{entry.EntryName()}, entry)
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
					processed, err := sdk.CreateMany(ctx, location, []string{}, entries)

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
					processed, err := sdk.CreateMany(ctx, location, []string{}, entries)

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
					processed, err := sdk.CreateMany(ctx, location, []string{}, entries)

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
				expected := append(existing, &MockEntryObject{Name: "4", Value: "D"})
				processed, err := sdk.UpdateMany(ctx, location, []string{}, existing, expected)

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveExactElements(expected))
				Expect(client.list()).To(HaveExactElements(expected))
			})

		})

		Context("when some entries are removed from the plan", func() {
			It("should properly remove deleted entries from the server and return back updated list", func() {
				stateEntries := []*MockEntryObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
				planEntries := []*MockEntryObject{{Name: "1", Value: "A"}, {Name: "3", Value: "C"}}
				processed, err := sdk.UpdateMany(ctx, location, []string{}, stateEntries, planEntries)

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(MatchEntries(planEntries))
				Expect(client.list()).To(MatchEntries(planEntries))
			})
		})
	})

	Context("Delete()", func() {
		Context("when entries from the plan are missing from the server", func() {
			It("should not delete anything from the server", func() {
				entries := []string{"4"}
				err := sdk.Delete(ctx, location, []string{}, entries)

				Expect(err).ToNot(HaveOccurred())
				Expect(client.list()).To(MatchEntries(existing))
			})
		})
	})
})
