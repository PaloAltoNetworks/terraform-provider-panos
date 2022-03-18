package panos

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname/IP address of the Palo Alto Networks firewall to connect to",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username (not used if the API key is set)",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password (not used if the API key is set)",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The api key of the firewall",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The protocol (https or http)",
				ValidateFunc: validateStringIn("https", "http", ""),
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If the port is non-standard for the protocol, the port number to use",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout for all communications with the firewall",
			},
			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target setting (NGFW serial number)",
			},
			"additional_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Additional HTTP headers to send with API calls",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"logging": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Logging options for the API connection",
			},
			"verify_certificate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "For HTTPS protocol connections, verify the certificate",
			},
			"json_config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Retrieve the provider configuration from this JSON file",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			// Shared data sources.
			"panos_address_object":                      dataSourceAddressObject(),
			"panos_address_objects":                     dataSourceAddressObjects(),
			"panos_anti_spyware_security_profile":       dataSourceAntiSpywareSecurityProfile(),
			"panos_anti_spyware_security_profiles":      dataSourceAntiSpywareSecurityProfiles(),
			"panos_antivirus_security_profile":          dataSourceAntivirusSecurityProfile(),
			"panos_antivirus_security_profiles":         dataSourceAntivirusSecurityProfiles(),
			"panos_api_key":                             dataSourceApiKey(),
			"panos_application_object":                  dataSourceApplicationObject(),
			"panos_application_objects":                 dataSourceApplicationObjects(),
			"panos_arp":                                 dataSourceArp(),
			"panos_arps":                                dataSourceArps(),
			"panos_audit_comment_history":               dataSourceAuditCommentHistory(),
			"panos_certificate_profile":                 dataSourceCertificateProfile(),
			"panos_certificate_profiles":                dataSourceCertificateProfiles(),
			"panos_custom_data_pattern_object":          dataSourceCustomDataPatternObject(),
			"panos_custom_data_pattern_objects":         dataSourceCustomDataPatternObjects(),
			"panos_custom_url_category":                 dataSourceCustomUrlCategory(),
			"panos_custom_url_categories":               dataSourceCustomUrlCategories(),
			"panos_data_filtering_security_profile":     dataSourceDataFilteringSecurityProfile(),
			"panos_data_filtering_security_profiles":    dataSourceDataFilteringSecurityProfiles(),
			"panos_decryption_rule":                     dataSourceDecryptionRule(),
			"panos_decryption_rules":                    dataSourceDecryptionRules(),
			"panos_device_group_parent":                 dataSourceDeviceGroupParent(),
			"panos_dos_protection_profile":              dataSourceDosProtectionProfile(),
			"panos_dos_protection_profiles":             dataSourceDosProtectionProfiles(),
			"panos_dynamic_user_group":                  dataSourceDynamicUserGroup(),
			"panos_dynamic_user_groups":                 dataSourceDynamicUserGroups(),
			"panos_edl":                                 dataSourceEdl(),
			"panos_edls":                                dataSourceEdls(),
			"panos_email_server_profile":                dataSourceEmailServerProfile(),
			"panos_email_server_profiles":               dataSourceEmailServerProfiles(),
			"panos_file_blocking_security_profile":      dataSourceFileBlockingSecurityProfile(),
			"panos_file_blocking_security_profiles":     dataSourceFileBlockingSecurityProfiles(),
			"panos_local_user_db_group":                 dataSourceLocalUserDbGroup(),
			"panos_local_user_db_groups":                dataSourceLocalUserDbGroups(),
			"panos_nat_rule":                            dataSourceNatRule(),
			"panos_nat_rules":                           dataSourceNatRules(),
			"panos_ospf":                                dataSourceOspf(),
			"panos_ospf_area":                           dataSourceOspfArea(),
			"panos_ospf_areas":                          dataSourceOspfAreas(),
			"panos_ospf_area_interface":                 dataSourceOspfAreaInterface(),
			"panos_ospf_area_interfaces":                dataSourceOspfAreaInterfaces(),
			"panos_ospf_area_virtual_link":              dataSourceOspfAreaVirtualLink(),
			"panos_ospf_area_virtual_links":             dataSourceOspfAreaVirtualLinks(),
			"panos_ospf_auth_profiles":                  dataSourceOspfAuthProfiles(),
			"panos_ospf_export":                         dataSourceOspfExport(),
			"panos_ospf_exports":                        dataSourceOspfExports(),
			"panos_pbf_rule":                            dataSourcePbfRule(),
			"panos_pbf_rules":                           dataSourcePbfRules(),
			"panos_plugin":                              dataSourcePlugin(),
			"panos_predefined_dlp_file_type":            dataSourcePredefinedDlpFileType(),
			"panos_predefined_tdb_file_type":            dataSourcePredefinedTdbFileType(),
			"panos_predefined_threat":                   dataSourcePredefinedThreat(),
			"panos_security_profile_group":              dataSourceSecurityProfileGroup(),
			"panos_security_profile_groups":             dataSourceSecurityProfileGroups(),
			"panos_security_rule":                       dataSourceSecurityRule(),
			"panos_security_rules":                      dataSourceSecurityRules(),
			"panos_ssl_decrypt":                         dataSourceSslDecrypt(),
			"panos_syslog_server_profile":               dataSourceSyslogServerProfile(),
			"panos_syslog_server_profiles":              dataSourceSyslogServerProfiles(),
			"panos_system_info":                         dataSourceSystemInfo(),
			"panos_tech_support_file":                   dataSourceTechSupportFile(),
			"panos_url_filtering_security_profile":      dataSourceUrlFilteringSecurityProfile(),
			"panos_url_filtering_security_profiles":     dataSourceUrlFilteringSecurityProfiles(),
			"panos_virtual_router":                      dataSourceVirtualRouter(),
			"panos_virtual_routers":                     dataSourceVirtualRouters(),
			"panos_vulnerability_security_profile":      dataSourceVulnerabilitySecurityProfile(),
			"panos_vulnerability_security_profiles":     dataSourceVulnerabilitySecurityProfiles(),
			"panos_wildfire_analysis_security_profile":  dataSourceWildfireAnalysisSecurityProfile(),
			"panos_wildfire_analysis_security_profiles": dataSourceWildfireAnalysisSecurityProfiles(),
			"panos_zone":                                dataSourceZone(),
			"panos_zones":                               dataSourceZones(),

			// Firewall data sources.
			"panos_dhcp_interface_info": dataSourceDhcpInterfaceInfo(),
			"panos_ip_tag":              dataSourceIpTag(),
			"panos_user_tag":            dataSourceUserTag(),

			// Panorama data sources.
			"panos_vm_auth_key":   dataSourceVmAuthKey(),
			"panos_device_group":  dataSourceDeviceGroup(),
			"panos_device_groups": dataSourceDeviceGroups(),

			// Aliases.
			"panos_panorama_plugin": dataSourcePlugin(),
		},

		ResourcesMap: map[string]*schema.Resource{
			// Shared resources.
			"panos_address_object":                     resourceAddressObject(),
			"panos_address_objects":                    resourceAddressObjects(),
			"panos_anti_spyware_security_profile":      resourceAntiSpywareSecurityProfile(),
			"panos_antivirus_security_profile":         resourceAntivirusSecurityProfile(),
			"panos_arp":                                resourceArp(),
			"panos_certificate_import":                 resourceCertificateImport(),
			"panos_certificate_profile":                resourceCertificateProfile(),
			"panos_custom_data_pattern_object":         resourceCustomDataPatternObject(),
			"panos_custom_url_category":                resourceCustomUrlCategory(),
			"panos_custom_url_category_entry":          resourceCustomUrlCategoryEntry(),
			"panos_data_filtering_security_profile":    resourceDataFilteringSecurityProfile(),
			"panos_decryption_rule_group":              resourceDecryptionRuleGroup(),
			"panos_dhcp_relay":                         resourceDhcpRelay(),
			"panos_dos_protection_profile":             resourceDosProtectionProfile(),
			"panos_dynamic_user_group":                 resourceDynamicUserGroup(),
			"panos_file_blocking_security_profile":     resourceFileBlockingSecurityProfile(),
			"panos_local_user_db_group":                resourceLocalUserDbGroup(),
			"panos_local_user_db_user":                 resourceLocalUserDbUser(),
			"panos_ospf":                               resourceOspf(),
			"panos_ospf_area":                          resourceOspfArea(),
			"panos_ospf_area_interface":                resourceOspfAreaInterface(),
			"panos_ospf_area_virtual_link":             resourceOspfAreaVirtualLink(),
			"panos_ospf_auth_profile":                  resourceOspfAuthProfile(),
			"panos_ospf_export":                        resourceOspfExport(),
			"panos_security_profile_group":             resourceSecurityProfileGroup(),
			"panos_ssl_decrypt":                        resourceSslDecrypt(),
			"panos_ssl_decrypt_trusted_root_ca_entry":  resourceSslDecryptTrustedRootCaEntry(),
			"panos_url_filtering_security_profile":     resourceUrlFilteringSecurityProfile(),
			"panos_vm_information_source":              resourceVmInformationSource(),
			"panos_vulnerability_security_profile":     resourceVulnerabilitySecurityProfile(),
			"panos_wildfire_analysis_security_profile": resourceWildfireAnalysisSecurityProfile(),

			// Panorama resources.
			"panos_device_group":                                  resourceDeviceGroup(),
			"panos_device_group_entry":                            resourceDeviceGroupEntry(),
			"panos_device_group_parent":                           resourceDeviceGroupParent(),
			"panos_panorama_address_group":                        resourcePanoramaAddressGroup(),
			"panos_panorama_address_object":                       resourcePanoramaAddressObject(),
			"panos_panorama_administrative_tag":                   resourcePanoramaAdministrativeTag(),
			"panos_panorama_aggregate_interface":                  resourcePanoramaAggregateInterface(),
			"panos_panorama_application_group":                    resourcePanoramaApplicationGroup(),
			"panos_panorama_application_object":                   resourcePanoramaApplicationObject(),
			"panos_panorama_application_signature":                resourcePanoramaApplicationSignature(),
			"panos_panorama_bfd_profile":                          resourcePanoramaBfdProfile(),
			"panos_panorama_bgp":                                  resourcePanoramaBgp(),
			"panos_panorama_bgp_aggregate":                        resourcePanoramaBgpAggregate(),
			"panos_panorama_bgp_aggregate_advertise_filter":       resourcePanoramaBgpAggregateAdvertiseFilter(),
			"panos_panorama_bgp_aggregate_suppress_filter":        resourcePanoramaBgpAggregateSuppressFilter(),
			"panos_panorama_bgp_auth_profile":                     resourcePanoramaBgpAuthProfile(),
			"panos_panorama_bgp_conditional_adv":                  resourcePanoramaBgpConditionalAdv(),
			"panos_panorama_bgp_conditional_adv_advertise_filter": resourcePanoramaBgpConditionalAdvAdvertiseFilter(),
			"panos_panorama_bgp_conditional_adv_non_exist_filter": resourcePanoramaBgpConditionalAdvNonExistFilter(),
			"panos_panorama_bgp_dampening_profile":                resourcePanoramaBgpDampeningProfile(),
			"panos_panorama_bgp_export_rule_group":                resourcePanoramaBgpExportRuleGroup(),
			"panos_panorama_bgp_import_rule_group":                resourcePanoramaBgpImportRuleGroup(),
			"panos_panorama_bgp_peer":                             resourcePanoramaBgpPeer(),
			"panos_panorama_bgp_peer_group":                       resourcePanoramaBgpPeerGroup(),
			"panos_panorama_bgp_redist_rule":                      resourcePanoramaBgpRedistRule(),
			"panos_panorama_edl":                                  resourcePanoramaEdl(),
			"panos_panorama_email_server_profile":                 resourcePanoramaEmailServerProfile(),
			"panos_panorama_ethernet_interface":                   resourcePanoramaEthernetInterface(),
			"panos_panorama_gcp_account":                          resourcePanoramaGcpAccount(),
			"panos_panorama_gke_cluster":                          resourcePanoramaGkeCluster(),
			"panos_panorama_gke_cluster_group":                    resourcePanoramaGkeClusterGroup(),
			"panos_panorama_gre_tunnel":                           resourcePanoramaGreTunnel(),
			"panos_panorama_http_server_profile":                  resourcePanoramaHttpServerProfile(),
			"panos_panorama_ike_crypto_profile":                   resourceIkeCryptoProfile(),
			"panos_panorama_ike_gateway":                          resourcePanoramaIkeGateway(),
			"panos_panorama_ipsec_crypto_profile":                 resourcePanoramaIpsecCryptoProfile(),
			"panos_panorama_ipsec_tunnel":                         resourcePanoramaIpsecTunnel(),
			"panos_panorama_ipsec_tunnel_proxy_id_ipv4":           resourcePanoramaIpsecTunnelProxyIdIpv4(),
			"panos_panorama_layer2_subinterface":                  resourcePanoramaLayer2Subinterface(),
			"panos_panorama_layer3_subinterface":                  resourcePanoramaLayer3Subinterface(),
			"panos_panorama_log_forwarding_profile":               resourcePanoramaLogForwardingProfile(),
			"panos_panorama_loopback_interface":                   resourcePanoramaLoopbackInterface(),
			"panos_panorama_management_profile":                   resourcePanoramaManagementProfile(),
			"panos_panorama_monitor_profile":                      resourcePanoramaMonitorProfile(),
			"panos_panorama_nat_rule":                             resourcePanoramaNatRule(),
			"panos_panorama_nat_rule_group":                       resourcePanoramaNatRuleGroup(),
			"panos_panorama_pbf_rule_group":                       resourcePanoramaPbfRuleGroup(),
			"panos_panorama_redistribution_profile_ipv4":          resourcePanoramaRedistributionProfileIpv4(),
			"panos_panorama_security_policy":                      resourcePanoramaSecurityPolicy(),
			"panos_panorama_security_rule_group":                  resourcePanoramaSecurityRuleGroup(),
			"panos_panorama_service_group":                        resourcePanoramaServiceGroup(),
			"panos_panorama_service_object":                       resourcePanoramaServiceObject(),
			"panos_panorama_snmptrap_server_profile":              resourcePanoramaSnmptrapServerProfile(),
			"panos_panorama_static_route_ipv4":                    resourcePanoramaStaticRouteIpv4(),
			"panos_panorama_syslog_server_profile":                resourcePanoramaSyslogServerProfile(),
			"panos_panorama_template":                             resourcePanoramaTemplate(),
			"panos_panorama_template_entry":                       resourcePanoramaTemplateEntry(),
			"panos_panorama_template_stack":                       resourcePanoramaTemplateStack(),
			"panos_panorama_template_stack_entry":                 resourcePanoramaTemplateStackEntry(),
			"panos_panorama_template_variable":                    resourcePanoramaTemplateVariable(),
			"panos_panorama_tunnel_interface":                     resourcePanoramaTunnelInterface(),
			"panos_panorama_virtual_router":                       resourcePanoramaVirtualRouter(),
			"panos_panorama_virtual_router_entry":                 resourcePanoramaVirtualRouterEntry(),
			"panos_panorama_vlan":                                 resourcePanoramaVlan(),
			"panos_panorama_vlan_entry":                           resourcePanoramaVlanEntry(),
			"panos_panorama_vlan_interface":                       resourcePanoramaVlanInterface(),
			"panos_panorama_zone":                                 resourcePanoramaZone(),
			"panos_panorama_zone_entry":                           resourcePanoramaZoneEntry(),
			"panos_vm_auth_key":                                   resourceVmAuthKey(),

			// Panorama aliases.
			"panos_panorama_nat_policy":            resourcePanoramaNatRule(),
			"panos_panorama_security_policies":     resourcePanoramaSecurityPolicy(),
			"panos_panorama_security_policy_group": resourcePanoramaSecurityRuleGroup(),
			"panos_panorama_device_group":          resourceDeviceGroup(),
			"panos_panorama_device_group_entry":    resourceDeviceGroupEntry(),

			// Firewall resources.
			"panos_address_group":                        resourceAddressGroup(),
			"panos_administrative_tag":                   resourceAdministrativeTag(),
			"panos_aggregate_interface":                  resourceAggregateInterface(),
			"panos_application_group":                    resourceApplicationGroup(),
			"panos_application_object":                   resourceApplicationObject(),
			"panos_application_signature":                resourceApplicationSignature(),
			"panos_aws_cloud_watch":                      resourceAwsCloudWatch(),
			"panos_bfd_profile":                          resourceBfdProfile(),
			"panos_bgp":                                  resourceBgp(),
			"panos_bgp_aggregate":                        resourceBgpAggregate(),
			"panos_bgp_aggregate_advertise_filter":       resourceBgpAggregateAdvertiseFilter(),
			"panos_bgp_aggregate_suppress_filter":        resourceBgpAggregateSuppressFilter(),
			"panos_bgp_auth_profile":                     resourceBgpAuthProfile(),
			"panos_bgp_conditional_adv":                  resourceBgpConditionalAdv(),
			"panos_bgp_conditional_adv_advertise_filter": resourceBgpConditionalAdvAdvertiseFilter(),
			"panos_bgp_conditional_adv_non_exist_filter": resourceBgpConditionalAdvNonExistFilter(),
			"panos_bgp_dampening_profile":                resourceBgpDampeningProfile(),
			"panos_bgp_export_rule_group":                resourceBgpExportRuleGroup(),
			"panos_bgp_import_rule_group":                resourceBgpImportRuleGroup(),
			"panos_bgp_peer":                             resourceBgpPeer(),
			"panos_bgp_peer_group":                       resourceBgpPeerGroup(),
			"panos_bgp_redist_rule":                      resourceBgpRedistRule(),
			"panos_dag_tags":                             resourceDagTags(),
			"panos_edl":                                  resourceEdl(),
			"panos_email_server_profile":                 resourceEmailServerProfile(),
			"panos_ethernet_interface":                   resourceEthernetInterface(),
			"panos_general_settings":                     resourceGeneralSettings(),
			"panos_gre_tunnel":                           resourceGreTunnel(),
			"panos_http_server_profile":                  resourceHttpServerProfile(),
			"panos_ike_crypto_profile":                   resourceIkeCryptoProfile(),
			"panos_ike_gateway":                          resourceIkeGateway(),
			"panos_ip_tag":                               resourceIpTag(),
			"panos_ipsec_crypto_profile":                 resourceIpsecCryptoProfile(),
			"panos_ipsec_tunnel":                         resourceIpsecTunnel(),
			"panos_ipsec_tunnel_proxy_id_ipv4":           resourceIpsecTunnelProxyIdIpv4(),
			"panos_layer2_subinterface":                  resourceLayer2Subinterface(),
			"panos_layer3_subinterface":                  resourceLayer3Subinterface(),
			"panos_license_api_key":                      resourceLicenseApiKey(),
			"panos_licensing":                            resourceLicensing(),
			"panos_log_forwarding_profile":               resourceLogForwardingProfile(),
			"panos_loopback_interface":                   resourceLoopbackInterface(),
			"panos_management_profile":                   resourceManagementProfile(),
			"panos_monitor_profile":                      resourceMonitorProfile(),
			"panos_nat_rule":                             resourceNatRule(),
			"panos_nat_rule_group":                       resourceNatRuleGroup(),
			"panos_pbf_rule_group":                       resourcePbfRuleGroup(),
			"panos_redistribution_profile_ipv4":          resourceRedistributionProfileIpv4(),
			"panos_security_policy":                      resourceSecurityPolicy(),
			"panos_security_rule_group":                  resourceSecurityRuleGroup(),
			"panos_service_group":                        resourceServiceGroup(),
			"panos_service_object":                       resourceServiceObject(),
			"panos_snmptrap_server_profile":              resourceSnmptrapServerProfile(),
			"panos_static_route_ipv4":                    resourceStaticRouteIpv4(),
			"panos_syslog_server_profile":                resourceSyslogServerProfile(),
			"panos_telemetry":                            resourceTelemetry(),
			"panos_tunnel_interface":                     resourceTunnelInterface(),
			"panos_user_tag":                             resourceUserTag(),
			"panos_userid_login":                         resourceUseridLogin(),
			"panos_virtual_router":                       resourceVirtualRouter(),
			"panos_virtual_router_entry":                 resourceVirtualRouterEntry(),
			"panos_vlan":                                 resourceVlan(),
			"panos_vlan_entry":                           resourceVlanEntry(),
			"panos_vlan_interface":                       resourceVlanInterface(),
			"panos_zone":                                 resourceZone(),
			"panos_zone_entry":                           resourceZoneEntry(),

			// Firewall aliases.
			"panos_nat_policy":            resourceNatRule(),
			"panos_security_policies":     resourceSecurityPolicy(),
			"panos_security_policy_group": resourceSecurityRuleGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		logging uint32
		err     error
	)

	lm := map[string]uint32{
		"quiet":                   pango.LogQuiet,
		"action":                  pango.LogAction,
		"query":                   pango.LogQuery,
		"op":                      pango.LogOp,
		"uid":                     pango.LogUid,
		"log":                     pango.LogLog,
		"export":                  pango.LogExport,
		"import":                  pango.LogImport,
		"xpath":                   pango.LogXpath,
		"send":                    pango.LogSend,
		"receive":                 pango.LogReceive,
		"osx_curl":                pango.LogOsxCurl,
		"curl_with_personal_data": pango.LogCurlWithPersonalData,
	}

	var hdrs map[string]string
	hconfig := d.Get("additional_headers").(map[string]interface{})
	if len(hconfig) > 0 {
		hdrs = make(map[string]string)
		for key, val := range hconfig {
			hdrs[key] = val.(string)
		}
	}

	if ll := d.Get("logging").([]interface{}); len(ll) > 0 {
		for i := range ll {
			s := ll[i].(string)
			if v, ok := lm[s]; !ok {
				return nil, fmt.Errorf("Unknown logging artifact requested: %s", s)
			} else {
				logging |= v
			}
		}
	}

	con, err := pango.ConnectUsing(
		pango.Client{
			Hostname:          d.Get("hostname").(string),
			Username:          d.Get("username").(string),
			Password:          d.Get("password").(string),
			ApiKey:            d.Get("api_key").(string),
			Protocol:          d.Get("protocol").(string),
			Port:              uint(d.Get("port").(int)),
			Timeout:           d.Get("timeout").(int),
			Target:            d.Get("target").(string),
			Headers:           hdrs,
			Logging:           logging,
			VerifyCertificate: d.Get("verify_certificate").(bool),
		},
		d.Get("json_config_file").(string),
		true,
	)

	if err != nil {
		return nil, err
	}

	return con, nil
}
