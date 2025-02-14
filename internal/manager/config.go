package manager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

type TFConfigObject[E any] interface {
	CopyToPango(context.Context, *map[string]types.String) (E, diag.Diagnostics)
	CopyFromPango(context.Context, E, *map[string]types.String) diag.Diagnostics
}

type SDKConfigService[C any, L ConfigLocation] interface {
	Create(context.Context, L, C) (C, error)
	Update(context.Context, L, C) (C, error)
	Read(context.Context, L, string) (C, error)
	Delete(context.Context, L, C) error
}

type ConfigLocation interface {
	Xpath(version.Number) ([]string, error)
}

type ConfigObjectManager[C any, L ConfigLocation, S SDKConfigService[C, L]] struct {
	service   S
	client    util.PangoClient
	specifier func(C) (any, error)
}

func NewConfigObjectManager[C any, L ConfigLocation, S SDKConfigService[C, L]](client util.PangoClient, service S, specifier func(C) (any, error)) *ConfigObjectManager[C, L, S] {
	return &ConfigObjectManager[C, L, S]{
		service:   service,
		client:    client,
		specifier: specifier,
	}
}

func (o *ConfigObjectManager[C, L, S]) Create(ctx context.Context, location L, config C) (C, error) {
	return o.service.Create(ctx, location, config)
}

func (o *ConfigObjectManager[C, L, S]) Update(ctx context.Context, location L, config C) (C, error) {
	return o.service.Update(ctx, location, config)
}

func (o *ConfigObjectManager[C, L, S]) Read(ctx context.Context, location L) (C, error) {
	obj, err := o.service.Read(ctx, location, "get")
	if err != nil && sdkerrors.IsObjectNotFound(err) {
		return obj, ErrObjectNotFound
	}

	return obj, err
}

func (o *ConfigObjectManager[C, L, S]) Delete(ctx context.Context, location L, config C) error {
	return o.service.Delete(ctx, location, config)
}
