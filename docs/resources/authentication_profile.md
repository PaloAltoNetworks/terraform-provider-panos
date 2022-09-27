---
page_title: "panos: panos_authentication_profile"
subcategory: "Device"
---

# panos_authentication_profile

Manages an authentication profile.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_authentication_profile" "example" {
    name = "fromTerraform"
    lockout_failed_attempts = "5"
    lockout_time = 4
    allow_list = ["aLocalUser"]
    type {
        local_database = true
    }
    multi_factor_authentication {
        enabled = true
        factors = [
            "first",
            "second",
            "third",
        ]
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
* `lockout_time` - Number of minutes to lock-out.
* `allow_list` - (List of string) List of allowed users or user groups.
* `type` - The type spec, as defined below.
* `username_modifier` - (PAN-OS 7.0+) Username modifier to handle user domain.  Valid values are `"%USERINPUT%"` (default), `"%USERINPUT%@%USERDOMAIN%"`, or `"%USERDOMAIN%\%USERINPUT%"`.
* `user_domain` - (PAN-OS 7.0+) Domain name to be used for authentication.
* `single_sign_on` - (PAN-OS 7.0+) Kerberos SSO settings spec, as defined below.
* `multi_factor_authentication` - (PAN-OS 8.0+) Specify MFA configuration spec, as defined below.

`type` supports the following arguments:

* `none` - (bool) No authentication.
* `local_database` - (bool) Local database authentication.
* `radius` - Radius authentication, as defined below.
* `ldap` - LDAP authenticatin, as defined below.

`type.radius` supports the following arguments:

* `server_profile` - (Required) Radius server profile object.
* `retrieve_user_group` - (bool, PAN-OS 7.0+) Retrieve user group from RADIUS.

`type.ldap` supports the following arguments:

* `server_profile` - (Required) LDAP server profile object.
* `login_attribute` - Login attribute in LDAP server to authenticate against.
* `password_expiry_warning` - Number of days prior to warning a user about password expiry (default: `"7"`).

`type.kerberos` supports the following arguments:

* `server_profile` - (Required) Kerberos server profile object.
* `realm` - (Required, PAN-OS 7.0+) Realm name to be used for authentication.

`type.tacacs_plus` supports the following arguments:

* `server_profile` - (Required) TACACS+ server profile object.
* `retrieve_user_group` - (bool, PAN-OS 8.0+) Retrieve user group from TACACS+.

`type.saml` supports the following arguments:

* `server_profile` - (Required) SAML IDP server profile object.
* `enable_single_logout` - (bool) Enable single logout.
* `request_signing_certificate` - (Signing certificate for SAML requests.
* `certificate_profile` - Certificate profile for IDP and SP.
* `username_attribute` - Attribute name for username to be extracted from SAML response (default: `"username"`).
* `user_group_attribute` - User group attribute.
* `admin_role_attribute` - Admin role attribute.
* `access_domain_attribute` - Access domain attribute.

`single_sign_on` supports the following arguments:

* `realm` - Kerberos realm to be used for authentication.
* `service_principal` - Kerberos service principal.
* `keytab` - Kerberos keytab.

`multi_factor_authentication` supports the following arguments:

* `enabled` - (bool) Enable additional authentication factors.
* `factors` - (List of strings) List of additional authentication factors.
