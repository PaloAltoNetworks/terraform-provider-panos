# A multicast PIM interface timer routing profile can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "my-template"
#       panorama_device = "localhost.localdomain"
#       ngfw_device     = "localhost.localdomain"
#     }
#   }
#
#   name = "example-pim-timer-profile"
# }
terraform import panos_multicast_pim_interface_timer_routing_profile.example $(echo '{"location":{"template":{"name":"my-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}},"name":"example-pim-timer-profile"}' | base64)
