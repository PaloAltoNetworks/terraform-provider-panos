---
page_title: "Provider: panos"
description: |-
    The Terraform provider for Palo Alto Networks PAN-OS.
---

# Provider panos

PAN-OS&reg; is the operating system for Palo Alto Networks&reg; NGFWs and Panorama&trade;. The panos provider allows you to manage various aspects of a firewall's or a Panorama's config, such as data interfaces and security policies.

Use the navigation to the left to read about the available Panorama and NGFW resources.

Refer to [the changelog](https://github.com/PaloAltoNetworks/terraform-provider-panos/blob/master/CHANGELOG.md) to see what's new.


## Versioning

The minimum version of PAN-OS is PAN-OS 10.1 and onwards.  In PAN-OS 10.1.5 a change was made with regards to the XML API that removed the provider's easy access to remove multiple config objects in a single API call.

To compensate for this loss of functionality, the provider now does multi-config API calls against PAN-OS when multiple `type=config` API calls must be run.

With regards to data sources and resources, the minimum version for any given data source or resource will be present in the documentation if the version is higher than PAN-OS 10.1.

If you need to work with multiple versions of PAN-OS where some versions have a new parameter and some don't, then make sure to use [conditionals](https://www.terraform.io/docs/configuration/expressions/conditionals.html) along with the `panos_system_info` data source to only set these variables when the version of PAN-OS is appropriate.


## List Filtering

Data sources that return a list of objects now supports filtering.  Please refer to the guide on the side for more information on how to use filtering.


## Provider Modes

The provider features two distinct modes:  API and local inspection.

The API mode is the standard mode where the provider connects to a running PAN-OS instance and executes API commands against it.  This mode can be used across both data sources and resources.

The local inspection mode in contrast allows users to use Terraform to inspect a locally saved XML schema that was previously exported from PAN-OS. Interaction with the exported XML schema is limited to data sources only, but there may still be some types of data sources that are incompatible with this (aka - User-ID data sources or exporting the tech support file). This combined with data source list filtering should enable users to quickly inspect previously saved configs from the comfort of Terraform. The versions of PAN-OS that [pango](https://github.com/PaloAltoNetworks/pango) supports determines the versions that the provider can successfully function in local inspection mode.


## Candidate vs Running Config

Most data sources (aka - those that are `type=config`) now support querying both the candidate config as well as the running config.  The default is to query candidate config (`get`), but if the data source supports the `action` param then this can be set to `show` to query the running config instead.

If the provider is in local inspection mode, then the `action` is ignored as the provider is getting answers directly from the given XML schema.


## Terraform / PAN-OS Interaction

### `lifecycle.create_before_destroy`

The order of operations that Terraform handles updates / deletes does not by default work the way that PAN-OS does things.  In order to make Terraform behave properly, inside of **each and every resource** you need to specify a [`lifecycle`](https://www.terraform.io/language/meta-arguments/lifecycle) block like so:

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

Terraform uses goroutines to speed up deployment, but the number of parallel operations is launches exceeds [what is recommended](https://docs.paloaltonetworks.com/pan-os/10-0/pan-os-panorama-api/pan-os-xml-api-request-types/apply-user-id-mapping-and-populate-dynamic-address-groups-api.html):

```
Limit the number of concurrent API calls to five. This limit ensures that there is no performance impact to the firewall web interface as the management plane web server handles requests from both the API and the web interface.
```

In order to accomplish this, make sure you set the [parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at or below this limit to prevent performance impacts.


## Commits

As of right now, Terraform does not provide native support for commits, so commits are handled out-of-band.  Please refer to the commit guide to the left for more information.


## AWS / GCP Considerations

If you are launching PAN-OS in AWS or GCP, there are additional considerations that you should be aware of with regards to initial configuration.  Please see the AWS / GCP Considerations guide off to the left.


## Example Provider Usage

```terraform
# Traditional provider example.
provider "panos" {
  hostname = "10.1.1.1"
  username = "admin"
  password = "secret"
}

# Local inspection mode provider example.
provider "panos" {
  config_file = file("/tmp/candidate-config.xml")

  # This is only used if a "detail-version" attribute is not present in
  # the exported XML schema. If it's there, this can be omitted.
  panos_version = "10.2.0"
}

terraform {
  required_providers {
    panos = {
      source  = "paloaltonetworks/terraform-provider-panos"
      version = "2.0.0"
    }
  }
}
```


## Provider Parameter Priority

There are multiple ways to specify the provider's parameters.  If overlapping values are configured for the provider, then this is the resolution order:

1. Directly in the `provider` block
2. Environment variable (where applicable)
3. From the JSON config file


<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `additional_headers` (Map of String) Additional HTTP headers to send with API calls Environment variable: `PANOS_HEADERS`. JSON config file variable: `additional_headers`.
- `api_key` (String) The API key for PAN-OS. Either specify this or give both username and password. Environment variable: `PANOS_API`. JSON config file variable: `api_key`.
- `api_key_in_request` (Boolean) Send the API key in the request body instead of using the authentication header. Environment variable: `PANOS_API_KEY_IN_REQUEST`. JSON config file variable: `api_key_in_request`.
- `auth_file` (String) Filesystem path to a JSON config file that specifies the provider's params.
- `config_file` (String) (Local inspection mode) The PAN-OS config file to load read in using `file()`
- `hostname` (String) The hostname or IP address of the PAN-OS instance (NGFW or Panorama). Environment variable: `PANOS_HOST`. JSON config file variable: `hostname`.
- `panos_version` (String) (Local inspection mode) The version of PAN-OS that exported the config file. This is only used if the root "config" block does not contain the "detail-version" attribute. Example: `10.2.3`.
- `password` (String, Sensitive) The password.  This is required if the api_key is not configured. Environment variable: `PANOS_PASSWORD`. JSON config file variable: `password`.
- `port` (Number) If the port is non-standard for the protocol, the port number to use. Environment variable: `PANOS_PORT`. JSON config file variable: `port`.
- `protocol` (String) The protocol (https or http). Default: `https`. Environment variable: `PANOS_PROTOCOL`. JSON config file variable: `protocol`.
- `skip_verify_certificate` (Boolean) (For https protocol) Skip verifying the HTTPS certificate. Environment variable: `PANOS_SKIP_VERIFY_CERTIFICATE`. JSON config file variable: `skip_verify_certificate`.
- `target` (String) Target setting (NGFW serial number). Environment variable: `PANOS_TARGET`. JSON config file variable: `target`.
- `username` (String) The username.  This is required if api_key is not configured. Environment variable: `PANOS_USERNAME`. JSON config file variable: `username`.

The following is an example of the contents of a JSON config file:

```json
{
    "hostname": "127.0.0.1",
    "api_key": "secret",
    "skip_verify_certificate": true
}
```


## Support

This template/script/solution is released “as-is”, with no warranty and no support. These should be seen as community supported and Palo Alto Networks may contribute its expertise at its discretion. Palo Alto Networks, including through its Authorized Support Centers (ASC) partners and backline support options, will not provide technical support or help in using or troubleshooting this template/script/solution. The underlying product used by this template/script/solution will still be supported in accordance with the product’s applicable support policy and the customer’s entitlements.
