---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "panos_ntp_settings Data Source - panos"
subcategory: ""
description: |-
  
---

# panos_ntp_settings (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `location` (Attributes) The location of this object. (see [below for nested schema](#nestedatt--location))

### Optional

- `ntp_servers` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers))

### Read-Only

- `encrypted_values` (Map of String, Sensitive)
- `tfid` (String) The Terraform ID.

<a id="nestedatt--location"></a>
### Nested Schema for `location`

Optional:

- `system` (Attributes) Located in a system settings. (see [below for nested schema](#nestedatt--location--system))
- `template` (Attributes) Located in a specific template. (see [below for nested schema](#nestedatt--location--template))
- `template_stack` (Attributes) Located in a specific template stack. (see [below for nested schema](#nestedatt--location--template_stack))

<a id="nestedatt--location--system"></a>
### Nested Schema for `location.system`

Optional:

- `ngfw_device` (String) The NGFW device.


<a id="nestedatt--location--template"></a>
### Nested Schema for `location.template`

Optional:

- `name` (String) The template.
- `ngfw_device` (String) The NGFW device.
- `panorama_device` (String) The panorama device.


<a id="nestedatt--location--template_stack"></a>
### Nested Schema for `location.template_stack`

Optional:

- `name` (String) The template stack.
- `ngfw_device` (String) The NGFW device.
- `panorama_device` (String) The panorama device.



<a id="nestedatt--ntp_servers"></a>
### Nested Schema for `ntp_servers`

Optional:

- `primary_ntp_server` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--primary_ntp_server))
- `secondary_ntp_server` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--secondary_ntp_server))

<a id="nestedatt--ntp_servers--primary_ntp_server"></a>
### Nested Schema for `ntp_servers.primary_ntp_server`

Optional:

- `authentication_type` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--primary_ntp_server--authentication_type))
- `ntp_server_address` (String) NTP Server IP Address or Domain Name

<a id="nestedatt--ntp_servers--primary_ntp_server--authentication_type"></a>
### Nested Schema for `ntp_servers.primary_ntp_server.authentication_type`

Optional:

- `autokey` (String) Autokey Authentication
- `none` (String) No NTP Authentication
- `symmetric_key` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--primary_ntp_server--authentication_type--symmetric_key))

<a id="nestedatt--ntp_servers--primary_ntp_server--authentication_type--symmetric_key"></a>
### Nested Schema for `ntp_servers.primary_ntp_server.authentication_type.symmetric_key`

Optional:

- `key_id` (Number) Symmetric Key Number
- `md5` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--primary_ntp_server--authentication_type--symmetric_key--md5))
- `sha1` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--primary_ntp_server--authentication_type--symmetric_key--sha1))

<a id="nestedatt--ntp_servers--primary_ntp_server--authentication_type--symmetric_key--md5"></a>
### Nested Schema for `ntp_servers.primary_ntp_server.authentication_type.symmetric_key.md5`

Optional:

- `authentication_key` (String, Sensitive) Symmetric Key MD5 String


<a id="nestedatt--ntp_servers--primary_ntp_server--authentication_type--symmetric_key--sha1"></a>
### Nested Schema for `ntp_servers.primary_ntp_server.authentication_type.symmetric_key.sha1`

Optional:

- `authentication_key` (String, Sensitive) Symmetric Key SHA1 Hexadecimal





<a id="nestedatt--ntp_servers--secondary_ntp_server"></a>
### Nested Schema for `ntp_servers.secondary_ntp_server`

Optional:

- `authentication_type` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--secondary_ntp_server--authentication_type))
- `ntp_server_address` (String) NTP Server IP Address or Domain Name

<a id="nestedatt--ntp_servers--secondary_ntp_server--authentication_type"></a>
### Nested Schema for `ntp_servers.secondary_ntp_server.authentication_type`

Optional:

- `autokey` (String) Autokey Authentication
- `none` (String) No NTP Authentication
- `symmetric_key` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--secondary_ntp_server--authentication_type--symmetric_key))

<a id="nestedatt--ntp_servers--secondary_ntp_server--authentication_type--symmetric_key"></a>
### Nested Schema for `ntp_servers.secondary_ntp_server.authentication_type.symmetric_key`

Optional:

- `key_id` (Number) Symmetric Key Number
- `md5` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--secondary_ntp_server--authentication_type--symmetric_key--md5))
- `sha1` (Attributes) (see [below for nested schema](#nestedatt--ntp_servers--secondary_ntp_server--authentication_type--symmetric_key--sha1))

<a id="nestedatt--ntp_servers--secondary_ntp_server--authentication_type--symmetric_key--md5"></a>
### Nested Schema for `ntp_servers.secondary_ntp_server.authentication_type.symmetric_key.md5`

Optional:

- `authentication_key` (String) Symmetric Key MD5 String


<a id="nestedatt--ntp_servers--secondary_ntp_server--authentication_type--symmetric_key--sha1"></a>
### Nested Schema for `ntp_servers.secondary_ntp_server.authentication_type.symmetric_key.sha1`

Optional:

- `authentication_key` (String) Symmetric Key SHA1 Hexadecimal
