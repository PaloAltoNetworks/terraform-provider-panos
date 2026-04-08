# A local user can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template_vsys = {
#       template        = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-user"
# }
terraform import panos_local_user.example $(echo '{"location":{"template_vsys":{"template":"example-template","panorama_device":"localhost.localdomain"}},"name":"example-user"}' | base64)
