---
page_title: "panos: panos_panorama_ike_gateway"
subcategory: "Panorama Networking"
---

# panos_panorama_ike_gateway

This resource allows you to add/update/delete Panorama IKE gateways
for both templates and template stacks.

## Example Usage

```hcl
# Note:  Only the template is a resource attribute variable, but all other
#   params (such as interface) should also be reference variables.
resource "panos_panorama_ike_gateway" "example" {
    template = panos_panorama_template.t.name
    name = "example"
    peer_ip_type = "dynamic"
    interface = "loopback.42"
    pre_shared_key = "secret"
    local_id_type = "ipaddr"
    local_id_value = "10.1.1.1"
    peer_id_type = "ipaddr"
    peer_id_value = "10.5.1.1"
    ikev1_crypto_profile = "myIkeProfile"
}

resource "panos_panorama_template" "t" {
    name = "my template"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The object's name
* `version` - (Optional, PAN-OS 7.0+) The IKE gateway version.  Valid values are
  `ikev1`, (the default), `ikev2`, or `ikev2-preferred`.  For PAN-OS 6.1, only
  `ikev1` is acceptable.
* `enable_ipv6` - (Optional, PAN-OS 7.0+, bool) Enable IPv6 or not.
* `disabled` - (Optional, PAN-OS 7.0+, bool) Set to `true` to disable.
* `peer_ip_type` - (Optional) The peer IP type.  Valid values are `ip`,
  `dynamic`, and `fqdn` (PANOS 8.1+).
* `peer_ip_value` - (Optional) The peer IP value.
* `interface` - (Required) The interface.
* `local_ip_address_type` - (Optional) The local IP address type.  Valid
  values for this are `ip`, `floating-ip`, or an empty string (the default)
  which is `None`.
* `local_ip_address_value` - (Optional) The IP address if `local_ip_address_type`
  is set to `ip`.
* `auth_type` - (Optional) The auth type.  Valid values are `pre-shared-key`
  (the default), or `certificate`.
* `pre_shared_key` - (Optional) The pre-shared key value.
* `local_id_type` - (Optional) The local ID type.  Valid values are `ipaddr`,
  `fqdn`, `ufqdn`, `keyid`, or `dn`.
* `local_id_value` - (Optional) The local ID value.
* `peer_id_type` - (Optional) The peer ID type.  Valid values are `ipaddr`,
  `fqdn`, `ufqdn`, `keyid`, or `dn`.
* `peer_id_value` - (Optional) The peer ID value.
* `peer_id_check` - (Optional) Enable peer ID wildcard match for certificate
  authentication.  Valid values are `exact` or `wildcard`.
* `local_cert` - (Optional) The local certificate name.
* `cert_enable_hash_and_url` - (Optional, PAN-OS 7.0+, bool) Set to `true` to use
  hash-and-url for local certificate.
* `cert_base_url` - (Optional) The host and directory part of URL for local
  certificates.
* `cert_use_management_as_source` - (Optional, PAN-OS 7.0+, bool) Set to `true` to
  use management interface IP as source to retrieve http certificates
* `cert_permit_payload_mismatch` - (Optional, bool) Set to `true` to permit
  peer identification and certificate payload identification mismatch.
* `cert_profile` - (Optional) Profile for certificate valdiation during IKE
  negotiation.
* `cert_enable_strict_validation` - (Optional, bool) Set to `true` to enable
  strict validation of peer's extended key use.
* `enable_passive_mode` - (Optional, bool) Set to `true` to enable passive
  mode (responder only).
* `enable_nat_traversal` - (Optional, bool) Set to `true` to enable NAT
  traversal.
* `nat_traversal_keep_alive` - (Optional, int) Sending interval for NAT
  keep-alive packets (in seconds).  For versions 6.1 - 8.1, this param, if specified,
  should be a multiple of 10 between 10 and 3600 to be valid.
* `nat_traversal_enable_udp_checksum` - (Optional, bool) Set to `true` to enable
  NAT traversal UDP checksum.
* `enable_fragmentation` - (Optional, bool) Set to `true` to enable fragmentation.
* `ikev1_exchange_mode` - (Optional) The IKEv1 exchange mode.
* `ikev1_crypto_profile` - (Optional) IKEv1 crypto profile.
* `enable_dead_peer_detection` - (Optional, bool) Set to `true` to enable dead
  peer detection.
* `dead_peer_detection_interval` - (Optional, int) The dead peer detection interval.
* `dead_peer_detection_retry` - (Optional, int) Number of retries before disconnection.
* `ikev2_crypto_profile` - (Optional, PAN-OS 7.0+) IKEv2 crypto profile.
* `ikev2_cookie_validation` - (Optional, PAN-OS 7.0+) Set to `true` to require cookie.
* `enable_liveness_check` - (Optional, , PAN-OS 7.0+bool) Set to `true` to
  enable sending empty information liveness check message.
* `liveness_check_interval` - (Optional, , PAN-OS 7.0+int) Delay interval before
  sending probing packets (in seconds).
