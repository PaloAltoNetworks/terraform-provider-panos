#!/bin/bash

# The entire PBF policy can be imported by providing the following base64 encoded object as the ID
# {
#     location = {
#         device_group = {
#         name = "example-device-group"
#         rulebase = "pre-rulebase"
#         panorama_device = "localhost.localdomain"
#         }
#     }
#
#
#     names = [
#         "route-guest-traffic", <- all rule names in the policy must be listed
#         "route-internal-traffic",
#     ]
# }
terraform import panos_pbf_policy.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["route-guest-traffic","route-internal-traffic"]}' | base64)
