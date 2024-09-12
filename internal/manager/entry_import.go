package manager

import (
	"context"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
)

type SDKImportableEntryService[E EntryObject, L EntryLocation, IL ImportLocation] interface {
	Create(context.Context, L, []IL, E) (E, error)
	Read(context.Context, L, string, string) (E, error)
	List(context.Context, L, string, string, string) ([]E, error)
	Update(context.Context, L, E, string) (E, error)
	Delete(context.Context, L, []IL, ...string) error
}

type ImportableEntryObjectManager[E EntryObject, L EntryLocation, IL ImportLocation, IS SDKImportableEntryService[E, L, IL]] struct {
	service   IS
	client    SDKClient
	specifier func(E) (any, error)
	matcher   func(E, E) bool
}

func NewImportableEntryObjectManager[E EntryObject, L EntryLocation, IL ImportLocation, IS SDKImportableEntryService[E, L, IL]](client SDKClient, service IS, specifier func(E) (any, error), matcher func(E, E) bool) *ImportableEntryObjectManager[E, L, IL, IS] {
	return &ImportableEntryObjectManager[E, L, IL, IS]{
		service:   service,
		client:    client,
		specifier: specifier,
		matcher:   matcher,
	}
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) ReadMany(ctx context.Context, location L, entries []E) ([]E, error) {
	return nil, &Error{err: ErrInternal, message: "called ReadMany on an importable singular resource"}
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) Read(ctx context.Context, location L, name string) (E, error) {
	object, err := o.service.Read(ctx, location, name, "get")
	if err != nil {
		return *new(E), ErrObjectNotFound
	}

	return object, nil
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) Create(ctx context.Context, location L, importLocs []IL, entry E) (E, error) {
	existing, err := o.service.List(ctx, location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		return *new(E), err
	}

	for _, elt := range existing {
		if elt.EntryName() == entry.EntryName() {
			return *new(E), ErrConflict
		}
	}

	obj, err := o.service.Create(ctx, location, importLocs, entry)
	return obj, err
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) Update(ctx context.Context, location L, entry E, name string) (E, error) {
	updated, err := o.service.Update(ctx, location, entry, name)
	if err != nil {
		return *new(E), &Error{err: err, message: "error during Update call"}
	}

	return updated, nil
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) Delete(ctx context.Context, location L, importLocations []IL, names []string, exhaustive ExhaustiveType) error {
	err := o.service.Delete(ctx, location, importLocations, names...)
	if err != nil {
		return &Error{err: err, message: "sdk error while deleting"}
	}
	return nil
}
