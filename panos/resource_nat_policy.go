package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNatPolicy() *schema.Resource {
	return &schema.Resource{
		Create: createNatPolicy,
		Read:   readNatPolicy,
		Update: updateNatPolicy,
		Delete: deleteNatPolicy,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this object in (default: vsys1)",
			},
			"rulebase": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rulebase",
				ForceNew:     true,
				Description:  "The rulebase (default: rulebase, pre-rulebase, post-rulebase)",
				ValidateFunc: validateStringIn("rulebase", "pre-rulebase", "post-rulebase"),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "NAT type (ipv4 default, nat64, or nptv6)",
				ValidateFunc: validateStringIn("ipv4", "nat64", "nptv6"),
			},
			"source_zone": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"to_interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
			},
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_address": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_address": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sat_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "none (default), dynamic-ip-and-port, dynamic-ip, or static-ip",
				ValidateFunc: validateStringIn("none", "dynamic-ip-and-port", "dynamic-ip", "static-ip"),
			},
			"sat_address_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "interface-address or translated-address",
				ValidateFunc: validateStringIn("interface-address", "translated-address"),
			},
			"sat_translated_address": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sat_interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_fallback_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("none", "interface-address", "translated-address"),
			},
			"sat_fallback_translated_address": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sat_fallback_interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_fallback_ip_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("ip", "floating"),
			},
			"sat_fallback_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_static_translated_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sat_static_bi_directional": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dat_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dat_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parseNatPolicy(d *schema.ResourceData) (string, string, nat.Entry) {
	vsys := d.Get("vsys").(string)
	rb := d.Get("rulebase").(string)

	o := nat.Entry{
		Name:                         d.Get("name").(string),
		Type:                         d.Get("type").(string),
		Description:                  d.Get("description").(string),
		SourceZone:                   asStringList(d, "source_zone"),
		DestinationZone:              d.Get("destination_zone").(string),
		ToInterface:                  d.Get("to_interface").(string),
		Service:                      d.Get("service").(string),
		SourceAddress:                asStringList(d, "source_address"),
		DestinationAddress:           asStringList(d, "destination_address"),
		SatType:                      d.Get("sat_type").(string),
		SatAddressType:               d.Get("sat_address_type").(string),
		SatTranslatedAddress:         asStringList(d, "sat_translated_address"),
		SatInterface:                 d.Get("sat_interface").(string),
		SatIpAddress:                 d.Get("sat_ip_address").(string),
		SatFallbackType:              d.Get("sat_fallback_type").(string),
		SatFallbackTranslatedAddress: asStringList(d, "sat_fallback_translated_address"),
		SatFallbackInterface:         d.Get("sat_fallback_interface").(string),
		SatFallbackIpType:            d.Get("sat_fallback_ip_type").(string),
		SatFallbackIpAddress:         d.Get("sat_fallback_ip_address").(string),
		SatStaticTranslatedAddress:   d.Get("sat_static_translated_address").(string),
		SatStaticBiDirectional:       d.Get("sat_static_bi_directional").(bool),
		DatAddress:                   d.Get("dat_address").(string),
		DatPort:                      d.Get("dat_port").(int),
		Disabled:                     d.Get("disabled").(bool),
		Tag:                          asStringList(d, "tags"),
	}

	return vsys, rb, o
}

func saveDataNatPolicy(d *schema.ResourceData, vsys, rb string, o nat.Entry) {
	var err error
	d.SetId(buildNatPolicyId(vsys, rb, o.Name))
	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("rulebase", rb)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	if err = d.Set("source_zone", o.SourceZone); err != nil {
		log.Printf("[WARN] Error setting 'source_zone' param for %q: %s", d.Id(), err)
	}
	d.Set("destination_zone", o.DestinationZone)
	d.Set("to_interface", o.ToInterface)
	d.Set("service", o.Service)
	if err = d.Set("source_address", o.SourceAddress); err != nil {
		log.Printf("[WARN] Error setting 'source_address' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("destination_address", o.DestinationAddress); err != nil {
		log.Printf("[WARN] Error setting 'destination_address' param for %q: %s", d.Id(), err)
	}
	d.Set("sat_type", o.SatType)
	d.Set("sat_address_type", o.SatAddressType)
	if err = d.Set("sat_translated_address", o.SatTranslatedAddress); err != nil {
		log.Printf("[WARN] Error setting 'sat_translated_address' param for %q: %s", d.Id(), err)
	}
	d.Set("sat_interface", o.SatInterface)
	d.Set("sat_ip_address", o.SatIpAddress)
	d.Set("sat_fallback_type", o.SatFallbackType)
	if err = d.Set("sat_fallback_translated_address", o.SatFallbackTranslatedAddress); err != nil {
		log.Printf("[WARN] Error setting 'sat_fallback_translated_address' param for %q: %s", d.Id(), err)
	}
	d.Set("sat_fallback_interface", o.SatFallbackInterface)
	d.Set("sat_fallback_ip_type", o.SatFallbackIpType)
	d.Set("sat_fallback_ip_address", o.SatFallbackIpAddress)
	d.Set("sat_static_translated_address", o.SatStaticTranslatedAddress)
	d.Set("sat_static_bi_directional", o.SatStaticBiDirectional)
	d.Set("dat_address", o.DatAddress)
	d.Set("dat_port", o.DatPort)
	d.Set("disabled", o.Disabled)
	if err = d.Set("tags", o.Tag); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
}

func parseNatPolicyId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildNatPolicyId(a, b, c string) string {
	return fmt.Sprintf("%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c)
}

func createNatPolicy(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb, o := parseNatPolicy(d)
	o.Defaults()

	if err := fw.Policies.Nat.Set(vsys, rb, o); err != nil {
		return err
	}

	saveDataNatPolicy(d, vsys, rb, o)
	return nil
}

func readNatPolicy(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb, name := parseNatPolicyId(d.Id())

	o, err := fw.Policies.Nat.Get(vsys, rb, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	saveDataNatPolicy(d, vsys, rb, o)
	return nil
}

func updateNatPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error
	fw := meta.(*pango.Firewall)
	vsys, rb, o := parseNatPolicy(d)
	o.Defaults()

	lo, err := fw.Policies.Nat.Get(vsys, rb, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	err = fw.Policies.Nat.Edit(vsys, rb, lo)

	if err == nil {
		saveDataNatPolicy(d, vsys, rb, o)
	}
	return err
}

func deleteNatPolicy(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb, name := parseNatPolicyId(d.Id())

	err := fw.Policies.Nat.Delete(vsys, rb, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
