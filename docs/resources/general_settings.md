---
page_title: "panos: panos_general_settings"
subcategory: "Device"
---

# panos_general_settings

This resource allows you to update the general device settings, such as DNS
or the hostname.

All params are optional for this resource.  If any options are not specified,
then whatever is already configured on the firewall is left as-is.  The
general device settings will always exist on the firewall, so `terraform
destroy` does not remove config from the firewall.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_general_settings" "example" {
    hostname = "ngfw220"
    dns_primary = "10.5.1.10"
    ntp_primary = "10.5.1.10"
    ntp_primary_auth_type = "none"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - Firewall hostname.
* `timezone` - The timezone (e.g. - `US/Pacific`).
* `domain` - The domain.
* `update_server` - The update server (Default: `updates.paloaltonetworks.com`).
* `verify_update_server` - Verify update server identity (Default: `true`).
* `login_banner` - Login banner that is shown during the login page
* `proxy_server` - (1.5+) Specify a proxy server.
* `proxy_port` - (int, 1.5+) Proxy's port number.
* `proxy_username` - (1.5+) Proxy's username.
* `proxy_password` - (1.5+) Proxy's password.
* `dns_primary` - Primary DNS server.
* `dns_secondary` - Secondary DNS server.
* `ntp_primary_address` - Primary NTP server.
* `ntp_primary_auth_type` - Primary NTP auth type.  This can be `none`,
  `autokey`, or `symmetric-key`.
* `ntp_primary_key_id` - Primary NTP `symmetric-key` key ID.
* `ntp_primary_algorithm` - Primary NTP `symmetric-key` algorithm.  This can be
  `sha1` or `md5`.
* `ntp_primary_auth_key` - Primary NTP `symmetric-key` auth key.  This is the
  SHA1 hash if the algorithm is `sha1`, or the md5sum if the algorithm is
  `md5`.
* `ntp_secondary_address` - Secondary NTP server.
* `ntp_secondary_auth_type` - Secondary NTP auth type.  This can be `none`,
  `autokey`, or `symmetric-key`.
* `ntp_secondary_key_id` - Secondary NTP `symmetric-key` key ID.
* `ntp_secondary_algorithm` - Secondary NTP `symmetric-key` algorithm.  This
  can be `sha1` or `md5`.
* `ntp_secondary_auth_key` - Secondary NTP `symmetric-key` auth key.  This is
  the SHA1 hash if the algorithm is `sha1`, or the md5sum if the algorithm is
  `md5`.
