package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv/filter/advertise"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBgpConditionalAdvAdvertiseFilter() *schema.Resource {
	return &schema.Resource{
		Create: createBgpConditionalAdvAdvertiseFilter,
		Read:   readBgpConditionalAdvAdvertiseFilter,
		Update: updateBgpConditionalAdvAdvertiseFilter,
		Delete: deleteBgpConditionalAdvAdvertiseFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpConditionalAdvAdvertiseFilterSchema(false),
	}
}

func bgpConditionalAdvAdvertiseFilterSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"bgp_conditional_adv": {
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
		"as_path_regex": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"community_regex": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"extended_community_regex": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"med": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"route_table": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", advertise.RouteTableUnicast, advertise.RouteTableMulticast, advertise.RouteTableBoth),
		},
		"address_prefixes": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"next_hops": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"from_peers": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if p {
		ans["template"] = templateSchema()
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData) (string, string, advertise.Entry) {
	vr := d.Get("virtual_router").(string)
	ca := d.Get("bgp_conditional_adv").(string)

	o := advertise.Entry{
		Name:                   d.Get("name").(string),
		Enable:                 d.Get("enable").(bool),
		AsPathRegex:            d.Get("as_path_regex").(string),
		CommunityRegex:         d.Get("community_regex").(string),
		ExtendedCommunityRegex: d.Get("extended_community_regex").(string),
		Med:                    d.Get("med").(string),
		RouteTable:             d.Get("route_table").(string),
		AddressPrefix:          asStringList(d.Get("address_prefixes").([]interface{})),
		NextHop:                asStringList(d.Get("next_hops").([]interface{})),
		FromPeer:               asStringList(d.Get("from_peers").([]interface{})),
	}

	return vr, ca, o
}

func saveBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, vr, ca string, o advertise.Entry) {
	d.Set("virtual_router", vr)
	d.Set("bgp_conditional_adv", ca)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("as_path_regex", o.AsPathRegex)
	d.Set("community_regex", o.CommunityRegex)
	d.Set("extended_community_regex", o.ExtendedCommunityRegex)
	d.Set("med", o.Med)
	d.Set("route_table", o.RouteTable)
	d.Set("address_prefixes", o.AddressPrefix)
	d.Set("next_hops", o.NextHop)
	d.Set("from_peers", o.FromPeer)
}

func parseBgpConditionalAdvAdvertiseFilterId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildBgpConditionalAdvAdvertiseFilterId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ca, o := parseBgpConditionalAdvAdvertiseFilter(d)

	if err = fw.Network.BgpConAdvAdvertiseFilter.Set(vr, ca, o); err != nil {
		return err
	}

	d.SetId(buildBgpConditionalAdvAdvertiseFilterId(vr, ca, o.Name))
	return readBgpConditionalAdvAdvertiseFilter(d, meta)
}

func readBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ca, name := parseBgpConditionalAdvAdvertiseFilterId(d.Id())

	o, err := fw.Network.BgpConAdvAdvertiseFilter.Get(vr, ca, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpConditionalAdvAdvertiseFilter(d, vr, ca, o)

	return nil
}

func updateBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ca, o := parseBgpConditionalAdvAdvertiseFilter(d)

	lo, err := fw.Network.BgpConAdvAdvertiseFilter.Get(vr, ca, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpConAdvAdvertiseFilter.Edit(vr, ca, lo); err != nil {
		return err
	}

	return readBgpConditionalAdvAdvertiseFilter(d, meta)
}

func deleteBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, ca, name := parseBgpConditionalAdvAdvertiseFilterId(d.Id())

	if err := fw.Network.BgpConAdvAdvertiseFilter.Delete(vr, ca, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
