## 1.8.0 (February 4, 2021)

NEW DATA SOURCES:

* `panos_api_key`

## 1.7.1 (February 4, 2021)

BUG FIXES:

* Fixing `panos_address_object` wanting to redeploy existing address objects.

## 1.7.0 (January 30, 2021)

Newly added resources and data sources now work with both NGFW and Panorama; there
is no separate `panos_panorama_` for Panorama.  If a data source or resource only
works with one, it will say as much in the documentation.  Additionally, the subheading
for the documentation will now be just "Objects" or "Network" for any of these new
style data sources/resources.

NEW DATA SOURCES:

* `panos_address_object` / `panos_address_objects`
* `panos_anti_spyware_security_profile` / `panos_anti_spyware_security_profiles`
* `panos_antivirus_security_profile` / `panos_antivirus_security_profiles`
* `panos_arp` / `panos_arps`
* `panos_custom_data_pattern_object` / `panos_custom_data_pattern_objects`
* `panos_data_filtering_security_profile` / `panos_data_filtering_security_profiles`
* `panos_device_group_parent`
* `panos_dos_protection_profile` / `panos_dos_protection_profiles`
* `panos_dynamic_user_group` / `panos_dynamic_user_groups`
* `panos_file_blocking_security_profile` / `panos_file_blocking_security_profiles`
* `panos_ospf`
* `panos_ospf_area` / `panos_ospf_areas`
* `panos_ospf_area_interface` / `panos_ospf_area_interfaces`
* `panos_ospf_area_virtual_link` / `panos_ospf_area_virtual_links`
* `panos_ospf_auth_profiles`
* `panos_ospf_export` / `panos_ospf_exports`
* `panos_predefined_dlp_file_type`
* `panos_predefined_tdb_file_type`
* `panos_predefined_threat`
* `panos_url_filtering_security_profile` / `panos_url_filtering_security_profiles`
* `panos_vulnerability_security_profile` / `panos_vulnerability_security_profiles`
* `panos_wildfire_analysis_security_profile` / `panos_wildfire_analysis_security_profiles`
* `panos_ip_tag`
* `panos_user_tag`
* `panos_vm_auth_key`

NEW RESOURCES:

* `panos_anti_spyware_security_profile`
* `panos_antivirus_security_profile`
* `panos_arp`
* `panos_custom_data_pattern_object`
* `panos_data_filtering_security_profile`
* `panos_dos_protection_profile`
* `panos_dynamic_user_group`
* `panos_file_blocking_security_profile`
* `panos_ip_tag`
* `panos_ospf`
* `panos_ospf_area`
* `panos_ospf_area_interface`
* `panos_ospf_area_virtual_link`
* `panos_ospf_auth_profile`
* `panos_ospf_export`
* `panos_url_filtering_security_profile`
* `panos_user_tag`
* `panos_vulnerability_security_profile`
* `panos_wildfire_analysis_security_profile`
* `panos_device_group_parent`
* `panos_vm_auth_key`

UPDATES:

* ECMP options added to `panos_virtual_router` / `panos_panorama_virtual_router`
* LACP, HA, and LLDP options added to both ethernet interfaces and aggregate interfaces
* `panos_panorama_plugin` has been renamed to `panos_plugin` and now also works with NGFW
  now.  `panos_panorama_plugin` still exists as an alias but will be removed in a
  future release, so please update your plan files accordingly.
* `panos_panorama_address_object` has been remade into the new "shared" style for
  data sources and resources as a kind of beta before touching any other existing
  resources.  This will cause extra unused params to exist in resources, but should
  not affect functionality.  Please let us know (GitHub issue) if this causes
  problems for you.  Otherwise the intent is to slowly retrofit resources into this
  new style.

## 1.6.3 (September 1, 2020)

* Release for Terraform Registry

## 1.6.2 (December 10, 2019)

UPDATES:

* Updated the provider to use the Terraform Plugin SDK ([#220](https://github.com/terraform-providers/terraform-provider-panos/issues/220))

## 1.6.1 (November 19, 2019)


UPDATES:

* The provider can now manage XFR PAN-OS releases ([#216](https://github.com/terraform-providers/terraform-provider-panos/issues/216))
* New optional provider param: `verify_certificate` ([#218](https://github.com/terraform-providers/terraform-provider-panos/issues/218))

## 1.6.0 (August 30, 2019)

NEW DATA SOURCES:

* `panos_panorama_plugin` ([#178](https://github.com/terraform-providers/terraform-provider-panos/issues/178))

NEW RESOURCES:

* `panos_gre_tunnel` / `panos_panorama_gre_tunnel` ([#162](https://github.com/terraform-providers/terraform-provider-panos/issues/162))
* `panos_monitor_profile` / `panos_panorama_monitor_profile` ([#182](https://github.com/terraform-providers/terraform-provider-panos/issues/182))
* `panos_panorama_gcp_account` ([#179](https://github.com/terraform-providers/terraform-provider-panos/issues/179))
* `panos_panorama_gke_cluster` ([#181](https://github.com/terraform-providers/terraform-provider-panos/issues/181))
* `panos_panorama_gke_cluster_group` ([#180](https://github.com/terraform-providers/terraform-provider-panos/issues/180))
* `panos_pbf_rule_group` / `panos_panorama_pbf_rule_group` ([#152](https://github.com/terraform-providers/terraform-provider-panos/issues/152))
* `panos_vlan` / `panos_panorama_vlan` ([#145](https://github.com/terraform-providers/terraform-provider-panos/issues/145))
* `panos_aggregate_interface` / `panos_panorama_aggregate_interface` ([#169](https://github.com/terraform-providers/terraform-provider-panos/issues/169))
* `panos_vlan_entry` / `panos_panorama_vlan_entry` ([#146](https://github.com/terraform-providers/terraform-provider-panos/issues/146))
* `panos_layer3_subinterface` / `panos_panorama_layer3_subinterface` ([#195](https://github.com/terraform-providers/terraform-provider-panos/issues/195))
* `panos_layer2_subinterface` / `panos_panorama_layer2_subinterface` ([#82](https://github.com/terraform-providers/terraform-provider-panos/issues/82))
* `panos_application_object` / `panos_panorama_application_object` ([#197](https://github.com/terraform-providers/terraform-provider-panos/issues/197))
* `panos_application_group` / `panos_panorama_application_group` ([#198](https://github.com/terraform-providers/terraform-provider-panos/issues/198))
* `panos_application_signature` / `panos_panorama_application_signature` ([#201](https://github.com/terraform-providers/terraform-provider-panos/issues/201))
* `panos_snmptrap_server_profile` / `panos_panorama_snmptrap_server_profile` ([#203](https://github.com/terraform-providers/terraform-provider-panos/issues/203))
* `panos_syslog_server_profile` / `panos_panorama_syslog_server_profile` ([#83](https://github.com/terraform-providers/terraform-provider-panos/issues/83))
* `panos_email_server_profile` / `panos_panorama_email_server_profile` ([#206](https://github.com/terraform-providers/terraform-provider-panos/issues/206))
* `panos_http_server_profile` / `panos_panorama_http_server_profile` ([#207](https://github.com/terraform-providers/terraform-provider-panos/issues/207))
* `panos_log_forwarding_profile` / `panos_panorama_log_forwarding_profile` ([#84](https://github.com/terraform-providers/terraform-provider-panos/issues/84))

NEW PARAMS ADDED TO:

* `panos_address_object` / `panos_panorama_address_object` ([#174](https://github.com/terraform-providers/terraform-provider-panos/issues/174))
* `panos_ethernet_interface` / `panos_panorama_ethernet_interface` ([#173](https://github.com/terraform-providers/terraform-provider-panos/issues/173))
* `panos_service_object` / `panos_panorama_service_object` ([#175](https://github.com/terraform-providers/terraform-provider-panos/issues/175))

BUG FIXES:

* The ordering for administrative tags on objects (such as address objects or address groups) is now preserved ([#161](https://github.com/terraform-providers/terraform-provider-panos/issues/161))
* The ordering for various list params in `panos_nat_rule_group` / `panos_panorama_nat_rule_group` is now ignored ([#143](https://github.com/terraform-providers/terraform-provider-panos/issues/143))
* New params added to `panos_nat_rule_group` / `panos_panorama_nat_rule_group` to to replace `static` and `dynamic` as `dynamic` is now a recursively reserved keyword in Terraform 0.12 ([#167](https://github.com/terraform-providers/terraform-provider-panos/issues/167))

## 1.5.2 (May 14, 2019)

* `panos_ike_gateway` / `panos_panorama_ike_gateway` - `floating-ip` is now an accepted value for `local_ip_address_type` ([#158](https://github.com/terraform-providers/terraform-provider-panos/issues/158))
* `panos_ike_gateway` / `panos_panorama_ike_gateway` - Updated documentation for `nat_traversal_keep_alive` ([#97](https://github.com/terraform-providers/terraform-provider-panos/issues/97))
* Fixed the acctest for `panos_panorama_nat_rule` ([#147](https://github.com/terraform-providers/terraform-provider-panos/issues/147))
* Properly handle when a list of strings is sent an empty string ([#153](https://github.com/terraform-providers/terraform-provider-panos/issues/153))
* Updated the commit script given on the main documentation page to include more ways to specify auth credentials ([#160](https://github.com/terraform-providers/terraform-provider-panos/issues/160))
* Updated to github.com/hashicorp/terraform@v0.12.0-rc1 ([#158](https://github.com/terraform-providers/terraform-provider-panos/issues/158))

## 1.5.1 (February 05, 2019)

The following resources can no longer be imported, as they have encrypted fields (thus there is no way to verify the plain text version of those fields) [[#139](https://github.com/terraform-providers/terraform-provider-panos/issues/139)]:

* `panos_bgp_auth_profile` / `panos_panorama_bgp_auth_profile`
* `panos_edl` / `panos_panorama_edl`
* `panos_ike_gateway` / `panos_panorama_ike_gateway`
* `panos_ipsec_tunnel` / `panos_panorama_ipsec_tunnel`

## 1.5.0 (February 04, 2019)

NEW RESOURCES:

* `panos_bfd_profile` / `panos_panorama_bfd_profile` ([#107](https://github.com/terraform-providers/terraform-provider-panos/issues/107))
* `panos_bgp` / `panos_panorama_bgp` ([#73](https://github.com/terraform-providers/terraform-provider-panos/issues/73))
* `panos_bgp_aggregate` / `panos_panorama_bgp_aggregate` ([#124](https://github.com/terraform-providers/terraform-provider-panos/issues/124))
* `panos_bgp_aggregate_advertise_filter` / `panos_panorama_bgp_aggregate_advertise_filter` ([#126](https://github.com/terraform-providers/terraform-provider-panos/issues/126))
* `panos_bgp_aggregate_suppress_filter` / `panos_panorama_bgp_aggregate_suppress_filter` ([#128](https://github.com/terraform-providers/terraform-provider-panos/issues/128))
* `panos_bgp_auth_profile` / `panos_panorama_bgp_auth_profile` ([#110](https://github.com/terraform-providers/terraform-provider-panos/issues/110))
* `panos_bgp_conditional_adv` / `panos_panorama_bgp_conditional_adv`, `panos_bgp_conditional_adv_advertise_filter` / `panos_panorama_bgp_conditional_adv_advertise_filter`, and `panos_bgp_conditional_adv_non_exist_filter` / `panos_panorama_bgp_conditional_adv_non_exist_filter` ([#122](https://github.com/terraform-providers/terraform-provider-panos/issues/122))
* `panos_bgp_dampening_profile` / `panos_panorama_bgp_dampening_profile` ([#111](https://github.com/terraform-providers/terraform-provider-panos/issues/111))
* `panos_bgp_export_rule_group` / `panos_panorama_bgp_export_rule_group` ([#120](https://github.com/terraform-providers/terraform-provider-panos/issues/120))
* `panos_bgp_import_rule_group` / `panos_panorama_bgp_import_rule_group` ([#118](https://github.com/terraform-providers/terraform-provider-panos/issues/118))
* `panos_bgp_peer` / `panos_panorama_bgp_peer` ([#116](https://github.com/terraform-providers/terraform-provider-panos/issues/116))
* `panos_bgp_peer_group` / `panos_panorama_bgp_peer_group` ([#114](https://github.com/terraform-providers/terraform-provider-panos/issues/114))
* `panos_bgp_redist_rule` / `panos_panorama_bgp_redist_rule` ([#130](https://github.com/terraform-providers/terraform-provider-panos/issues/130))
* `panos_nat_rule_group` / `panos_panorama_nat_rule_group` ([#78](https://github.com/terraform-providers/terraform-provider-panos/issues/78))
* `panos_redistribution_profile_ivp4` / `panos_panorama_redistribution_profile_ipv4` ([#92](https://github.com/terraform-providers/terraform-provider-panos/issues/92))

ENHANCEMENTS:

* Almost every resource can now be imported ([#86](https://github.com/terraform-providers/terraform-provider-panos/issues/86))
* Added proxy params to `panos_general_settings` ([#96](https://github.com/terraform-providers/terraform-provider-panos/issues/96))

DEPRECATED RESOURCES:

* `panos_nat_rule` / `panos_panorama_nat_rule` are both deprecated.  Please use `panos_nat_rule_group` / `panos_panorama_nat_rule_group` instead.

## 1.4.1 (October 26, 2018)

NEW RESOURCES:

* `panos_virtual_router_entry` and `panos_panorama_virtual_router_entry` ([#71](https://github.com/terraform-providers/terraform-provider-panos/issues/71))
* `panos_zone_entry` and `panos_panorama_zone_entry` ([#74](https://github.com/terraform-providers/terraform-provider-panos/issues/74))

BUG FIXES:

* Panorama device groups no longer require a description. ([#81](https://github.com/terraform-providers/terraform-provider-panos/issues/81))
* Panorama template stacks can now define a `default_vsys` ([#85](https://github.com/terraform-providers/terraform-provider-panos/issues/85))

## 1.4.0 (August 27, 2018)

NEW FEATURES:

* Support for both templates and template stacks has been added to the provider.  When defining your resource, use either the `template` variable if you want to attach it to a template, or `template_stack` if you want to attach it to a template stack.

NEW DATA SOURCES:

* `panos_dhcp_interface_info` ([#35](https://github.com/terraform-providers/terraform-provider-panos/issues/35))

NEW RESOURCES:

* `panos_ike_crypto_profile` and `panos_panorama_ike_crypto_profile` ([#37](https://github.com/terraform-providers/terraform-provider-panos/issues/37))
* `panos_ipsec_crypto_profile` and `panos_panorama_ipsec_crypto_profile` ([#38](https://github.com/terraform-providers/terraform-provider-panos/issues/38))
* `panos_tunnel_interface` and `panos_panorama_tunnel_interface` ([#42](https://github.com/terraform-providers/terraform-provider-panos/issues/42))
* `panos_ike_gateway` and `panos_panorama_ike_gateway` ([#39](https://github.com/terraform-providers/terraform-provider-panos/issues/39))
* `panos_ipsec_tunnel`, `panos_ipsec_tunnel_proxy_id_ipv4`, `panos_panorama_ipsec_tunnel`, and `panos_panorama_ipsec_tunnel_proxy_id_ipv4` ([#28](https://github.com/terraform-providers/terraform-provider-panos/issues/28))
* `panos_edl` and `panos_panorama_edl` ([#27](https://github.com/terraform-providers/terraform-provider-panos/issues/27))
* `panos_loopback_interface` and `panos_panorama_loopback_interface` ([#41](https://github.com/terraform-providers/terraform-provider-panos/issues/41))
* `panos_vlan_interface` and `panos_panorama_vlan_interface` ([#40](https://github.com/terraform-providers/terraform-provider-panos/issues/40))
* `panos_static_route_ipv4` and `panos_panorama_static_route_ipv4` ([#30](https://github.com/terraform-providers/terraform-provider-panos/issues/30))
* `panos_panorama_template`, `panos_panorama_template_entry`, `panos_panorama_template_stack`, `panos_panorama_template_stack_entry`, and `panos_panorama_template_variable` ([#43](https://github.com/terraform-providers/terraform-provider-panos/issues/43))
* `panos_license_api_key` and `panos_licensing` ([#24](https://github.com/terraform-providers/terraform-provider-panos/issues/24))
* `panos_panorama_management_profile` ([#58](https://github.com/terraform-providers/terraform-provider-panos/issues/58))
* `panos_panorama_ethernet_interface` ([#60](https://github.com/terraform-providers/terraform-provider-panos/issues/60))
* `panos_panorama_zone` ([#62](https://github.com/terraform-providers/terraform-provider-panos/issues/62))
* `panos_panorama_virtual_router` ([#64](https://github.com/terraform-providers/terraform-provider-panos/issues/64))

## 1.3.0 (June 27, 2018)

RENAMED RESOURCES:

The following resources have been renamed for clarity from their original names.  Both the old name and the new name will work right now, but please update your plans to use the new names as the original names may be removed / repurposed in the future.

* `panos_nat_policy` is now `panos_nat_rule` ([#34](https://github.com/terraform-providers/terraform-provider-panos/issues/34))
* `panos_security_policies` is now `panos_security_policy` ([#34](https://github.com/terraform-providers/terraform-provider-panos/issues/34))
* `panos_security_policy_group` is now `panos_security_rule_group` ([#34](https://github.com/terraform-providers/terraform-provider-panos/issues/34))
* `panos_panorama_nat_policy` is now `panos_panorama_nat_rule` ([#34](https://github.com/terraform-providers/terraform-provider-panos/issues/34))
* `panos_panorama_security_policies` is now `panos_panorama_security_policy` ([#34](https://github.com/terraform-providers/terraform-provider-panos/issues/34))
* `panos_panorama_security_policy_group` is now `panos_panorama_security_rule_group` ([#34](https://github.com/terraform-providers/terraform-provider-panos/issues/34))

## 1.2.0 (June 19, 2018)

FEATURES:

* **New Resource:** `panos_telemetry` ([#31](https://github.com/terraform-providers/terraform-provider-panos/issues/31))
* **New Resource:** `panos_security_policy_group` ([#20](https://github.com/terraform-providers/terraform-provider-panos/issues/20)] [[#32](https://github.com/terraform-providers/terraform-provider-panos/issues/32))
* **New Resource:** `panos_panorama_security_policy_group` ([#20](https://github.com/terraform-providers/terraform-provider-panos/issues/20)] [[#32](https://github.com/terraform-providers/terraform-provider-panos/issues/32))

NOTES:

* The new `DatType` param is now required if you are doing destination address translation in your NAT policies.  This applies to both `panos_nat_policy` and `panos_panorama_nat_policy`.  Please update your plan files accordingly.

ENHANCEMENTS:

* `panos_nat_policy` and `panos_panorama_nat_policy` now support PAN-OS 8.1's dynamic destination NAT address type ([#25](https://github.com/terraform-providers/terraform-provider-panos/issues/25)] [[#33](https://github.com/terraform-providers/terraform-provider-panos/issues/33))

FIXES:

* Creating Panorama service objects in device groups ([#26](https://github.com/terraform-providers/terraform-provider-panos/issues/26)] [[#29](https://github.com/terraform-providers/terraform-provider-panos/issues/29))

## 1.1.0 (April 26, 2018)

FEATURES:

* **New Feature:** Added Panorama support ([#3](https://github.com/terraform-providers/terraform-provider-panos/issues/3))
* **New Feature:** Added support for credentials file for provider config ([#5](https://github.com/terraform-providers/terraform-provider-panos/issues/5))
* **New Resource:** `panos_panorama_address_group`
* **New Resource:** `panos_panorama_address_object`
* **New Resource:** `panos_panorama_administrative_tag`
* **New Resource:** `panos_panorama_device_group`
* **New Resource:** `panos_panorama_device_group_entry`
* **New Resource:** `panos_panorama_nat_policy`
* **New Resource:** `panos_panorama_security_policies`
* **New Resource:** `panos_panorama_service_group`
* **New Resource:** `panos_panorama_service_object`

ENHANCEMENTS:

* `panos_nat_policy`: The `rulebase` parameter has been deprecated.  You can safely remove this from your plan files.
* `panos_security_policies`: The `rulebase` parameter has been deprecated.  You can safely remove this from your plan files.

## 1.0.0 (January 31, 2018)

FEATURES:

* **New Data Source:** `panos_system_info`
* **New Resource:** `panos_address_group`
* **New Resource:** `panos_address_object`
* **New Resource:** `panos_administrative_tag`
* **New Resource:** `panos_dag_tags`
* **New Resource:** `panos_ethernet_interface`
* **New Resource:** `panos_general_settings`
* **New Resource:** `panos_management_profile`
* **New Resource:** `panos_nat_policy`
* **New Resource:** `panos_security_policies`
* **New Resource:** `panos_service_group`
* **New Resource:** `panos_service_object`
* **New Resource:** `panos_virtual_router`
* **New Resource:** `panos_zone`
