#!/bin/bash
# An OSPF authentication routing profile can be imported by providing the following base64 encoded object as the ID

# Import from an NGFW device
# {
#   location = {
#     ngfw = {
#       ngfw_device = "localhost.localdomain"
#     }
#   }
#
#   name = "ospf-simple-password"
# }
terraform import panos_ospf_auth_routing_profile.example $(echo '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"name":"ospf-simple-password"}' | base64)

# Import from a Panorama template
# {
#   location = {
#     template = {
#       name            = "ospf-routing-template"
#       panorama_device = "localhost.localdomain"
#       ngfw_device     = "localhost.localdomain"
#     }
#   }
#
#   name = "ospf-md5-auth"
# }
terraform import panos_ospf_auth_routing_profile.example $(echo '{"location":{"template":{"name":"ospf-routing-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"ospf-md5-auth"}' | base64)

# Import from a Panorama template stack
# {
#   location = {
#     template_stack = {
#       name            = "ospf-routing-stack"
#       panorama_device = "localhost.localdomain"
#       ngfw_device     = "localhost.localdomain"
#     }
#   }
#
#   name = "ospf-md5-auth"
# }
terraform import panos_ospf_auth_routing_profile.example $(echo '{"location":{"template_stack":{"name":"ospf-routing-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"ospf-md5-auth"}' | base64)
