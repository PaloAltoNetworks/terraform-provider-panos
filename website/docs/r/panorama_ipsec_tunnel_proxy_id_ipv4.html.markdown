---
layout: "panos"
page_title: "panos: panos_panorama_ipsec_tunnel_proxy_id_ipv4"
sidebar_current: "docs-panos-panorama-resource-ipsec-tunnel-proxy-id-ipv4"
description: |-
  Manages Panorama IPv4 proxy IDs for auto key IPSec tunnels.
---

# panos_panorama_ipsec_tunnel_proxy_id_ipv4

This resource allows you to add/update/delete Panorama IPSec tunnel proxy IDs
to a parent auto key IPSec tunnel for both templates and template stacks.

## Example Usage

```hcl
resource "panos_panorama_ipsec_tunnel_proxy_id_ipv4" "example" {
    template = "my template"
    ipsec_tunnel = "myIpsecTunnel"
    name = "example"
    local = "10.1.1.1"
    remote = "10.2.1.1"
    protocol_any = true
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

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
