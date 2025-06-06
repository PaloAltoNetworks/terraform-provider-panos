package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAdminRoleDevice(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: adminRoleResource1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-admin-role", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("cli"),
						knownvalue.StringExact("superuser"),
					),
					// device.restapi.device
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("email_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("http_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("ldap_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("snmp_trap_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("email_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("syslog_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("virtual_systems"),
						knownvalue.StringExact("read-only"),
					),
					// device.restapi.network
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("aggregate_ethernet_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("bfd_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("bgp_routing_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("dhcp_relays"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("dhcp_servers"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("dns_proxies"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("ethernet_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_clientless_app_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_clientless_apps"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_gateways"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_ipsec_crypto_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_mdm_servers"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_portals"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("gre_tunnels"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("ike_crypto_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("ike_gateway_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("interface_management_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("ipsec_crypto_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("ipsec_tunnels"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("lldp"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("lldp_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("logical_routers"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("loopback_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("qos_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("qos_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("sdwan_interface_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("sdwan_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("tunnel_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("tunnel_monitor_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("virtual_routers"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("virtual_wires"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("vlan_interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("vlans"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("zone_protection_network_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("zones"),
						knownvalue.StringExact("read-only"),
					),
					// device.restapi.objects
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("address_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("addresses"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("anti_spyware_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("antivirus_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("application_filters"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("application_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("applications"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("authentication_enforcements"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_data_patterns"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_spyware_signatures"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_url_categories"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_vulnerability_signatures"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("data_filtering_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("decryption_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("devices"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("dos_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("dynamic_user_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("external_dynamic_lists"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("file_blocking_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("globalprotect_hip_objects"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("globalprotect_hip_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("gtp_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("log_forwarding_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("packet_broker_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("regions"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("schedules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sctp_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_error_correction_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_path_quality_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_saas_quality_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_traffic_distribution_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("security_profile_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("service_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("services"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("tags"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("url_filtering_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("vulnerability_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("wildfire_analysis_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					// device.restapi.policies
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("application_override_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("authentication_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("decryption_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("dos_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("nat_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("network_packet_broker_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("policy_based_forwarding_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("qos_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("sdwan_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("security_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("tunnel_inspection_rules"),
						knownvalue.StringExact("read-only"),
					),
					// device.restapi.system
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("restapi").
							AtMapKey("system").
							AtMapKey("configuration"),
						knownvalue.StringExact("read-only"),
					),
					// device.restapi.webui
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("acc"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("dashboard"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("tasks"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("validate"),
						knownvalue.StringExact("enable"),
					),
					// device.restapi.webui.commit
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("commit").
							AtMapKey("commit_for_other_admins"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("commit").
							AtMapKey("device"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("commit").
							AtMapKey("object_level_changes"),
						knownvalue.StringExact("enable"),
					),
					// device.restapi.webui.device
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("admin_roles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("administrators"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("authentication_profile"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("authentication_sequence"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("block_pages"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("config_audit"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("data_redistribution"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("device_quarantine"),
						knownvalue.StringExact("read-only"),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_admin_role.role",
					// 	tfjsonpath.New("role").
					// 		AtMapKey("device").
					// 		AtMapKey("webui").
					// 		AtMapKey("device").
					// 		AtMapKey("dhcp_syslog_server"),
					// 	knownvalue.StringExact("read-only"),
					// ),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("dynamic_updates"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("global_protect_client"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("high_availability"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("licenses"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_fwd_card"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("master_key"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("plugins"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("scheduled_log_export"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("shared_gateways"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("software"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("support"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("troubleshooting"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("user_identification"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("virtual_systems"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("vm_info_source"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.device.webui.device.certificate_management
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("certificate_profile"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("certificates"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ocsp_responder"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("scep"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ssh_service_profile"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ssl_decryption_exclusion"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ssl_tls_service_profile"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.device.webui.device.local_user_database
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("local_user_database").
							AtMapKey("user_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("local_user_database").
							AtMapKey("users"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.device.webui.device.log_settings
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("cc_alarm"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("config"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("correlation"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("globalprotect"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("hipmatch"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("iptag"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("manage_log"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("system"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("user_id"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.device.webui.device.policy_recommendations
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("policy_recommendations").
							AtMapKey("iot"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("policy_recommendations").
							AtMapKey("saas"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.device.webui.device.server_profile
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("dns"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("email"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("http"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("kerberos"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("ldap"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("mfa"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("netflow"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("radius"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("saml_idp"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("scp"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("snmp_trap"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("syslog"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("tacplus"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.device.webui.device.setup
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("content_id"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("hsm"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("management"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("operations"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("services"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("session"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("telemetry"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("wildfire"),
						knownvalue.StringExact("read-only"),
					),
					// device.webui.global
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("global").
							AtMapKey("system_alarms"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.monitor
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("app_scope"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("application_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("block_ip_list"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("botnet"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("external_logs"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("gtp_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("packet_capture"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("sctp_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("session_browser"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("threat_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("traffic_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("url_filtering_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("view_custom_reports"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.monitor.automated_correlation_engine
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("automated_correlation_engine").
							AtMapKey("correlated_events"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("automated_correlation_engine").
							AtMapKey("correlation_objects"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.monitor.custom_reports
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("application_statistics"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("auth"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("data_filtering_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("decryption_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("decryption_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("globalprotect"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("gtp_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("gtp_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("hipmatch"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("iptag"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("sctp_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("sctp_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("threat_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("threat_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("traffic_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("traffic_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("tunnel_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("tunnel_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("url_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("url_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("userid"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("wildfire_log"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.monitor.logs
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("alarm"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("authentication"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("configuration"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("data_filtering"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("decryption"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("globalprotect"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("gtp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("hipmatch"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("iptag"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("sctp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("system"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("threat"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("traffic"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("tunnel"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("url"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("userid"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("wildfire"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.monitor.pdf_reports
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("email_scheduler"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("manage_pdf_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("pdf_summary_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("report_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("saas_application_usage_report"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("user_activity_report"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.network
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("dhcp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("dns_proxy"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("gre_tunnels"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("interfaces"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("ipsec_tunnels"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("lldp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("qos"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("sdwan_interface_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("secure_web_gateway"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("virtual_routers"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("virtual_wires"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("vlans"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("zones"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.network.global_protect
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("clientless_app_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("clientless_apps"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("gateways"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("mdm"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("portals"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.network.network_profiles
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("bfd_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("gp_app_ipsec_crypto"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("ike_crypto"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("ike_gateways"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("interface_mgmt"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("ipsec_crypto"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("lldp_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("qos_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("tunnel_monitor"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("network_profiles").
							AtMapKey("zone_protection"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.network.routing
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("logical_routers"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.network.routing.routing_profiles
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("bfd"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("bgp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("filters"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("multicast"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("ospf"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("ospfv3"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("routing").
							AtMapKey("routing_profiles").
							AtMapKey("ripv2"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.objects
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("address_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("addresses"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("application_filters"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("application_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("applications"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("authentication"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("devices"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("dynamic_block_lists"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("dynamic_user_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("log_forwarding"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("packet_broker_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("regions"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("schedules"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profile_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("service_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("services"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("tags"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.objects.custom_objects
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("data_patterns"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("spyware"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("url_category"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("vulnerability"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.objects.decryption
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("decryption").
							AtMapKey("decryption_profile"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.objects.global_protect
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("global_protect").
							AtMapKey("hip_objects"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("global_protect").
							AtMapKey("hip_profiles"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.objects.sdwan
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_dist_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_error_correction_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_saas_quality_profile"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.objects.security_profiles
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("anti_spyware"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("antivirus"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("data_filtering"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("dos_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("file_blocking"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("gtp_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("sctp_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("url_filtering"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("vulnerability_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("wildfire_analysis"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.operations
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("download_core_files"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("download_pcap_files"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("generate_stats_dump_file"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("generate_tech_support_file"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("reboot"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.policies
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("application_override_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("authentication_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("dos_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("nat_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("network_packet_broker_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("pbf_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("qos_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("rule_hit_count_reset"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("sdwan_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("security_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("ssl_decryption_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("tunnel_inspect_rulebase"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.privacy
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("privacy").
							AtMapKey("show_full_ip_addresses"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("privacy").
							AtMapKey("show_user_names_in_logs_and_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("privacy").
							AtMapKey("view_pcap_files"),
						knownvalue.StringExact("disable"),
					),
					// device.webui.save
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("save").
							AtMapKey("object_level_changes"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("save").
							AtMapKey("partial_save"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("webui").
							AtMapKey("save").
							AtMapKey("save_for_other_admins"),
						knownvalue.StringExact("disable"),
					),
					// device.xmlapi
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("commit"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("config"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("export"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("import"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("iot"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("op"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("report"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("device").
							AtMapKey("xmlapi").
							AtMapKey("user_id"),
						knownvalue.StringExact("disable"),
					),
				},
			},
		},
	})
}

func TestAccAdminRoleVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: adminRoleResource2,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-admin-role", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("cli"),
						knownvalue.StringExact("vsysadmin"),
					),
					// vsys.restapi.device
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("email_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("http_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("ldap_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("snmp_trap_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("email_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("syslog_server_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("device").
							AtMapKey("virtual_systems"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.restapi.network
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_clientless_app_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_clientless_apps"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_gateways"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_mdm_servers"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("globalprotect_portals"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("sdwan_interface_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("network").
							AtMapKey("zones"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.restapi.objects
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("address_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("addresses"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("anti_spyware_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("antivirus_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("application_filters"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("application_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("applications"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("authentication_enforcements"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_data_patterns"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_spyware_signatures"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_url_categories"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("custom_vulnerability_signatures"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("data_filtering_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("decryption_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("devices"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("dos_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("dynamic_user_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("external_dynamic_lists"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("file_blocking_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("globalprotect_hip_objects"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("globalprotect_hip_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("gtp_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("log_forwarding_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("packet_broker_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("regions"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("schedules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sctp_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_error_correction_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_path_quality_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_saas_quality_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("sdwan_traffic_distribution_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("security_profile_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("service_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("services"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("tags"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("url_filtering_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("vulnerability_protection_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("objects").
							AtMapKey("wildfire_analysis_security_profiles"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.restapi.policies
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("application_override_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("authentication_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("decryption_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("dos_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("nat_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("network_packet_broker_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("policy_based_forwarding_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("qos_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("sdwan_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("security_rules"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("policies").
							AtMapKey("tunnel_inspection_rules"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.restapi.system
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("restapi").
							AtMapKey("system").
							AtMapKey("configuration"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.restapi.webui
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("acc"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("dashboard"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("tasks"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("validate"),
						knownvalue.StringExact("enable"),
					),
					// vsys.restapi.webui.commit
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("commit").
							AtMapKey("commit_for_other_admins"),
						knownvalue.StringExact("enable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("commit").
							AtMapKey("virtual_systems"),
						knownvalue.StringExact("disable"),
					),
					// vsys.restapi.webui.device
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("administrators"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("authentication_profile"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("authentication_sequence"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("block_pages"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("data_redistribution"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("device_quarantine"),
						knownvalue.StringExact("read-only"),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_admin_role.role",
					// 	tfjsonpath.New("role").
					// 		AtMapKey("vsys").
					// 		AtMapKey("webui").
					// 		AtMapKey("device").
					// 		AtMapKey("dhcp_syslog_server"),
					// 	knownvalue.StringExact("read-only"),
					// ),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("troubleshooting"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("user_identification"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("vm_info_source"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.device.webui.device.certificate_management
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("certificate_profile"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("certificates"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ocsp_responder"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("scep"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ssh_service_profile"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ssl_decryption_exclusion"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("certificate_management").
							AtMapKey("ssl_tls_service_profile"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.device.webui.device.local_user_database
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("local_user_database").
							AtMapKey("user_groups"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("local_user_database").
							AtMapKey("users"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.device.webui.device.log_settings
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("config"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("correlation"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("globalprotect"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("hipmatch"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("iptag"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("system"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("log_settings").
							AtMapKey("user_id"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.device.webui.device.policy_recommendations
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("policy_recommendations").
							AtMapKey("iot"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("policy_recommendations").
							AtMapKey("saas"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.device.webui.device.server_profile
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("dns"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("email"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("http"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("kerberos"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("ldap"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("mfa"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("netflow"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("radius"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("saml_idp"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("scp"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("snmp_trap"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("syslog"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("server_profile").
							AtMapKey("tacplus"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.device.webui.device.setup
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("content_id"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("hsm"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("interfaces"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("management"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("operations"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("services"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("session"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("telemetry"),
						knownvalue.StringExact("read-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("device").
							AtMapKey("setup").
							AtMapKey("wildfire"),
						knownvalue.StringExact("read-only"),
					),
					// vsys.webui.monitor
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("app_scope"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("block_ip_list"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("external_logs"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("session_browser"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("view_custom_reports"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.monitor.automated_correlation_engine
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("automated_correlation_engine").
							AtMapKey("correlated_events"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("automated_correlation_engine").
							AtMapKey("correlation_objects"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.monitor.custom_reports
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("application_statistics"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("auth"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("data_filtering_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("decryption_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("decryption_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("globalprotect"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("gtp_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("gtp_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("hipmatch"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("iptag"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("sctp_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("sctp_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("threat_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("threat_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("traffic_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("traffic_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("tunnel_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("tunnel_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("url_log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("url_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("userid"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("custom_reports").
							AtMapKey("wildfire_log"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.monitor.logs
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("authentication"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("data_filtering"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("decryption"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("globalprotect"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("gtp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("hipmatch"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("iptag"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("sctp"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("threat"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("traffic"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("tunnel"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("url"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("userid"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("logs").
							AtMapKey("wildfire"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.monitor.pdf_reports
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("email_scheduler"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("manage_pdf_summary"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("pdf_summary_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("report_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("saas_application_usage_report"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("monitor").
							AtMapKey("pdf_reports").
							AtMapKey("user_activity_report"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.network
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("sdwan_interface_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("zones"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.network.global_protect
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("clientless_app_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("clientless_apps"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("gateways"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("mdm"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("network").
							AtMapKey("global_protect").
							AtMapKey("portals"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.objects
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("address_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("addresses"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("application_filters"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("application_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("applications"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("authentication"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("devices"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("dynamic_block_lists"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("dynamic_user_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("log_forwarding"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("packet_broker_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("regions"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("schedules"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profile_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("service_groups"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("services"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("tags"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.objects.custom_objects
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("data_patterns"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("spyware"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("url_category"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("custom_objects").
							AtMapKey("vulnerability"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.objects.decryption
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("decryption").
							AtMapKey("decryption_profile"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.objects.global_protect
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("global_protect").
							AtMapKey("hip_objects"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("global_protect").
							AtMapKey("hip_profiles"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.objects.sdwan
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_dist_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_error_correction_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_profile"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("sdwan").
							AtMapKey("sdwan_saas_quality_profile"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.objects.security_profiles
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("anti_spyware"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("antivirus"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("data_filtering"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("dos_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("file_blocking"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("gtp_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("sctp_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("url_filtering"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("vulnerability_protection"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("objects").
							AtMapKey("security_profiles").
							AtMapKey("wildfire_analysis"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.operations
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("download_core_files"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("download_pcap_files"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("generate_stats_dump_file"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("generate_tech_support_file"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("operations").
							AtMapKey("reboot"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.policies
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("application_override_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("authentication_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("dos_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("nat_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("network_packet_broker_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("pbf_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("qos_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("rule_hit_count_reset"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("sdwan_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("security_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("ssl_decryption_rulebase"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("policies").
							AtMapKey("tunnel_inspect_rulebase"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.privacy
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("privacy").
							AtMapKey("show_full_ip_addresses"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("privacy").
							AtMapKey("show_user_names_in_logs_and_reports"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("privacy").
							AtMapKey("view_pcap_files"),
						knownvalue.StringExact("disable"),
					),
					// vsys.webui.save
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("save").
							AtMapKey("object_level_changes"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("save").
							AtMapKey("partial_save"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("webui").
							AtMapKey("save").
							AtMapKey("save_for_other_admins"),
						knownvalue.StringExact("disable"),
					),
					// vsys.xmlapi
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("commit"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("config"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("export"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("import"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("iot"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("log"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("op"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("report"),
						knownvalue.StringExact("disable"),
					),
					statecheck.ExpectKnownValue(
						"panos_admin_role.role",
						tfjsonpath.New("role").
							AtMapKey("vsys").
							AtMapKey("xmlapi").
							AtMapKey("user_id"),
						knownvalue.StringExact("disable"),
					),
				},
			},
		},
	})
}

const adminRoleResource1 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_admin_role" "role" {
  location = { template = { name = panos_template.template.name } }

  name = format("%s-admin-role", var.prefix)

  description = "admin role description"

  role = {
    device = {
      cli = "superuser"
      restapi = {
        device = {
          email_server_profiles = "read-only"
          http_server_profiles = "read-only"
          ldap_server_profiles = "read-only"
          log_interface_setting = "read-only"
          snmp_trap_server_profiles = "read-only"
          syslog_server_profiles = "read-only"
          virtual_systems = "read-only"
        }
        network = {
          aggregate_ethernet_interfaces = "read-only"
          bfd_network_profiles = "read-only"
          bgp_routing_profiles = "read-only"
          dhcp_relays = "read-only"
          dhcp_servers = "read-only"
          dns_proxies = "read-only"
          ethernet_interfaces = "read-only"
          globalprotect_clientless_app_groups = "read-only"
          globalprotect_clientless_apps = "read-only"
          globalprotect_gateways = "read-only"
          globalprotect_ipsec_crypto_network_profiles = "read-only"
          globalprotect_mdm_servers = "read-only"
          globalprotect_portals = "read-only"
          gre_tunnels = "read-only"
          ike_crypto_network_profiles = "read-only"
          ike_gateway_network_profiles = "read-only"
          interface_management_network_profiles = "read-only"
          ipsec_crypto_network_profiles = "read-only"
          ipsec_tunnels = "read-only"
          lldp = "read-only"
          lldp_network_profiles = "read-only"
          logical_routers = "read-only"
          loopback_interfaces = "read-only"
          qos_interfaces = "read-only"
          qos_network_profiles = "read-only"
          sdwan_interface_profiles = "read-only"
          sdwan_interfaces = "read-only"
          tunnel_interfaces = "read-only"
          tunnel_monitor_network_profiles = "read-only"
          virtual_routers = "read-only"
          virtual_wires = "read-only"
          vlan_interfaces = "read-only"
          vlans = "read-only"
          zone_protection_network_profiles = "read-only"
          zones = "read-only"
        }
        objects = {
          address_groups = "read-only"
          addresses = "read-only"
          anti_spyware_security_profiles = "read-only"
          antivirus_security_profiles = "read-only"
          application_filters = "read-only"
          application_groups = "read-only"
          applications = "read-only"
          authentication_enforcements = "read-only"
          custom_data_patterns = "read-only"
          custom_spyware_signatures = "read-only"
          custom_url_categories = "read-only"
          custom_vulnerability_signatures = "read-only"
          data_filtering_security_profiles = "read-only"
          decryption_profiles = "read-only"
          devices = "read-only"
          dos_protection_security_profiles = "read-only"
          dynamic_user_groups = "read-only"
          external_dynamic_lists = "read-only"
          file_blocking_security_profiles = "read-only"
          globalprotect_hip_objects = "read-only"
          globalprotect_hip_profiles = "read-only"
          gtp_protection_security_profiles = "read-only"
          log_forwarding_profiles = "read-only"
          packet_broker_profiles = "read-only"
          regions = "read-only"
          schedules = "read-only"
          sctp_protection_security_profiles = "read-only"
          sdwan_error_correction_profiles = "read-only"
          sdwan_path_quality_profiles = "read-only"
          sdwan_saas_quality_profiles = "read-only"
          sdwan_traffic_distribution_profiles = "read-only"
          security_profile_groups = "read-only"
          service_groups = "read-only"
          services = "read-only"
          tags = "read-only"
          url_filtering_security_profiles = "read-only"
          vulnerability_protection_security_profiles = "read-only"
          wildfire_analysis_security_profiles = "read-only"
        }
        policies = {
          application_override_rules = "read-only"
          authentication_rules = "read-only"
          decryption_rules = "read-only"
          dos_rules = "read-only"
          nat_rules = "read-only"
          network_packet_broker_rules = "read-only"
          policy_based_forwarding_rules = "read-only"
          qos_rules = "read-only"
          sdwan_rules = "read-only"
          security_rules = "read-only"
          tunnel_inspection_rules = "read-only"
        }
        system = {
          configuration = "read-only"
        }
      }
      webui = {
        acc = "enable"
        dashboard = "enable"
        tasks = "enable"
        validate = "enable"
        commit = {
          commit_for_other_admins = "enable"
          device = "disable"
          object_level_changes = "enable"
        }
        device = {
          access_domain = "read-only"
          admin_roles = "read-only"
          administrators = "read-only"
          authentication_profile = "read-only"
          authentication_sequence = "read-only"
          block_pages = "read-only"
          config_audit = "enable"
          data_redistribution = "read-only"
          device_quarantine = "read-only"
          # dhcp_syslog_server = "read-only"
          dynamic_updates = "read-only"
          global_protect_client = "read-only"
          high_availability = "read-only"
          licenses = "read-only"
          log_fwd_card = "read-only"
          master_key = "read-only"
          plugins = "disable"
          scheduled_log_export = "enable"
          shared_gateways = "read-only"
          software = "read-only"
          support = "read-only"
          troubleshooting = "read-only"
          user_identification = "read-only"
          virtual_systems = "read-only"
          vm_info_source = "read-only"
          certificate_management = {
            certificate_profile = "read-only"
            certificates = "read-only"
            ocsp_responder = "read-only"
            scep = "read-only"
            ssh_service_profile = "read-only"
            ssl_decryption_exclusion = "read-only"
            ssl_tls_service_profile = "read-only"
          }
          local_user_database = {
            user_groups = "read-only"
            users = "read-only"
          }
          log_settings = {
            cc_alarm = "read-only"
            config = "read-only"
            correlation = "read-only"
            globalprotect = "read-only"
            hipmatch = "read-only"
            iptag = "read-only"
            manage_log = "read-only"
            system = "read-only"
            user_id = "read-only"
          }
          policy_recommendations = {
            iot = "read-only"
            saas = "read-only"
          }
          server_profile = {
            dns = "read-only"
            email = "read-only"
            http = "read-only"
            kerberos = "read-only"
            ldap = "read-only"
            mfa = "read-only"
            netflow = "read-only"
            radius = "read-only"
            saml_idp = "read-only"
            scp = "read-only"
            snmp_trap = "read-only"
            syslog = "read-only"
            tacplus = "read-only"
          }
          setup = {
            content_id = "read-only"
            hsm = "read-only"
            interfaces = "read-only"
            management = "read-only"
            operations = "read-only"
            services = "read-only"
            session = "read-only"
            telemetry = "read-only"
            wildfire = "read-only"
          }
        }
        global = {
          system_alarms = "disable"
        }
        monitor = {
          app_scope = "disable"
          application_reports = "disable"
          block_ip_list = "disable"
          botnet = "disable"
          external_logs = "disable"
          gtp_reports = "disable"
          packet_capture = "disable"
          sctp_reports = "disable"
          session_browser = "disable"
          threat_reports = "disable"
          traffic_reports = "disable"
          url_filtering_reports = "disable"
          view_custom_reports = "disable"
          automated_correlation_engine = {
            correlated_events = "disable"
            correlation_objects = "disable"
          }
          custom_reports = {
            application_statistics = "disable"
            auth = "disable"
            data_filtering_log = "disable"
            decryption_log = "disable"
            decryption_summary = "disable"
            globalprotect = "disable"
            gtp_log = "disable"
            gtp_summary = "disable"
            hipmatch = "disable"
            iptag = "disable"
            sctp_log = "disable"
            sctp_summary = "disable"
            threat_log = "disable"
            threat_summary = "disable"
            traffic_log = "disable"
            traffic_summary = "disable"
            tunnel_log = "disable"
            tunnel_summary = "disable"
            url_log = "disable"
            url_summary = "disable"
            userid = "disable"
            wildfire_log = "disable"
          }
          logs = {
            alarm = "disable"
            authentication = "disable"
            configuration = "disable"
            data_filtering = "disable"
            decryption = "disable"
            globalprotect = "disable"
            gtp = "disable"
            hipmatch = "disable"
            iptag = "disable"
            sctp = "disable"
            system = "disable"
            threat = "disable"
            traffic = "disable"
            tunnel = "disable"
            url = "disable"
            userid = "disable"
            wildfire = "disable"
          }
          pdf_reports = {
            email_scheduler = "disable"
            manage_pdf_summary = "disable"
            pdf_summary_reports = "disable"
            report_groups = "disable"
            saas_application_usage_report = "disable"
            user_activity_report = "disable"
          }
        }
        network = {
          dhcp = "disable"
          dns_proxy = "disable"
          gre_tunnels = "disable"
          interfaces = "disable"
          ipsec_tunnels = "disable"
          lldp = "disable"
          qos = "disable"
          sdwan_interface_profile = "disable"
          secure_web_gateway = "disable"
          virtual_routers = "disable"
          virtual_wires = "disable"
          vlans = "disable"
          zones = "disable"
          global_protect = {
            clientless_app_groups = "disable"
            clientless_apps = "disable"
            gateways = "disable"
            mdm = "disable"
            portals = "disable"
          }
          network_profiles = {
            bfd_profile = "disable"
            gp_app_ipsec_crypto = "disable"
            ike_crypto = "disable"
            ike_gateways = "disable"
            interface_mgmt = "disable"
            ipsec_crypto = "disable"
            lldp_profile = "disable"
            qos_profile = "disable"
            tunnel_monitor = "disable"
            zone_protection = "disable"
          }
          routing = {
            logical_routers = "disable"
            routing_profiles = {
              bfd = "disable"
              bgp = "disable"
              filters = "disable"
              multicast = "disable"
              ospf = "disable"
              ospfv3 = "disable"
              ripv2 = "disable"
            }
          }
        }
        objects = {
          address_groups = "disable"
          addresses = "disable"
          application_filters = "disable"
          application_groups = "disable"
          applications = "disable"
          authentication = "disable"
          devices = "disable"
          dynamic_block_lists = "disable"
          dynamic_user_groups = "disable"
          log_forwarding = "disable"
          packet_broker_profile = "disable"
          regions = "disable"
          schedules = "disable"
          security_profile_groups = "disable"
          service_groups = "disable"
          services = "disable"
          tags = "disable"
          custom_objects = {
            data_patterns = "disable"
            spyware = "disable"
            url_category = "disable"
            vulnerability = "disable"
          }
          decryption = {
            decryption_profile = "disable"
          }
          global_protect = {
            hip_objects = "disable"
            hip_profiles = "disable"
          }
          sdwan = {
            sdwan_dist_profile = "disable"
            sdwan_error_correction_profile = "disable"
            sdwan_profile = "disable"
            sdwan_saas_quality_profile = "disable"
          }
          security_profiles = {
            anti_spyware = "disable"
            antivirus = "disable"
            data_filtering = "disable"
            dos_protection = "disable"
            file_blocking = "disable"
            gtp_protection = "disable"
            sctp_protection = "disable"
            url_filtering = "disable"
            vulnerability_protection = "disable"
            wildfire_analysis = "disable"
          }
        }
        operations = {
          download_core_files = "disable"
          download_pcap_files = "disable"
          generate_stats_dump_file = "disable"
          generate_tech_support_file = "disable"
          reboot = "disable"
        }
        policies = {
          application_override_rulebase = "disable"
          authentication_rulebase = "disable"
          dos_rulebase = "disable"
          nat_rulebase = "disable"
          network_packet_broker_rulebase = "disable"
          pbf_rulebase = "disable"
          qos_rulebase = "disable"
          rule_hit_count_reset = "disable"
          sdwan_rulebase = "disable"
          security_rulebase = "disable"
          ssl_decryption_rulebase = "disable"
          tunnel_inspect_rulebase = "disable"
        }
        privacy = {
          show_full_ip_addresses = "disable"
          show_user_names_in_logs_and_reports = "disable"
          view_pcap_files = "disable"
        }
        save = {
          object_level_changes = "disable"
          partial_save = "disable"
          save_for_other_admins = "disable"
        }
      }
      xmlapi = {
        commit = "disable"
        config = "disable"
        export = "disable"
        import = "disable"
        iot = "disable"
        log = "disable"
        op = "disable"
        report = "disable"
        user_id = "disable"
      }
    }
  }
}
`

const adminRoleResource2 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_admin_role" "role" {
  location = { template = { name = panos_template.template.name } }

  name = format("%s-admin-role", var.prefix)

  description = "admin role description"

  role = {
    vsys = {
      cli = "vsysadmin"
      restapi = {
        device = {
          email_server_profiles = "read-only"
          http_server_profiles = "read-only"
          ldap_server_profiles = "read-only"
          log_interface_setting = "read-only"
          snmp_trap_server_profiles = "read-only"
          syslog_server_profiles = "read-only"
          virtual_systems = "read-only"
        }
        network = {
          globalprotect_clientless_app_groups = "read-only"
          globalprotect_clientless_apps = "read-only"
          globalprotect_gateways = "read-only"
          globalprotect_mdm_servers = "read-only"
          globalprotect_portals = "read-only"
          sdwan_interface_profiles = "read-only"
          zones = "read-only"
        }
        objects = {
          address_groups = "read-only"
          addresses = "read-only"
          anti_spyware_security_profiles = "read-only"
          antivirus_security_profiles = "read-only"
          application_filters = "read-only"
          application_groups = "read-only"
          applications = "read-only"
          authentication_enforcements = "read-only"
          custom_data_patterns = "read-only"
          custom_spyware_signatures = "read-only"
          custom_url_categories = "read-only"
          custom_vulnerability_signatures = "read-only"
          data_filtering_security_profiles = "read-only"
          decryption_profiles = "read-only"
          devices = "read-only"
          dos_protection_security_profiles = "read-only"
          dynamic_user_groups = "read-only"
          external_dynamic_lists = "read-only"
          file_blocking_security_profiles = "read-only"
          globalprotect_hip_objects = "read-only"
          globalprotect_hip_profiles = "read-only"
          gtp_protection_security_profiles = "read-only"
          log_forwarding_profiles = "read-only"
          packet_broker_profiles = "read-only"
          regions = "read-only"
          schedules = "read-only"
          sctp_protection_security_profiles = "read-only"
          sdwan_error_correction_profiles = "read-only"
          sdwan_path_quality_profiles = "read-only"
          sdwan_saas_quality_profiles = "read-only"
          sdwan_traffic_distribution_profiles = "read-only"
          security_profile_groups = "read-only"
          service_groups = "read-only"
          services = "read-only"
          tags = "read-only"
          url_filtering_security_profiles = "read-only"
          vulnerability_protection_security_profiles = "read-only"
          wildfire_analysis_security_profiles = "read-only"
        }
        policies = {
          application_override_rules = "read-only"
          authentication_rules = "read-only"
          decryption_rules = "read-only"
          dos_rules = "read-only"
          nat_rules = "read-only"
          network_packet_broker_rules = "read-only"
          policy_based_forwarding_rules = "read-only"
          qos_rules = "read-only"
          sdwan_rules = "read-only"
          security_rules = "read-only"
          tunnel_inspection_rules = "read-only"
        }
        system = {
          configuration = "read-only"
        }
      }
      webui = {
        acc = "enable"
        dashboard = "enable"
        tasks = "enable"
        validate = "enable"
        commit = {
          commit_for_other_admins = "enable"
          virtual_systems = "disable"
        }
        device = {
          administrators = "read-only"
          authentication_profile = "read-only"
          authentication_sequence = "read-only"
          block_pages = "read-only"
          data_redistribution = "read-only"
          device_quarantine = "read-only"
          # dhcp_syslog_server = "read-only"
          troubleshooting = "read-only"
          user_identification = "read-only"
          vm_info_source = "read-only"
          certificate_management = {
            certificate_profile = "read-only"
            certificates = "read-only"
            ocsp_responder = "read-only"
            scep = "read-only"
            ssh_service_profile = "read-only"
            ssl_decryption_exclusion = "read-only"
            ssl_tls_service_profile = "read-only"
          }
          local_user_database = {
            user_groups = "read-only"
            users = "read-only"
          }
          log_settings = {
            config = "read-only"
            correlation = "read-only"
            globalprotect = "read-only"
            hipmatch = "read-only"
            iptag = "read-only"
            system = "read-only"
            user_id = "read-only"
          }
          policy_recommendations = {
            iot = "read-only"
            saas = "read-only"
          }
          server_profile = {
            dns = "read-only"
            email = "read-only"
            http = "read-only"
            kerberos = "read-only"
            ldap = "read-only"
            mfa = "read-only"
            netflow = "read-only"
            radius = "read-only"
            saml_idp = "read-only"
            scp = "read-only"
            snmp_trap = "read-only"
            syslog = "read-only"
            tacplus = "read-only"
          }
          setup = {
            content_id = "read-only"
            hsm = "read-only"
            interfaces = "read-only"
            management = "read-only"
            operations = "read-only"
            services = "read-only"
            session = "read-only"
            telemetry = "read-only"
            wildfire = "read-only"
          }
        }
        monitor = {
          app_scope = "disable"
          block_ip_list = "disable"
          external_logs = "disable"
          session_browser = "disable"
          view_custom_reports = "disable"
          automated_correlation_engine = {
            correlated_events = "disable"
            correlation_objects = "disable"
          }
          custom_reports = {
            application_statistics = "disable"
            auth = "disable"
            data_filtering_log = "disable"
            decryption_log = "disable"
            decryption_summary = "disable"
            globalprotect = "disable"
            gtp_log = "disable"
            gtp_summary = "disable"
            hipmatch = "disable"
            iptag = "disable"
            sctp_log = "disable"
            sctp_summary = "disable"
            threat_log = "disable"
            threat_summary = "disable"
            traffic_log = "disable"
            traffic_summary = "disable"
            tunnel_log = "disable"
            tunnel_summary = "disable"
            url_log = "disable"
            url_summary = "disable"
            userid = "disable"
            wildfire_log = "disable"
          }
          logs = {
            authentication = "disable"
            data_filtering = "disable"
            decryption = "disable"
            globalprotect = "disable"
            gtp = "disable"
            hipmatch = "disable"
            iptag = "disable"
            sctp = "disable"
            threat = "disable"
            traffic = "disable"
            tunnel = "disable"
            url = "disable"
            userid = "disable"
            wildfire = "disable"
          }
          pdf_reports = {
            email_scheduler = "disable"
            manage_pdf_summary = "disable"
            pdf_summary_reports = "disable"
            report_groups = "disable"
            saas_application_usage_report = "disable"
            user_activity_report = "disable"
          }
        }
        network = {
          sdwan_interface_profile = "disable"
          zones = "disable"
          global_protect = {
            clientless_app_groups = "disable"
            clientless_apps = "disable"
            gateways = "disable"
            mdm = "disable"
            portals = "disable"
          }
        }
        objects = {
          address_groups = "disable"
          addresses = "disable"
          application_filters = "disable"
          application_groups = "disable"
          applications = "disable"
          authentication = "disable"
          devices = "disable"
          dynamic_block_lists = "disable"
          dynamic_user_groups = "disable"
          log_forwarding = "disable"
          packet_broker_profile = "disable"
          regions = "disable"
          schedules = "disable"
          security_profile_groups = "disable"
          service_groups = "disable"
          services = "disable"
          tags = "disable"
          custom_objects = {
            data_patterns = "disable"
            spyware = "disable"
            url_category = "disable"
            vulnerability = "disable"
          }
          decryption = {
            decryption_profile = "disable"
          }
          global_protect = {
            hip_objects = "disable"
            hip_profiles = "disable"
          }
          sdwan = {
            sdwan_dist_profile = "disable"
            sdwan_error_correction_profile = "disable"
            sdwan_profile = "disable"
            sdwan_saas_quality_profile = "disable"
          }
          security_profiles = {
            anti_spyware = "disable"
            antivirus = "disable"
            data_filtering = "disable"
            dos_protection = "disable"
            file_blocking = "disable"
            gtp_protection = "disable"
            sctp_protection = "disable"
            url_filtering = "disable"
            vulnerability_protection = "disable"
            wildfire_analysis = "disable"
          }
        }
        operations = {
          download_core_files = "disable"
          download_pcap_files = "disable"
          generate_stats_dump_file = "disable"
          generate_tech_support_file = "disable"
          reboot = "disable"
        }
        policies = {
          application_override_rulebase = "disable"
          authentication_rulebase = "disable"
          dos_rulebase = "disable"
          nat_rulebase = "disable"
          network_packet_broker_rulebase = "disable"
          pbf_rulebase = "disable"
          qos_rulebase = "disable"
          rule_hit_count_reset = "disable"
          sdwan_rulebase = "disable"
          security_rulebase = "disable"
          ssl_decryption_rulebase = "disable"
          tunnel_inspect_rulebase = "disable"
        }
        privacy = {
          show_full_ip_addresses = "disable"
          show_user_names_in_logs_and_reports = "disable"
          view_pcap_files = "disable"
        }
        save = {
          object_level_changes = "disable"
          partial_save = "disable"
          save_for_other_admins = "disable"
        }
      }
      xmlapi = {
        commit = "disable"
        config = "disable"
        export = "disable"
        import = "disable"
        iot = "disable"
        log = "disable"
        op = "disable"
        report = "disable"
        user_id = "disable"
      }
    }
  }
}
`
