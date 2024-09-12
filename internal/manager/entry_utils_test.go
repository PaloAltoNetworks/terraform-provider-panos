package manager_test

import (
	"container/list"
	"context"
	"net/http"
	"net/url"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

type MockEntryObject struct {
	Name  string
	Value string
}

func (o *MockEntryObject) EntryName() string {
	return o.Name
}

func (o *MockEntryObject) SetEntryName(name string) {
	o.Name = name
}

func (o *MockEntryObject) EntryUuid() *string {
	panic("mock entry object EntryUuid() called")
}

func (o *MockEntryObject) SetEntryUuid(uuid *string) {
	panic("mock entry object SetEntryUuid() called")
}

func (o *MockEntryObject) DeepCopy() any {
	return &MockUuidObject{
		Name:  o.Name,
		Value: o.Value,
	}
}

type MockEntryClient[E manager.UuidObject] struct {
	Initial          *list.List
	MultiConfigOpers []MultiConfigOper
}

func NewMockEntryClient[E manager.UuidObject](initial []E) *MockEntryClient[E] {
	l := list.New()
	for _, elt := range initial {
		l.PushBack(elt)
	}

	return &MockEntryClient[E]{
		Initial: l,
	}
}

func (o *MockEntryClient[E]) GetTarget() string {
	return "mocked-target"
}

func (o *MockEntryClient[E]) Versioning() version.Number {
	v, _ := version.New("10.0.0")
	return v
}

type entryState string

const (
	entryDeleted entryState = "delete"
	entryOk      entryState = "ok"
)

func (o *MockEntryClient[E]) MultiConfig(ctx context.Context, updates *xmlapi.MultiConfig, arg1 bool, arg2 url.Values) ([]byte, *http.Response, *xmlapi.MultiConfigResponse, error) {
	o.MultiConfigOpers, _ = MultiConfig[E](updates, &o.Initial, multiConfigEntry, 0)

	return nil, nil, nil, nil
}

func (o *MockEntryClient[E]) list() []E {
	var entries []E
	for e := o.Initial.Front(); e != nil; e = e.Next() {
		entries = append(entries, e.Value.(E))
	}

	return entries
}

type MockEntryService[E manager.UuidObject, L manager.EntryLocation] struct {
	client *MockEntryClient[E]
}

func (o *MockEntryService[E, L]) Create(ctx context.Context, location L, entry E) (E, error) {
	o.client.Initial.PushBack(entry)

	return entry, nil
}

func (o *MockEntryService[E, L]) Update(ctx context.Context, location L, entry E, name string) (E, error) {
	for e := o.client.Initial.Front(); e != nil; e = e.Next() {
		eltEntry := e.Value.(E)
		if entry.EntryName() == eltEntry.EntryName() {
			e.Value = entry
			return entry, nil
		}
	}

	return *new(E), sdkerrors.Panos{Code: 7, Msg: "Object not found"}
}

func (o *MockEntryService[E, L]) List(ctx context.Context, location L, action string, filter string, quote string) ([]E, error) {
	return o.client.list(), nil
}

func (o *MockEntryService[E, L]) Read(ctx context.Context, location L, name string, action string) (E, error) {
	for e := o.client.Initial.Front(); e != nil; e = e.Next() {
		entry := e.Value.(E)
		if entry.EntryName() == name {
			return entry, nil
		}
	}

	return *new(E), sdkerrors.Panos{Code: 7, Msg: "Object not found"}
}

func (o *MockEntryService[E, L]) Delete(ctx context.Context, location L, names ...string) error {
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

func NewMockEntryService[E manager.UuidObject, L manager.EntryLocation](client *MockEntryClient[E]) *MockEntryService[E, L] {
	return &MockEntryService[E, L]{
		client: client,
	}
}

func MockEntrySpecifier(entry *MockEntryObject) (any, error) {
	return entry, nil
}

func MockEntryMatcher(entry *MockEntryObject, other *MockEntryObject) bool {
	return entry.Value == other.Value
}
