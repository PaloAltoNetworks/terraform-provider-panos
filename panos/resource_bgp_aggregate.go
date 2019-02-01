package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBgpAggregate() *schema.Resource {
	return &schema.Resource{
		Create: createBgpAggregate,
		Read:   readBgpAggregate,
		Update: updateBgpAggregate,
		Delete: deleteBgpAggregate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAggregateSchema(false),
	}
}

func bgpAggregateSchema(p bool) map[string]*schema.Schema {
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
		"prefix": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"summary": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"as_set": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"local_preference": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"med": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"weight": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"next_hop": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"origin": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      aggregate.OriginIncomplete,
			ValidateFunc: validateStringIn("", aggregate.OriginIgp, aggregate.OriginEgp, aggregate.OriginIncomplete),
		},
		"as_path_limit": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"as_path_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      aggregate.AsPathTypeNone,
			ValidateFunc: validateStringIn("", aggregate.AsPathTypeNone, aggregate.AsPathTypeRemove, aggregate.AsPathTypePrepend, aggregate.AsPathTypeRemoveAndPrepend),
		},
		"as_path_value": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"community_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      aggregate.CommunityTypeNone,
			ValidateFunc: validateStringIn("", aggregate.CommunityTypeNone, aggregate.CommunityTypeRemoveAll, aggregate.CommunityTypeRemoveRegex, aggregate.CommunityTypeAppend, aggregate.CommunityTypeOverwrite),
		},
		"community_value": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"extended_community_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      aggregate.CommunityTypeNone,
			ValidateFunc: validateStringIn("", aggregate.CommunityTypeNone, aggregate.CommunityTypeRemoveAll, aggregate.CommunityTypeRemoveRegex, aggregate.CommunityTypeAppend, aggregate.CommunityTypeOverwrite),
		},
		"extended_community_value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema()
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgpAggregate(d *schema.ResourceData) (string, aggregate.Entry) {
	vr := d.Get("virtual_router").(string)

	o := aggregate.Entry{
		Name:                   d.Get("name").(string),
		Prefix:                 d.Get("prefix").(string),
		Enable:                 d.Get("enable").(bool),
		Summary:                d.Get("summary").(bool),
		AsSet:                  d.Get("as_set").(bool),
		LocalPreference:        d.Get("local_preference").(string),
		Med:                    d.Get("med").(string),
		Weight:                 d.Get("weight").(int),
		NextHop:                d.Get("next_hop").(string),
		Origin:                 d.Get("origin").(string),
		AsPathLimit:            d.Get("as_path_limit").(int),
		AsPathType:             d.Get("as_path_type").(string),
		AsPathValue:            d.Get("as_path_value").(string),
		CommunityType:          d.Get("community_type").(string),
		CommunityValue:         d.Get("community_value").(string),
		ExtendedCommunityType:  d.Get("extended_community_type").(string),
		ExtendedCommunityValue: d.Get("extended_community_value").(string),
	}

	return vr, o
}

func saveBgpAggregate(d *schema.ResourceData, vr string, o aggregate.Entry) {
	d.Set("virtual_router", vr)

	d.Set("name", o.Name)
	d.Set("prefix", o.Prefix)
	d.Set("enable", o.Enable)
	d.Set("summary", o.Summary)
	d.Set("as_set", o.AsSet)
	d.Set("local_preference", o.LocalPreference)
	d.Set("med", o.Med)
	d.Set("weight", o.Weight)
	d.Set("next_hop", o.NextHop)
	d.Set("origin", o.Origin)
	d.Set("as_path_limit", o.AsPathLimit)
	d.Set("as_path_type", o.AsPathType)
	d.Set("as_path_value", o.AsPathValue)
	d.Set("community_type", o.CommunityType)
	d.Set("community_value", o.CommunityValue)
	d.Set("extended_community_type", o.ExtendedCommunityType)
	d.Set("extended_community_value", o.ExtendedCommunityValue)
}

func parseBgpAggregateId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildBgpAggregateId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpAggregate(d)

	if err = fw.Network.BgpAggregate.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildBgpAggregateId(vr, o.Name))
	return readBgpAggregate(d, meta)
}

func readBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpAggregateId(d.Id())

	o, err := fw.Network.BgpAggregate.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpAggregate(d, vr, o)

	return nil
}

func updateBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpAggregate(d)

	lo, err := fw.Network.BgpAggregate.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpAggregate.Edit(vr, lo); err != nil {
		return err
	}

	return readBgpAggregate(d, meta)
}

func deleteBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpAggregateId(d.Id())

	if err := fw.Network.BgpAggregate.Delete(vr, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
