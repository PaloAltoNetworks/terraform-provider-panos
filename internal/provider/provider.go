package provider

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	sdk "github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the provider implementation interface is sound.
var (
	_ provider.Provider              = &PanosProvider{}
	_ provider.ProviderWithFunctions = &PanosProvider{}
	_ provider.ProviderWithActions   = &PanosProvider{}
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
	MultiConfigBatchSize  types.Int64  `tfsdk:"multi_config_batch_size"`
	PanosVersion          types.String `tfsdk:"panos_version"`
	Password              types.String `tfsdk:"password"`
	Port                  types.Int64  `tfsdk:"port"`
	Protocol              types.String `tfsdk:"protocol"`
	SdkLogCategories      types.String `tfsdk:"sdk_log_categories"`
	SdkLogLevel           types.String `tfsdk:"sdk_log_level"`
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
					"PANOS_HOSTNAME",
					"hostname",
				),
				Optional: true,
			},
			"multi_config_batch_size": schema.Int64Attribute{
				Description: ProviderParamDescription(
					"Number of operations to send as part of a single MultiConfig update",
					"500",
					"PANOS_MULTI_CONFIG_BATCH_SIZE",
					"multi_config_batch_size",
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
			"sdk_log_categories": schema.StringAttribute{
				Description: ProviderParamDescription(
					"Log categories to configure for the PAN-OS SDK library",
					"",
					"PANOS_LOG_CATEGORIES",
					"sdk_log_categories",
				),
				Optional: true,
			},
			"sdk_log_level": schema.StringAttribute{
				Description: ProviderParamDescription(
					"SDK logging Level for categories",
					"INFO",
					"PANOS_LOG_LEVEL",
					"sdk_log_level",
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

type ProviderData struct {
	Client               *sdk.Client
	MultiConfigBatchSize int
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
		var logCategories sdk.LogCategory
		if !config.SdkLogCategories.IsNull() {
			categories := strings.Split(config.SdkLogCategories.ValueString(), ",")
			var err error
			logCategories, err = sdk.LogCategoryFromStrings(categories)
			if err != nil {
				resp.Diagnostics.AddError("Failed to configure Terraform provider", err.Error())
				return
			}
		}

		var logLevel slog.Level
		if !config.SdkLogLevel.IsNull() {
			levelStr := config.SdkLogLevel.ValueString()
			err := logLevel.UnmarshalText([]byte(levelStr))
			if err != nil {
				resp.Diagnostics.AddError("Failed to configure Terraform provider", fmt.Sprintf("Invalid Log Level: %s", levelStr))
			}
		} else {
			logLevel = slog.LevelInfo
		}

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
			Logging: sdk.LoggingInfo{
				LogLevel:      logLevel,
				LogCategories: logCategories,
			},
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

	batchSize := config.MultiConfigBatchSize.ValueInt64()
	if batchSize == 0 {
		batchSize = 500
	} else if batchSize < 0 || batchSize > 10000 {
		resp.Diagnostics.AddError("Failed to configure Terraform provider", fmt.Sprintf("multi_config_batch_size must be between 1 and 10000, value: %d", batchSize))
		return
	}

	providerData := &ProviderData{
		Client:               con,
		MultiConfigBatchSize: int(batchSize),
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
	resp.EphemeralResourceData = providerData
	resp.ActionData = providerData

	// Done.
	tflog.Info(ctx, "Configured client", map[string]any{"success": true})
}

// DataSources defines the data sources for this provider.
func (p *PanosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAdminRoleDataSource,
		NewAuthenticationProfileDataSource,
		NewCertificateDataSource,
		NewDnsSettingsDataSource,
		NewDynamicUpdatesDataSource,
		NewGeneralSettingsDataSource,
		NewGlobalprotectGatewayDataSource,
		NewConfigLogSettingsDataSource,
		NewCorrelationLogSettingsDataSource,
		NewGlobalprotectLogSettingsDataSource,
		NewHipmatchLogSettingsDataSource,
		NewIptagLogSettingsDataSource,
		NewSystemLogSettingsDataSource,
		NewUseridLogSettingsDataSource,
		NewNtpSettingsDataSource,
		NewLdapProfileDataSource,
		NewSyslogProfileDataSource,
		NewProxySettingsDataSource,
		NewSslDecryptDataSource,
		NewDhcpDataSource,
		NewGlobalprotectPortalDataSource,
		NewIkeGatewayDataSource,
		NewAggregateInterfaceDataSource,
		NewEthernetInterfaceDataSource,
		NewLoopbackInterfaceDataSource,
		NewTunnelInterfaceDataSource,
		NewVlanInterfaceDataSource,
		NewLogicalRouterDataSource,
		NewAntiSpywareSecurityProfileDataSource,
		NewInterfaceManagementProfileDataSource,
		NewBgpAddressFamilyRoutingProfileDataSource,
		NewBgpAuthRoutingProfileDataSource,
		NewBgpDampeningRoutingProfileDataSource,
		NewBgpFilteringRoutingProfileDataSource,
		NewBgpRedistributionRoutingProfileDataSource,
		NewBgpTimerRoutingProfileDataSource,
		NewFiltersAccessListRoutingProfileDataSource,
		NewFiltersAsPathAccessListRoutingProfileDataSource,
		NewFiltersBgpRouteMapRoutingProfileDataSource,
		NewFiltersCommunityListRoutingProfileDataSource,
		NewFiltersPrefixListRoutingProfileDataSource,
		NewFiltersRouteMapsRedistributionRoutingProfileDataSource,
		NewAggregateLayer3SubinterfaceDataSource,
		NewEthernetLayer3SubinterfaceDataSource,
		NewIpsecTunnelDataSource,
		NewVirtualRouterStaticRouteIpv4DataSource,
		NewVirtualRouterStaticRoutesIpv4DataSource,
		NewVirtualRouterStaticRouteIpv6DataSource,
		NewVirtualRouterStaticRoutesIpv6DataSource,
		NewVirtualRouterDataSource,
		NewZoneDataSource,
		NewAddressGroupDataSource,
		NewAddressDataSource,
		NewAddressesDataSource,
		NewAdministrativeTagDataSource,
		NewApplicationGroupDataSource,
		NewApplicationDataSource,
		NewCustomUrlCategoryDataSource,
		NewExternalDynamicListDataSource,
		NewAntivirusSecurityProfileDataSource,
		NewCertificateProfileDataSource,
		NewFileBlockingSecurityProfileDataSource,
		NewIkeCryptoProfileDataSource,
		NewIpsecCryptoProfileDataSource,
		NewLogForwardingProfileDataSource,
		NewSecurityProfileGroupDataSource,
		NewSslTlsServiceProfileDataSource,
		NewUrlFilteringSecurityProfileDataSource,
		NewVulnerabilitySecurityProfileDataSource,
		NewWildfireAnalysisSecurityProfileDataSource,
		NewServiceGroupDataSource,
		NewServiceDataSource,
		NewDeviceGroupParentDataSource,
		NewDeviceGroupDataSource,
		NewTemplateStackDataSource,
		NewTemplateVariableDataSource,
		NewTemplateDataSource,
		NewDecryptionPolicyDataSource,
		NewDecryptionPolicyRulesDataSource,
		NewDefaultSecurityPolicyDataSource,
		NewNatPolicyDataSource,
		NewNatPolicyRulesDataSource,
		NewSecurityPolicyDataSource,
		NewSecurityPolicyRulesDataSource,
	}
}

// Resources defines the data sources for this provider.
func (p *PanosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAdminRoleResource,
		NewAuthenticationProfileResource,
		NewCertificateImportResource,
		NewDnsSettingsResource,
		NewDynamicUpdatesResource,
		NewGeneralSettingsResource,
		NewGlobalprotectGatewayResource,
		NewConfigLogSettingsResource,
		NewCorrelationLogSettingsResource,
		NewGlobalprotectLogSettingsResource,
		NewHipmatchLogSettingsResource,
		NewIptagLogSettingsResource,
		NewSystemLogSettingsResource,
		NewUseridLogSettingsResource,
		NewNtpSettingsResource,
		NewLdapProfileResource,
		NewSyslogProfileResource,
		NewProxySettingsResource,
		NewSslDecryptResource,
		NewDhcpResource,
		NewGlobalprotectPortalResource,
		NewIkeGatewayResource,
		NewAggregateInterfaceResource,
		NewEthernetInterfaceResource,
		NewLoopbackInterfaceResource,
		NewTunnelInterfaceResource,
		NewVlanInterfaceResource,
		NewLogicalRouterResource,
		NewAntiSpywareSecurityProfileResource,
		NewInterfaceManagementProfileResource,
		NewBgpAddressFamilyRoutingProfileResource,
		NewBgpAuthRoutingProfileResource,
		NewBgpDampeningRoutingProfileResource,
		NewBgpFilteringRoutingProfileResource,
		NewBgpRedistributionRoutingProfileResource,
		NewBgpTimerRoutingProfileResource,
		NewFiltersAccessListRoutingProfileResource,
		NewFiltersAsPathAccessListRoutingProfileResource,
		NewFiltersBgpRouteMapRoutingProfileResource,
		NewFiltersCommunityListRoutingProfileResource,
		NewFiltersPrefixListRoutingProfileResource,
		NewFiltersRouteMapsRedistributionRoutingProfileResource,
		NewAggregateLayer3SubinterfaceResource,
		NewEthernetLayer3SubinterfaceResource,
		NewIpsecTunnelResource,
		NewVirtualRouterStaticRouteIpv4Resource,
		NewVirtualRouterStaticRoutesIpv4Resource,
		NewVirtualRouterStaticRouteIpv6Resource,
		NewVirtualRouterStaticRoutesIpv6Resource,
		NewVirtualRouterResource,
		NewZoneResource,
		NewAddressGroupResource,
		NewAddressResource,
		NewAddressesResource,
		NewAdministrativeTagResource,
		NewApplicationGroupResource,
		NewApplicationResource,
		NewCustomUrlCategoryResource,
		NewExternalDynamicListResource,
		NewAntivirusSecurityProfileResource,
		NewCertificateProfileResource,
		NewFileBlockingSecurityProfileResource,
		NewIkeCryptoProfileResource,
		NewIpsecCryptoProfileResource,
		NewLogForwardingProfileResource,
		NewSecurityProfileGroupResource,
		NewSslTlsServiceProfileResource,
		NewUrlFilteringSecurityProfileResource,
		NewVulnerabilitySecurityProfileResource,
		NewWildfireAnalysisSecurityProfileResource,
		NewServiceGroupResource,
		NewServiceResource,
		NewDeviceGroupParentResource,
		NewDeviceGroupResource,
		NewTemplateStackResource,
		NewTemplateVariableResource,
		NewTemplateResource,
		NewDecryptionPolicyResource,
		NewDecryptionPolicyRulesResource,
		NewDefaultSecurityPolicyResource,
		NewNatPolicyResource,
		NewNatPolicyRulesResource,
		NewSecurityPolicyResource,
		NewSecurityPolicyRulesResource,
	}
}

func (p *PanosProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewApiKeyResource,
		NewVmAuthKeyResource,
	}
}

func (p *PanosProvider) Actions(_ context.Context) []func() action.Action {
	return []func() action.Action{
		NewCommitAction,
		NewCommitAllAction,
	}
}

func (p *PanosProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewAddressValueFunction,
		NewCreateImportIdFunction,
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

type CreateResourceIdFunc func(context.Context, types.Object) ([]byte, error)

type resourceFuncs struct {
	CreateImportId CreateResourceIdFunc
}

var resourceFuncMap = map[string]resourceFuncs{
	"panos_address": resourceFuncs{
		CreateImportId: AddressImportStateCreator,
	},
	"panos_address_group": resourceFuncs{
		CreateImportId: AddressGroupImportStateCreator,
	},
	"panos_addresses": resourceFuncs{
		CreateImportId: AddressesImportStateCreator,
	},
	"panos_admin_role": resourceFuncs{
		CreateImportId: AdminRoleImportStateCreator,
	},
	"panos_administrative_tag": resourceFuncs{
		CreateImportId: AdministrativeTagImportStateCreator,
	},
	"panos_aggregate_interface": resourceFuncs{
		CreateImportId: AggregateInterfaceImportStateCreator,
	},
	"panos_aggregate_layer3_subinterface": resourceFuncs{
		CreateImportId: AggregateLayer3SubinterfaceImportStateCreator,
	},
	"panos_anti_spyware_security_profile": resourceFuncs{
		CreateImportId: AntiSpywareSecurityProfileImportStateCreator,
	},
	"panos_antivirus_security_profile": resourceFuncs{
		CreateImportId: AntivirusSecurityProfileImportStateCreator,
	},
	"panos_application": resourceFuncs{
		CreateImportId: ApplicationImportStateCreator,
	},
	"panos_application_group": resourceFuncs{
		CreateImportId: ApplicationGroupImportStateCreator,
	},
	"panos_authentication_profile": resourceFuncs{
		CreateImportId: AuthenticationProfileImportStateCreator,
	},
	"panos_bgp_address_family_routing_profile": resourceFuncs{
		CreateImportId: BgpAddressFamilyRoutingProfileImportStateCreator,
	},
	"panos_bgp_auth_routing_profile": resourceFuncs{
		CreateImportId: BgpAuthRoutingProfileImportStateCreator,
	},
	"panos_bgp_dampening_routing_profile": resourceFuncs{
		CreateImportId: BgpDampeningRoutingProfileImportStateCreator,
	},
	"panos_bgp_filtering_routing_profile": resourceFuncs{
		CreateImportId: BgpFilteringRoutingProfileImportStateCreator,
	},
	"panos_bgp_redistribution_routing_profile": resourceFuncs{
		CreateImportId: BgpRedistributionRoutingProfileImportStateCreator,
	},
	"panos_bgp_timer_routing_profile": resourceFuncs{
		CreateImportId: BgpTimerRoutingProfileImportStateCreator,
	},
	"panos_certificate_profile": resourceFuncs{
		CreateImportId: CertificateProfileImportStateCreator,
	},
	"panos_config_log_settings": resourceFuncs{
		CreateImportId: ConfigLogSettingsImportStateCreator,
	},
	"panos_correlation_log_settings": resourceFuncs{
		CreateImportId: CorrelationLogSettingsImportStateCreator,
	},
	"panos_custom_url_category": resourceFuncs{
		CreateImportId: CustomUrlCategoryImportStateCreator,
	},
	"panos_decryption_policy": resourceFuncs{
		CreateImportId: DecryptionPolicyImportStateCreator,
	},
	"panos_decryption_policy_rules": resourceFuncs{
		CreateImportId: DecryptionPolicyRulesImportStateCreator,
	},
	"panos_default_security_policy": resourceFuncs{
		CreateImportId: DefaultSecurityPolicyImportStateCreator,
	},
	"panos_device_group": resourceFuncs{
		CreateImportId: DeviceGroupImportStateCreator,
	},
	"panos_dhcp": resourceFuncs{
		CreateImportId: DhcpImportStateCreator,
	},
	"panos_ethernet_interface": resourceFuncs{
		CreateImportId: EthernetInterfaceImportStateCreator,
	},
	"panos_ethernet_layer3_subinterface": resourceFuncs{
		CreateImportId: EthernetLayer3SubinterfaceImportStateCreator,
	},
	"panos_external_dynamic_list": resourceFuncs{
		CreateImportId: ExternalDynamicListImportStateCreator,
	},
	"panos_file_blocking_security_profile": resourceFuncs{
		CreateImportId: FileBlockingSecurityProfileImportStateCreator,
	},
	"panos_filters_access_list_routing_profile": resourceFuncs{
		CreateImportId: FiltersAccessListRoutingProfileImportStateCreator,
	},
	"panos_filters_as_path_access_list_routing_profile": resourceFuncs{
		CreateImportId: FiltersAsPathAccessListRoutingProfileImportStateCreator,
	},
	"panos_filters_bgp_route_map_routing_profile": resourceFuncs{
		CreateImportId: FiltersBgpRouteMapRoutingProfileImportStateCreator,
	},
	"panos_filters_community_list_routing_profile": resourceFuncs{
		CreateImportId: FiltersCommunityListRoutingProfileImportStateCreator,
	},
	"panos_filters_prefix_list_routing_profile": resourceFuncs{
		CreateImportId: FiltersPrefixListRoutingProfileImportStateCreator,
	},
	"panos_filters_route_maps_redistribution_routing_profile": resourceFuncs{
		CreateImportId: FiltersRouteMapsRedistributionRoutingProfileImportStateCreator,
	},
	"panos_globalprotect_gateway": resourceFuncs{
		CreateImportId: GlobalprotectGatewayImportStateCreator,
	},
	"panos_globalprotect_log_settings": resourceFuncs{
		CreateImportId: GlobalprotectLogSettingsImportStateCreator,
	},
	"panos_globalprotect_portal": resourceFuncs{
		CreateImportId: GlobalprotectPortalImportStateCreator,
	},
	"panos_hipmatch_log_settings": resourceFuncs{
		CreateImportId: HipmatchLogSettingsImportStateCreator,
	},
	"panos_ike_crypto_profile": resourceFuncs{
		CreateImportId: IkeCryptoProfileImportStateCreator,
	},
	"panos_ike_gateway": resourceFuncs{
		CreateImportId: IkeGatewayImportStateCreator,
	},
	"panos_interface_management_profile": resourceFuncs{
		CreateImportId: InterfaceManagementProfileImportStateCreator,
	},
	"panos_ipsec_crypto_profile": resourceFuncs{
		CreateImportId: IpsecCryptoProfileImportStateCreator,
	},
	"panos_ipsec_tunnel": resourceFuncs{
		CreateImportId: IpsecTunnelImportStateCreator,
	},
	"panos_iptag_log_settings": resourceFuncs{
		CreateImportId: IptagLogSettingsImportStateCreator,
	},
	"panos_ldap_profile": resourceFuncs{
		CreateImportId: LdapProfileImportStateCreator,
	},
	"panos_log_forwarding_profile": resourceFuncs{
		CreateImportId: LogForwardingProfileImportStateCreator,
	},
	"panos_logical_router": resourceFuncs{
		CreateImportId: LogicalRouterImportStateCreator,
	},
	"panos_loopback_interface": resourceFuncs{
		CreateImportId: LoopbackInterfaceImportStateCreator,
	},
	"panos_nat_policy": resourceFuncs{
		CreateImportId: NatPolicyImportStateCreator,
	},
	"panos_nat_policy_rules": resourceFuncs{
		CreateImportId: NatPolicyRulesImportStateCreator,
	},
	"panos_security_policy": resourceFuncs{
		CreateImportId: SecurityPolicyImportStateCreator,
	},
	"panos_security_policy_rules": resourceFuncs{
		CreateImportId: SecurityPolicyRulesImportStateCreator,
	},
	"panos_security_profile_group": resourceFuncs{
		CreateImportId: SecurityProfileGroupImportStateCreator,
	},
	"panos_service": resourceFuncs{
		CreateImportId: ServiceImportStateCreator,
	},
	"panos_service_group": resourceFuncs{
		CreateImportId: ServiceGroupImportStateCreator,
	},
	"panos_ssl_tls_service_profile": resourceFuncs{
		CreateImportId: SslTlsServiceProfileImportStateCreator,
	},
	"panos_syslog_profile": resourceFuncs{
		CreateImportId: SyslogProfileImportStateCreator,
	},
	"panos_system_log_settings": resourceFuncs{
		CreateImportId: SystemLogSettingsImportStateCreator,
	},
	"panos_template": resourceFuncs{
		CreateImportId: TemplateImportStateCreator,
	},
	"panos_template_stack": resourceFuncs{
		CreateImportId: TemplateStackImportStateCreator,
	},
	"panos_template_variable": resourceFuncs{
		CreateImportId: TemplateVariableImportStateCreator,
	},
	"panos_tunnel_interface": resourceFuncs{
		CreateImportId: TunnelInterfaceImportStateCreator,
	},
	"panos_url_filtering_security_profile": resourceFuncs{
		CreateImportId: UrlFilteringSecurityProfileImportStateCreator,
	},
	"panos_userid_log_settings": resourceFuncs{
		CreateImportId: UseridLogSettingsImportStateCreator,
	},
	"panos_virtual_router": resourceFuncs{
		CreateImportId: VirtualRouterImportStateCreator,
	},
	"panos_virtual_router_static_route_ipv4": resourceFuncs{
		CreateImportId: VirtualRouterStaticRouteIpv4ImportStateCreator,
	},
	"panos_virtual_router_static_route_ipv6": resourceFuncs{
		CreateImportId: VirtualRouterStaticRouteIpv6ImportStateCreator,
	},
	"panos_virtual_router_static_routes_ipv4": resourceFuncs{
		CreateImportId: VirtualRouterStaticRoutesIpv4ImportStateCreator,
	},
	"panos_virtual_router_static_routes_ipv6": resourceFuncs{
		CreateImportId: VirtualRouterStaticRoutesIpv6ImportStateCreator,
	},
	"panos_vlan_interface": resourceFuncs{
		CreateImportId: VlanInterfaceImportStateCreator,
	},
	"panos_vulnerability_security_profile": resourceFuncs{
		CreateImportId: VulnerabilitySecurityProfileImportStateCreator,
	},
	"panos_wildfire_analysis_security_profile": resourceFuncs{
		CreateImportId: WildfireAnalysisSecurityProfileImportStateCreator,
	},
	"panos_zone": resourceFuncs{
		CreateImportId: ZoneImportStateCreator,
	},
}
