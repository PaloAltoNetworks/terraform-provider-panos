package provider_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	sdkClient *pango.Client

	testAccProviders = map[string]func() (tfprotov6.ProviderServer, error){
		"panos": providerserver.NewProtocol6WithError(provider.New(version)()),
	}
)

func init() {
	sdkClient = &pango.Client{
		CheckEnvironment: true,
	}

	ctx := context.Background()

	if err := sdkClient.Setup(); err != nil {
		slog.Error("setting up pango client: ", slog.String("error", err.Error()))
	}

	if err := sdkClient.Initialize(ctx); err != nil {
		slog.Error("initialization pango client: ", slog.String("error", err.Error()))
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("PANOS_HOSTNAME") == "" {
		t.Fatal("PANOS_HOSTNAME must be set for acceptance tests")
	}

	if os.Getenv("PANOS_API") != "" {
		return
	}

	if os.Getenv("PANOS_USERNAME") == "" {
		t.Fatal("PANOS_USERNAME must be set for acceptance tests")
	}

	if os.Getenv("PANOS_PASSWORD") == "" {
		t.Fatal("PANOS_PASSWORD must be set for acceptance tests")
	}
}
