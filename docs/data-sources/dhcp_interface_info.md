---
page_title: "panos: panos_dhcp_interface_info"
subcategory: "Operational State"
---

# panos_dhcp_interface_info

Use this data source to retrieve DHCP client information about the given
firewall interface.


## PAN-OS

NGFW


## Example Usage

```hcl
data "panos_dhcp_interface_info" "example" {
    interface = "ethernet1/1"
}

output "eth1_ip" {
    value = data.panos_dhcp_interface_info.example.ip
}
```

## Argument Reference

* `interface` - (Required) The data interface to get DHCP information for.

## Attribute Reference

* `state` - The interface's state.
* `ip` - DHCP IP address.
* `gateway` - The default gateway assigned.
* `server` - The DHCP server IP
* `server_id` - DHCP server ID
* `primary_dns` - Primary DNS server
* `secondary_dns` - Secondary DNS server
* `primary_wins` - Primary WINS server
* `secondary_wins` - Secondary WINS
* `primary_nis` - Primary NIS
* `secondary_nis` - Secondary NIS
* `primary_ntp` - Primary NTP
* `secondary_ntp` - Secondary NTP
* `pop3_server` - POP3 Server
* `smtp_server` - SMTP Server
* `dns_suffix` - DNS Suffix
