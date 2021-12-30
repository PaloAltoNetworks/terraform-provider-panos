package namespace

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

/*
Plugin is a namespace struct for config that exists in PAN-OS as a plugin.
*/
type Plugin struct {
	Common
}

// Set performs a SET to configure one or more objects.
func (n *Plugin) Set(pather Pather, specs []PluginSpecifier) error {
	data := make([]interface{}, 0, len(specs))
	names := make([]string, 0, len(specs))

	tally := make(map[string]int)
	for _, s := range specs {
		name, val, err := s.Specify(n.Client.Plugins())
		if err != nil {
			return err
		}
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
func (n *Plugin) Edit(pather Pather, spec PluginSpecifier) error {
	name, data, err := spec.Specify(n.Client.Plugins())
	if err != nil {
		return err
	}

	if n.Plural != "" {
		n.Client.LogAction("(edit) %s: %s", n.Singular, name)
	} else {
		n.Client.LogAction("(edit) %s", n.Singular)
	}

	path, pErr := pather([]string{name})
	if pErr != nil {
		return pErr
	}

	_, err = n.Client.Edit(path, data, nil, nil)
	return err
}

// Delete performs a DELETE to remove config.
func (n *Plugin) Delete(pather Pather, names []string, nErr error) error {
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
