# Testing Infrastructure

Here is where we keep the code of infrastructure to be used as a target for running acceptance tests of this provider.
[PANOS NGFW is available](https://docs.paloaltonetworks.com/compatibility-matrix/vm-series-firewalls.html)
in many forms on many clouds and you _should_ be able to run tests against any of these.
We codified the following:

 - AWS

## Initial Configuration

It is necessary to configure the firewall before running acceptance tests against it.
This part is a bit tricky to automate as PANOS can only be configured via custom prompt,
which means the only way to automate this is to pipe commands via stdin.

Provisioning via stdin is [currently not supported](https://github.com/hashicorp/terraform/issues/16800)
in Terraform. This is why there's a small utility `panosconfig`
which does exactly that (allows "scriptable" configuration of the firewall).

## How

```sh
cd aws
terraform apply -var=panos_version=9.0.0 -var=panos_username=terraform -var=panos_password=terraformpass
```

Export variables

```sh
export PANOS_HOSTNAME=$(terraform output hostname)
export PANOS_USERNAME=$(terraform output username)
export PANOS_PASSWORD=$(terraform output password)
export PANOS_SSH_PRIVATE_KEY="$(terraform output ssh_private_key)"
```

It takes aproximately 13 minutes (after boot) until the instance
is ready to be **configured** and prompt is available.

```sh
go run ../panosconfig
```

Run acceptance tests after the instance was configured.

```sh
cd ../..
make testacc
```

```sh
cd test-infra/aws
terraform destroy
```
