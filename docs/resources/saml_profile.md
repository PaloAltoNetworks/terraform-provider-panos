---
page_title: "panos: panos_saml_profile"
subcategory: "Device"
---

# panos_saml_profile

Manages a SAML IDP profile.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_saml_profiles" "example" {
    name = "fromTerraform"
    base_dn = "baseDn"
    bind_dn = "bindDn"
    password = "secret"
    bind_timeout = 5
    search_timeout = 7
    retry_interval = 120
    server {
        name = "first"
        server = "first.example.com"
    }
    server {
        name = "second"
        server = "192.168.0.5"
        port = 23430
    }

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.

NGFW / Panorama:

* `vsys` - The vsys (default: `shared`).

The following arguments are supported:

* `name` - (Required) The name.
* `admin_use_only` - (bool) Administrator use only.
* `identity_provider_id` - (Required) Unique identifier for SAML IdP.
* `identity_provider_certificate` - (Required) Object name of IdP signing certificate.
* `sso_url` - (Required) The single sign on service URL for the IdP server.
* `sso_binding` - SAML HTTP binding for SSO requests to IdP. Valid values are `"post"` (default) or `"redirect"`.
* `slo_url` - The single logout service URL for the IdP server.
* `slo_binding` - SAML HTTP binding for SLO requests to IdP. Valid values are `"post"` (default) or `"redirect"`.
* `validate_identity_provider_certificate` - (bool) Validate identity provider certificate (default: `true`).
* `sign_saml_message` - (bool) Sign SAML message to IdP (default: `true`).
* `max_clock_skew` - (int) Maximum allowed clock skew in seconds between SAML entities (default: `60`).
