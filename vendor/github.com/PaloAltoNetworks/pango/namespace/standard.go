package namespace

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

/*
Standard is a namespace struct for config that is not imported into a vsys.
*/
type Standard struct {
	Common
}

// Set performs a SET to configure one or more objects.
func (n *Standard) Set(pather Pather, specs []Specifier) error {
	v := n.Client.Versioning()
	data := make([]interface{}, 0, len(specs))
	names := make([]string, 0, len(specs))

	tally := make(map[string]int)
	for _, s := range specs {
		name, val := s.Specify(v)
		tally[name] = tally[name] + 1
		if tally[name] > 1 {
			return fmt.Errorf("%s is defined multiple times: %q", n.Singular, name)
		}
		data = append(data, val)
		names = append(names, name)
	}

	path, pErr := pather(names)
	if pErr != nil {
		return pErr
	}

	if n.Plural != "" {
		n.Client.LogAction("(set) %s: %v", n.Plural, names)
	} else {
		n.Client.LogAction("(set) %s", n.Singular)
	}

	if len(data) == 0 {
		return nil
	}

	elm := util.BulkElement{
		XMLName: xml.Name{Local: path[len(path)-2]},
		Data:    data,
	}

	if len(data) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	_, err := n.Client.Set(path, elm.Config(), nil, nil)
	return err
}

// Edit performs an EDIT to modify a single object.
func (n *Standard) Edit(pather Pather, spec Specifier) error {
	name, data := spec.Specify(n.Client.Versioning())

	if n.Plural != "" {
		n.Client.LogAction("(edit) %s: %s", n.Singular, name)
	} else {
		n.Client.LogAction("(edit) %s", n.Singular)
	}

	path, pErr := pather([]string{name})
	if pErr != nil {
		return pErr
	}

	_, err := n.Client.Edit(path, data, nil, nil)
	return err
}

// Delete performs a DELETE to remove config.
func (n *Standard) Delete(pather Pather, names []string, nErr error) error {
	if nErr != nil {
		return nErr
	}

	if n.Plural != "" {
		n.Client.LogAction("(delete) %s: %v", n.Plural, names)
		if len(names) == 0 {
			return nil
		}
	} else {
		n.Client.LogAction("(delete) %s", n.Singular)
	}

	path, pErr := pather(names)
	if pErr != nil {
		return pErr
	}

	_, err := n.Client.Delete(path, nil, nil)
	return err
}

// MoveGroup places a logical group of objects in the desired location (rulebase
// objects).
//
// The `movement` param should be one of the Move constants in the util package.
//
// The `rule` param is the other rule the `movement` param is referencing.  If
// this is an empty string, then the first rule in the group isn't moved
// anywhere, but all other rules will still be moved to be grouped with the
// first one.
func (n *Standard) MoveGroup(pather Pather, lister MoveLister, movement int, rule string, grp []string) error {
	var err error

	fIdx, oIdx := -1, -1

	n.Client.LogAction("(move) %s group", n.Singular)
	if len(grp) == 0 {
		return fmt.Errorf("Requires at least one %s", n.Singular)
	} else if rule == grp[0] {
		return fmt.Errorf("Can't position %q in relation to itself", rule)
	} else if !util.ValidMovement(movement) {
		return fmt.Errorf("Invalid movement specified: %d", movement)
	} else if util.RelativeMovement(movement) && rule == "" {
		return fmt.Errorf("Specify 'rule' in order to perform relative group positioning")
	}

	path, pErr := pather([]string{grp[0]})
	if pErr != nil {
		return pErr
	}

	if movement != util.MoveSkip {
		// Get the list of current security policies.
		curList, err := lister()
		if err != nil {
			return err
		} else if len(curList) == 0 {
			return fmt.Errorf("No policies found")
		}

		switch movement {
		case util.MoveTop:
			_, em := n.Client.Move(path, "top", "", nil, nil)
			if em != nil {
				if em.Error() != "already at the top" {
					err = em
				}
			}
		case util.MoveBottom:
			_, em := n.Client.Move(path, "bottom", "", nil, nil)
			if em != nil {
				if em.Error() != "already at the bottom" {
					err = em
				}
			}
		default:
			// Find the indexes of the first security policy and the ref policy.
			for i, v := range curList {
				if v == grp[0] {
					fIdx = i
				} else if v == rule {
					oIdx = i
				}
				if fIdx != -1 && oIdx != -1 {
					break
				}
			}

			// Sanity check:  both rules should be present.
			if fIdx == -1 {
				return fmt.Errorf("First %s in group %q does not exist", n.Singular, grp[0])
			} else if oIdx == -1 {
				return fmt.Errorf("Reference %s %q does not exist", n.Singular, rule)
			}

			// Perform the move of the first security policy, if needed.
			if (movement == util.MoveBefore && fIdx > oIdx) || (movement == util.MoveDirectlyBefore && fIdx+1 != oIdx) {
				_, err = n.Client.Move(path, "before", rule, nil, nil)
			} else if (movement == util.MoveAfter && fIdx < oIdx) || (movement == util.MoveDirectlyAfter && fIdx != oIdx+1) {
				_, err = n.Client.Move(path, "after", rule, nil, nil)
			}
		}

		// If we moved something, make sure it worked.
		if err != nil {
			return err
		}
	}

	// Now move all other rules under the first.
	li := len(path) - 1
	for i := 1; i < len(grp); i++ {
		path[li] = util.AsEntryXpath([]string{grp[i]})
		_, err = n.Client.Move(path, "after", grp[i-1], nil, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
