package manager_test

import (
	"container/list"
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

type MockEntryObject struct {
	Name  string
	Value string
}

func (o *MockEntryObject) EntryUuid() *string {
	panic("called EntryUuid on MockEntryObject")
}

func (o *MockEntryObject) SetEntryUuid(value *string) {
	panic("called SetEntryUuid on MockEntryObject")
}

func (o *MockEntryObject) EntryName() string {
	return o.Name
}

func (o *MockEntryObject) SetEntryName(name string) {
	o.Name = name
}

func (o *MockEntryObject) DeepCopy() any {
	return &MockEntryObject{
		Name:  o.Name,
		Value: o.Value,
	}
}

type MockEntryClient[E manager.UuidObject] struct {
	Initial          *list.List
	Current          *list.List
	MultiConfigOpers []MultiConfigOper
}

func NewMockEntryClient[E manager.UuidObject](initial []E) *MockEntryClient[E] {
	l := list.New()
	c := list.New()

	for _, elt := range initial {
		entry := interface{}(elt).(*MockEntryObject)
		l.PushBack(entry)
		c.PushBack(entry.DeepCopy())
	}

	return &MockEntryClient[E]{
		Initial: l,
		Current: c,
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

func (o *MockEntryClient[E]) ChunkedMultiConfig(ctx context.Context, updates *xmlapi.MultiConfig, strict bool, extras url.Values) ([]xmlapi.ChunkedMultiConfigResponse, error) {
	data, httpResponse, mcResponse, err := o.MultiConfig(ctx, updates, strict, extras)
	if err != nil {
		return nil, err
	}

	return []xmlapi.ChunkedMultiConfigResponse{{Data: data, HttpResponse: httpResponse, MultiConfigResponse: mcResponse}}, nil
}

func (o *MockEntryClient[E]) MultiConfig(ctx context.Context, updates *xmlapi.MultiConfig, arg1 bool, arg2 url.Values) ([]byte, *http.Response, *xmlapi.MultiConfigResponse, error) {
	o.MultiConfigOpers, _ = MultiConfig[E](updates, &o.Current, multiConfigEntry, 0)

	return nil, nil, nil, nil
}

func (o *MockEntryClient[E]) list() []E {
	var entries []E
	slog.Debug("MockEntryClient list()", "o.Current", o.Current, "o.Current.Front()", o.Current.Front())
	for e := o.Current.Front(); e != nil; e = e.Next() {
		slog.Debug("MockEntryClient list()", "entry", e.Value.(E).EntryName())
		entries = append(entries, e.Value.(E))
	}

	return entries
}

type MockEntryService[E manager.UuidObject, L manager.EntryLocation] struct {
	client *MockEntryClient[E]
}

func (o *MockEntryService[E, L]) CreateWithXpath(ctx context.Context, xpath string, entry E) error {
	_, err := o.Create(ctx, *new(L), entry)
	return err
}

func (o *MockEntryService[E, L]) Create(ctx context.Context, location L, entry E) (E, error) {
	o.client.Initial.PushBack(entry)

	return entry, nil
}

func (o *MockEntryService[E, L]) UpdateWithXpath(ctx context.Context, xpath string, entry E, name string) error {
	for e := o.client.Initial.Front(); e != nil; e = e.Next() {
		eltEntry := e.Value.(E)
		if entry.EntryName() == eltEntry.EntryName() {
			e.Value = entry
			return nil
		}
	}

	return sdkerrors.Panos{Code: 7, Msg: "Object not found"}
}

func (o *MockEntryService[E, L]) ListWithXpath(ctx context.Context, xpath string, action string, filter string, quote string) ([]E, error) {
	return o.client.list(), nil
}

func (o *MockEntryService[E, L]) ReadWithXpath(ctx context.Context, xpath string, action string) (E, error) {
	components := strings.Split(xpath, "/")
	name := components[len(components)-1]
	name = strings.TrimPrefix(name, "entry[@name='")
	name = strings.TrimSuffix(name, "']")
	return o.Read(ctx, *new(L), name, action)
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
