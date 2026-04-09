#!/bin/bash
# A logical router can be imported by providing the following base64 encoded object as the ID

# Import from an NGFW vsys
# {
#   location = {
#     vsys = {
#       name        = "vsys1"
#       ngfw_device = "localhost.localdomain"
#     }
#   }
#
#   name = "lr-basic"
# }
terraform import panos_logical_router.example $(echo '{"location":{"vsys":{"name":"vsys1","ngfw_device":"localhost.localdomain"}},"name":"lr-basic"}' | base64)

# Import from a Panorama template
# {
#   location = {
#     template = {
#       name            = "my-template"
#       panorama_device = "localhost.localdomain"
#       ngfw_device     = "localhost.localdomain"
#     }
#   }
#
#   name = "lr-basic"
# }
terraform import panos_logical_router.example $(echo '{"location":{"template":{"name":"my-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"lr-basic"}' | base64)

# Import from a Panorama template stack
# {
#   location = {
#     template_stack = {
#       name            = "my-template-stack"
#       panorama_device = "localhost.localdomain"
#       ngfw_device     = "localhost.localdomain"
#     }
#   }
#
#   name = "lr-basic"
# }
terraform import panos_logical_router.example $(echo '{"location":{"template_stack":{"name":"my-template-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"lr-basic"}' | base64)

# Import from an NGFW device (device-level configuration)
# {
#   location = {
#     ngfw = {
#       ngfw_device = "localhost.localdomain"
#     }
#   }
#
#   name = "lr-basic"
# }
terraform import panos_logical_router.example $(echo '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"name":"lr-basic"}' | base64)
