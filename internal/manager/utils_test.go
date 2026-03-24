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

type MockLocation struct {
	Filter string
}

func (o MockLocation) IsValid() error {
	panic("unimplemented")
}

func (o MockLocation) XpathWithComponents(version version.Number, components ...string) ([]string, error) {
	return []string{"some", "location", components[0]}, nil
}

func (o MockLocation) XpathWithUuid(version version.Number, uuid string) ([]string, error) {
	panic("unimplemented")
}

func (o MockLocation) Xpath(version version.Number) ([]string, error) {
	panic("unimplemented")
}

func (o MockLocation) LocationFilter() *string {
	if o.Filter == "" {
		return nil
	}

	return &o.Filter
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
	MultiConfigOperMove   MultiConfigOperType = "move"
)

type MultiConfigOper struct {
	Operation   MultiConfigOperType
	EntryName   string
	EntryUuid   *string
	Where       string
	Destination string
	NewName     string
	Value       string
}

func entryNameFromXpath(xpath string) string {
	xpathPrefix := "/some/location/"

	return xpath[len(xpathPrefix):]
}

func findElementByXpath[E sdkmanager.EntryObject](existing *list.List, xpath string) *list.Element {
	needleEntryName := entryNameFromXpath(xpath)
	var next *list.Element
	for e := existing.Front(); e != nil; e = next {
		next = e.Next()
		entry := e.Value.(E)
		if entry.EntryName() == needleEntryName {
			return e
		}
	}

	return nil
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
			State: entryInitial,
			Idx:   idx,
		}
		idx += 1
	}
	existingLen := idx

	fixIndices := func(pivot int) {
		for name, elt := range entriesByName {
			if elt.Idx > pivot {
				elt.Idx -= 1
				entriesByName[name] = elt
			}
		}
	}

	entries := make([]*entryWithIdx, len(updates.Operations)+existing.Len())
	for _, oper := range updates.Operations {
		xpathParts := strings.Split(oper.Xpath, "/")
		entryName := xpathParts[len(xpathParts)-1]
		entryName = strings.TrimPrefix(entryName, "entry[@name='")
		entryName = strings.TrimSuffix(entryName, "']")
		op := oper.XMLName.Local

		operEntry := MultiConfigOper{
			Operation:   MultiConfigOperType(op),
			EntryName:   entryName,
			Destination: oper.Destination,
			Where:       oper.Where,
			NewName:     oper.NewName,
		}

		slog.Debug("MultiConfig", "operEntry", operEntry)

		switch MultiConfigOperType(op) {
		case MultiConfigOperSet, MultiConfigOperEdit:
			slog.Debug("MultiConfig() OperSet/OperEdit")
			entry := oper.Data.(E)

			if existing, found := entriesByName[entryName]; found {
				if objectType == multiConfigUuid {
					entry.SetEntryUuid(existing.Entry.EntryUuid())
				}
				existing.Entry = entry
				entriesByName[entryName] = existing
			} else {
				entryUuid := fmt.Sprintf("%05d", uuid)
				if objectType == multiConfigUuid {
					entry.SetEntryUuid(&entryUuid)
				}
				entriesByName[entryName] = entryWithIdx{
					Entry: entry,
					State: entryUpdated,
					Idx:   idx,
				}

				uuid += 1
				idx += 1
			}
		case MultiConfigOperRename:
			_, found := entriesByName[oper.NewName]
			if found {
				panic(fmt.Sprintf("FIXME: should propagate back error from MultiConfig"))
			}

			entry, found := entriesByName[entryName]
			if !found {
				panic(fmt.Sprintf("FIXME: should propagate back error from MultiConfig"))
			}

			delete(entriesByName, entryName)
			entry.Entry.SetEntryName(oper.NewName)
			entry.State = entryUpdated
			entriesByName[oper.NewName] = entry

			operEntry.NewName = oper.NewName
		case MultiConfigOperDelete:
			entry, found := entriesByName[entryName]
			if !found {
				continue
			}

			fixIndices(entry.Idx)
			delete(entriesByName, entryName)
			idx -= 1
		case MultiConfigOperMove:
			pivot, found := entriesByName[oper.Destination]
			if !found && oper.Destination != "top" && oper.Destination != "bottom" {
				panic(fmt.Sprintf("could not find pivot element for move action: xpath: %s, where: %s, destination: %s", oper.Xpath, oper.Where, oper.Destination))
			}

			moved, found := entriesByName[entryNameFromXpath(oper.Xpath)]
			if !found {
				panic(fmt.Sprintf("could not find moved element for move action: xpath: %s, where: %s, destination: %s", oper.Xpath, oper.Where, oper.Destination))
			}

			switch oper.Where {
			case "after":
				movedOldIdx := moved.Idx
				moved.Idx = pivot.Idx + 1
				if movedOldIdx < pivot.Idx {
					for _, elt := range entriesByName {
						if elt.Idx > movedOldIdx && elt.Idx < pivot.Idx {
							elt.Idx -= 1
						}
						entriesByName[elt.Entry.EntryName()] = elt
					}
				} else if movedOldIdx > pivot.Idx {
					for _, elt := range entriesByName {
						if elt.Idx >= moved.Idx {
							elt.Idx += 1
						}
						entriesByName[elt.Entry.EntryName()] = elt
					}
				}

				entriesByName[moved.Entry.EntryName()] = moved
			case "before":
				movedOldIdx := moved.Idx
				moved.Idx = pivot.Idx
				if movedOldIdx < pivot.Idx {
					for _, elt := range entriesByName {
						if elt.Idx < movedOldIdx && elt.Idx > pivot.Idx {
							elt.Idx -= 1
						}
						entriesByName[elt.Entry.EntryName()] = elt
					}
				} else if movedOldIdx > pivot.Idx {
					for _, elt := range entriesByName {
						if elt.Idx >= moved.Idx {
							elt.Idx += 1
						}
						entriesByName[elt.Entry.EntryName()] = elt
					}
				}

				entriesByName[moved.Entry.EntryName()] = moved
			case "bottom":
				movedOldIdx := moved.Idx
				moved.Idx = existingLen - 1
				for _, elt := range entriesByName {
					if elt.Idx > movedOldIdx {
						elt.Idx -= 1
					}
					entriesByName[elt.Entry.EntryName()] = elt
				}

				entriesByName[moved.Entry.EntryName()] = moved
			case "top":
				movedOldIdx := moved.Idx
				moved.Idx = 0
				for _, elt := range entriesByName {
					if elt.Idx < movedOldIdx && elt.Entry.EntryName() != moved.Entry.EntryName() {
						elt.Idx += 1
					}
					entriesByName[elt.Entry.EntryName()] = elt
				}

				entriesByName[moved.Entry.EntryName()] = moved
			default:
				panic(fmt.Sprintf("Unknown move where: %s", oper.Where))
			}

			operEntry.Where = oper.Where
			operEntry.Destination = oper.Destination
		default:
			panic(fmt.Sprintf("UNKNOWN OPERATION: %s", op))
		}

		opers = append(opers, operEntry)
	}

	for idx := range entries {
		entries[idx] = nil
	}

	for _, elt := range entriesByName {
		if elt.State == entryDeleted {
			continue
		}

		if entries[elt.Idx] != nil {
			var formattedEntries []string
			idx := 1
			for _, elt := range entriesByName {
				formattedEntries = append(
					formattedEntries,
					fmt.Sprintf("%d:{Entry:%s State:%s Idx:%d}", idx, elt.Entry.EntryName(), elt.State, elt.Idx),
				)
				idx += 1
			}
			formattedString := fmt.Sprintf("map[%s]", strings.Join(formattedEntries, " "))
			slog.Debug("Seen elements with duplicated indices", "entries", formattedString)
			panic(fmt.Sprintf("element with idx %d already seen, problem with movement logic", elt.Idx))
		}
		entries[elt.Idx] = &elt
	}

	transformed := list.New()
	for _, elt := range entries {
		if elt != nil && elt.State != entryDeleted {
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
		if typed, ok := o.expected.([]*MockEntryObject); !ok {
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
				if elt.Location != "" && elt.Location != entries[idx].Location {
					return false, nil
				}
			}
		}
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
				if elt.Location != entries[idx].Location {
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
