package manager

import (
	"context"
	"errors"
	"fmt"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
)

type SDKImportableEntryService[E EntryObject, L EntryLocation, IL ImportLocation] interface {
	Create(context.Context, L, []IL, E) (E, error)
	CreateWithXpath(context.Context, string, E) error
	Read(context.Context, L, string, string) (E, error)
	ReadWithXpath(context.Context, string, string) (E, error)
	List(context.Context, L, string, string, string) ([]E, error)
	Update(context.Context, L, E, string) (E, error)
	Delete(context.Context, L, []IL, ...string) error
	ImportToLocations(context.Context, L, []IL, string) error
	UnimportFromLocations(context.Context, L, []IL, []string) error
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

func (o *ImportableEntryObjectManager[E, L, IL, IS]) Read(ctx context.Context, location L, components []string) (E, error) {
	xpath, err := location.XpathWithComponents(o.client.Versioning(), components...)
	if err != nil {
		return *new(E), err
	}

	object, err := o.service.ReadWithXpath(ctx, util.AsXpath(xpath), "get")
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			return *new(E), ErrObjectNotFound
		}
		return *new(E), &Error{err: err}
	}

	return object, nil
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) Create(ctx context.Context, location L, components []string, entry E) (E, error) {
	_, err := o.Read(ctx, location, components)
	if err == nil {
		return *new(E), &Error{err: ErrConflict, message: fmt.Sprintf("entry '%s' already exists", entry.EntryName())}
	}

	if err != nil && !errors.Is(err, ErrObjectNotFound) {
		return *new(E), err
	}

	xpath, err := location.XpathWithComponents(o.client.Versioning(), components...)
	if err != nil {
		return *new(E), err
	}

	err = o.service.CreateWithXpath(ctx, util.AsXpath(xpath[:len(xpath)-1]), entry)
	if err != nil {
		return *new(E), err
	}

	return o.Read(ctx, location, components)
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

func (o *ImportableEntryObjectManager[E, L, IL, IS]) ImportToLocations(ctx context.Context, location L, importLocs []IL, entry string) error {
	return o.service.ImportToLocations(ctx, location, importLocs, entry)
}

func (o *ImportableEntryObjectManager[E, L, IL, IS]) UnimportFromLocations(ctx context.Context, location L, importLocs []IL, entry string) error {
	return o.service.UnimportFromLocations(ctx, location, importLocs, []string{entry})
}
