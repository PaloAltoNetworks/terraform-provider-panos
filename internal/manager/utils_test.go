package manager_test

import (
	"container/list"
	"fmt"
	"log/slog"
	"strings"

	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

var _ = slog.Debug

type MockLocation struct{}

func (o MockLocation) IsValid() error {
	panic("unimplemented")
}

func (o MockLocation) XpathWithEntryName(version version.Number, name string) ([]string, error) {
	return []string{"some", "location", name}, nil
}

func (o MockLocation) XpathWithUuid(version version.Number, uuid string) ([]string, error) {
	panic("unimplemented")
}

func (o MockLocation) Xpath(version version.Number) ([]string, error) {
	panic("unimplemented")
}

type multiConfigType string

const (
	multiConfigEntry multiConfigType = "entry"
	multiConfigUuid  multiConfigType = "uuid"
)

type MultiConfigOperType string

const (
	MultiConfigOperSet    MultiConfigOperType = "set"
	MultiConfigOperEdit   MultiConfigOperType = "edit"
	MultiConfigOperAdd    MultiConfigOperType = "add"
	MultiConfigOperRemove MultiConfigOperType = "remove"
	MultiConfigOperModify MultiConfigOperType = "modify"
	MultiConfigOperDelete MultiConfigOperType = "delete"
	MultiConfigOperRename MultiConfigOperType = "rename"
)

type MultiConfigOper struct {
	Operation MultiConfigOperType
	EntryName string
	EntryUuid *string
}

func MultiConfig[E sdkmanager.UuidObject](updates *xmlapi.MultiConfig, existingPtr **list.List, objectType multiConfigType, uuid int) ([]MultiConfigOper, int) {
	type entryWithIdx struct {
		Entry E
		State entryState
		Idx   int
	}

	existing := *existingPtr
	var opers []MultiConfigOper

	entriesByName := make(map[string]entryWithIdx, existing.Len())
	idx := 0
	for e := existing.Front(); e != nil; e = e.Next() {
		entry := e.Value.(E)
		entriesByName[entry.EntryName()] = entryWithIdx{
			Entry: entry,
			State: entryOk,
			Idx:   idx,
		}
		idx += 1
	}

	fixIndices := func(pivot int) {
		for name, elt := range entriesByName {
			if elt.Idx > pivot {
				elt.Idx -= 1
				entriesByName[name] = elt
			}
		}
	}

	for _, oper := range updates.Operations {
		xpathParts := strings.Split(oper.Xpath, "/")
		entryName := xpathParts[len(xpathParts)-1]
		op := oper.XMLName.Local

		operEntry := MultiConfigOper{
			Operation: MultiConfigOperType(op),
			EntryName: entryName,
		}

		switch MultiConfigOperType(op) {
		case MultiConfigOperSet, MultiConfigOperEdit:
			entry := oper.Data.(E)

			if existing, found := entriesByName[entryName]; found {
				if objectType == multiConfigUuid {
					entry.SetEntryUuid(existing.Entry.EntryUuid())
				}
				existing.Entry = entry
			} else {
				entryUuid := fmt.Sprintf("%05d", uuid)
				if objectType == multiConfigUuid {
					entry.SetEntryUuid(&entryUuid)
				}
				entriesByName[entryName] = entryWithIdx{
					Entry: entry,
					State: entryOk,
					Idx:   idx,
				}

				uuid += 1
				idx += 1
			}
		case MultiConfigOperDelete:
			fixIndices(entriesByName[entryName].Idx)
			delete(entriesByName, entryName)
			idx -= 1
		default:
			panic(fmt.Sprintf("UNKNOWN OPERATION: %s", op))
		}

		opers = append(opers, operEntry)

	}

	entries := make([]entryWithIdx, len(updates.Operations)+existing.Len())
	for _, elt := range entriesByName {
		if elt.State == entryOk {
			entries[elt.Idx] = elt
		}
	}

	transformed := list.New()
	for _, elt := range entries {
		if elt.State == entryOk {
			transformed.PushBack(elt.Entry)
		}
	}

	*existingPtr = transformed

	return opers, uuid
}

type Copyable interface {
	DeepCopy() any
}

type equalEntries struct {
	expected any
}

func MatchEntries(expected any) types.GomegaMatcher {
	return &equalEntries{
		expected: expected,
	}
}

func (o *equalEntries) Match(actual any) (bool, error) {
	switch entries := actual.(type) {
	case []*MockEntryObject:
		panic("unimplemented 1")
	case []*MockUuidObject:
		if typed, ok := o.expected.([]*MockUuidObject); !ok {
			return false, fmt.Errorf("Expected %T to match %T", o.expected, actual)
		} else {
			if len(entries) != len(typed) {
				return false, nil
			}
			for idx, elt := range typed {
				if elt.Name != entries[idx].Name {
					return false, nil
				}
				if elt.Value != entries[idx].Value {
					return false, nil
				}
			}
		}
	default:
		return false, fmt.Errorf("invalid type: %T", entries)
	}

	return true, nil
}

func (o *equalEntries) FailureMessage(actual any) string {
	return format.Message(actual, "to equal", o.expected)
}

func (o *equalEntries) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to equal", o.expected)
}
