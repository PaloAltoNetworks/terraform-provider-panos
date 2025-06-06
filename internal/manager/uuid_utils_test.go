package manager_test

import (
	"container/list"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/PaloAltoNetworks/pango/movement"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

var _ = slog.LevelDebug

type MockUuidObject struct {
	Name  string
	Uuid  *string
	Value string
}

func (o MockUuidObject) EntryUuid() *string {
	return o.Uuid
}

func (o *MockUuidObject) SetEntryUuid(value *string) {
	o.Uuid = value
}

func (o *MockUuidObject) SetEntryUuidValue(value string) {
	o.Uuid = &value
}

func (o MockUuidObject) EntryName() string {
	return o.Name
}

func (o *MockUuidObject) SetEntryName(name string) {
	o.Name = name
}

func (o *MockUuidObject) DeepCopy() any {
	return &MockUuidObject{
		Name:  o.Name,
		Uuid:  o.Uuid,
		Value: o.Value,
	}
}

type MockUuidClient[E manager.UuidObject] struct {
	Uuid             int
	Initial          *list.List
	Current          *list.List
	MultiConfigOpers []MultiConfigOper
}

func NewMockUuidClient[E manager.UuidObject](initial []E) *MockUuidClient[E] {
	l := list.New()
	c := list.New()
	uuid := 1
	for _, elt := range initial {
		entry := interface{}(elt).(*MockUuidObject)
		entry.SetEntryUuidValue(fmt.Sprintf("%05d", uuid))
		uuid += 1
		l.PushBack(entry)
		c.PushBack(entry.DeepCopy())
	}

	return &MockUuidClient[E]{
		Uuid:    uuid,
		Initial: l,
		Current: c,
	}
}

func (o *MockUuidClient[E]) GetTarget() string {
	return "mocked-target"
}

func (o *MockUuidClient[E]) Versioning() version.Number {
	v, _ := version.New("10.0.0")
	return v
}

func (o *MockUuidClient[E]) ChunkedMultiConfig(ctx context.Context, updates *xmlapi.MultiConfig, strict bool, extras url.Values) ([]xmlapi.ChunkedMultiConfigResponse, error) {
	data, httpResponse, mcResponse, err := o.MultiConfig(ctx, updates, strict, extras)
	if err != nil {
		return nil, err
	}

	return []xmlapi.ChunkedMultiConfigResponse{{Data: data, HttpResponse: httpResponse, MultiConfigResponse: mcResponse}}, nil
}

func (o *MockUuidClient[E]) MultiConfig(ctx context.Context, updates *xmlapi.MultiConfig, arg1 bool, arg2 url.Values) ([]byte, *http.Response, *xmlapi.MultiConfigResponse, error) {
	o.MultiConfigOpers, o.Uuid = MultiConfig[E](updates, &o.Current, multiConfigUuid, o.Uuid)

	return nil, nil, nil, nil
}

func (o *MockUuidClient[E]) list() []E {
	var entries []E
	for e := o.Current.Front(); e != nil; e = e.Next() {
		entries = append(entries, e.Value.(E))
	}
	return entries
}

type MockUuidService[E manager.UuidObject, L any] struct {
	// used to verify the correctness of MoveGroup() calls
	moveGroupEntries []*MockUuidObject

	client *MockUuidClient[E]
}

func MockUuidSpecifier(entry *MockUuidObject) (any, error) {
	return &MockUuidObject{
		Name:  entry.Name,
		Uuid:  entry.Uuid,
		Value: entry.Value,
	}, nil
}

func MockUuidMatcher(entry *MockUuidObject, other *MockUuidObject) bool {
	return entry.Value == other.Value
}

func NewMockUuidService[E manager.UuidObject, T any](client *MockUuidClient[E]) *MockUuidService[E, T] {
	return &MockUuidService[E, T]{
		client: client,
	}
}

func (o *MockUuidService[E, T]) List(ctx context.Context, location MockLocation, action string, filter string, quote string) ([]E, error) {
	return o.client.list(), nil
}

func (o *MockUuidService[E, T]) Create(ctx context.Context, location MockLocation, entry *MockUuidObject) (*MockUuidObject, error) {
	panic("unreachable")
}

func (o *MockUuidService[E, T]) Update(ctx context.Context, location MockLocation, entry *MockUuidObject, name string) (*MockUuidObject, error) {
	panic("unreachable")
}

func (o *MockUuidService[E, T]) Delete(ctx context.Context, location MockLocation, names ...string) error {
	namesMap := make(map[string]struct{}, len(names))
	for _, elt := range names {
		namesMap[elt] = struct{}{}
	}

	var next *list.Element
	for e := o.client.Initial.Front(); e != nil; e = next {
		next = e.Next()
		entry := e.Value.(E)
		if _, found := namesMap[entry.EntryName()]; found {
			o.client.Initial.Remove(e)
		}
	}

	return nil
}

func (o *MockUuidService[E, L]) removeEntriesFromCurrent(entries []*MockUuidObject) int {
	entriesByName := make(map[string]*MockUuidObject)
	for _, elt := range entries {
		entriesByName[elt.EntryName()] = elt
	}

	firstIdx := -1
	idx := 0
	var next *list.Element
	for e := o.client.Current.Front(); e != nil; e = next {
		next = e.Next()
		entry := e.Value.(E)
		if _, found := entriesByName[entry.EntryName()]; found {
			if firstIdx == -1 {
				firstIdx = idx
			}
			entriesByName[entry.EntryName()].SetEntryUuid(entry.EntryUuid())
			o.client.Current.Remove(e)
		}
		idx += 1
	}

	return firstIdx
}

func (o *MockUuidService[E, T]) MoveGroup(ctx context.Context, location MockLocation, position movement.Position, entries []*MockUuidObject, batchSize int) error {
	o.moveGroupEntries = entries

	if len(entries) == 0 {
		return nil
	}

	var existing []*MockUuidObject
	var next *list.Element
	for e := o.client.Current.Front(); e != nil; e = next {
		next = e.Next()
		entry := e.Value.(*MockUuidObject)
		existing = append(existing, entry)
	}

	movements, err := movement.MoveGroup(position, entries, existing)
	if err != nil {
		return err
	}

	updates := xmlapi.NewMultiConfig(len(movements))
	for _, elt := range movements {
		path, err := location.XpathWithComponents(o.client.Versioning(), elt.Movable.EntryName())
		if err != nil {
			return err
		}

		switch elt.Where {
		case movement.ActionWhereFirst, movement.ActionWhereLast:
			updates.Add(&xmlapi.Config{
				Action:      "move",
				Xpath:       util.AsXpath(path),
				Where:       string(elt.Where),
				Destination: string(elt.Where),
				Target:      o.client.GetTarget(),
			})
		case movement.ActionWhereBefore, movement.ActionWhereAfter:
			updates.Add(&xmlapi.Config{
				Action:      "move",
				Xpath:       util.AsXpath(path),
				Where:       string(elt.Where),
				Destination: elt.Destination.EntryName(),
				Target:      o.client.GetTarget(),
			})
		}
	}

	if len(updates.Operations) > 0 {
		_, _, _, err = o.client.MultiConfig(ctx, updates, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
