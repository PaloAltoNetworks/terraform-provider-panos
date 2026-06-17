package provider

import (
	"context"

	"github.com/PaloAltoNetworks/pango/util"
)

func administratorCreatePasswordHash(ctx context.Context, client util.PangoClient, password string) (string, error) {
	return client.RequestPasswordHash(ctx, password)
}
