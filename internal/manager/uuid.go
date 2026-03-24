package manager

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/movement"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type ExhaustiveType int

const (
	Exhaustive    ExhaustiveType = iota
	NonExhaustive ExhaustiveType = iota
)

type UuidObject interface {
	EntryName() string
	SetEntryName(name string)
	EntryUuid() *string
	SetEntryUuid(value *string)
	GetMiscAttributes() []xml.Attr
}

type UuidLocation interface {
	LocationFilter() *string
	XpathWithComponents(version.Number, ...string) ([]string, error)
}

type SDKUuidService[E UuidObject, L UuidLocation] interface {
	Create(context.Context, L, E) (E, error)
	List(context.Context, L, string, string, string) ([]E, error)
	Delete(context.Context, L, ...string) error
	MoveGroup(context.Context, L, movement.Position, []E, int) error
}

type uuidObjectWithState[E EntryObject] struct {
	Entry    E
	State    entryState
	StateIdx int
	NewName  string
}

type TFUuidObject[E any] interface {
	EntryName() string
	EntryUuid() *string
	CopyToPango(context.Context, *map[string]types.String) (E, diag.Diagnostics)
	CopyFromPango(context.Context, E, *map[string]types.String) diag.Diagnostics
}

type UuidObjectManager[E UuidObject, L UuidLocation, S SDKUuidService[E, L]] struct {
	batchSize int
	service   S
	client    SDKClient
	specifier func(E) (any, error)
	matcher   func(E, E) bool
}

func NewUuidObjectManager[E UuidObject, L UuidLocation, S SDKUuidService[E, L]](client SDKClient, service S, batchSize int, specifier func(E) (any, error), matcher func(E, E) bool) *UuidObjectManager[E, L, S] {
	return &UuidObjectManager[E, L, S]{
		batchSize: batchSize,
		service:   service,
		client:    client,
		specifier: specifier,
		matcher:   matcher,
	}
}

func (o *UuidObjectManager[E, L, S]) entriesByName(entries []E, initialState entryState) map[string]uuidObjectWithState[E] {
	entriesByName := make(map[string]uuidObjectWithState[E], len(entries))
	for idx, elt := range entries {
		entriesByName[elt.EntryName()] = uuidObjectWithState[E]{
			Entry:    elt,
			StateIdx: idx,
			State:    initialState,
		}
	}

	return entriesByName
}

func (o *UuidObjectManager[E, L, S]) entriesProperlySorted(existing []E, planEntriesByName map[string]uuidObjectWithState[E]) (bool, error) {
	// All entries returned from the server are filtered out, to gain a list
	// of entries that are part of the plan. For every entry, we calculate its
	// actual index in the plan so we can compare it with the expected one.
	existingEntriesByName := make(map[string]uuidObjectWithState[E], len(existing))
	managedEntriesByName := make(map[string]uuidObjectWithState[E], len(planEntriesByName))
	idx := 0

	for existingIdx, elt := range existing {
		name := elt.EntryName()
		existingEntriesByName[name] = uuidObjectWithState[E]{
			Entry:    existing[idx],
			StateIdx: existingIdx,
		}

		if planEntry, found := planEntriesByName[name]; found {
			if elt.EntryUuid() == nil {
				return false, ErrMissingUuid
			}

			// moveNonExhausitve is called just after we've executed MultiConfig which could have created
			// some new entries on the server. If so, we want to make sure new UUIDs are stored in the state.
			planEntry.Entry = elt
			managedEntriesByName[name] = uuidObjectWithState[E]{
				Entry: existing[idx],
				// managedEntriesByName StateIdx reflects index of the entry on the server.
				StateIdx: idx,
			}
		}
		idx++
	}

	var movementRequired bool

	var previousManagedEntry, previousPlanEntry *uuidObjectWithState[E]
	// First step is to check whether any of the elements from the plan are out
	// of order on the server.
	for _, elt := range managedEntriesByName {
		// plannedEntriesByName is a map of entries from the plan indexed by their
		// name, where each element has StateIdx, indicating its position in the plan.
		// Index of the entry in the plan is compared against its remote index (relative
		// to other entries from the plan).
		planEntry := planEntriesByName[(elt.Entry).EntryName()]
		if planEntry.StateIdx != elt.StateIdx {
			movementRequired = true
			break
		}

		// If this is the first entry to compare, store it for future reference for
		// both this managed entry and the plan entry. It will be used later to calculate
		// distance between two entries on the list.
		if previousManagedEntry == nil {
			previousManagedEntry = &elt
			previousPlanEntry = &planEntry
			continue
		}

		serverDistance := elt.StateIdx - previousManagedEntry.StateIdx
		planDistance := planEntry.StateIdx - previousPlanEntry.StateIdx

		// If the distance between previous and the current element differs between
		// PAN-OS and Terraform plan, we need to move entries on the server.
		if serverDistance != planDistance {
			movementRequired = true
			break
		}

		previousManagedEntry = &elt
		previousPlanEntry = &planEntry
	}

	return movementRequired, nil
}

func (o *UuidObjectManager[E, L, S]) filterEntriesByLocation(location L, entries []E) []E {
	filter := location.LocationFilter()
	if filter == nil {
		return entries
	}

	getLocAttribute := func(entry E) *string {
		for _, elt := range entry.GetMiscAttributes() {
			if elt.Name.Local == "loc" {
				return &elt.Value
			}
		}
		return nil
	}

	var filtered []E
	for _, elt := range entries {
		location := getLocAttribute(elt)
		if location == nil || *location == *filter {
			filtered = append(filtered, elt)
		}
	}

	return filtered
}

func (o *UuidObjectManager[E, L, S]) moveExhaustive(ctx context.Context, location L, entriesByName map[string]uuidObjectWithState[E], position movement.Position) error {
	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && err.Error() != "Object not found" {
		return &Error{err: err, message: "Failed to list existing entries"}
	}

	existing = o.filterEntriesByLocation(location, existing)

	movementRequired, err := o.entriesProperlySorted(existing, entriesByName)
	if err != nil {
		return err
	}

	if movementRequired {
		entries := make([]E, len(entriesByName))
		for _, elt := range entriesByName {
			entries[elt.StateIdx] = elt.Entry
		}

		err = o.service.MoveGroup(ctx, location, position, entries, o.batchSize)
		if err != nil {
			return &Error{err: err, message: "Failed to move group of entries"}
		}
	}

	return nil
}

type PositionWhereType string

const (
	PositionWhereFirst  PositionWhereType = "first"
	PositionWhereLast   PositionWhereType = "last"
	PositionWhereBefore PositionWhereType = "before"
	PositionWhereAfter  PositionWhereType = "after"
)

type position struct {
	Where      PositionWhereType
	PivotEntry string
	Directly   bool
}

// moveNonExhaustive implements algorithm for executing movements for a subset of a list entries
//
// When moveNonExhaustive is called, the given list is not entirely managed by the Terraform resource.
// In that case a care has to be taken to only execute movement on a subset of entries, those that
// are under Terraform control.
func (o *UuidObjectManager[E, L, S]) moveNonExhaustive(ctx context.Context, location L, planEntries []E, planEntriesByName map[string]uuidObjectWithState[E], sdkPosition movement.Position) error {
	entries := make([]E, len(planEntriesByName))
	for _, elt := range planEntriesByName {
		entries[elt.StateIdx] = elt.Entry
	}

	err := o.service.MoveGroup(ctx, location, sdkPosition, entries, o.batchSize)
	if err != nil {
		return &Error{err: err, message: "Failed to move group of entries"}
	}

	return nil
}

func (o *UuidObjectManager[E, L, S]) CreateMany(ctx context.Context, location L, parentComponents []string, planEntries []E, exhaustive ExhaustiveType, sdkPosition movement.Position) ([]E, error) {
	var diags diag.Diagnostics

	planEntriesByName := o.entriesByName(planEntries, entryUnknown)
	if len(planEntriesByName) != len(planEntries) {
		return nil, ErrPlanConflict
	}

	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, fmt.Errorf("Failed to list remote entries: %w", err)
	}

	existing = o.filterEntriesByLocation(location, existing)

	updates := xmlapi.NewChunkedMultiConfig(len(planEntries), o.batchSize)

	switch exhaustive {
	case Exhaustive:
		for _, elt := range existing {
			components := append(parentComponents, util.AsEntryXpath(elt.EntryName()))
			path, err := location.XpathWithComponents(o.client.Versioning(), components...)
			if err != nil {
				return nil, err
			}

			updates.Add(&xmlapi.Config{
				Action: "delete",
				Xpath:  util.AsXpath(path),
				Target: o.client.GetTarget(),
			})
		}
	case NonExhaustive:
		for _, elt := range existing {
			if _, found := planEntriesByName[elt.EntryName()]; found {
				return nil, ErrConflict
			}
		}
	}

	for _, elt := range planEntries {
		path, err := location.XpathWithComponents(o.client.Versioning(), util.AsEntryXpath(elt.EntryName()))
		if err != nil {
			return nil, err
		}

		xmlEntry, err := o.specifier(elt)
		if err != nil {
			diags.AddError("Failed to transform entry into XML document", err.Error())
			return nil, ErrMarshaling
		}

		updates.Add(&xmlapi.Config{
			Action:  "edit",
			Xpath:   util.AsXpath(path),
			Element: xmlEntry,
			Target:  o.client.GetTarget(),
		})
	}

	if len(updates.Operations) > 0 {
		_, err := o.client.ChunkedMultiConfig(ctx, updates, false, nil)
		if err != nil {
			cleanupErr := o.cleanUpIncompleteUpdate(ctx, location, parentComponents, planEntriesByName, exhaustive)
			return nil, errors.Join(fmt.Errorf("failed to create entries on the server: %w", err), cleanupErr)
		}
	}

	switch exhaustive {
	case Exhaustive:
		err := o.moveExhaustive(ctx, location, planEntriesByName, sdkPosition)
		if err != nil {
			return nil, err
		}
	case NonExhaustive:
		err := o.moveNonExhaustive(ctx, location, planEntries, planEntriesByName, sdkPosition)
		if err != nil {
			return nil, err
		}
	}

	existing, err = o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, fmt.Errorf("Failed to list remote entries: %w", err)
	}

	existing = o.filterEntriesByLocation(location, existing)

	entries := make([]E, len(planEntries))
	for _, elt := range existing {
		if planEntry, found := planEntriesByName[elt.EntryName()]; found {
			entries[planEntry.StateIdx] = elt
		}
	}

	return entries, nil
}

func (o *UuidObjectManager[E, L, S]) UpdateMany(ctx context.Context, location L, parentComponents []string, stateEntries []E, planEntries []E, exhaustive ExhaustiveType, position movement.Position) ([]E, error) {
	stateEntriesByName := o.entriesByName(stateEntries, entryUnknown)
	planEntriesByName := o.entriesByName(planEntries, entryUnknown)
	if len(planEntriesByName) != len(planEntries) {
		return nil, ErrPlanConflict
	}

	findMatchingStateEntry := func(entry E) (E, bool) {
		var found bool
		var foundEntry E
		for _, elt := range stateEntriesByName {
			if o.matcher(entry, elt.Entry) {
				foundEntry = elt.Entry
				found = true
				break
			}
		}

		if !found {
			return *new(E), false
		}

		// If matched entry already exists in the plan, this is not a rename
		// but adding a missing entry.
		if _, ok := planEntriesByName[foundEntry.EntryName()]; ok {
			return *new(E), false
		}

		return foundEntry, true
	}

	renamedEntries := make(map[string]struct{})
	processedStateEntries := make(map[string]uuidObjectWithState[E])
	// First, we compare plan with the state to create a map of all entries with their
	// state.
	for idx, elt := range planEntries {
		var processedEntry uuidObjectWithState[E]

		if stateEntry, found := stateEntriesByName[elt.EntryName()]; !found {
			// If any given entry from the plan is not found in the state, check if there
			// is another entry in the state that matches it, without name and uuid. If so,
			// this is renamed entry: reuse index from the state and set its entryState to
			// entryRenamed.
			if renamedStateEntry, found := findMatchingStateEntry(elt); found {
				processedEntry = uuidObjectWithState[E]{
					Entry:    renamedStateEntry,
					State:    entryRenamed,
					StateIdx: stateEntriesByName[renamedStateEntry.EntryName()].StateIdx,
					NewName:  elt.EntryName(),
				}
				renamedEntries[elt.EntryName()] = struct{}{}
			} else {
				processedEntry = uuidObjectWithState[E]{
					Entry:    elt,
					State:    entryMissing,
					StateIdx: idx,
				}
			}

			processedEntryName := processedEntry.Entry.EntryName()
			// If there is already a processed entry with entryMissing state, it means we've
			// encountered a new entry with the name matching renamedStateEntry old name.
			// It will have entryOutdated state because its spec didn't match spec of the
			// entry about to be renamed.
			// This processed entry is actually a new entry instead, so change its state to
			// entryMissing and update index to match index from the plan.
			if previousEntry, found := processedStateEntries[processedEntryName]; found {
				if previousEntry.State != entryOutdated {
					return nil, &Error{err: ErrInternal, message: fmt.Sprintf("previousEntry.State '%s' != entryOutdated", previousEntry.State)}
				}
			}
			processedStateEntries[processedEntryName] = processedEntry
		} else {
			// If entry from the plan was found in the state check if they match and set the
			// processedEntry State accordingly.
			elt.SetEntryUuid(stateEntry.Entry.EntryUuid())
			processedEntry = uuidObjectWithState[E]{
				Entry:    elt,
				StateIdx: idx,
			}

			if o.matcher(elt, stateEntry.Entry) {
				processedEntry.State = entryOk
			} else {
				processedEntry.State = entryOutdated
			}

			processedStateEntries[elt.EntryName()] = processedEntry
		}
	}

	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && err.Error() != "Object not found" {
		return nil, &Error{err: err, message: "sdk error while listing resources"}
	}

	existing = o.filterEntriesByLocation(location, existing)

	updates := xmlapi.NewChunkedMultiConfig(len(existing), o.batchSize)

	// Next, we compare existing entries from the server against entries processed in the previous
	// state to find any required updates.
	for _, existingEntry := range existing {
		existingEntryName := existingEntry.EntryName()
		path, err := location.XpathWithComponents(o.client.Versioning(), util.AsEntryXpath(existingEntryName))
		if err != nil {
			return nil, &Error{err: err, message: "failed to create xpath for an existing entry"}
		}

		// We don't support "taking over" any existing entries, so if any existing entry was not already
		// part of the plan, and it's about to be either created or renamed, then error out with ErrConflict.
		_, foundInState := stateEntriesByName[existingEntryName]
		_, foundInRenamed := renamedEntries[existingEntryName]
		_, foundInPlan := planEntriesByName[existingEntryName]

		if !foundInState && (foundInRenamed || foundInPlan) {
			return nil, &Error{err: ErrConflict, message: fmt.Sprintf("entry '%s' already exists, created outside of terraform", existingEntryName)}
		}

		// If the existing entry name matches new name for the renamed entry,
		// we delete it before adding rename commands.
		if _, found := renamedEntries[existingEntryName]; found {
			updates.Add(&xmlapi.Config{
				Action: "delete",
				Xpath:  util.AsXpath(path),
				Target: o.client.GetTarget(),
			})
			continue
		}

		processedEntry, found := processedStateEntries[existingEntryName]
		if !found {
			if exhaustive == Exhaustive {
				// If existing entry is not found in the processedStateEntries map,
				// entry is not managed by terraform. If Exhaustive update has been
				// called, delete it from the server.
				updates.Add(&xmlapi.Config{
					Action: "delete",
					Xpath:  util.AsXpath(path),
					Target: o.client.GetTarget(),
				})
			}
			continue
		}

		existingEntryUuid := existingEntry.EntryUuid()
		if existingEntryUuid == nil {
			return nil, &Error{err: ErrInternal, message: "existing entry uuid is nil"}
		}
		processedEntryUuid := processedEntry.Entry.EntryUuid()

		if found && processedEntryUuid != nil && *processedEntryUuid == *existingEntryUuid {
			// If uuid match but the processedEntry is being renamed, don't compare them
			if processedEntry.State == entryRenamed {
				continue
			}

			if o.matcher(processedEntry.Entry, existingEntry) {
				processedEntry.State = entryOutdated
			} else {
				processedEntry.State = entryOk
			}
		}
	}

	// Finally, we check if any state entries are not in the plan, and mark them for deletion.
	for name, elt := range stateEntriesByName {
		if _, processedEntryFound := processedStateEntries[name]; !processedEntryFound {
			elt.State = entryDeleted
			processedStateEntries[name] = elt
		}
	}

	createOps := make([]*xmlapi.Config, len(planEntries))

	for _, elt := range processedStateEntries {
		components := append(parentComponents, util.AsEntryXpath(elt.Entry.EntryName()))
		path, err := location.XpathWithComponents(o.client.Versioning(), components...)
		if err != nil {
			return nil, &Error{err: err, message: "failed to create xpath for an existing entry"}
		}

		xmlEntry, err := o.specifier(elt.Entry)
		if err != nil {
			return nil, &Error{err: err, message: "failed to transform entry into xml document"}
		}

		switch elt.State {
		case entryMissing:
			createOps[elt.StateIdx] = &xmlapi.Config{
				Action:  "edit",
				Xpath:   util.AsXpath(path),
				Element: xmlEntry,
				Target:  o.client.GetTarget(),
			}
		case entryOutdated:
			updates.Add(&xmlapi.Config{
				Action:  "edit",
				Xpath:   util.AsXpath(path),
				Element: xmlEntry,
				Target:  o.client.GetTarget(),
			})
		case entryRenamed:
			updates.Add(&xmlapi.Config{
				Action:  "rename",
				Xpath:   util.AsXpath(path),
				NewName: elt.NewName,
				Target:  o.client.GetTarget(),
			})

			// If existing entry is found in our plan with state entryRenamed,
			// we move entry in processedStateEntries from old name to the new name,
			// indicating it has been renamed. This is used later to properly assign
			// uuids to all entries that are being saved to the state.
			delete(processedStateEntries, elt.Entry.EntryName())
			elt.Entry.SetEntryName(elt.NewName)
			processedStateEntries[elt.NewName] = elt
		case entryDeleted:
			updates.Add(&xmlapi.Config{
				Action: "delete",
				Xpath:  util.AsXpath(path),
				Target: o.client.GetTarget(),
			})
		case entryUnknown:
			return nil, &Error{err: ErrInternal, message: "some entries were not processed"}
		case entryOk:
			// Nothing to do for entries that have no changes
		}
	}

	for _, elt := range createOps {
		if elt != nil {
			updates.Add(elt)
		}
	}

	if len(updates.Operations) > 0 {
		if _, err := o.client.ChunkedMultiConfig(ctx, updates, false, nil); err != nil {
			cleanupErr := o.cleanUpIncompleteUpdate(ctx, location, parentComponents, processedStateEntries, exhaustive)
			return nil, errors.Join(&Error{err: err, message: "failed to execute MultiConfig command"}, cleanupErr)
		}
	}

	switch exhaustive {
	case Exhaustive:
		err := o.moveExhaustive(ctx, location, processedStateEntries, position)
		if err != nil {
			return nil, err
		}
	case NonExhaustive:
		err := o.moveNonExhaustive(ctx, location, planEntries, planEntriesByName, position)
		if err != nil {
			return nil, err
		}
	}

	existing, err = o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, fmt.Errorf("Failed to list remote entries: %w", err)
	}

	existing = o.filterEntriesByLocation(location, existing)

	entries := make([]E, len(planEntries))
	for _, elt := range existing {
		if planEntry, found := planEntriesByName[elt.EntryName()]; found {
			entries[planEntry.StateIdx] = elt
		}
	}

	return entries, nil
}

func (o *UuidObjectManager[E, L, S]) ReadMany(ctx context.Context, location L, stateEntries []E, exhaustive ExhaustiveType, position movement.Position) ([]E, bool, error) {
	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			return nil, false, ErrObjectNotFound
		}
		return nil, false, &Error{err: err, message: "failed to list remote entries"}
	}

	existing = o.filterEntriesByLocation(location, existing)

	if exhaustive == Exhaustive {
		// For resources that take sole ownership of a given list, Read()
		// will return all existing entries from the server.
		return existing, false, nil
	}

	// For resources that only manage a subset of items, Read() must
	// only return entries that are already part of the state.
	stateEntriesByName := make(map[string]uuidObjectWithState[E], len(stateEntries))
	for idx, elt := range stateEntries {
		stateEntriesByName[elt.EntryName()] = uuidObjectWithState[E]{
			Entry:    elt,
			State:    entryMissing,
			StateIdx: idx,
		}
	}

	commonCount := 0
	for _, elt := range existing {
		if stateEntry, found := stateEntriesByName[elt.EntryName()]; found {
			stateEntry.State = entryOk
			stateEntry.Entry = elt
			stateEntry.StateIdx = commonCount
			stateEntriesByName[elt.EntryName()] = stateEntry
			commonCount += 1
		}
	}

	common := make([]E, commonCount)
	var stateEntriesMissing bool
	for _, elt := range stateEntriesByName {
		if elt.State == entryOk {
			common[elt.StateIdx] = elt.Entry
		}

		if elt.State == entryMissing {
			stateEntriesMissing = true
		}
	}

	if stateEntriesMissing {
		return common, false, nil
	}

	// If position is nil, there is a lifecycle { ignore_changes = [position] } in place
	// and the state now keep null position instead of an actual position attribute.
	//
	// In that case we return movementRequired = true, but this will be ignored by terraform
	// anyway until lifecycly ignore_changes is updated to remove position. When that happens
	// terraform will detect a drift (actual position -> null) and generate a valid plan.
	if position == nil {
		return common, true, nil
	}

	actions, err := movement.MoveGroup(position, stateEntries, existing)
	if err != nil {
		if errors.Is(err, movement.ErrSlicesNotEqualLength) {
			return nil, false, fmt.Errorf("Not all entries found on the server: %w", err)
		}
		return nil, false, err
	}

	if len(actions) > 0 {
		return common, true, nil
	}

	return common, false, nil
}

func (o *UuidObjectManager[E, L, S]) Delete(ctx context.Context, location L, parentComponents []string, names []string, exhaustive ExhaustiveType) error {
	deletes := xmlapi.NewChunkedMultiConfig(o.batchSize, len(names))

	for _, elt := range names {
		components := append(parentComponents, util.AsEntryXpath(elt))
		xpath, err := location.XpathWithComponents(o.client.Versioning(), components...)
		if err != nil {
			return err
		}

		deletes.Add(&xmlapi.Config{
			Action: "delete",
			Xpath:  util.AsXpath(xpath),
			Target: o.client.GetTarget(),
		})
	}

	_, _, _, err := o.client.MultiConfig(ctx, deletes, false, nil)
	if err != nil {
		return &Error{err: err, message: "sdk error while deleting"}
	}

	return nil
}

func (o *UuidObjectManager[E, L, S]) cleanUpIncompleteUpdate(ctx context.Context, location L, parentComponents []string, entries map[string]uuidObjectWithState[E], exhaustive ExhaustiveType) error {
	var names []string
	for _, elt := range entries {
		if elt.State == entryMissing {
			names = append(names, elt.Entry.EntryName())
		}
	}

	return o.Delete(ctx, location, parentComponents, names, exhaustive)
}
