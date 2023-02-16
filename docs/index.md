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

Refer to
[the changelog](https://github.com/PaloAltoNetworks/terraform-provider-panos/blob/master/CHANGELOG.md)
to see what's new.


## Provider v2

The `panos` provider has been around a while, and most of the releases have been
concerned with trying to cover a split between missing PAN-OS features and PAN-OS
updates that cause regressions in the provider.  Provider code itself has also
undergone a few shifts in philosophy.  In addition to this, the `panos` provider
is still using v1 of the
[terraform provider sdk](https://github.com/hashicorp/terraform-plugin-sdk),
which affects users ability to control provider logging as well as being the source 
[of issues](https://github.com/PaloAltoNetworks/terraform-provider-panos/issues/359)
with the provider itself.

All of this means that it's time for the `panos` provider to undergo a v2 update.  This
update will:

* Update the underlying HashiCorp libraries
* Include major performance improvements for all `_group` resources
* Converge the firewall and panorama resources across the whole provider
* Remove old resource aliases and deprecated parameters

These updates will take time.  Once the development effort begines, a new branch will
be created for the new v2 provider, allowing users to download the code, build it locally,
and provide feedback to the development team.

Please report any issues found with the v2 provider in github as usual.


## Terraform / PAN-OS Interaction

### `lifecycle.create_before_destroy`

The order of operations that Terraform handles updates / deletes does not by
default work the way that PAN-OS does things.  In order to make Terraform behave
properly, inside of **each and every resource** you need to specify a
[`lifecycle`](https://www.terraform.io/language/meta-arguments/lifecycle) block
like so:

```hcl
resource "panos_address_object" "example" {
    name = "web server 1"
    # continue with the rest of the definition
    ...

    lifecycle {
        create_before_destroy = true
    }
}
```


### Parallelism

Terraform uses goroutines to speed up deployment, but the number of parallel
operations is launches exceeds
[what is recommended](https://docs.paloaltonetworks.com/pan-os/10-0/pan-os-panorama-api/pan-os-xml-api-request-types/apply-user-id-mapping-and-populate-dynamic-address-groups-api.html):

```
Limit the number of concurrent API calls to five. This limit ensures that there is no performance impact to the firewall web interface as the management plane web server handles requests from both the API and the web interface.
```

In order to accomplish this, make sure you set the
[parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at or
below this limit to prevent performance impacts.


## Resource / Data Source Naming

In earlier releases of the `panos` provider, resources for the NGFW had a
`panos_` prefix, while resources intended for Panorama had a `panos_panorama_`
prefix.  This has become a bit of a stumbling block for new users of the `panos`
provider.  So starting in `panos` provider v1.7, we are doing away with this
distinction for any resources or data sources added, and will be slowing working to
retrofit the existing resources to behave this way.


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

* `hostname` - (env:`PANOS_HOSTNAME`) The hostname / IP address of PAN-OS.
* `username` - (env:`PANOS_USERNAME`) The PAN-OS username.  This is ignored
  if the `api_key` is given.
* `password` - (env:`PANOS_PASSWORD`) The PAN-OS password.  This is ignored
  if the `api_key` is given.
* `api_key` - (env:`PANOS_API_KEY`) The API key for the firewall.
* `protocol` - (env:`PANOS_PROTOCOL`) The communication protocol.  Valid
  values are `https` (the default) or `http`.
* `port` - (int, env:`PANOS_PORT`) If the port number is non-standard for
  the desired protocol, then the port number to use.
* `timeout` - (int, env:`PANOS_TIMEOUT`) The timeout for all communications
  with PAN-OS (default: `10`).
* `target` - (env:`PANOS_TARGET`) The firewall serial number to target
  configuration commands to (the `hostname` should be a Panorama PAN-OS).
* `additional_headers` - (env:`PANOS_HEADERS`, added in v1.9.0) Mapping of
  any additional headers to send with all API requests to PAN-OS.
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
* `log` - (v1.9.0) Log `log`.
* `export` - (v1.9.0) Log `export`.
* `import` - (v1.9.0) Log `import`.
* `xpath` - Log the XPATH associated with various actions.
* `send` - Log the raw request sent to the device.  This is probably
  only useful in development of the provider itself.
* `receive` - Log the raw response sent back from the device.  This is probably
  only useful in development of the provider itself.
* `osx_curl` - (v1.9.0) Output the API calls as OSX cURL calls.  Using the provider
  may uncover issues with PAN-OS itself.  If you believe you have encountered a bug
  with PAN-OS, enable cURL logging and give TAC the output of that, as the provider
  itself is still community supported.
* `curl_with_personal_data` - (v1.9.0) Without this specified, any curl logging will
  replace the hostname with `HOST`, the API key with `APIKEY`, and will not include
  any additional headers specified by the `headers` provider param.  If this logging type
  is specified in the logging, then these modifications are not done, essentially
  allowing you to copy/paste the cURL command and execute it yourself in your
  environment.


The following is an example of the contents of a JSON config file:

```json
{
    "hostname": "127.0.0.1",
    "api_key": "secret",
    "timeout": 10,
    "logging": ["action", "op", "uid", "osx_curl"],
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
