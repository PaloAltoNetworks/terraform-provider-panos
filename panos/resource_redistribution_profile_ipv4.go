package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/profile/redist/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRedistributionProfileIpv4() *schema.Resource {
	return &schema.Resource{
		Create: createRedistributionProfileIpv4,
		Read:   readRedistributionProfileIpv4,
		Update: updateRedistributionProfileIpv4,
		Delete: deleteRedistributionProfileIpv4,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: redistributionProfileIpv4Schema(false),
	}
}

func redistributionProfileIpv4Schema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"virtual_router": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"priority": &schema.Schema{
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validateIntInRange(1, 255),
		},
		"action": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      ipv4.ActionRedist,
			ValidateFunc: validateStringIn(ipv4.ActionRedist, ipv4.ActionNoRedist),
		},
		"types": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateStringIn(ipv4.TypeBgp, ipv4.TypeConnect, ipv4.TypeOspf, ipv4.TypeRip, ipv4.TypeStatic),
			},
		},
		"interfaces": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"destinations": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"next_hops": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ospf_path_types": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateStringIn(ipv4.OspfPathTypeIntraArea, ipv4.OspfPathTypeInterArea, ipv4.OspfPathTypeExt1, ipv4.OspfPathTypeExt2),
			},
		},
		"ospf_areas": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ospf_tags": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"bgp_communities": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"bgp_extended_communities": &schema.Schema{
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

func parseRedistributionProfileIpv4(d *schema.ResourceData) (string, ipv4.Entry) {
	vr := d.Get("virtual_router").(string)

	o := ipv4.Entry{
		Name:                   d.Get("name").(string),
		Priority:               d.Get("priority").(int),
		Action:                 d.Get("action").(string),
		Types:                  asStringList(d.Get("types").([]interface{})),
		Interfaces:             asStringList(d.Get("interfaces").([]interface{})),
		Destinations:           asStringList(d.Get("destinations").([]interface{})),
		NextHops:               asStringList(d.Get("next_hops").([]interface{})),
		OspfPathTypes:          asStringList(d.Get("ospf_path_types").([]interface{})),
		OspfAreas:              asStringList(d.Get("ospf_areas").([]interface{})),
		OspfTags:               asStringList(d.Get("ospf_tags").([]interface{})),
		BgpCommunities:         asStringList(d.Get("bgp_communities").([]interface{})),
		BgpExtendedCommunities: asStringList(d.Get("bgp_extended_communities").([]interface{})),
	}

	return vr, o
}

func saveRedistributionProfileIpv4(d *schema.ResourceData, vr string, o ipv4.Entry) {
	var err error

	d.Set("virtual_router", vr)

	d.Set("name", o.Name)
	d.Set("priority", o.Priority)
	d.Set("action", o.Action)
	if err = d.Set("types", o.Types); err != nil {
		log.Printf("[WARN] Error setting 'types' for %q: %s", d.Id(), err)
	}
	if err = d.Set("interfaces", o.Interfaces); err != nil {
		log.Printf("[WARN] Error setting 'interfaces' for %q: %s", d.Id(), err)
	}
	if err = d.Set("destinations", o.Destinations); err != nil {
		log.Printf("[WARN] Error setting 'destinations' for %q: %s", d.Id(), err)
	}
	if err = d.Set("next_hops", o.NextHops); err != nil {
		log.Printf("[WARN] Error setting 'next_hops' for %q: %s", d.Id(), err)
	}
	if err = d.Set("ospf_path_types", o.OspfPathTypes); err != nil {
		log.Printf("[WARN] Error setting 'ospf_path_types' for %q: %s", d.Id(), err)
	}
	if err = d.Set("ospf_areas", o.OspfAreas); err != nil {
		log.Printf("[WARN] Error setting 'ospf_areas' for %q: %s", d.Id(), err)
	}
	if err = d.Set("ospf_tags", o.OspfTags); err != nil {
		log.Printf("[WARN] Error setting 'ospf_tags' for %q: %s", d.Id(), err)
	}
	if err = d.Set("bgp_communities", o.BgpCommunities); err != nil {
		log.Printf("[WARN] Error setting 'bgp_communities' for %q: %s", d.Id(), err)
	}
	if err = d.Set("bgp_extended_communities", o.BgpExtendedCommunities); err != nil {
		log.Printf("[WARN] Error setting 'bgp_extended_communities' for %q: %s", d.Id(), err)
	}
}

func parseRedistributionProfileIpv4Id(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildRedistributionProfileIpv4Id(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseRedistributionProfileIpv4(d)

	if err := fw.Network.RedistributionProfile.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildRedistributionProfileIpv4Id(vr, o.Name))
	return readRedistributionProfileIpv4(d, meta)
}

func readRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, name := parseRedistributionProfileIpv4Id(d.Id())

	o, err := fw.Network.RedistributionProfile.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveRedistributionProfileIpv4(d, vr, o)

	return nil
}

func updateRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseRedistributionProfileIpv4(d)

	lo, err := fw.Network.RedistributionProfile.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.RedistributionProfile.Edit(vr, lo); err != nil {
		return err
	}

	return readRedistributionProfileIpv4(d, meta)
}

func deleteRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseRedistributionProfileIpv4Id(d.Id())

	err := fw.Network.RedistributionProfile.Delete(vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
