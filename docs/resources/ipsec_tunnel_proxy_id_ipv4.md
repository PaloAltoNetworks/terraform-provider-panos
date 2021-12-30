---
page_title: "panos: panos_ipsec_tunnel_proxy_id_ipv4"
subcategory: "Network"
---

# panos_ipsec_tunnel_proxy_id_ipv4

This resource allows you to add/update/delete IPSec tunnel proxy IDs to
a parent auto key IPSec tunnel.


## PAN-OS

NGFW


## Import Name

```
<ipsec_tunnel>:<name>
```


## Example Usage

```hcl
resource "panos_ipsec_tunnel_proxy_id_ipv4" "example" {
    ipsec_tunnel = panos_ipsec_tunnel.t1.name
    name = "example"
    local = "10.1.1.1"
    remote = "10.2.1.1"
    protocol_any = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name
* `ipsec_tunnel` - (Required) The auto key IPSec tunnel to attach this 
  proxy ID to.
* `local` - (Optional) IP subnet or IP address represents local network.
* `remote` - (Optional) IP subnet or IP address represents remote network.
* `protocol_any` - (Optional, bool) Set to `true` for any IP protocol.
* `protocol_number` - (Optional, int) IP protocol number.
* `protocol_tcp_local` - (Optional, int) Local TCP port number.
* `protocol_tcp_remote` - (Optional, int) Remote TCP port number.
* `protocol_udp_local` - (Optional, int) Local UDP port number.
* `protocol_udp_remote` - (Optional, int) Remote UDP port number.
