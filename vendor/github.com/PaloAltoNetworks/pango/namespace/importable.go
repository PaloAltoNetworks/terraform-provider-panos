package namespace

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

/*
Importable is a namespace struct for config that is imported into a vsys.

The current list of importable config is as follows:
    * interfaces
    * virtual routers
    * virtual wires
    * vlans

The ImportPath param should be set to any of the valid Import constants in the util
package.
*/
type Importable struct {
	Common
	ImportPath string
}

/*
Set performs a SET to configure one or more objects.

As this is an importable config, first all objects are unimported, then
everything is configured, and finally the config is imported into the
specified vsys.
*/
func (n *Importable) Set(tmpl, ts, vsys string, pather Pather, specs []ImportSpecifier) error {
	var err error

	// Sanity check: Import path should be configured.
	if n.ImportPath == "" {
		return fmt.Errorf("Namespace did not configure 'ImportPath'")
	}

	v := n.Client.Versioning()
	data := make([]interface{}, 0, len(specs))
	names := make([]string, 0, len(specs))
	imports := make([]string, 0, len(specs))

	tally := make(map[string]int)
	for _, s := range specs {
		name, iName, val := s.Specify(v)
		tally[name] = tally[name] + 1
		if tally[name] > 1 {
			return fmt.Errorf("%s is defined multiple times: %q", n.Singular, name)
		}
		data = append(data, val)
		names = append(names, name)
		if iName != "" {
			imports = append(imports, name)
		}
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

	// Unimport first.
	if err = n.Client.VsysUnimport(n.ImportPath, tmpl, ts, names); err != nil {
		return err
	}

	// Now configure the objects.
	if _, err = n.Client.Set(path, elm.Config(), nil, nil); err != nil {
		return err
	}

	// Finally import all valid importables.
	return n.Client.VsysImport(n.ImportPath, tmpl, ts, vsys, imports)
}

/*
Edit performs an EDIT to configure one object.

As this is an importable config, first the object is unimported, then the object
is configured, and finally the object is imported, if applicable.
*/
func (n *Importable) Edit(tmpl, ts, vsys string, pather Pather, spec ImportSpecifier) error {
	var err error

	// Sanity check: Import path should be configured.
	if n.ImportPath == "" {
		return fmt.Errorf("Namespace did not configure 'ImportPath'")
	}

	name, iName, data := spec.Specify(n.Client.Versioning())

	n.Client.LogAction("(edit) %s: %s", n.Singular, name)

	path, pErr := pather([]string{name})
	if pErr != nil {
		return pErr
	}

	// Unimport first.
	if err = n.Client.VsysUnimport(n.ImportPath, tmpl, ts, []string{name}); err != nil {
		return err
	}

	// Now configure the object.
	if _, err = n.Client.Edit(path, data, nil, nil); err != nil {
		return err
	}

	// Finally import if applicable.
	if iName != "" {
		if err = n.Client.VsysImport(n.ImportPath, tmpl, ts, vsys, []string{name}); err != nil {
			return err
		}
	}

	return nil
}

/*
Delete performs a DELETE to remove one or more objects.

As this is an importable config, first all objects are unimported, then all
objects are deleted from the config.
*/
func (n *Importable) Delete(tmpl, ts string, pather Pather, names []string, nErr error) error {
	var err error

	// Sanity check: Import path should be configured.
	if n.ImportPath == "" {
		return fmt.Errorf("Namespace did not configure 'ImportPath'")
	}

	if nErr != nil {
		return nErr
	}

	n.Client.LogAction("(delete) %s: %v", n.Plural, names)

	if len(names) == 0 {
		return nil
	}

	path, pErr := pather(names)
	if pErr != nil {
		return pErr
	}

	// Unimport first.
	if err = n.Client.VsysUnimport(n.ImportPath, tmpl, ts, names); err != nil {
		return err
	}

	// Delete the config.
	_, err = n.Client.Delete(path, nil, nil)
	return err
}
