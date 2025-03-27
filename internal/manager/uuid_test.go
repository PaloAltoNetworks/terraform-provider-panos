package manager_test

import (
	"context"
	"log"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PaloAltoNetworks/pango/movement"
	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

var _ = log.Printf
var _ = Expect
var _ = slog.Debug

type MockConfigObject struct {
	Value string
}

func (o MockConfigObject) EntryName() string {
	panic("unimplemented")
}

var _ = Describe("Entry", func() {

})

var _ = Describe("Server", func() {
	var initial []*MockUuidObject
	var manager *sdkmanager.UuidObjectManager[*MockUuidObject, MockLocation, sdkmanager.SDKUuidService[*MockUuidObject, MockLocation]]
	var client *MockUuidClient[*MockUuidObject]
	var service sdkmanager.SDKUuidService[*MockUuidObject, MockLocation]
	var mockService *MockUuidService[*MockUuidObject, MockLocation]
	var location MockLocation
	var ctx context.Context

	var position movement.Position
	var entries []*MockUuidObject
	var mode sdkmanager.ExhaustiveType

	BeforeEach(func() {
		location = MockLocation{}
		ctx = context.Background()
		initial = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
		client = NewMockUuidClient(initial)
		service = NewMockUuidService[*MockUuidObject, MockLocation](client)
		var ok bool
		if mockService, ok = service.(*MockUuidService[*MockUuidObject, MockLocation]); !ok {
			panic("failed to cast service to mockService")
		}
		manager = sdkmanager.NewUuidObjectManager(client, service, MockUuidSpecifier, MockUuidMatcher)
	})

	Describe("Creating new resources on the server", func() {
		Context("When server has no entries yet", func() {
			BeforeEach(func() {
				initial := []*MockUuidObject{}
				client = NewMockUuidClient(initial)
				service = NewMockUuidService[*MockUuidObject, MockLocation](client)
				manager = sdkmanager.NewUuidObjectManager(client, service, MockUuidSpecifier, MockUuidMatcher)
			})

			It("CreateMany() should create new entries on the server, and return them with uuid set", func() {
				entries := []*MockUuidObject{{Name: "1", Value: "A"}}
				processed, err := manager.CreateMany(ctx, location, entries, sdkmanager.Exhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(1))
				Expect(processed).To(MatchEntries(entries))

				current := client.list()
				Expect(current).To(HaveLen(1))
				Expect(current).To(MatchEntries(entries))
			})
		})

		Context("When server already has some entries", func() {
			BeforeEach(func() {
				entries = []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}}
			})

			Context("and entries with the same name are being created in NonExhaustive mode", func() {
				BeforeEach(func() {
					entries = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "4", Value: "D"}}
					mode = sdkmanager.NonExhaustive
				})

				It("should not create any entries and return an error", func() {
					processed, err := manager.CreateMany(ctx, location, entries, mode, position)

					Expect(err).To(MatchError(sdkmanager.ErrConflict))
					Expect(processed).To(BeNil())

					Expect(client.list()).To(Equal(initial))
				})
			})

			Context("and all entries being created are new to the server", func() {
				It("should create those entries in the correct position", func() {
					processed, err := manager.CreateMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(2))

					Expect(processed).To(Equal(entries))
					Expect(mockService.moveGroupEntries).To(Equal(entries))

					current := client.list()
					Expect(current[0:2]).To(MatchEntries(processed))
				})
			})

			Context("and entries are created in Exhaustive mode", func() {
				BeforeEach(func() {
					entries = []*MockUuidObject{{Name: "1", Value: "A'"}, {Name: "3", Value: "C"}}
					mode = sdkmanager.Exhaustive
					position = movement.PositionFirst{}
				})

				It("should not return any error and overwrite all entries on the server", func() {
					processed, err := manager.CreateMany(ctx, location, entries, mode, position)

					Expect(err).ToNot(HaveOccurred())

					// We don't want to mutate the provided list of entries, but we have to pass
					// them via pointer to satisfy generic type. Make sure uuid is still nil.
					Expect(entries[0].Uuid).To(BeNil())

					Expect(client.MultiConfigOpers).To(HaveExactElements([]MultiConfigOper{
						{Operation: MultiConfigOperDelete, EntryName: "1"},
						{Operation: MultiConfigOperDelete, EntryName: "2"},
						{Operation: MultiConfigOperDelete, EntryName: "3"},
						{Operation: MultiConfigOperEdit, EntryName: "1"},
						{Operation: MultiConfigOperEdit, EntryName: "3"},
					}))

					Expect(processed).To(MatchEntries(entries))

					current := client.list()
					Expect(current).To(HaveLen(2))
					Expect(current).To(MatchEntries(entries))
				})
			})
		})
	})

	Context("updating existing entries", func() {
		Context("when some of the entries are out of order", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
				client = NewMockUuidClient(initial)
				service = NewMockUuidService[*MockUuidObject, MockLocation](client)
				var ok bool
				if mockService, ok = service.(*MockUuidService[*MockUuidObject, MockLocation]); !ok {
					panic("failed to cast service to mockService")
				}
				manager = sdkmanager.NewUuidObjectManager(client, service, MockUuidSpecifier, MockUuidMatcher)

			})
			It("should move the entries in order", func() {
				entries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "3", Value: "C"}, {Name: "2", Value: "B"}}

				processed, err := manager.ReadMany(ctx, location, entries, sdkmanager.NonExhaustive)

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(processed).NotTo(MatchEntries(entries))

				processed, err = manager.UpdateMany(ctx, location, entries, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))

				Expect(processed).To(Equal(entries))
			})
		})
	})

	Context("initially has some entries", func() {
		Context("when creating new entries with NonExhaustive type", func() {
			Context("and position is set to first", func() {
				It("should create new entries on the top of the list", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					processed, err := manager.CreateMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})
					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))

					Expect(processed[0]).To(Equal(entries[0]))
					Expect(processed[1]).To(Equal(entries[1]))
					Expect(processed[2]).To(Equal(entries[2]))

					clientEntries := client.list()
					Expect(clientEntries).To(HaveLen(6))

					Expect(mockService.moveGroupEntries).To(Equal(entries))

					Expect(clientEntries[0]).To(Equal(entries[0]))
					Expect(clientEntries[1]).To(Equal(entries[1]))
					Expect(clientEntries[2]).To(Equal(entries[2]))

				})
			})
			Context("and position is set to last", func() {
				It("should create new entries on the bottom of the list", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					processed, err := manager.CreateMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionLast{})
					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))

					Expect(processed[0]).To(Equal(entries[0]))
					Expect(processed[1]).To(Equal(entries[1]))
					Expect(processed[2]).To(Equal(entries[2]))

					clientEntries := client.list()
					Expect(clientEntries).To(HaveLen(6))

					Expect(mockService.moveGroupEntries).To(Equal(entries))

					Expect(clientEntries[3]).To(Equal(entries[0]))
					Expect(clientEntries[4]).To(Equal(entries[1]))
					Expect(clientEntries[5]).To(Equal(entries[2]))

				})
			})
			Context("and position is set to directly after first element", func() {
				It("should create new entries directly after first existing element", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					processed, err := manager.CreateMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionAfter{Directly: true, Pivot: initial[0].Name})

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))

					Expect(processed[0]).To(Equal(entries[0]))
					Expect(processed[1]).To(Equal(entries[1]))
					Expect(processed[2]).To(Equal(entries[2]))

					clientEntries := client.list()
					Expect(clientEntries).To(HaveLen(6))

					Expect(clientEntries[1]).To(Equal(entries[0]))
					Expect(clientEntries[2]).To(Equal(entries[1]))
					Expect(clientEntries[3]).To(Equal(entries[2]))

					Expect(clientEntries[0]).To(Equal(initial[0]))
					Expect(clientEntries[4]).To(Equal(initial[1]))
					Expect(clientEntries[5]).To(Equal(initial[2]))

					Expect(mockService.moveGroupEntries).To(Equal(entries))
				})
			})
			Context("and position is set to directly before last element", func() {
				It("should create new entries directly before last element", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					pivot := initial[2].Name // "3"
					position = movement.PositionBefore{Directly: true, Pivot: pivot}
					processed, err := manager.CreateMany(ctx, location, entries, sdkmanager.NonExhaustive, position)

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))
					Expect(processed).To(MatchEntries(entries))

					current := client.list()
					Expect(current).To(HaveLen(6))
					Expect(current[2:5]).To(MatchEntries(entries))

					Expect(current[0:1]).To(MatchEntries(initial[0:1]))
					Expect(current[5:5]).To(MatchEntries(initial[2:2]))

					Expect(mockService.moveGroupEntries).To(Equal(entries))
				})
			})
			Context("and there is a duplicate entry within a list", func() {
				It("should properly raise an error", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "4", Value: "D"}}
					_, err := manager.CreateMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})

					Expect(err).To(MatchError(sdkmanager.ErrPlanConflict))
				})
			})
		})
	})
})
