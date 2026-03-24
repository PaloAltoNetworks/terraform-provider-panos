#!/bin/bash

# This script is used to import a panos_globalprotect_gateway resource.
#
# The following arguments are required:
#   - name: the name of the resource
#   - location: the location of the resource (e.g. "template", "template-stack", "vsys")
#
# The following arguments are optional:
#   - panorama_device: the panorama device name (default: "localhost.localdomain")
#   - template: the template name
#   - template_stack: the template stack name
#   - ngfw_device: the ngfw device name (default: "localhost.localdomain")
#   - vsys: the vsys name (default: "vsys1")
#
# Example:
# ./import.sh -name "gw-example" -location "template" -template "tf-example-gw"

# a string with all the arguments
ARGS=$@

terraform import panos_globalprotect_gateway.gw $ARGS
