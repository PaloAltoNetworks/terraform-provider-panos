package provider

import (
	"github.com/PaloAltoNetworks/pango/errors"
)

var InspectionModeError = "Resources are unavailable when the provider is in inspection mode.  Resources are only available in API mode."

func IsObjectNotFound(e error) bool {
	e2, ok := e.(errors.Panos)
	if ok && e2.ObjectNotFound() {
		return true
	}

	return false
}
