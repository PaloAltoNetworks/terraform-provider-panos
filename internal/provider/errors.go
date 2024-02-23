package provider

import (
    "github.com/PaloAltoNetworks/pango/errors"
)

func IsObjectNotFound(e error) bool {
	e2, ok := e.(errors.Panos)
	if ok && e2.ObjectNotFound() {
		return true
	}

	return false
}
