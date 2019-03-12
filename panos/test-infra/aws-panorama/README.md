# aws-panorama

Infrastructure for Panorama.

## How to Use

Create infrastructure

```sh
terraform apply -var=panorama_version=8.1.2
```

Export variables

```sh
export PANOS_HOSTNAME=$(terraform output hostname)
export PANOS_SSH_PRIVATE_KEY="$(terraform output ssh_private_key)"
```

It takes aproximately 25 minutes (after boot) until the instance
is ready to be **configured** and prompt is available.

```sh
PANOS_USERNAME=terraform PANOS_PASSWORD=terraformpass go run ../panosconfig
```

Run acceptance tests after the instance was configured.

```sh
INFRA_DIR=$(pwd)
cd ../../..
make testacc
cd $INFRA_DIR
```

Destroy infrastructure

```sh
terraform destroy
```
