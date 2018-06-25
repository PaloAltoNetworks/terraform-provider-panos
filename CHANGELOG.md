## 1.2.1 (Unreleased)

RENAMED RESOURCES:

The following resources have been renamed for clarity from their original names.  Both the old name and the new name will work right now, but please update your plans to use the new names as the original names may be removed / repurposed in the future.

* `panos_nat_policy` is now `panos_nat_rule` [GH-34]
* `panos_security_policies` is now `panos_security_policy` [GH-34]
* `panos_security_policy_group` is now `panos_security_rule_group` [GH-34]
* `panos_panorama_nat_policy` is now `panos_panorama_nat_rule` [GH-34]
* `panos_panorama_security_policies` is now `panos_panorama_security_policy` [GH-34]
* `panos_panorama_security_policy_group` is now `panos_panorama_security_rule_group` [GH-34]

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
