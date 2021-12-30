package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaNatRule() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaNatRule,
		Read:   readPanoramaNatRule,
		Update: updatePanoramaNatRule,
		Delete: deletePanoramaNatRule,

		DeprecationMessage: "Please use panos_panorama_nat_rule_group instead",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
				ForceNew: true,
			},
			"rulebase": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      util.PreRulebase,
				ValidateFunc: validateStringIn(util.Rulebase, util.PreRulebase, util.PostRulebase),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ipv4",
				Description:  "NAT type (ipv4 default, nat64, or nptv6)",
				ValidateFunc: validateStringIn("ipv4", "nat64", "nptv6"),
			},
			"source_zones": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
			},
			"source_addresses": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_addresses": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sat_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				Description:  "none (default), dynamic-ip-and-port, dynamic-ip, or static-ip",
				ValidateFunc: validateStringIn("none", "dynamic-ip-and-port", "dynamic-ip", "static-ip"),
			},
			"sat_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "interface-address or translated-address",
				ValidateFunc: validateStringIn("interface-address", "translated-address"),
			},
			"sat_translated_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sat_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_fallback_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("none", "interface-address", "translated-address"),
			},
			"sat_fallback_translated_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sat_fallback_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_fallback_ip_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("ip", "floating"),
			},
			"sat_fallback_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_static_translated_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_static_bi_directional": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dat_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateStringIn("static", "dynamic"),
				Optional:     true,
			},
			"dat_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dat_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dat_dynamic_distribution": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target":        targetSchema(false),
			"negate_target": negateTargetSchema(),
		},
	}
}

func parsePanoramaNatRule(d *schema.ResourceData) (string, string, nat.Entry) {
	dg := d.Get("device_group").(string)
	rb := d.Get("rulebase").(string)

	o := nat.Entry{
		Name:                           d.Get("name").(string),
		Type:                           d.Get("type").(string),
		Description:                    d.Get("description").(string),
		SourceZones:                    asStringList(d.Get("source_zones").([]interface{})),
		DestinationZone:                d.Get("destination_zone").(string),
		ToInterface:                    d.Get("to_interface").(string),
		Service:                        d.Get("service").(string),
		SourceAddresses:                asStringList(d.Get("source_addresses").([]interface{})),
		DestinationAddresses:           asStringList(d.Get("destination_addresses").([]interface{})),
		SatType:                        d.Get("sat_type").(string),
		SatAddressType:                 d.Get("sat_address_type").(string),
		SatTranslatedAddresses:         asStringList(d.Get("sat_translated_addresses").([]interface{})),
		SatInterface:                   d.Get("sat_interface").(string),
		SatIpAddress:                   d.Get("sat_ip_address").(string),
		SatFallbackType:                d.Get("sat_fallback_type").(string),
		SatFallbackTranslatedAddresses: asStringList(d.Get("sat_fallback_translated_addresses").([]interface{})),
		SatFallbackInterface:           d.Get("sat_fallback_interface").(string),
		SatFallbackIpType:              d.Get("sat_fallback_ip_type").(string),
		SatFallbackIpAddress:           d.Get("sat_fallback_ip_address").(string),
		SatStaticTranslatedAddress:     d.Get("sat_static_translated_address").(string),
		SatStaticBiDirectional:         d.Get("sat_static_bi_directional").(bool),
		DatAddress:                     d.Get("dat_address").(string),
		DatPort:                        d.Get("dat_port").(int),
		DatDynamicDistribution:         d.Get("dat_dynamic_distribution").(string),
		Disabled:                       d.Get("disabled").(bool),
		Tags:                           asStringList(d.Get("tags").([]interface{})),
		Targets:                        loadTarget(d.Get("target")),
		NegateTarget:                   d.Get("negate_target").(bool),
	}

	switch d.Get("dat_type").(string) {
	case "static":
		o.DatType = nat.DatTypeStatic
	case "dynamic":
		o.DatType = nat.DatTypeDynamic
	}

	return dg, rb, o
}

func parsePanoramaNatRuleId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaNatRuleId(a, b, c string) string {
	return fmt.Sprintf("%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c)
}

func createPanoramaNatRule(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, rb, o := parsePanoramaNatRule(d)

	if err := pano.Policies.Nat.Set(dg, rb, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaNatRuleId(dg, rb, o.Name))
	return readPanoramaNatRule(d, meta)
}

func readPanoramaNatRule(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, name := parsePanoramaNatRuleId(d.Id())

	o, err := pano.Policies.Nat.Get(dg, rb, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("device_group", dg)
	d.Set("rulebase", rb)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	if err = d.Set("source_zones", o.SourceZones); err != nil {
		log.Printf("[WARN] Error setting 'source_zones' param for %q: %s", d.Id(), err)
	}
	d.Set("destination_zone", o.DestinationZone)
	d.Set("to_interface", o.ToInterface)
	d.Set("service", o.Service)
	if err = d.Set("source_addresses", o.SourceAddresses); err != nil {
		log.Printf("[WARN] Error setting 'source_addresses' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("destination_addresses", o.DestinationAddresses); err != nil {
		log.Printf("[WARN] Error setting 'destination_addresses' param for %q: %s", d.Id(), err)
	}
	d.Set("sat_type", o.SatType)
	d.Set("sat_address_type", o.SatAddressType)
	if err = d.Set("sat_translated_addresses", o.SatTranslatedAddresses); err != nil {
		log.Printf("[WARN] Error setting 'sat_translated_addresses' param for %q: %s", d.Id(), err)
	}
	d.Set("sat_interface", o.SatInterface)
	d.Set("sat_ip_address", o.SatIpAddress)
	d.Set("sat_fallback_type", o.SatFallbackType)
	if err = d.Set("sat_fallback_translated_addresses", o.SatFallbackTranslatedAddresses); err != nil {
		log.Printf("[WARN] Error setting 'sat_fallback_translated_addresses' param for %q: %s", d.Id(), err)
	}
	d.Set("sat_fallback_interface", o.SatFallbackInterface)
	d.Set("sat_fallback_ip_type", o.SatFallbackIpType)
	d.Set("sat_fallback_ip_address", o.SatFallbackIpAddress)
	d.Set("sat_static_translated_address", o.SatStaticTranslatedAddress)
	d.Set("sat_static_bi_directional", o.SatStaticBiDirectional)
	switch o.DatType {
	case nat.DatTypeStatic:
		d.Set("dat_type", "static")
	case nat.DatTypeDynamic:
		d.Set("dat_type", "dynamic")
	}
	d.Set("dat_address", o.DatAddress)
	d.Set("dat_port", o.DatPort)
	d.Set("dat_dynamic_distribution", o.DatDynamicDistribution)
	d.Set("disabled", o.Disabled)
	if err = d.Set("tags", o.Tags); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("target", dumpTarget(o.Targets)); err != nil {
		log.Printf("[WARN] Error setting 'target' param for %q: %s", d.Id(), err)
	}
	d.Set("negate_target", o.NegateTarget)

	return nil
}

func updatePanoramaNatRule(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, o := parsePanoramaNatRule(d)

	lo, err := pano.Policies.Nat.Get(dg, rb, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Policies.Nat.Edit(dg, rb, lo); err != nil {
		return err
	}

	return readPanoramaNatRule(d, meta)
}

func deletePanoramaNatRule(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, rb, name := parsePanoramaNatRuleId(d.Id())

	err := pano.Policies.Nat.Delete(dg, rb, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
