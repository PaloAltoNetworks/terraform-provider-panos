---
page_title: "Provider: panos"
---

# Provider panos

PAN-OS&reg; is the operating system for Palo Alto Networks&reg; NGFWs and
Panorama&trade;. The panos provider allows you to manage various aspects
of a firewall's or a Panorama's config, such as data interfaces and security
policies.

Use the navigation to the left to read about the available Panorama and NGFW
resources.

## Resource Naming

In earlier releases of the `panos` provider, resources for the NGFW had a
`panos_` prefix, while resources intended for Panorama had a `panos_panorama_`
prefix.  This has become a bit of a stumbling block for new users of the `panos`
provider.  So starting in `panos` provider v1.7, we are doing away with this
distinction for any resources or data sources added, and will be slowing working to
retrofit the existing resources to behave this way.

One resource has been retrofitted in this regard, and that's the `panos_address_object`
resource.  In order to not break existing plans / deployments, the provider has an alias
for `panos_panorama_address_object` that uses this new approach.  This means that as of
1.7, you can add an address object to Panorama using the name `panos_address_object` or
`panos_panorama_address_object`.  This is the only pre-1.7 resource that has been
modified like this as a trial run.  If you encounter any issues, please open a GitHub
issue against the provider.

This merging also has an affect on the ID for a resource.  Since NGFW and Panorama
often require slightly different params, the names will be slightly different based
on if the provider is acting on a NGFW or a Panorama.  The ID is important when you
try to import an existing config into Terraform.  So in the documentation for these
newer style resources, there will be one format for NGFW and another format for
Panorama.

## Versioning

In general, the panos provider has support for PAN-OS 6.1 onwards.  Data
sources or resources that have minimum PAN-OS version requirements will
specify their version requirements in their documentation.

Some resources may contain variables that are only applicable for newer
versions of PAN-OS.  If you need to work with multiple versions of PAN-OS
where some versions have a new parameter and some don't, then make sure to use
[conditionals](https://www.terraform.io/docs/configuration/expressions/conditionals.html)
along with the `panos_system_info` data source to only set these variables
when the version of PAN-OS is appropriate.

One such resource is `panos_ethernet_interface` and the `ipv4_mss_adjust`
parameter.  Doing the following is one way to correctly configure this
parameter only when it's applicable:

```hcl
data "panos_system_info" "config" {}

data "panos_ethernet_interface" "eth1" {
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
    adjust_tcp_mss = true
    ipv4_mss_adjust = "${data.panos_system_info.config.version_major >= 8 ? 42 : 0}"
    # ...
}
```

## Commits

As of right now, Terraform does not provide native support for commits, so
commits are handled out-of-band.  Please refer to the commit guide to the left
for more information.

## AWS / GCP Considerations

If you are launching PAN-OS in AWS or GCP, there are additional considerations
that you should be aware of with regards to initial configuration.  Please see
the AWS / GCP Considerations guide off to the left.

## Importing Resources

Many resources support being imported.  Any resource that supports `terraform
import` will have a "Import Name" section in the documentation.  The variables
given in this section directly match up with the resource params you would
specify in your plan file.  Thus if you were importing an ethernet interface
whose import name is `<vsys>:<name>`, your import name would be something like
`vsys1:ethernet1/1`.

Of special note is the Panorama resources.  The templated resources often
have both the template and the template stack in the resource name, however
only one of these can ever be present.  Thus, the one that isn't being used
should just be an empty string.  For example, if you were trying to import
a Panorama IPv4 static route whose import name is
`<template>:<template_stack>:<virtual_router>:<name>` that resides in a
template, your import name would be something like
`myTemplate::myVirtualRouter:myStaticRouteName`.

If the resource is a combined resource (see [Resource
Naming](#resource-naming) above), then the documentation for that resource
or data source will have separate entries for NGFW and Panorama naming.

## Example Provider Usage

```hcl
# Configure the panos provider
provider "panos" {
    hostname = "127.0.0.1"
    json_config_file = "../panos-creds.json"
}
```

## Argument Reference

Arguments can be given in any or all of the following ways.  A parameter value
is taken from the highest priority source, with lower priority sources being
ignored.  From highest to lowest priority, these ways are:

1. Directly in the `provider` block
2. Environment variable setting (where applicable)
3. From the JSON config file


The following arguments are supported:

* `hostname` - (Optional, env:`PANOS_HOSTNAME`) The hostname / IP address of PAN-OS.
* `username` - (Optional, env:`PANOS_USERNAME`) The PAN-OS username.  This is ignored
  if the `api_key` is given.
* `password` - (Optional, env:`PANOS_PASSWORD`) The PAN-OS password.  This is ignored
  if the `api_key` is given.
* `api_key` - (Optional, env:`PANOS_API_KEY`) The API key for the firewall.
* `protocol` - (Optional, env:`PANOS_PROTOCOL`) The communication protocol.  Valid
  values are `https` (the default) or `http`.
* `port` - (Optional, int, env:`PANOS_PORT`) If the port number is non-standard for
  the desired protocol, then the port number to use.
* `timeout` - (Optional, int, env:`PANOS_TIMEOUT`) The timeout for all communications
  with PAN-OS (default: `10`).
* `target` - (Optional, env:`PANOS_TARGET`) The firewall serial number to target
  configuration commands to (the `hostname` should be a Panorama PAN-OS).
* `logging` - (Optional, env:`PANOS_LOGGING`) List of logging options for the
  provider's connection to the API.  If this is unspecified, then it defaults to
  `["action", "uid"]`.  If this is being specified as an environment variable,
  then it should be a CSV list.
* `verify_certificate` - (Optional, bool, env:`PANOS_VERIFY_CERTIFICATE`) For HTTPS
  protocol connections, verify that the certificate is valid.
* `json_config_file` - (Optional) The path to a JSON configuration file that
  contains any number of the provider's parameters.  If specified, the params
  present act as a last resort for any other provider param that has not been
  specified yet.  Params in the JSON config file match what the provider block
  supports, both in naming convention and data types.  See below for an example of
  the JSON config file contents.

The list of strings supported for `logging` are as follows:

* `quiet` - Disables logging if only `quiet` is specified.
* `action` - Log `set` / `edit` / `delete`.
* `query` - Log `get`.
* `op` - Log `op`.
* `uid` - Log user-id envocations.
* `xpath` - Log the XPATH associated with various actions.
* `send` - Log the raw request sent to the device.  This is probably
  only useful in development of the provider itself.
* `receive` - Log the raw response sent back from the device.  This is probably
  only useful in development of the provider itself.

The following is an example of the contents of a JSON config file:

```json
{
    "hostname": "127.0.0.1",
    "api_key": "secret",
    "timeout": 10,
    "logging": ["action", "op", "uid"],
    "verify_certificate": false
}
```

## Support

This template/solution are released under an as-is, best effort, support
policy. These scripts should be seen as community supported and Palo Alto
Networks will contribute our expertise as and when possible. We do not
provide technical support or help in using or troubleshooting the components
of the project through our normal support options such as Palo Alto Networks
support teams, or ASC (Authorized Support Centers) partners and backline
support options. The underlying product used (the VM-Series firewall) by the
scripts or templates are still supported, but the support is only for the
product functionality and not for help in deploying or using the template or
script itself. Unless explicitly tagged, all projects or work posted in our
GitHub repository (at https://github.com/PaloAltoNetworks) or sites other
than our official Downloads page on https://support.paloaltonetworks.com
are provided under the best effort policy.
