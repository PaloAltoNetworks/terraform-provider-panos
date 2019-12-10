package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate/filter/suppress"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpAggregateSuppressFilter() *schema.Resource {
	return &schema.Resource{
		Create: createBgpAggregateSuppressFilter,
		Read:   readBgpAggregateSuppressFilter,
		Update: updateBgpAggregateSuppressFilter,
		Delete: deleteBgpAggregateSuppressFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAggregateSuppressFilterSchema(false),
	}
}

func bgpAggregateSuppressFilterSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"bgp_aggregate": {
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
			ValidateFunc: validateStringIn("", suppress.RouteTableUnicast, suppress.RouteTableMulticast, suppress.RouteTableBoth),
		},
		"address_prefix": {
			Type:     schema.TypeSet,
			Optional: true,
			Set:      resourceMatchAddressPrefixHash,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"prefix": {
						Type:     schema.TypeString,
						Required: true,
					},
					"exact": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
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

func parseBgpAggregateSuppressFilter(d *schema.ResourceData) (string, string, suppress.Entry) {
	vr := d.Get("virtual_router").(string)
	ag := d.Get("bgp_aggregate").(string)

	o := suppress.Entry{
		Name:                   d.Get("name").(string),
		Enable:                 d.Get("enable").(bool),
		AsPathRegex:            d.Get("as_path_regex").(string),
		CommunityRegex:         d.Get("community_regex").(string),
		ExtendedCommunityRegex: d.Get("extended_community_regex").(string),
		Med:                    d.Get("med").(string),
		RouteTable:             d.Get("route_table").(string),
		NextHop:                asStringList(d.Get("next_hops").([]interface{})),
		FromPeer:               asStringList(d.Get("from_peers").([]interface{})),
	}

	sl := d.Get("address_prefix").(*schema.Set).List()
	if len(sl) > 0 {
		o.AddressPrefix = make(map[string]bool)
		for i := range sl {
			sli := sl[i].(map[string]interface{})
			o.AddressPrefix[sli["prefix"].(string)] = sli["exact"].(bool)
		}
	}

	return vr, ag, o
}

func saveBgpAggregateSuppressFilter(d *schema.ResourceData, vr, ag string, o suppress.Entry) {
	d.Set("virtual_router", vr)
	d.Set("bgp_aggregate", ag)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("as_path_regex", o.AsPathRegex)
	d.Set("community_regex", o.CommunityRegex)
	d.Set("extended_community_regex", o.ExtendedCommunityRegex)
	d.Set("med", o.Med)
	d.Set("route_table", o.RouteTable)
	d.Set("next_hops", o.NextHop)
	d.Set("from_peers", o.FromPeer)

	aps := &schema.Set{
		F: resourceMatchAddressPrefixHash,
	}
	for k, v := range o.AddressPrefix {
		aps.Add(map[string]interface{}{
			"prefix": k,
			"exact":  v,
		})
	}
	if err := d.Set("address_prefix", aps); err != nil {
		log.Printf("[WARN] Error setting `address_prefix` for %q: %s", d.Id(), err)
	}
}

func parseBgpAggregateSuppressFilterId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildBgpAggregateSuppressFilterId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ag, o := parseBgpAggregateSuppressFilter(d)

	if err = fw.Network.BgpAggSuppressFilter.Set(vr, ag, o); err != nil {
		return err
	}

	d.SetId(buildBgpAggregateSuppressFilterId(vr, ag, o.Name))
	return readBgpAggregateSuppressFilter(d, meta)
}

func readBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, ag, name := parseBgpAggregateSuppressFilterId(d.Id())

	o, err := fw.Network.BgpAggSuppressFilter.Get(vr, ag, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpAggregateSuppressFilter(d, vr, ag, o)

	return nil
}

func updateBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, ag, o := parseBgpAggregateSuppressFilter(d)

	lo, err := fw.Network.BgpAggSuppressFilter.Get(vr, ag, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpAggSuppressFilter.Edit(vr, ag, lo); err != nil {
		return err
	}

	return readBgpAggregateSuppressFilter(d, meta)
}

func deleteBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, ag, name := parseBgpAggregateSuppressFilterId(d.Id())

	err := fw.Network.BgpAggSuppressFilter.Delete(vr, ag, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
