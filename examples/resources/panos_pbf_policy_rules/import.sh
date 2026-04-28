#!/bin/bash

# A set of PBF policy rules can be imported by providing the following base64 encoded object as the ID
# {
#     location = {
#         device_group = {
#         name = "example-device-group"
#         rulebase = "pre-rulebase"
#         panorama_device = "localhost.localdomain"
#         }
#     }
#
#     position = { where = "after", directly = true, pivot = "route-voip-traffic" }
#
#     names = [
#         "route-backup-traffic",
#         "block-suspicious-traffic"
#     ]
# }
terraform import panos_pbf_policy_rules.application_routing $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["route-backup-traffic","block-suspicious-traffic"],"position":{"directly":true,"pivot":"route-voip-traffic","where":"after"}}' | base64)
