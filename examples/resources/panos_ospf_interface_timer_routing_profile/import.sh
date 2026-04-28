#!/bin/bash
# Import an OSPF interface timer routing profile from a template
location='{"template":{"name":"ospf-routing-template","panorama_device":"localhost.localdomain","ngfw_device":"localhost.localdomain"}}'
encoded_location=$(echo -n "$location" | base64)
terraform import "panos_ospf_interface_timer_routing_profile.custom_timers" "$encoded_location:custom-if-timer-profile"
