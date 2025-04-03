package manager

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type entryObjectWithState[E EntryObject] struct {
	Entry    E
	State    entryState
	StateIdx int
	NewName  string
}

type TFEntryObject[E any] interface {
	EntryName() string
	CopyToPango(context.Context, *map[string]types.String) (E, diag.Diagnostics)
	CopyFromPango(context.Context, E, *map[string]types.String) diag.Diagnostics
}

type EntryObject interface {
	EntryName() string
	SetEntryName(name string)
}

type EntryLocation interface {
	XpathWithEntryName(version.Number, string) ([]string, error)
}

type SDKEntryService[E EntryObject, L EntryLocation] interface {
	Create(context.Context, L, E) (E, error)
	List(context.Context, L, string, string, string) ([]E, error)
	Read(context.Context, L, string, string) (E, error)
	Update(context.Context, L, E, string) (E, error)
	Delete(context.Context, L, ...string) error
}

type EntryObjectManager[E EntryObject, L EntryLocation, S SDKEntryService[E, L]] struct {
	batchSize int
	service   S
	client    SDKClient
	specifier func(E) (any, error)
	matcher   func(E, E) bool
}

func NewEntryObjectManager[E EntryObject, L EntryLocation, S SDKEntryService[E, L]](client SDKClient, service S, batchSize int, specifier func(E) (any, error), matcher func(E, E) bool) *EntryObjectManager[E, L, S] {
	return &EntryObjectManager[E, L, S]{
		batchSize: batchSize,
		service:   service,
		client:    client,
		specifier: specifier,
		matcher:   matcher,
	}
}

func (o *EntryObjectManager[E, L, S]) ReadMany(ctx context.Context, location L, entries []E) ([]E, error) {
	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			return nil, ErrObjectNotFound
		} else {
			return nil, &Error{err: err, message: "Failed to read entries from the server"}
		}
	}

	entriesByName := o.entriesByName(entries, entryOk)

	filtered := make([]E, 0, len(entries))
	for _, elt := range existing {
		if _, found := entriesByName[elt.EntryName()]; found {
			filtered = append(filtered, elt)
		}
	}

	return filtered, nil
}

func (o *EntryObjectManager[E, L, S]) Read(ctx context.Context, location L, name string) (E, error) {
	object, err := o.service.Read(ctx, location, name, "get")
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			return *new(E), ErrObjectNotFound
		}
		return *new(E), &Error{err: err}
	}

	return object, nil
}

func (o *EntryObjectManager[E, T, S]) Create(ctx context.Context, location T, entry E) (E, error) {
	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return *new(E), err
	}

	for _, elt := range existing {
		if elt.EntryName() == entry.EntryName() {
			return *new(E), ErrConflict
		}
	}

	obj, err := o.service.Create(ctx, location, entry)
	return obj, err
}

func (o *EntryObjectManager[E, L, S]) CreateMany(ctx context.Context, location L, entries []E) ([]E, error) {
	// First, check if none of the entries from the plan already exist on the server
	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, &Error{err: err, message: "failed to list existing entries on the server"}
	}

	entriesByName := o.entriesByName(entries, entryMissing)

	for _, elt := range existing {
		if _, found := entriesByName[elt.EntryName()]; found {
			return nil, &Error{err: ErrConflict, message: fmt.Sprintf("entry '%s' already exists", elt.EntryName())}
		}
	}

	updates := xmlapi.NewChunkedMultiConfig(len(existing), o.batchSize)

	for _, elt := range entries {
		path, err := location.XpathWithEntryName(o.client.Versioning(), elt.EntryName())
		if err != nil {
			return nil, &Error{err: err, message: "failed to create xpath for an existing entry"}
		}

		xmlEntry, err := o.specifier(elt)
		if err != nil {
			return nil, &Error{err: err, message: "failed to marshal entry into XML document"}
		}

		updates.Add(&xmlapi.Config{
			Action:  "edit",
			Xpath:   util.AsXpath(path),
			Element: xmlEntry,
			Target:  o.client.GetTarget(),
		})
	}

	if len(updates.Operations) > 0 {
		if _, err := o.client.ChunkedMultiConfig(ctx, updates, false, nil); err != nil {
			return nil, &Error{err: err, message: "Failed to execute MultiConfig command"}
		}
	}

	existing, err = o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, &Error{err: err, message: "failed to list existing entries on the server"}
	}

	var filtered []E
	for _, elt := range existing {
		if _, found := entriesByName[elt.EntryName()]; !found {
			continue
		}
		filtered = append(filtered, elt)
	}

	return filtered, nil
}

func (o *EntryObjectManager[E, T, S]) Delete(ctx context.Context, location T, names []string) error {
	err := o.service.Delete(ctx, location, names...)
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			return ErrObjectNotFound
		} else {
			return &Error{err: err, message: "sdk error while deleting"}
		}
	}
	return nil
}

func (o *EntryObjectManager[E, L, S]) Update(ctx context.Context, location L, entry E, name string) (E, error) {
	updated, err := o.service.Update(ctx, location, entry, name)
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			return *new(E), ErrObjectNotFound
		} else {
			return *new(E), &Error{err: err, message: "sdk error while updating"}
		}
	}

	return updated, nil
}

func (o *EntryObjectManager[E, L, S]) UpdateMany(ctx context.Context, location L, stateEntries []E, planEntries []E) ([]E, error) {
	stateEntriesByName := o.entriesByName(stateEntries, entryUnknown)
	planEntriesByName := o.entriesByName(planEntries, entryUnknown)
	renamedEntries := make(map[string]struct{})
	processedStateEntriesByName := make(map[string]entryObjectWithState[E])

	findMatchingStateEntry := func(entry E) (*entryObjectWithState[E], bool) {
		var found *entryObjectWithState[E]
		for _, elt := range stateEntriesByName {
			if o.matcher(entry, elt.Entry) {
				found = &elt
			}
		}

		if found == nil {
			return nil, false
		}

		// If matched entry already exists in the plan, found entry
		// is not a rename, but an missing entry about to be added.
		if _, ok := planEntriesByName[found.Entry.EntryName()]; ok {
			return nil, false
		}

		return found, true
	}

	for idx, elt := range planEntries {
		eltEntryName := elt.EntryName()
		var processedEntry *entryObjectWithState[E]

		if stateElt, found := stateEntriesByName[eltEntryName]; !found {
			// If given plan entry is not found in the state, check if there
			// is another entry that matches it without name. If so, this plan
			// entry is a rename: keep the renamedEntry index, adn set its state
			// to entryRename.
			if stateEntry, found := findMatchingStateEntry(elt); found {
				processedEntry = &entryObjectWithState[E]{
					Entry:   stateEntry.Entry,
					State:   entryRenamed,
					NewName: eltEntryName,
				}
				renamedEntries[eltEntryName] = struct{}{}
			} else {
				processedEntry = &entryObjectWithState[E]{
					Entry: elt,
					State: entryMissing,
				}
			}

			// If there is already a processed entry with state entryMissing, it means
			// we've encountered a new entry with the name matching renamedEntry old name.
			// It will have state entryOutdated because its spec didn't match spec of the
			// entry about to be renamed.
			// Change its state to entryMissing instead, and update its index to match
			// index from the plan.
			if previousEntry, found := processedStateEntriesByName[processedEntry.Entry.EntryName()]; found {
				if previousEntry.State != entryOutdated && previousEntry.State != entryMissing {
					return nil, &Error{err: ErrInternal, message: "invalid entry state when building update list"}
				}
			}
			processedStateEntriesByName[processedEntry.Entry.EntryName()] = *processedEntry
		} else {
			processedEntry = &entryObjectWithState[E]{
				Entry:    elt,
				StateIdx: idx,
			}

			if o.matcher(elt, stateElt.Entry) {
				processedEntry.State = entryOk
			} else {
				processedEntry.State = entryOutdated
			}
			processedStateEntriesByName[elt.EntryName()] = *processedEntry
		}
	}

	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, &Error{err: err, message: "failed to get a list of existing entries from the server"}
	}

	for name, elt := range stateEntriesByName {
		if _, processedEntryFound := processedStateEntriesByName[name]; !processedEntryFound {
			elt.State = entryDeleted
			processedStateEntriesByName[name] = elt
		}
	}

	updates := xmlapi.NewChunkedMultiConfig(len(planEntries), o.batchSize)

	for _, existingEntry := range existing {
		existingEntryName := existingEntry.EntryName()
		_, foundInState := stateEntriesByName[existingEntryName]
		_, foundInRenamed := renamedEntries[existingEntryName]
		_, foundInPlan := planEntriesByName[existingEntryName]

		if !foundInState && (foundInRenamed || foundInPlan) {
			return nil, &Error{err: ErrConflict, message: "entry with a matching name already exists"}
		}

		path, err := location.XpathWithEntryName(o.client.Versioning(), existingEntryName)
		if err != nil {
			return nil, &Error{err: err, message: "failed to create xpath for an existing entry"}
		}

		// If the existing entry name matches new name for the renamed entry,
		// we delete it before adding Renamed commands.
		if _, found := renamedEntries[existingEntryName]; found {
			updates.Add(&xmlapi.Config{
				Action: "delete",
				Xpath:  util.AsXpath(path),
				Target: o.client.GetTarget(),
			})
			continue
		}

		processedElt, found := processedStateEntriesByName[existingEntryName]
		if found {
			// XXX: If entry from the plan is in process of being renamed, and its content
			// differs from what exists on the server we should switch its state to entryOutdated
			// instead.
			if processedElt.State == entryRenamed {
				continue
			}

			if !o.matcher(processedElt.Entry, existingEntry) {
				processedElt.State = entryOutdated
			} else {
				processedElt.State = entryOk
			}
		}

	}

	for _, elt := range processedStateEntriesByName {
		path, err := location.XpathWithEntryName(o.client.Versioning(), elt.Entry.EntryName())
		if err != nil {
			return nil, &Error{err: err, message: "failed to create xpath for entry"}
		}

		xmlEntry, err := o.specifier(elt.Entry)
		if err != nil {
			return nil, &Error{err: err, message: "failed to marshal entry into XML document"}
		}

		switch elt.State {
		case entryMissing, entryOutdated:
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
			// we move entry in processedEntries from old name to the new name,
			// indicating it has been renamed.
			// This is used later when we assign uuids to all entries.
			delete(processedStateEntriesByName, elt.Entry.EntryName())
			elt.Entry.SetEntryName(elt.NewName)
			processedStateEntriesByName[elt.NewName] = elt
		case entryDeleted:
			updates.Add(&xmlapi.Config{
				Action: "delete",
				Xpath:  util.AsXpath(path),
				Target: o.client.GetTarget(),
			})
		case entryUnknown:
			slog.Warn("Entry state is still unknown after reconciliation", "Name", elt.Entry.EntryName())
		case entryOk:
			// Nothing to do for entries that have no changes
		}

	}

	if len(updates.Operations) > 0 {
		if _, err := o.client.ChunkedMultiConfig(ctx, updates, false, nil); err != nil {
			return nil, &Error{err: err, message: "Failed to execute MultiConfig command"}
		}
	}

	existing, err = o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return nil, fmt.Errorf("Failed to list remote entries: %w", err)
	}

	entries := make([]E, len(planEntries))
	for _, elt := range existing {
		if planEntry, found := planEntriesByName[elt.EntryName()]; found {
			entries[planEntry.StateIdx] = elt
		}
	}

	return entries, nil
}

func (o *EntryObjectManager[E, L, S]) entriesByName(entries []E, initialState entryState) map[string]entryObjectWithState[E] {
	entriesByName := make(map[string]entryObjectWithState[E], len(entries))
	for idx, elt := range entries {
		entriesByName[elt.EntryName()] = entryObjectWithState[E]{
			Entry:    elt,
			StateIdx: idx,
			State:    initialState,
		}
	}

	return entriesByName
}
