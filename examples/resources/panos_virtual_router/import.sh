#!/bin/bash
# Import a virtual router from a template
location='{"template":{"name":"example-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}}'
encoded_location=$(echo -n "$location" | base64)
terraform import "panos_virtual_router.example" "$encoded_location:production-vr"
