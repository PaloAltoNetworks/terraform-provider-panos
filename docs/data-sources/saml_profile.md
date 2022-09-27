---
page_title: "panos: panos_saml_profile"
subcategory: "Device"
---

# panos_saml_profile

Gets information on a SAML IDP profile.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_saml_profiles" "example" {
    name = "fromTerraform"
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


## Attribute reference

The following attributes are supported:

* `admin_use_only` - (bool) Administrator use only.
* `identity_provider_id` - Unique identifier for SAML IdP.
* `identity_provider_certificate` - Object name of IdP signing certificate.
* `sso_url` - The single sign on service URL for the IdP server.
* `sso_binding` - SAML HTTP binding for SSO requests to IdP.
* `slo_url` - The single logout service URL for the IdP server.
* `slo_binding` - SAML HTTP binding for SLO requests to IdP.
* `validate_identity_provider_certificate` - (bool) Validate identity provider certificate.
* `sign_saml_message` - (bool) Sign SAML message to IdP.
* `max_clock_skew` - (int) Maximum allowed clock skew in seconds between SAML entities.
