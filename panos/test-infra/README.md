# Testing Infrastructure

Here is where we keep the code of infrastructure to be used as a target for running acceptance tests of this provider.

Both [PAN-OS](https://docs.paloaltonetworks.com/compatibility-matrix/vm-series-firewalls.html) and
[Panorama](https://docs.paloaltonetworks.com/panorama/8-1/panorama-admin/set-up-panorama/set-up-the-panorama-virtual-appliance/install-the-panorama-virtual-appliance.html)
are available in many forms on many clouds and you _should_ be able to run tests against any of these.

We codified the following:

 - Panorama @ AWS
 - PAN-OS @ AWS

## Initial Configuration (`panosconfig`)

It is necessary to configure the firewall before running acceptance tests against it.
This part is a bit tricky to automate as PANOS can only be configured via custom prompt,
which means the only way to automate this is to pipe commands via stdin.

Provisioning via stdin is [currently not supported](https://github.com/hashicorp/terraform/issues/16800)
in Terraform. This is why there's a small utility `panosconfig`
which does exactly that (allows "scriptable" configuration of the firewall).

## How

Follow guides inside each directory.
