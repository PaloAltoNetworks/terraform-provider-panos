package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type Error struct {
	message string
	err     error
}

func (o *Error) Error() string {
	if o.err != nil {
		return fmt.Sprintf("%s: %s", o.message, o.err)
	}

	return o.message
}

func (o *Error) Unwrap() error {
	return o.err
}

var (
	ErrConflict          = errors.New("entry from the plan already exists on the server")
	ErrMissingUuid       = errors.New("entry is missing required uuid")
	ErrMarshaling        = errors.New("failed to marshal entry to XML document")
	ErrInvalidPosition   = errors.New("position is not valid")
	ErrMissingPivotPoint = errors.New("provided pivot entry does not exist")
	ErrInternal          = errors.New("internal provider error")
	ErrObjectNotFound    = errors.New("Object not found")
)

type entryState string

const (
	entryUnknown  entryState = "unknown"
	entryMissing  entryState = "missing"
	entryOutdated entryState = "outdated"
	entryRenamed  entryState = "renamed"
	entryDeleted  entryState = "deleted"
	entryOk       entryState = "ok"
)

type SDKClient interface {
	Versioning() version.Number
	GetTarget() string
	MultiConfig(context.Context, *xmlapi.MultiConfig, bool, url.Values) ([]byte, *http.Response, *xmlapi.MultiConfigResponse, error)
}

type ImportLocation interface {
	XpathForLocation(version.Number, util.ILocation) ([]string, error)
	MarshalPangoXML([]string) (string, error)
	UnmarshalPangoXML([]byte) ([]string, error)
}
