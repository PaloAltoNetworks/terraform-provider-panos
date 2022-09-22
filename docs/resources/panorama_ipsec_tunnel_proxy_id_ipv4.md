---
page_title: "panos: panos_panorama_ipsec_tunnel_proxy_id_ipv4"
subcategory: "Network"
---

# panos_panorama_ipsec_tunnel_proxy_id_ipv4

This resource allows you to add/update/delete Panorama IPSec tunnel proxy IDs
to a parent auto key IPSec tunnel for templates.


## PAN-OS

Panorama


## Import Name

```shell
<template>:<template_stack>:<ipsec_tunnel>:<name>
```


## Example Usage

```hcl
# NOTE: ipsec_tunnel should be an attribute resource variable (like how the
#  template param is referenced) in practice.
resource "panos_panorama_ipsec_tunnel_proxy_id_ipv4" "example" {
    template = panos_panorama_template.t.name
    ipsec_tunnel = "myIpsecTunnel"
    name = "example"
    local = "10.1.1.1"
    remote = "10.2.1.1"
    protocol_any = true

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "t" {
    name = "my template"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The template name.
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
