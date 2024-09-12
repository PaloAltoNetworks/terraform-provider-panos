Terraform Provider for Palo Alto Networks PANOS
===============================================

- Website: https://www.terraform.io
- Documentation: https://www.terraform.io/docs/providers/panos/index.html

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 1.8.x or higher
- [Go](https://golang.org/doc/install) 1.22 (to build the provider plugin)

Building The Provider
---------------------

1. Install [Go](https://go.dev/dl)

2. Clone the SDK repo:

```sh
git clone https://github.com/paloaltonetworks/pango
```

3. Clone this repo:

```sh
git clone https://github.com/paloaltonetworks/terraform-provider-panos
```

4. Build the provider:

```sh
cd terraform-provider-panos
go build .
```

5. Specify the `dev_overrides` configuration per the next section below. This tells Terraform where to find the provider you just built. The directory to specify is the full path to the cloned provider repo.


Developing the Provider
-----------------------

With Terraform v1 and later, [development overrides for provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers) can be leveraged in order to use the provider built from source.

To do this, populate a Terraform CLI configuration file (`~/.terraformrc` for all platforms other than Windows; `terraform.rc` in the `%APPDATA%` directory when using Windows) with at least the following options:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/paloaltonetworks-local/panos" = "/directory/containing/the/provider/binary/here"
  }

  direct {}
}
```

Then when referencing the locally built provider, use the local name in the `terraform` configuration block like so:

```hcl
terraform {
    required_providers {
        panos = {
            source = "paloaltonetworks-local/panos"
            version = "2.0.0"
        }
    }
}
```