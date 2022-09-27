---
page_title: "panos: panos_ipsec_tunnel"
subcategory: "Network"
---

# panos_ipsec_tunnel

This resource allows you to add/update/delete IPSec tunnels.

A large number of params have prefixes:

* `ak` - Auto key
* `mk` - Manual key
* `gps` - GlobalProtect Satellite


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_ipsec_tunnel" "example" {
    name = "example"
    tunnel_interface = "tunnel.7"
    anti_replay = true
    ak_ike_gateway = "myIkeGateway"
    ak_ipsec_crypto_profile = "myIkeProfile"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name
* `tunnel_interface` - (Required) The tunnel interface.
* `anti_replay` - (Optional, bool) Set to `true` to enable Anti-Replay check
  on this tunnel.
* `enable_ipv6` - (Optional, PAN-OS 7.0+, bool) Set to `true` to enable IPv6.
* `copy_tos` - (Optional, bool) Set to `true` to copy IP TOS bits from inner
  packet to IPSec packet (not recommended).
* `copy_flow_label` - (Optional, PAN-OS 7.0+, bool) Set to `true` to copy IPv6
  flow label for 6in6 tunnel from inner packet to IPSec packet (not recommended).
* `disabled` - (Optional, PAN-OS 7.0+, bool) Set to `true` to disable this
  IPSec tunnel.
* `type` - (Optional) The type.  Valid values are `auto-key` (the default),
  `manual-key`, or `global-protect-satellite`.
* `ak_ike_gateway` - (Optional) IKE gateway name.
* `ak_ipsec_crypto_profile` - (Optional) IPSec crypto profile name.
* `mk_local_spi` - (Optional) Outbound SPI, hex format.
* `mk_remote_spi` - (Optional) Inbound SPI, hex format.
* `mk_local_address_ip` - (Optional) Specify exact IP address if interface
  has multiple addresses.
* `mk_local_address_floating_ip` - (Optional) Floating IP address in HA
  Active-Active configuration.
* `mk_protocol` - (Optional) Manual key protocol.  Valid valies are
  `esp` or `ah`.
* `mk_auth_type` - (Optional) Authentication algorithm.  Valid values are
  `md5`, `sha1`, `sha256`, `sha384`, `sha512`, or `none`.
* `mk_auth_key` - (Optional) The auth key for the given auth type.
* `mk_esp_encryption_type` - (Optional) The encryption algorithm.  Valid values
  are `des`, `3des`, `aes-128-cbc`, `aes-192-cbc`, `aes-256-cbc`, or `null`.
* `mk_esp_encryption_key` - (Optional) The encryption key.
* `gps_interface` - (Optional) Interface to communicate with portal.
* `gps_portal_address` - (Optional) GlobalProtect portal address.
* `gps_prefer_ipv6` - (Optional, PAN-OS 8.0+, bool) Prefer to register the
  portal in IPv6. Only applicable to FQDN portal-address.
* `gps_interface_ip_ipv4` - (Optional) specify exact IP address if interface
  has multiple addresses (IPv4).
* `gps_interface_ip_ipv6` - (Optional, PAN-OS 8.0+) specify exact IP address if interface
  has multiple addresses (IPv6).
* `gps_interface_floating_ip_ipv4` - (Optional, PAN-OS 7.0+) Floating IPv4
  address in HA Active-Active configuration.
* `gps_interface_floating_ip_ipv6` - (Optional, PAN-OS 8.0+) Floating IPv6
  address in HA Active-Active configuration.
* `gps_publish_connected_routes` - (Optional, bool) Set to `true` to to publish
  connected and static routes.
* `gps_publish_routes` - (Optional, list) Specify list of routes to publish
  to Global Protect Gateway.
* `gps_local_certificate` - (Optional) GlobalProtect satellite certificate
  file name.
* `gps_certificate_profile` - (Optional) Profile for authenticating
  GlobalProtect gateway certificates.
* `enable_tunnel_monitor` - (Optional, bool) Enable tunnel monitoring on this tunnel.
* `tunnel_monitor_destination_ip` - (Optional) Destination IP to send ICMP probe.
* `tunnel_monitor_source_ip` - (Optional) Source IP to send ICMP probe
* `tunnel_monitor_profile` - (Optional) Tunnel monitor profile.
* `tunnel_monitor_proxy_id` - (Optional, PAN-OS 7.0+) Which proxy-id (or
  proxy-id-v6) the monitoring traffic will use.
