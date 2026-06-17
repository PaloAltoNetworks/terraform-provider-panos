#!/bin/bash
# Import a BFD Network Profile from a Panorama template.
# The location must be base64-encoded JSON before passing to terraform import.

location='{"template":{"name":"bfd-network-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}}'
encoded_location=$(echo -n "$location" | base64)
terraform import "panos_bfd_network_profile.active" "$encoded_location:bfd-active-fast"
