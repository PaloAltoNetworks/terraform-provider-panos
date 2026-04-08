# A log export schedule can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "ftp-export-schedule"
# }
terraform import panos_log_export_schedule.ftp_example $(echo '{"location":{"template":{"name":"example-template","panorama_device":"localhost.localdomain"}},"name":"ftp-export-schedule"}' | base64)
