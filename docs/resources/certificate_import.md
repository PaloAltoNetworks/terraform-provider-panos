---
page_title: "panos: panos_certificate_import"
subcategory: "Device"
---

# panos_certificate_import

This resource allows you to import/update/delete a PEM or PKCS12 certificate.

-> **NOTE:**  Importing into a Template vsys that isn't shared does not work right now.


## PAN-OS

NGFW and Panorama.


## Import Name

Encrypted fields prevent Terraform from importing this resource.


## Example Usage

```hcl
# A PEM style cert.
resource "panos_certificate_import" "example" {
    name = "tfcert"
    pem {
        certificate = file("cert.pem")
        private_key = file("key.pem")
        passphrase = "secret"
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

```hcl
# A PKCS12 style cert.
resource "panos_certificate_import" "example2" {
    name = "tfcert2"
    pkcs12 {
        certificate = file("cert.pfx")
        passphrase = "foobar"
    }

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `template` - The template.

NGFW / Panorama:

* `vsys` - The vsys (default: `shared`).

The following arguments are supported:

* `name` - (Required) The name.
* `pem` - A PEM style certificate, as defined below. Conflicts with `pkcs12`.
* `pkcs12` - A PKCS12 style certificate, as defined below. Conflicts with `pem`.

`pem` supports the following arguments:

* `certificate` - (Required) The contents of the certificate file.
* `certificate_filename` - The certificate filename for uploading to
  PAN-OS (default: `cert.pem`).
* `private_key` - (Required) The contents of the private key file.
* `private_key_filename` - The private key filename for uplaoding to
  PAN-OS (default: `key.pem`).
* `passphrase` - (Required) The private key file passphrase.

`pkcs12` supports the following arguments:

* `certificate` - (Required) The contents of the certificate file.
* `certificate_filename` - The certificate filename for uploading to
  PAN-OS (default: `cert.pfx`).
* `passphrase` - (Required) The private key file passphrase.


## Attribute Reference

The following attributes are supported:

* `cert_format` - The certificate format.
* `common_name` - The common name.
* `algorithm` - The algorithm.
* `ca` - The CA.
* `not_valid_after` - Certificate is not valid after this date.
* `not_valid_before` - Certificate is not valid before this date.
* `expiry_epoch` - Expiry ephoch.
* `subject` - Subject.
* `subject_hash` - The subject hash.
* `issuer` - Certificate issuer.
* `issuer_hash` - Hash of the issuer.
* `csr` - The CSR.
* `public_key` - Public key.
* `private_key` - Encrypted private key.
* `private_key_on_hsm` - (bool) Private key on HSM.
* `status` - The status.
* `revoke_date_epoch` - Revoke date epoch.
