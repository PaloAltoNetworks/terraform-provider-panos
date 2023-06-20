package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/profile/dampening"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpDampeningProfile() *schema.Resource {
	return &schema.Resource{
		Create: createBgpDampeningProfile,
		Read:   readBgpDampeningProfile,
		Update: updateBgpDampeningProfile,
		Delete: deleteBgpDampeningProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpDampeningProfileSchema(false),
	}
}

func bgpDampeningProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"enable": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"cutoff": &schema.Schema{
			Type:     schema.TypeFloat,
			Optional: true,
			Default:  1.25,
		},
		"reuse": &schema.Schema{
			Type:     schema.TypeFloat,
			Optional: true,
			Default:  0.5,
		},
		"max_hold_time": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  900,
		},
		"decay_half_life_reachable": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  300,
		},
		"decay_half_life_unreachable": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  900,
		},
	}

	if p {
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func saveBgpDampeningProfile(d *schema.ResourceData, vr string, o dampening.Entry) {
	d.Set("virtual_router", vr)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("cutoff", o.Cutoff)
	d.Set("reuse", o.Reuse)
	d.Set("max_hold_time", o.MaxHoldTime)
	d.Set("decay_half_life_reachable", o.DecayHalfLifeReachable)
	d.Set("decay_half_life_unreachable", o.DecayHalfLifeUnreachable)
}

func parseBgpDampeningProfile(d *schema.ResourceData) (string, dampening.Entry) {
	vr := d.Get("virtual_router").(string)
	o := dampening.Entry{
		Name:                     d.Get("name").(string),
		Enable:                   d.Get("enable").(bool),
		Cutoff:                   d.Get("cutoff").(float64),
		Reuse:                    d.Get("reuse").(float64),
		MaxHoldTime:              d.Get("max_hold_time").(int),
		DecayHalfLifeReachable:   d.Get("decay_half_life_reachable").(int),
		DecayHalfLifeUnreachable: d.Get("decay_half_life_unreachable").(int),
	}

	return vr, o
}

func parseBgpDampeningProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildBgpDampeningProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseBgpDampeningProfile(d)

	if err := fw.Network.BgpDampeningProfile.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildBgpDampeningProfileId(vr, o.Name))
	return readBgpDampeningProfile(d, meta)
}

func readBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, name := parseBgpDampeningProfileId(d.Id())

	o, err := fw.Network.BgpDampeningProfile.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpDampeningProfile(d, vr, o)

	return nil
}

func updateBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpDampeningProfile(d)

	lo, err := fw.Network.BgpDampeningProfile.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpDampeningProfile.Edit(vr, lo); err != nil {
		return err
	}

	return readBgpDampeningProfile(d, meta)
}

func deleteBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpDampeningProfileId(d.Id())

	err := fw.Network.BgpDampeningProfile.Delete(vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
