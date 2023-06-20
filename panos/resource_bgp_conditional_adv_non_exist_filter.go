package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/conadv/filter/nonexist"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpConditionalAdvNonExistFilter() *schema.Resource {
	return &schema.Resource{
		Create: createBgpConditionalAdvNonExistFilter,
		Read:   readBgpConditionalAdvNonExistFilter,
		Update: updateBgpConditionalAdvNonExistFilter,
		Delete: deleteBgpConditionalAdvNonExistFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpConditionalAdvNonExistFilterSchema(false),
	}
}

func bgpConditionalAdvNonExistFilterSchema(p bool) map[string]*schema.Schema {
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
			ValidateFunc: validateStringIn("", nonexist.RouteTableUnicast, nonexist.RouteTableMulticast, nonexist.RouteTableBoth),
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
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func saveBgpConditionalAdvNonExistFilter(d *schema.ResourceData, vr, ca string, o nonexist.Entry) {
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

func parseBgpConditionalAdvNonExistFilter(d *schema.ResourceData) (string, string, nonexist.Entry) {
	vr := d.Get("virtual_router").(string)
	ca := d.Get("bgp_conditional_adv").(string)

	o := nonexist.Entry{
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

func parseBgpConditionalAdvNonExistFilterId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildBgpConditionalAdvNonExistFilterId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ca, o := parseBgpConditionalAdvNonExistFilter(d)

	if err = fw.Network.BgpConAdvNonExistFilter.Set(vr, ca, o); err != nil {
		return err
	}

	d.SetId(buildBgpConditionalAdvNonExistFilterId(vr, ca, o.Name))
	return readBgpConditionalAdvNonExistFilter(d, meta)
}

func readBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ca, name := parseBgpConditionalAdvNonExistFilterId(d.Id())

	o, err := fw.Network.BgpConAdvNonExistFilter.Get(vr, ca, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpConditionalAdvNonExistFilter(d, vr, ca, o)

	return nil
}

func updateBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ca, o := parseBgpConditionalAdvNonExistFilter(d)

	lo, err := fw.Network.BgpConAdvNonExistFilter.Get(vr, ca, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpConAdvNonExistFilter.Edit(vr, ca, lo); err != nil {
		return err
	}

	return readBgpConditionalAdvNonExistFilter(d, meta)
}

func deleteBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, ca, name := parseBgpConditionalAdvNonExistFilterId(d.Id())

	if err := fw.Network.BgpConAdvNonExistFilter.Delete(vr, ca, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
