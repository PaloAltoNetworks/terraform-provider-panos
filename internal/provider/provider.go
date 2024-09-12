package provider

import (
	"context"

	sdk "github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the provider implementation interface is sound.
var (
	_ provider.Provider = &PanosProvider{}
)

// PanosProvider is the provider implementation.
type PanosProvider struct {
	version string
}

// PanosProviderModel maps provider schema data to a Go type.
type PanosProviderModel struct {
	AdditionalHeaders     types.Map    `tfsdk:"additional_headers"`
	ApiKey                types.String `tfsdk:"api_key"`
	ApiKeyInRequest       types.Bool   `tfsdk:"api_key_in_request"`
	AuthFile              types.String `tfsdk:"auth_file"`
	ConfigFile            types.String `tfsdk:"config_file"`
	Hostname              types.String `tfsdk:"hostname"`
	PanosVersion          types.String `tfsdk:"panos_version"`
	Password              types.String `tfsdk:"password"`
	Port                  types.Int64  `tfsdk:"port"`
	Protocol              types.String `tfsdk:"protocol"`
	SkipVerifyCertificate types.Bool   `tfsdk:"skip_verify_certificate"`
	Target                types.String `tfsdk:"target"`
	Username              types.String `tfsdk:"username"`
}

// Metadata returns the provider type name.
func (p *PanosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "panos"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *PanosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider to interact with Palo Alto Networks PAN-OS.",
		Attributes: map[string]schema.Attribute{
			"additional_headers": schema.MapAttribute{
				Description: ProviderParamDescription(
					"Additional HTTP headers to send with API calls",
					"",
					"PANOS_HEADERS",
					"additional_headers",
				),
				Optional:    true,
				ElementType: types.StringType,
			},
			"api_key": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The API key for PAN-OS. Either specify this or give both username and password.",
					"",
					"PANOS_API_KEY",
					"api_key",
				),
				Optional: true,
			},
			"api_key_in_request": schema.BoolAttribute{
				Description: ProviderParamDescription(
					"Send the API key in the request body instead of using the authentication header.",
					"",
					"PANOS_API_KEY_IN_REQUEST",
					"api_key_in_request",
				),
				Optional: true,
			},
			"auth_file": schema.StringAttribute{
				Description: ProviderParamDescription(
					"Filesystem path to a JSON config file that specifies the provider's params.",
					"",
					"",
					"auth_file",
				),
				Optional: true,
			},
			"config_file": schema.StringAttribute{
				Description: ProviderParamDescription(
					"(Local inspection mode) The PAN-OS config file to load read in using `file()`",
					"",
					"",
					"config_file",
				),
				Optional: true,
			},
			"hostname": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The hostname or IP address of the PAN-OS instance (NGFW or Panorama).",
					"",
					"PANOS_HOST",
					"hostname",
				),
				Optional: true,
			},
			"panos_version": schema.StringAttribute{
				Description: ProviderParamDescription(
					"(Local inspection mode) The version of PAN-OS that exported the config file. This is only used if the root 'config' block does not contain the 'detail-version' attribute. Example: `10.2.3`.",
					"",
					"",
					"panos_version",
				),
				Optional: true,
			},
			"password": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The password.  This is required if the api_key is not configured.",
					"",
					"PANOS_PASSWORD",
					"password",
				),
				Optional:  true,
				Sensitive: true,
			},
			"port": schema.Int64Attribute{
				Description: ProviderParamDescription(
					"If the port is non-standard for the protocol, the port number to use.",
					"",
					"PANOS_PORT",
					"port",
				),
				Optional: true,
			},
			"protocol": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The protocol (https or http).",
					"https",
					"PANOS_PROTOCOL",
					"protocol",
				),
				Optional: true,
			},
			"skip_verify_certificate": schema.BoolAttribute{
				Description: ProviderParamDescription(
					"(For https protocol) Skip verifying the HTTPS certificate.",
					"",
					"PANOS_SKIP_VERIFY_CERTIFICATE",
					"skip_verify_certificate",
				),
				Optional: true,
			},
			"target": schema.StringAttribute{
				Description: ProviderParamDescription(
					"Target setting (NGFW serial number).",
					"",
					"PANOS_TARGET",
					"target",
				),
				Optional: true,
			},
			"username": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The username.  This is required if api_key is not configured.",
					"",
					"PANOS_USERNAME",
					"username",
				),
				Optional: true,
			},
		},
	}
}

// Configure prepares the provider.
func (p *PanosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring the provider client...")

	var config PanosProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var con *sdk.Client

	if config.ConfigFile.ValueStringPointer() != nil {
		tflog.Info(ctx, "Configuring client for local inspection mode")
		con = &sdk.Client{}
		if err := con.SetupLocalInspection(config.ConfigFile.ValueString(), config.PanosVersion.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error setting up local inspection mode", err.Error())
			return
		}
	} else {
		tflog.Info(ctx, "Configuring client for API mode")
		con = &sdk.Client{
			Hostname:        config.Hostname.ValueString(),
			Username:        config.Username.ValueString(),
			Password:        config.Password.ValueString(),
			ApiKey:          config.ApiKey.ValueString(),
			Protocol:        config.Protocol.ValueString(),
			Port:            int(config.Port.ValueInt64()),
			Target:          config.Target.ValueString(),
			ApiKeyInRequest: config.ApiKeyInRequest.ValueBool(),
			// Headers from AdditionalHeaders
			SkipVerifyCertificate: config.SkipVerifyCertificate.ValueBool(),
			AuthFile:              config.AuthFile.ValueString(),
			CheckEnvironment:      true,
			//Agent:            fmt.Sprintf("Terraform/%s Provider/scm Version/%s", req.TerraformVersion, p.version),
		}

		if err := con.Setup(); err != nil {
			resp.Diagnostics.AddError("Provider parameter value error", err.Error())
			return
		}

		//con.HttpClient.Transport = sdkapi.NewTransport(con.HttpClient.Transport, con)

		if err := con.Initialize(ctx); err != nil {
			resp.Diagnostics.AddError("Initialization error", err.Error())
			return
		}
	}

	resp.DataSourceData = con
	resp.ResourceData = con

	// Done.
	tflog.Info(ctx, "Configured client", map[string]any{"success": true})
}

// DataSources defines the data sources for this provider.
func (p *PanosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTfidDataSource,
		NewDnsSettingsDataSource,
		NewNtpSettingsDataSource,
		NewEthernetInterfaceDataSource,
		NewLoopbackInterfaceDataSource,
		NewInterfaceManagementProfileDataSource,
		NewVirtualRouterDataSource,
		NewZoneDataSource,
		NewAddressGroupDataSource,
		NewAddressesDataSource,
		NewAdministrativeTagDataSource,
		NewCustomUrlCategoryDataSource,
		NewServiceGroupDataSource,
		NewServiceDataSource,
		NewDeviceGroupParentDataSource,
		NewDeviceGroupDataSource,
		NewTemplateStackDataSource,
		NewTemplateVariableDataSource,
		NewTemplateDataSource,
		NewSecurityPolicyDataSource,
		NewSecurityPolicyRulesDataSource,
	}
}

// Resources defines the data sources for this provider.
func (p *PanosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDnsSettingsResource,
		NewNtpSettingsResource,
		NewEthernetInterfaceResource,
		NewLoopbackInterfaceResource,
		NewInterfaceManagementProfileResource,
		NewVirtualRouterResource,
		NewZoneResource,
		NewAddressGroupResource,
		NewAddressesResource,
		NewAdministrativeTagResource,
		NewCustomUrlCategoryResource,
		NewServiceGroupResource,
		NewServiceResource,
		NewDeviceGroupParentResource,
		NewDeviceGroupResource,
		NewTemplateStackResource,
		NewTemplateVariableResource,
		NewTemplateResource,
		NewSecurityPolicyResource,
		NewSecurityPolicyRulesResource,
	}
}

func (p *PanosProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewAddressValueFunction,
	}
}

// New is a helper function to get the provider implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PanosProvider{
			version: version,
		}
	}
}
