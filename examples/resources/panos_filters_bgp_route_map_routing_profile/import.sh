# A BGP route map routing profile can be imported by providing the following base64 encoded object as the ID

# Import from an NGFW device
# {
#   location = {
#     ngfw = {
#       ngfw_device = "localhost.localdomain"
#     }
#   }
#
#   name = "my-route-map"
# }
terraform import panos_filters_bgp_route_map_routing_profile.example $(echo '{"location":{"ngfw":{"ngfw_device":"localhost.localdomain"}},"name":"my-route-map"}' | base64)

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
#   name = "my-route-map"
# }
terraform import panos_filters_bgp_route_map_routing_profile.example $(echo '{"location":{"template":{"name":"my-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"my-route-map"}' | base64)

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
#   name = "my-route-map"
# }
terraform import panos_filters_bgp_route_map_routing_profile.example $(echo '{"location":{"template_stack":{"name":"my-template-stack","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"my-route-map"}' | base64)
