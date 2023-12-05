package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/redist"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBgpRedistRule() *schema.Resource {
	return &schema.Resource{
		Create: createBgpRedistRule,
		Read:   readBgpRedistRule,
		Update: updateBgpRedistRule,
		Delete: deleteBgpRedistRule,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpRedistRuleSchema(false),
	}
}

func bgpRedistRuleSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"address_family": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      redist.AddressFamilyIpv4,
			ValidateFunc: validateStringIn("", redist.AddressFamilyIpv4, redist.AddressFamilyIpv6),
		},
		"route_table": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", redist.RouteTableUnicast, redist.RouteTableMulticast, redist.RouteTableBoth),
		},
		"metric": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"set_origin": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      redist.SetOriginIncomplete,
			ValidateFunc: validateStringIn("", redist.SetOriginIgp, redist.SetOriginEgp, redist.SetOriginIncomplete),
		},
		"set_med": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"set_local_preference": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"set_as_path_limit": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"set_communities": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"set_extended_communities": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if p {
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgpRedistRule(d *schema.ResourceData) (string, redist.Entry) {
	vr := d.Get("virtual_router").(string)

	o := redist.Entry{
		Name:                 d.Get("name").(string),
		Enable:               d.Get("enable").(bool),
		AddressFamily:        d.Get("address_family").(string),
		RouteTable:           d.Get("route_table").(string),
		Metric:               d.Get("metric").(int),
		SetOrigin:            d.Get("set_origin").(string),
		SetMed:               d.Get("set_med").(string),
		SetLocalPreference:   d.Get("set_local_preference").(string),
		SetAsPathLimit:       d.Get("set_as_path_limit").(int),
		SetCommunity:         asStringList(d.Get("set_communities").([]interface{})),
		SetExtendedCommunity: asStringList(d.Get("set_extended_communities").([]interface{})),
	}

	return vr, o
}

func saveBgpRedistRule(d *schema.ResourceData, vr string, o redist.Entry) {
	var err error

	d.Set("virtual_router", vr)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("address_family", o.AddressFamily)
	d.Set("route_table", o.RouteTable)
	d.Set("metric", o.Metric)
	d.Set("set_origin", o.SetOrigin)
	d.Set("set_med", o.SetMed)
	d.Set("set_local_preference", o.SetLocalPreference)
	d.Set("set_as_path_limit", o.SetAsPathLimit)

	if err = d.Set("set_communities", o.SetCommunity); err != nil {
		log.Printf("[WARN] Error setting 'set_communities' for %q: %s", d.Id(), err)
	}

	if err = d.Set("set_extended_communities", o.SetExtendedCommunity); err != nil {
		log.Printf("[WARN] Error setting 'set_extended_communities' for %q: %s", d.Id(), err)
	}
}

func parseBgpRedistRuleId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildBgpRedistRuleId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpRedistRule(d)

	if err = fw.Network.BgpRedistRule.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildBgpRedistRuleId(vr, o.Name))
	return readBgpRedistRule(d, meta)
}

func readBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpRedistRuleId(d.Id())

	o, err := fw.Network.BgpRedistRule.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpRedistRule(d, vr, o)

	return nil
}

func updateBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpRedistRule(d)

	lo, err := fw.Network.BgpRedistRule.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpRedistRule.Edit(vr, lo); err != nil {
		return err
	}

	return readBgpRedistRule(d, meta)
}

func deleteBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpRedistRuleId(d.Id())

	if err := fw.Network.BgpRedistRule.Delete(vr, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
