#!/bin/bash

# An OSPFv3 redistribution routing profile can be imported by providing a base64 encoded JSON object as the ID

# Import from an NGFW device
# {
#   "location": {
#     "ngfw": {
#       "ngfw_device": "localhost.localdomain"
#     }
#   },
#   "name": "ospfv3-redistribute-connected"
# }
terraform import panos_ospfv3_redistribution_routing_profile.connected_routes $(echo '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"name":"ospfv3-redistribute-connected"}' | base64)

# Import from a Panorama template
# {
#   "location": {
#     "template": {
#       "name": "production-template",
#       "panorama_device": "localhost.localdomain",
#       "ngfw_device": "localhost.localdomain"
#     }
#   },
#   "name": "ospfv3-redistribute-connected"
# }
terraform import panos_ospfv3_redistribution_routing_profile.connected_routes $(echo '{"location":{"template":{"name":"production-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"ospfv3-redistribute-connected"}' | base64)

# Import from a Panorama template stack
# {
#   "location": {
#     "template_stack": {
#       "name": "production-stack",
#       "panorama_device": "localhost.localdomain",
#       "ngfw_device": "localhost.localdomain"
#     }
#   },
#   "name": "ospfv3-redistribute-connected"
# }
terraform import panos_ospfv3_redistribution_routing_profile.connected_routes $(echo '{"location":{"template_stack":{"name":"production-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"ospfv3-redistribute-connected"}' | base64)
