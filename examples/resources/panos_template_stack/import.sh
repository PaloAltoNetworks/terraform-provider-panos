# A template stack can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     panorama = {
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-template-stack"
# }
terraform import panos_template_stack.example $(echo '{"location":{"panorama":{"panorama_device":"localhost.localdomain"}},"name":"example-template-stack"}' | base64)
