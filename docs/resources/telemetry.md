---
page_title: "panos: panos_telemetry"
subcategory: "Device"
---

# panos_telemetry

This resource allows you to add/update/delete telemetry sharing.

Join other Palo Alto Networks customers in a global sharing community, helping
to raise the bar against the latest attack techniques. Your participation
allows us to deliver new threat prevention controls across the attack
lifecycle. Choose the type of data you share across applications, threat
intelligence, and device health information to improve the fidelity of the
protections we deliver. This is an opt-in feature controlled with granular
policy, and we encourage you to join the community.


## PAN-OS

NGFW


## Import Name

```shell
<provider hostname>
```


## Example Usage

```hcl
resource "panos_telemetry" "example" {
    threat_prevention_reports = true
    threat_prevention_data = true
    threat_prevention_packet_captures = true
}
```

## Argument Reference

The following arguments are supported:

* `application_reports` - (Bool, optional) Application reports.
* `threat_prevention_reports` - (Bool, optional) Threat reports.
* `url_reports` - (Bool, optional) URL reports.
* `file_type_identification_reports` - (Bool, optional) File type identification
  reports.
* `threat_prevention_data` - (Bool, optional) Threat prevention data.
* `threat_prevention_packet_captures` - (Bool, optional) Enable sending packet-
  captures with threat prevention information. This requires that
  `threat_prevention_data` also be enabled.
* `product_usage_stats` - (Bool, optional) Health and performance reports.
* `passive_dns_monitoring` - (Bool, optional) Passive DNS monitoring.
