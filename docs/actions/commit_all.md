---
page_title: "panos_commit_all Action - panos"
subcategory: ""
description: |-
  Commits pending changes and pushes them to all managed devices in Panorama (commit-all operation).
---

# panos_commit_all (Action)

The `panos_commit_all` action performs a commit-all operation on Panorama, which commits pending configuration changes and pushes them to all managed firewalls.

## Overview

This action is **Panorama-specific** and pushes committed configuration changes to all managed devices (firewalls).

**Important**: This action does NOT commit changes to Panorama. It only pushes already-committed configuration. You must run `panos_commit` first to commit your changes to Panorama, then run `panos_commit_all` to push those committed changes to managed devices.

### Typical Workflow

1. Apply Terraform resources (creates pending changes in Panorama)
2. Run `panos_commit` action (commits pending changes to Panorama)
3. Run `panos_commit_all` action (pushes committed configuration to devices)

This is equivalent to:
- Running "Commit" in Panorama (step 2)
- Then clicking "Push to Devices" (step 3)

## When to Use

Use `panos_commit_all` when:

- You're managing a Panorama instance (not a standalone firewall)
- You want to activate configuration changes on managed firewalls
- You need to push device group or template configurations to devices
- You're deploying changes that should be applied across your firewall fleet

## Comparison with panos_commit

| Action | Purpose | Target | What it does |
|--------|---------|--------|--------------|
| `panos_commit` | Commits pending changes | Panorama or Firewall | Saves configuration changes to the device |
| `panos_commit_all` | Pushes committed config | Panorama only | Distributes saved configuration to managed devices |

**Key Differences**:
- `panos_commit`: Takes pending changes and commits them to Panorama (or firewall)
- `panos_commit_all`: Takes committed configuration and pushes it to device groups

**Important**:
- You typically run `panos_commit` BEFORE `panos_commit_all`
- Running `panos_commit_all` without pending committed changes will succeed but push nothing
- Running `panos_commit_all` on a standalone firewall will fail (Panorama only)

## Prerequisites

- Provider must be configured to connect to a Panorama instance
- User must have commit permissions in Panorama
- At least one device group must exist in Panorama

## Configuration

### device_groups (Optional)

You can optionally specify which device groups to push configuration to:

```terraform
action "panos_commit_all" "push_to_specific" {
  config {
    device_groups = ["Production", "DMZ"]
  }
}
```

- **Type**: List of strings
- **Default**: If not specified, pushes to all device groups in Panorama
- **Example**: `["Production", "DMZ", "Branch-Offices"]`

**When to use**:
- Push to specific device groups only
- Avoid pushing to all device groups when testing
- Target specific environments (production vs staging)

## Example Usage

### Declarative Action Block (Recommended)

For users who prefer to define actions in their Terraform configuration instead of using CLI commands:

```terraform
terraform {
  required_providers {
    panos = {
      source  = "PaloAltoNetworks/panos"
      version = "~> 2.0"
    }
  }
}

provider "panos" {
  hostname = "panorama.example.com"
  username = "admin"
  password = var.panos_password
}

# Define your resources
resource "panos_device_group" "production" {
  location = {
    panorama = {}
  }
  name        = "Production-Firewalls"
  description = "Production environment firewall group"
}

resource "panos_security_policy" "allow_web" {
  location = {
    device_group = {
      name = panos_device_group.production.name
    }
  }

  rule {
    name                  = "Allow-Web-Traffic"
    source_zones          = ["trust"]
    destination_zones     = ["untrust"]
    source_addresses      = ["any"]
    destination_addresses = ["any"]
    applications          = ["web-browsing", "ssl"]
    services              = ["application-default"]
    action                = "allow"
  }
}

# Define the commit-all action in your configuration
action "panos_commit_all" "push_to_devices" {
  # No config needed - pushes to all device groups
}

# Or push to specific device groups
action "panos_commit_all" "push_to_production" {
  config {
    device_groups = ["Production-Firewalls"]
  }
}

# Now you can invoke it declaratively:
# terraform apply
# terraform action invoke push_to_devices
```

### CLI-Based Invocation with terraform apply

You can invoke actions as part of the apply operation using the `-invoke` flag:

```terraform
provider "panos" {
  hostname = "panorama.example.com"
  username = "admin"
  password = var.panos_password
}

resource "panos_address" "web_server" {
  location = {
    device_group = {
      name = "Production"
    }
  }

  name       = "web-server-01"
  ip_netmask = "192.168.1.100/32"
}

# Define the action
action "panos_commit_all" "deploy" {
}

# After applying resources, invoke during apply:
# terraform apply -invoke=action.panos_commit_all.deploy
```

### Targeted Device Group Push

```terraform
provider "panos" {
  hostname = "panorama.example.com"
  username = "admin"
  api_key  = var.panos_api_key
}

# Configure resources for DMZ device group
resource "panos_address" "dmz_server" {
  location = {
    device_group = {
      name = "DMZ-Firewalls"
    }
  }

  name       = "dmz-web-server"
  ip_netmask = "192.168.100.10/32"
}

# Define actions - one for all groups, one for DMZ only
action "panos_commit" "commit_to_panorama" {
}

action "panos_commit_all" "push_all" {
  # No config = pushes to all device groups
}

action "panos_commit_all" "push_dmz_only" {
  config {
    device_groups = ["DMZ-Firewalls"]
  }
}

# Workflow:
# terraform apply
# terraform action invoke commit_to_panorama
# terraform action invoke push_dmz_only  # Only pushes to DMZ
```

### Multi-Stage Deployment Workflow

This example shows the recommended workflow: commit to Panorama first, then push to devices.

```terraform
# Create your resources
resource "panos_address_group" "web_servers" {
  location = {
    shared = {}
  }

  name           = "Web-Servers"
  static_members = ["web1", "web2", "web3"]
}

resource "panos_security_policy" "branch_policy" {
  location = {
    device_group = {
      name = "Branch-Offices"
    }
  }

  rule {
    name                  = "Branch-Internet-Access"
    source_zones          = ["inside"]
    destination_zones     = ["outside"]
    source_addresses      = ["any"]
    destination_addresses = ["any"]
    applications          = ["web-browsing"]
    action                = "allow"
  }
}

# Define actions for the deployment workflow
action "panos_commit" "commit_to_panorama" {
}

action "panos_commit_all" "push_to_devices" {
}

# Deployment workflow:
# 1. terraform apply                              # Apply resources (creates pending changes)
# 2. terraform action invoke commit_to_panorama   # Commit changes to Panorama
# 3. (Optional) Review in Panorama UI
# 4. terraform action invoke push_to_devices      # Push to managed devices
```

### Automated Workflow with Lifecycle Hooks

This example uses Terraform lifecycle hooks to automatically trigger both commit and push actions:

```terraform
provider "panos" {
  hostname = "panorama.example.com"
  username = "admin"
  password = var.panos_password
}

resource "panos_security_policy" "app_policy" {
  location = {
    device_group = {
      name = "Production"
    }
  }

  rules = [{
    name                  = "Allow-App-Traffic"
    source_zones          = ["trust"]
    destination_zones     = ["untrust"]
    source_addresses      = ["any"]
    destination_addresses = ["any"]
    applications          = ["web-browsing"]
    action                = "allow"
  }]

  # Automatically commit and push after create or update
  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [
        action.panos_commit.commit_to_panorama,
        action.panos_commit_all.push_to_devices
      ]
    }
  }
}

action "panos_commit" "commit_to_panorama" {
}

action "panos_commit_all" "push_to_devices" {
}

# With this setup, running terraform apply will:
# 1. Create/update the security policy
# 2. Automatically commit to Panorama
# 3. Automatically push to devices
```

### Template Stack Push Example

```terraform
provider "panos" {
  hostname = "panorama.example.com"
  username = "admin"
  password = var.panos_password

  # Target template stack
  target = "Branch-Template-Stack"
}

# Configure network settings in template
resource "panos_ethernet_interface" "wan" {
  location = {
    template = {
      name = "Branch-Template"
    }
  }

  name = "ethernet1/1"
  mode = "layer3"

  layer3 {
    ipv4 {
      static {
        ip_addresses = ["10.0.0.1/24"]
      }
    }
  }
}

# Define actions for manual invocation
action "panos_commit" "commit_template" {
}

action "panos_commit_all" "push_template" {
}

# Manual workflow:
# terraform apply
# terraform action invoke commit_template
# terraform action invoke push_template
```

### Conditional Push Based on Environment

```terraform
variable "environment" {
  description = "Deployment environment"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "production"], var.environment)
    error_message = "Environment must be dev, staging, or production."
  }
}

provider "panos" {
  hostname = var.environment == "production" ? "panorama-prod.example.com" : "panorama-dev.example.com"
  username = "admin"
  api_key  = var.panos_api_key
  target   = "${var.environment}-firewalls"
}

resource "panos_security_policy" "app_policy" {
  location = {
    device_group = {
      name = "${var.environment}-firewalls"
    }
  }

  rule {
    name                  = "Allow-App-Traffic"
    source_zones          = ["trust"]
    destination_zones     = ["untrust"]
    source_addresses      = ["any"]
    destination_addresses = ["any"]
    applications          = ["web-browsing"]
    action                = "allow"
  }
}

# Environment-specific action
action "panos_commit_all" "push_${var.environment}" {
}

# Usage:
# terraform apply -var="environment=dev"
# terraform action invoke push_dev
```

## Usage in CI/CD Pipelines

### Declarative Approach in CI/CD

```hcl
# main.tf
terraform {
  required_providers {
    panos = {
      source  = "PaloAltoNetworks/panos"
      version = "~> 2.0"
    }
  }
}

provider "panos" {
  hostname = var.panorama_host
  api_key  = var.panorama_api_key
}

# Your resources here
resource "panos_address" "servers" {
  # ... configuration ...
}

# Define commit action
action "panos_commit_all" "deploy" {
}
```

```bash
#!/bin/bash
# deploy.sh - Automated deployment script

set -e

# Step 1: Initialize Terraform
terraform init

# Step 2: Plan the changes
terraform plan -out=tfplan

# Step 3: Apply configuration changes to Panorama (creates pending changes)
terraform apply tfplan

# Step 4: Commit changes to Panorama
echo "Committing configuration to Panorama..."
terraform action invoke commit_to_panorama

# Step 5: Push committed configuration to managed devices
echo "Pushing configuration to managed devices..."
terraform action invoke push_to_devices

echo "Deployment complete!"
```

### CLI-Based CI/CD Approach

```bash
#!/bin/bash
# deploy.sh - Production deployment with manual review checkpoint

set -e

# Step 1: Plan the changes
terraform plan -out=tfplan

# Step 2: Apply configuration changes to Panorama (creates pending changes)
terraform apply tfplan

# Step 3: Commit to Panorama
echo "Committing changes to Panorama..."
terraform action invoke commit_to_panorama

# Step 4: Manual review checkpoint for production
if [ "$ENVIRONMENT" = "production" ]; then
  echo "Configuration committed to Panorama."
  echo "Review changes in Panorama UI before pushing to devices."
  echo "Continue with push? (yes/no)"
  read -r response
  if [ "$response" != "yes" ]; then
    echo "Push cancelled. Changes are committed to Panorama but not pushed to devices."
    exit 0
  fi
fi

# Step 5: Push committed configuration to all managed devices
echo "Pushing configuration to devices..."
terraform action invoke push_to_devices

echo "Deployment complete!"
```

## Invocation Methods Comparison

Both methods require defining an `action` block in your Terraform configuration first.

| Method | Command | When to Use |
|--------|---------|-------------|
| **Standalone Invoke** | `terraform action invoke <name>` | Invoke action independently after resources are already applied |
| **Apply + Invoke** | `terraform apply -invoke=action.panos_commit_all.<name>` | Apply resources and invoke action in a single command |

**Note**: Both methods reference the same action definition in your configuration.

## Behavior and Notes

### Synchronous Operation

The action waits for the commit-all job to complete before returning. For large deployments with many devices, this can take several minutes.

```
Polling interval: 2 seconds
Timeout: Follows client configuration (typically 10 minutes)
```

### Error Handling

If the commit-all operation fails on any device:
- The action returns an error with details from the Panorama API
- Partial commits may have succeeded on some devices
- Check Panorama's commit status to see which devices succeeded/failed

### Target Scope

The `target` parameter in the provider configuration determines the scope:

| Target Value | Behavior |
|--------------|----------|
| Not specified | Pushes all pending changes to all device groups/templates |
| `"DeviceGroup-Name"` | Pushes only changes for that specific device group |
| `"Template-Name"` | Pushes only changes for that specific template |

### Best Practices

1. **Test in Non-Production First**: Always test commit-all operations in a development/staging Panorama before production

2. **Review Before Push**: Review pending changes in Panorama UI before invoking the action

3. **Use Declarative Actions**: Define actions in your Terraform configuration for better version control and reusability

4. **Commit Locally First**: Consider running `panos_commit` first to commit to Panorama, then `panos_commit_all` to push to devices

5. **Monitor Device Status**: After the action completes, verify device commit status in Panorama

6. **Handle Failures Gracefully**: Implement error handling in automation scripts to catch and report commit failures

7. **Use Maintenance Windows**: Schedule commit-all operations during maintenance windows to minimize impact

## Troubleshooting

### "Object not found" Error

**Cause**: Trying to run commit-all on a standalone firewall

**Solution**: Only use this action with Panorama. Use `panos_commit` for standalone firewalls.

### "No changes to commit" Warning

**Cause**: No pending configuration changes in Panorama

**Solution**: This is informational. The action completes successfully with no changes pushed.

### Timeout Errors

**Cause**: Commit-all taking longer than the configured timeout (many devices, slow network)

**Solution**: Increase the provider timeout or check device connectivity in Panorama.

### Partial Commit Failures

**Cause**: Some devices failed to receive the configuration push

**Solution**:
1. Check device connectivity in Panorama
2. Review commit logs in Panorama → Commit Status
3. Re-run the action after resolving device issues

### Action Not Found

**Cause**: Using `terraform action invoke` with older Terraform versions

**Solution**: Ensure you're using Terraform 1.9+ which supports the action invocation syntax, or use CLI-based invocation instead.

## API Details

The action generates the following PAN-OS XML API call:

```
GET/POST https://panorama.example.com/api/
  ?type=commit
  &action=all
  &cmd=<commit-all></commit-all>
  &target=<device-group-or-template>
  &key=<api-key>
```

Response includes a job ID that is polled until completion.

## Schema

This action takes no input parameters. All configuration is inherited from the provider.

<!-- schema generated by tfplugindocs -->
## Attributes

This action has no configurable attributes. The action inherits its configuration from the provider block, including:

- Connection details (hostname, credentials)
- Target device group or template
- Timeout settings

## See Also

- [`panos_commit`](commit.md) - Commit changes locally without pushing to devices
- [Panorama Administrator's Guide - Commit and Push](https://docs.paloaltonetworks.com/panorama) - Official documentation
- [Provider Configuration](../index.md#target) - How to configure the target parameter
- [Terraform Actions Documentation](https://developer.hashicorp.com/terraform/language/actions) - Understanding Terraform actions
