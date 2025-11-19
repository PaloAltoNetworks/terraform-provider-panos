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
	var batchSize int

	BeforeEach(func() {
		initial = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
		location = MockLocation{}
	})

	JustBeforeEach(func() {
		batchSize = 500
		ctx = context.Background()

		client = NewMockUuidClient(initial)
		service = NewMockUuidService[*MockUuidObject, MockLocation](client)
		var ok bool
		if mockService, ok = service.(*MockUuidService[*MockUuidObject, MockLocation]); !ok {
			panic("failed to cast service to mockService")
		}
		manager = sdkmanager.NewUuidObjectManager(client, service, batchSize, MockUuidSpecifier, MockUuidMatcher)
	})

	Describe("Reading entries from the server", func() {
		BeforeEach(func() {
			initial = []*MockUuidObject{
				{Name: "1", Value: "A", Location: "parent"},
				{Name: "2", Value: "B", Location: "child"},
				{Name: "3", Value: "C", Location: "child"},
			}
			location = MockLocation{Filter: "child"}
		})

		Context("with exhaustive mode", func() {
			Context("when entries have mixed device group location", func() {
				It("should only return entries from the specific device group", func() {
					entries := []*MockUuidObject{{Name: "2", Value: "B", Location: "child"}, {Name: "3", Value: "C", Location: "child"}}
					processed, movementRequired, err := manager.ReadMany(ctx, location, entries, sdkmanager.Exhaustive, movement.PositionFirst{})
					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(MatchEntries(entries))
					Expect(movementRequired).To(Equal(false))
				})
			})
		})
		Context("with non-exhaustive mode", func() {
			It("should return only entries managed by the state", func() {
				entries := []*MockUuidObject{{Name: "2", Value: "B", Location: "child"}}
				processed, movementRequired, err := manager.ReadMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionLast{})
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(MatchEntries(entries))
				// initial is [1, 2, 3] and state is [2]. 1 is from a different location, 3 is unmanaged
				// by the state, 2 is supposed to be at the end so movementRequired is true.
				Expect(movementRequired).To(Equal(true))
			})
		})
	})

	Describe("Creating new resources on the server", func() {
		Context("When server has no entries yet", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{}
			})

			It("CreateMany() should create new entries on the server, and return them with uuid set", func() {
				entries := []*MockUuidObject{{Name: "1", Value: "A"}}
				processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.Exhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(1))
				Expect(processed).To(MatchEntries(entries))

				current := client.list()
				Expect(current).To(HaveLen(1))
				Expect(current).To(MatchEntries(entries))
			})
		})

		Context("When server already has some entries", func() {
			Context("and entries with the same name are being created in NonExhaustive mode", func() {
				It("should not create any entries and return an error", func() {
					entries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "4", Value: "D"}}
					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})

					Expect(err).To(MatchError(sdkmanager.ErrConflict))
					Expect(processed).To(BeNil())

					Expect(client.list()).To(Equal(initial))
				})
			})

			Context("and all entries being created are new to the server", func() {
				It("should create those entries in the correct position", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}}
					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(2))

					Expect(processed).To(MatchEntries(entries))
					Expect(mockService.moveGroupEntries).To(Equal(entries))

					current := client.list()
					Expect(current[0:2]).To(MatchEntries(processed))
				})
			})

			Context("and entries are created in Exhaustive mode", func() {
				It("should not return any error and overwrite all entries on the server", func() {
					entries := []*MockUuidObject{{Name: "1", Value: "A'"}, {Name: "3", Value: "C"}}
					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.Exhaustive, movement.PositionFirst{})

					Expect(err).ToNot(HaveOccurred())

					// We don't want to mutate the provided list of entries, but we have to pass
					// them via pointer to satisfy generic type. Make sure uuid is still nil.
					Expect(entries[0].Uuid).To(BeNil())

					Expect(client.MultiConfigOpers[0]).To(HaveExactElements([]MultiConfigOper{
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

	Context("Updating resource on the server", func() {
		Context("with NonExhaustive mode, when entries are out of order", func() {
			It("should reorder the entries on the server", func() {
				stateEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "3", Value: "C"}, {Name: "2", Value: "B"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, stateEntries, planEntries, sdkmanager.NonExhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(client.list()).To(MatchEntries(planEntries))

				Expect(client.MultiConfigOpers).To(HaveLen(1))
				Expect(client.MultiConfigOpers[0]).To(ContainElements([]MultiConfigOper{
					{Operation: "move", EntryName: "3", Where: "after", Destination: "1"},
				}))
			})
		})

		Context("with Exhaustive mode, when entries are out of order", func() {
			It("should reorder the entries on the server", func() {
				stateEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "3", Value: "C"}, {Name: "2", Value: "B"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, stateEntries, planEntries, sdkmanager.Exhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(client.list()).To(MatchEntries(planEntries))

				Expect(client.MultiConfigOpers).To(HaveLen(1))
				Expect(client.MultiConfigOpers[0]).To(ContainElements([]MultiConfigOper{
					{Operation: "move", EntryName: "3", Where: "after", Destination: "1"},
				}))
			})
		})

		Context("when adding a new entry", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}}
			})

			It("should add the entry to the server", func() {
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}
				processed, err := manager.UpdateMany(ctx, location, []string{}, initial, planEntries, sdkmanager.NonExhaustive, movement.PositionLast{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(client.list()).To(MatchEntries(planEntries))
			})
		})

		Context("when deleting an entry", func() {
			It("should delete the entry from the server", func() {
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "3", Value: "C"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, initial, planEntries, sdkmanager.NonExhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(2))
				Expect(client.list()).To(MatchEntries(planEntries))
			})
		})

		Context("when modifying an existing entry", func() {
			It("should update the entry on the server", func() {
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B_modified"}, {Name: "3", Value: "C"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, initial, planEntries, sdkmanager.NonExhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(client.list()).To(MatchEntries(planEntries))
			})
		})

		Context("when renaming an existing entry", func() {
			It("should rename the entry on the server", func() {
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "two", Value: "B"}, {Name: "3", Value: "C"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, initial, planEntries, sdkmanager.NonExhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(client.list()).To(MatchEntries(planEntries))

				Expect(client.MultiConfigOpers).To(HaveLen(1))
				var renameOp bool
				for _, op := range client.MultiConfigOpers[0] {
					if op.Operation == MultiConfigOperRename && op.EntryName == "2" {
						renameOp = true
						break
					}
				}
				Expect(renameOp).To(BeTrue(), "Expected a rename operation for entry '2'")
			})
		})

		Context("in exhaustive mode with unmanaged entries on server", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{{Name: "0", Value: "A"}, {Name: "99", Value: "Z"}}
			})

			It("should delete unmanaged entries", func() {
				stateEntries := []*MockUuidObject{}
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, stateEntries, planEntries, sdkmanager.Exhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(2))
				Expect(client.list()).To(MatchEntries(planEntries))
			})
		})

		Context("when plan conflicts with an unmanaged entry", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "conflict", Value: "Z"}}
			})

			It("should return a conflict error", func() {
				stateEntries := []*MockUuidObject{{Name: "1", Value: "A"}}
				planEntries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "conflict", Value: "new"}}
				_, err := manager.UpdateMany(ctx, location, []string{}, stateEntries, planEntries, sdkmanager.NonExhaustive, movement.PositionFirst{})

				Expect(err).To(MatchError(sdkmanager.ErrConflict))
			})
		})

		Context("with combined operations", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{
					{Name: "1", Value: "A"},
					{Name: "2", Value: "B"},
					{Name: "3", Value: "C"},
					{Name: "4", Value: "D"},
					{Name: "99", Value: "ZZ"},
				}
			})

			It("should correctly apply all changes", func() {
				stateEntries := []*MockUuidObject{
					{Name: "1", Value: "A"},
					{Name: "2", Value: "B"},
					{Name: "3", Value: "C"},
					{Name: "4", Value: "D"},
				}

				planEntries := []*MockUuidObject{
					// { Name: "2", Value: "B"},      // deleted
					{Name: "5", Value: "E"},          // added
					{Name: "3", Value: "C"},          // reordered
					{Name: "four", Value: "D"},       // renamed
					{Name: "1", Value: "A_modified"}, // modified
				}

				expectedFinalState := []*MockUuidObject{
					{Name: "5", Value: "E"},
					{Name: "3", Value: "C"},
					{Name: "four", Value: "D"},
					{Name: "1", Value: "A_modified"},
				}

				expectedFinalServerObjects := append(expectedFinalState, &MockUuidObject{Name: "99", Value: "ZZ"})

				processed, err := manager.UpdateMany(ctx, location, []string{}, stateEntries, planEntries, sdkmanager.NonExhaustive, movement.PositionFirst{})

				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(4))
				Expect(processed).To(MatchEntries(expectedFinalState))
				Expect(client.list()).To(MatchEntries(expectedFinalServerObjects))

				Expect(client.MultiConfigOpers[0]).To(ContainElements([]MultiConfigOper{
					{Operation: "rename", EntryName: "4", NewName: "four"},
					{Operation: "edit", EntryName: "1"},
					{Operation: "delete", EntryName: "2"},
					{Operation: "edit", EntryName: "5"},
				}))

				Expect(client.MultiConfigOpers[1]).To(HaveExactElements([]MultiConfigOper{
					{Operation: "move", EntryName: "5", Where: "top", Destination: "top"},
					{Operation: "move", EntryName: "3", Where: "after", Destination: "5"},
					{Operation: "move", EntryName: "four", Where: "after", Destination: "3"},
				}))
			})
		})

		Context("when multiple entries on the server are equal", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{
					{Name: "1", Value: ""},
					{Name: "2", Value: ""},
					{Name: "3", Value: "C"},
					{Name: "4", Value: ""},
				}
			})

			It("should properly handle the requested update", func() {
				state := []*MockUuidObject{{Name: "1", Value: ""}, {Name: "2", Value: ""}, {Name: "3", Value: "C"}}
				plan := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}

				processed, err := manager.UpdateMany(ctx, location, []string{}, state, plan, sdkmanager.NonExhaustive, movement.PositionLast{})
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(MatchEntries(plan))
				Expect(client.list()).To(MatchEntries(append([]*MockUuidObject{{Name: "4", Value: ""}}, plan...)))
			})
		})
	})

	Context("Delete()", func() {
		When("deleting entries that exist", func() {
			It("should delete the entries from the server", func() {
				entries := []string{"1", "3"}
				err := manager.Delete(ctx, location, []string{}, entries, sdkmanager.NonExhaustive)

				Expect(err).ToNot(HaveOccurred())

				remaining := client.list()
				Expect(remaining).To(HaveLen(1))
				Expect(remaining[0].EntryName()).To(Equal("2"))
			})
		})

		When("deleting entries that are missing from the server", func() {
			It("should not change the list of entries on the server", func() {
				entries := []string{"4"}
				err := manager.Delete(ctx, location, []string{}, entries, sdkmanager.NonExhaustive)

				Expect(err).ToNot(HaveOccurred())
				Expect(client.list()).To(MatchEntries(initial))
			})
		})

		When("deleting a mix of existing and missing entries", func() {
			It("should only delete the existing entries from the server", func() {
				entries := []string{"1", "4"}
				err := manager.Delete(ctx, location, []string{}, entries, sdkmanager.NonExhaustive)

				Expect(err).ToNot(HaveOccurred())

				remaining := client.list()
				Expect(remaining).To(HaveLen(2))
				Expect(remaining[0].EntryName()).To(Equal("2"))
				Expect(remaining[1].EntryName()).To(Equal("3"))
			})
		})

		Context("when some of the entries were removed from the server", func() {
			BeforeEach(func() {
				initial = []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "3", Value: "C"}}
				client = NewMockUuidClient(initial)
				service = NewMockUuidService[*MockUuidObject, MockLocation](client)
				var ok bool
				if mockService, ok = service.(*MockUuidService[*MockUuidObject, MockLocation]); !ok {
					panic("failed to cast service to mockService")
				}
				manager = sdkmanager.NewUuidObjectManager(client, service, batchSize, MockUuidSpecifier, MockUuidMatcher)

			})

			It("should recreate missing entries on the server based on the state", func() {
				entries := []*MockUuidObject{{Name: "1", Value: "A"}, {Name: "2", Value: "B"}, {Name: "3", Value: "C"}}

				processed, moveRequired, err := manager.ReadMany(ctx, location, entries, sdkmanager.NonExhaustive, movement.PositionLast{})

				Expect(err).ToNot(HaveOccurred())
				Expect(moveRequired).To(BeFalse())
				Expect(processed).To(HaveLen(2))

				processed, err = manager.UpdateMany(ctx, location, []string{}, processed, entries, sdkmanager.NonExhaustive, movement.PositionLast{})
				Expect(client.list()).To(HaveLen(3))
				Expect(err).ToNot(HaveOccurred())
				Expect(processed).To(HaveLen(3))
				Expect(processed).To(MatchEntries(entries))
			})
		})
	})

	Context("initially has some entries", func() {
		Context("when creating new entries with NonExhaustive type", func() {
			Context("and position is set to first", func() {
				It("should create new entries on the top of the list", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})
					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))

					Expect(processed).To(MatchEntries(entries))

					clientEntries := client.list()
					Expect(clientEntries).To(HaveLen(6))

					Expect(mockService.moveGroupEntries).To(Equal(entries))

					Expect(clientEntries[0:3]).To(MatchEntries(entries))
				})
			})
			Context("and position is set to last", func() {
				It("should create new entries on the bottom of the list", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, movement.PositionLast{})
					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))

					Expect(processed).To(MatchEntries(entries))

					clientEntries := client.list()
					Expect(clientEntries).To(HaveLen(6))

					Expect(mockService.moveGroupEntries).To(Equal(entries))

					Expect(clientEntries[3:]).To(MatchEntries(entries))
				})
			})
			Context("and position is set to directly after first element", func() {
				It("should create new entries directly after first existing element", func() {
					entries := []*MockUuidObject{{Name: "4", Value: "D"}, {Name: "5", Value: "E"}, {Name: "6", Value: "F"}}

					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, movement.PositionAfter{Directly: true, Pivot: initial[0].Name})

					Expect(err).ToNot(HaveOccurred())
					Expect(processed).To(HaveLen(3))

					Expect(processed).To(MatchEntries(entries))

					clientEntries := client.list()
					Expect(clientEntries).To(HaveLen(6))

					Expect(clientEntries[1:4]).To(MatchEntries(entries))

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
					position := movement.PositionBefore{Directly: true, Pivot: pivot}
					processed, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, position)

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
					_, err := manager.CreateMany(ctx, location, []string{}, entries, sdkmanager.NonExhaustive, movement.PositionFirst{})

					Expect(err).To(MatchError(sdkmanager.ErrPlanConflict))
				})
			})
		})
	})

})
