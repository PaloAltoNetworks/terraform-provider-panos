package namespace

import (
	"github.com/PaloAltoNetworks/pango/plugin"
	"github.com/PaloAltoNetworks/pango/version"
)

// Namer returns the names of objects returned from PAN-OS.
type Namer interface {
	Names() []string
}

// MovePather returns an xpath given the name.
type MovePather func(string) []string

// MoveLister returns a list of current rules.
type MoveLister func() ([]string, error)

// Pather returns an xpath for the given list of names.
type Pather func([]string) ([]string, error)

// Specifier is an object that has a Specify function given the current version number.
//
// There are two items returned:
//
// 1) the unique name of this config element
// 2) a struct specific to this version PAN-OS representing the desired config
type Specifier interface {
	Specify(version.Number) (string, interface{})
}

// ImportSpecifier is an object that has a Specify function given the current version number.
//
// There are three items returned:
//
// 1) the unique name of this config element
// 2) the unique name to be used when importing; an empty string means "do not import"
// 3) a struct specific to this version PAN-OS representing the desired config
type ImportSpecifier interface {
	Specify(version.Number) (string, string, interface{})
}

/*
PluginSpecifier is an object that has a Specify function given a list of plugins.

There are three items returned:

1) the unique name of this config element
2) a struct specific to this version PAN-OS representing the desired config
3) an error if there is a mismatch between the plugins installed and what is supported
*/
type PluginSpecifier interface {
	Specify([]plugin.Info) (string, interface{}, error)
}
