#!/bin/bash
# A VLAN can be imported by providing the following base64 encoded object as the ID.

# Import from a Panorama template
# {
#   "location": {
#     "template": {
#       "name":            "branch-office-template",
#       "panorama_device": "localhost.localdomain",
#       "ngfw_device":     "localhost.localdomain"
#     }
#   },
#   "name": "vlan-10-production"
# }
terraform import panos_vlan.production $(echo -n '{"location":{"template":{"name":"branch-office-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"vlan-10-production"}' | base64)

# Import from a Panorama template stack
# {
#   "location": {
#     "template_stack": {
#       "name":            "branch-office-stack",
#       "panorama_device": "localhost.localdomain",
#       "ngfw_device":     "localhost.localdomain"
#     }
#   },
#   "name": "vlan-10-production"
# }
terraform import panos_vlan.production $(echo -n '{"location":{"template_stack":{"name":"branch-office-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"vlan-10-production"}' | base64)

# Import from a standalone NGFW
# {
#   "location": {
#     "ngfw": {
#       "ngfw_device": "localhost.localdomain"
#     }
#   },
#   "name": "vlan-10-production"
# }
terraform import panos_vlan.production $(echo -n '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"name":"vlan-10-production"}' | base64)
