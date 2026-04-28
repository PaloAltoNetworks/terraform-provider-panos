# An OSPFv3 auth routing profile can be imported by providing the following base64 encoded object as the ID

# Import from an NGFW device
# {
#   location = {
#     ngfw = {
#       ngfw_device = "localhost.localdomain"
#     }
#   }
#
#   name = "ospfv3-ah-sha256-profile"
# }
terraform import panos_ospfv3_auth_routing_profile.example $(echo '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"name":"ospfv3-ah-sha256-profile"}' | base64)

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
#   name = "ospfv3-esp-secure-profile"
# }
terraform import panos_ospfv3_auth_routing_profile.example $(echo '{"location":{"template":{"name":"my-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"ospfv3-esp-secure-profile"}' | base64)

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
#   name = "ospfv3-esp-encrypt-only"
# }
terraform import panos_ospfv3_auth_routing_profile.example $(echo '{"location":{"template_stack":{"name":"my-template-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"ospfv3-esp-encrypt-only"}' | base64)
