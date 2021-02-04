---
page_title: "panos: panos_api_key"
subcategory: "Operational State"
---

# panos_api_key

Use this data source to retrieve the API key associated with the provider
authentication credentials supplied.

Using the API key instead of classic username / password may result in some
speed gains when using the `panos` provider.


## Example Usage

```hcl
data "panos_api_key" "example" {}
```


## Attribute Reference

The following attributes are supported:

* `api_key` - The API key
