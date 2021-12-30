package namespace

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
)

// Common contains the shared methods every namespace has.
type Common struct {
	Singular   string
	Plural     string
	Client     util.XapiClient
	Predefined bool
}

/*
Object returns a single object's config.

cmd should be util.Get or util.Show.
pather creates the xpath to use.
name is used for logging only, but is the name of the object to return.
ans is an interface to unmarshal the response into.
*/
func (n *Common) Object(cmd string, pather Pather, name string, ans interface{}) error {
	path, pErr := pather([]string{name})
	if err := n.retrieve(cmd, path, true, name, false, false, ans, pErr); err != nil {
		return err
	}

	return nil
}

/*
Objects returns multiple object's config.

cmd should be util.Get or util.Show.
path is the xpath.
ans is an interface to unmarshal the response into.
*/
func (n *Common) Objects(cmd string, pather Pather, ans interface{}) error {
	path, pErr := pather(nil)
	if err := n.retrieve(cmd, path, false, "", true, false, ans, pErr); err != nil {
		return err
	}

	return nil
}

/*
Listing returns a list of names.

cmd should be util.Get or util.Show.
path is the xpath.
ans is an interface to unmarshal the response into.
*/
func (n *Common) Listing(cmd string, pather Pather, ans Namer) ([]string, error) {
	path, pErr := pather(nil)
	if err := n.retrieve(cmd, path, false, "", true, true, ans, pErr); err != nil {
		return nil, err
	}

	return ans.Names(), nil
}

// retrieve does either a GET or SHOW to retrieve config.
func (n *Common) retrieve(cmd string, path []string, singular bool, singleDesc string, plural, namesOnly bool, ans interface{}, pErr error) error {
	var err error
	var data []byte
	var tag string

	// Sanity checks.
	if cmd != util.Get && cmd != util.Show {
		return fmt.Errorf("invalid cmd: %s", cmd)
	}
	if ans == nil {
		return fmt.Errorf("ans must be specified")
	}

	if pErr != nil {
		return pErr
	}

	// Do logging and determine the actual path to query.
	var predef string
	if n.Predefined {
		predef = "predefined "
	}
	if singular {
		if singleDesc != "" {
			n.Client.LogQuery("(%s) %s%s: %s", cmd, predef, n.Singular, singleDesc)
		} else {
			n.Client.LogQuery("(%s) %s%s", cmd, predef, n.Singular)
		}
	} else if plural {
		tag = path[len(path)-2]
		if cmd == util.Show {
			path = path[:len(path)-1]
		}
		if namesOnly {
			if cmd == util.Get {
				path = append(path, "@name")
			}
			n.Client.LogQuery("(%s) %s%s names", cmd, predef, n.Singular)
		} else {
			if cmd == util.Get {
				path = path[:len(path)-1]
			}
			n.Client.LogQuery("(%s) list of %s%s", cmd, predef, n.Plural)
		}
	}

	// Perform the query.
	switch cmd {
	case util.Get:
		data, err = n.Client.Get(path, nil, nil)
	case util.Show:
		data, err = n.Client.Show(path, nil, nil)
	}
	if err != nil {
		if plural && (err.Error() == "No such node" || err.Error() == "Object not found") {
			return nil
		}
		return err
	}

	// Unmarshal the response into the given struct.
	data = util.StripPanosPackaging(data, tag)
	return UnpackageXmlInto(data, ans)
}

/*
AllFromPanosConfig returns multiple objects' config from the configuration
retrieved from PAN-OS and stored in the client's config tree.

path is the xpath function.
ans is an interface to unmarshal the found config into.
*/
func (n *Common) AllFromPanosConfig(pather Pather, ans interface{}) error {
	path, pErr := pather(nil)
	return n.loadConfig(path, ans, pErr)
}

/*
FromPanosConfig returns a single object's config from the configuration
retrieved from PAN-OS and stored in the client's config tree.

path is the xpath function.
name is the name of the object to retrieve.
ans is an interface to unmarshal the found config into.
*/
func (n *Common) FromPanosConfig(pather Pather, name string, ans Namer) error {
	path, pErr := pather([]string{name})
	if err := n.loadConfig(path, ans, pErr); err != nil {
		return err
	}

	names := ans.Names()
	if len(names) == 0 {
		return errors.ObjectNotFound()
	}

	return nil
}

func (n *Common) loadConfig(path []string, ans interface{}, pErr error) error {
	if pErr != nil {
		return pErr
	}

	config := n.Client.ConfigTree()
	if config == nil {
		return fmt.Errorf("no config tree loaded")
	}

	elm := util.FindXmlNodeInTree(path, config)
	if elm == nil {
		return nil
	}

	data, err := xml.Marshal(elm)
	if err != nil {
		return err
	}

	return UnpackageXmlInto(data, ans)
}
