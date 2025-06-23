---
page_title: 'Migration Guide'
---

Version 2.0.0 of the Terraform provider introduces significant breaking changes to the schemas, and there is no automatic state upgrade available.

Most resources support importing, even if the documentation does not include an import section. We are currently updating the documentation, and all resources will have an import section included in future versions.

We recommend using [configuration generation](https://developer.hashicorp.com/terraform/language/import/generating-configuration), which has been available since Terraform v1.5 as an experimental feature.

The import ID of a resource is the base64-encoded version of the Terraform resource configuration.

## Migration Steps

1. Backup your existing Terraform configuration and state files
2. Create a variable or a local value to reflect what you are importing

   ```hcl
    locals {
      import_address = {
        location = {
          device_group = {
            name            = "import-address"
            panorama_device = "localhost.localdomain"
          }
        }

        name = "example-address"
      }
    }
   ```

For shared objects, set location to this both for the import and resource:

   ```hcl
    location = { shared = {} }
   ```

3. Add the `import` block
   ```hcl
    import {
      id = base64encode(jsonencode(local.import_address))
      to = panos_address.example
    }
   ```
4. Plan and generate configuration

   ```bash
    terraform plan -generate-config-out=generated.tf
    ...
    Terraform will perform the following actions:

      # panos_address.example will be imported
      # (config will be generated)
        resource "panos_address" "example" {
            description = "example address 1"
            ip_netmask  = "10.0.0.1/32"
            location    = {
                device_group = {
                    name            = "import-address"
                    panorama_device = "localhost.localdomain"
                }
            }
            name        = "example-address"
        }

    Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
    ╷
    │ Warning: Config generation is experimental
    │
    │ Generating configuration during import is currently experimental, and the generated configuration format may change in future
    │ versions.
   ```

5. Review generated configuration

   ```hcl
    # __generated__ by Terraform
    # Please review these resources and move them into your main configuration files.

    # __generated__ by Terraform from "eyJsb2NhdGlvbiI6eyJkZXZpY2VfZ3JvdXAiOnsibmFtZSI6ImltcG9ydC1hZGRyZXNzIiwicGFub3JhbWFfZGV2aWNlIjoibG9jYWxob3N0LmxvY2FsZG9tYWluIn19LCJuYW1lIjoiZXhhbXBsZS1hZGRyZXNzIn0="
    resource "panos_address" "example" {
      description      = "example address 1"
      disable_override = null
      fqdn             = null
      ip_netmask       = "10.0.0.1/32"
      ip_range         = null
      ip_wildcard      = null
      location = {
        device_group = {
          name            = "import-address"
          panorama_device = "localhost.localdomain"
        }
        shared = null
        vsys   = null
      }
      name = "example-address"
      tags = null
    }
   ```

6. Apply

   ```hcl
    terraform apply
    ....
    Terraform will perform the following actions:

      # panos_address.example will be imported
        resource "panos_address" "example" {
            description = "example address 1"
            ip_netmask  = "10.0.0.1/32"
            location    = {
                device_group = {
                    name            = "import-address"
                    panorama_device = "localhost.localdomain"
                }
            }
            name        = "example-address"
        }

    Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
    panos_address.example: Importing... [id=eyJsb2NhdGlvbiI6eyJkZXZpY2VfZ3JvdXAiOnsibmFtZSI6ImltcG9ydC1hZGRyZXNzIiwicGFub3JhbWFfZGV2aWNlIjoibG9jYWxob3N0LmxvY2FsZG9tYWluIn19LCJuYW1lIjoiZXhhbXBsZS1hZGRyZXNzIn0=]
    panos_address.example: Import complete [id=eyJsb2NhdGlvbiI6eyJkZXZpY2VfZ3JvdXAiOnsibmFtZSI6ImltcG9ydC1hZGRyZXNzIiwicGFub3JhbWFfZGV2aWNlIjoibG9jYWxob3N0LmxvY2FsZG9tYWluIn19LCJuYW1lIjoiZXhhbXBsZS1hZGRyZXNzIn0=]

   ```

## Example Import IDs

### Importing an Address

An address can be imported by providing the following base64 encoded object as the ID

```hcl
{
  location = {
    device_group = {
      name            = "example-device-group"
      panorama_device = "localhost.localdomain"
  }
 }

  name = "addr1"
}
```

```bash
terraform import panos_address.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"name":"addr1"}' | base64)
```

### Importing the entire security policy rule base

The entire policy can be imported by providing the following base64 encoded object as the ID. In this instance you only have to specify the name of the first rule in the policy to import all the rules.

```hcl
{
  location = {
    device_group = {
      name = "example-device-group"
      rulebase = "pre-rulebase"
      panorama_device = "localhost.localdomain"
    }
  }

  names = [
    "rule-1", # the first rule in the policy
  ]
}
```

```bash
terraform import panos_security_policy.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["rule-1"]}' | base64)
```
