#!/bin/bash
# A VLAN static MAC entry can be imported by providing the following
# base64 encoded object as the ID.
# The import ID must include both the parent VLAN name (vlan) and the
# static MAC address (name).

# Import from a Panorama template
# {
#   "location": {
#     "template": {
#       "name":            "branch-office-template",
#       "panorama_device": "localhost.localdomain",
#       "ngfw_device":     "localhost.localdomain"
#     }
#   },
#   "vlan": "vlan-10-production",
#   "name": "00:1a:2b:3c:4d:5e"
# }
terraform import panos_vlan_entry.server_web $(echo -n '{"location":{"template":{"name":"branch-office-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"vlan":"vlan-10-production","name":"00:1a:2b:3c:4d:5e"}' | base64)

# Import from a Panorama template stack
# {
#   "location": {
#     "template_stack": {
#       "name":            "branch-office-stack",
#       "panorama_device": "localhost.localdomain",
#       "ngfw_device":     "localhost.localdomain"
#     }
#   },
#   "vlan": "vlan-10-production",
#   "name": "00:1a:2b:3c:4d:5e"
# }
terraform import panos_vlan_entry.server_web $(echo -n '{"location":{"template_stack":{"name":"branch-office-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"vlan":"vlan-10-production","name":"00:1a:2b:3c:4d:5e"}' | base64)

# Import from a standalone NGFW
# {
#   "location": {
#     "ngfw": {
#       "ngfw_device": "localhost.localdomain"
#     }
#   },
#   "vlan": "vlan-10-production",
#   "name": "00:1a:2b:3c:4d:5e"
# }
terraform import panos_vlan_entry.server_web $(echo -n '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"vlan":"vlan-10-production","name":"00:1a:2b:3c:4d:5e"}' | base64)
