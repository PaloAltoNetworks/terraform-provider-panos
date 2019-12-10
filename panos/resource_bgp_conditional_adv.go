package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpConditionalAdv() *schema.Resource {
	return &schema.Resource{
		Create: createBgpConditionalAdv,
		Read:   readBgpConditionalAdv,
		Update: updateBgpConditionalAdv,
		Delete: deleteBgpConditionalAdv,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpConditionalAdvSchema(false),
	}
}

func bgpConditionalAdvSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"used_by": {
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

func parseBgpConditionalAdvId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildBgpConditionalAdvId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseBgpConditionalAdv(d *schema.ResourceData) (string, conadv.Entry) {
	vr := d.Get("virtual_router").(string)

	o := conadv.Entry{
		Name:   d.Get("name").(string),
		Enable: d.Get("enable").(bool),
		UsedBy: asStringList(d.Get("used_by").([]interface{})),
	}

	return vr, o
}

func saveBgpConditionalAdv(d *schema.ResourceData, vr string, o conadv.Entry) {
	d.Set("virtual_router", vr)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("used_by", o.UsedBy)
}

func createBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpConditionalAdv(d)

	if err = fw.Network.BgpConditionalAdv.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildBgpConditionalAdvId(vr, o.Name))
	return readBgpConditionalAdv(d, meta)
}

func readBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpConditionalAdvId(d.Id())

	o, err := fw.Network.BgpConditionalAdv.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpConditionalAdv(d, vr, o)

	return nil
}

func updateBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseBgpConditionalAdv(d)

	lo, err := fw.Network.BgpConditionalAdv.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpConditionalAdv.Edit(vr, lo); err != nil {
		return err
	}

	return readBgpConditionalAdv(d, meta)
}

func deleteBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpConditionalAdvId(d.Id())

	err := fw.Network.BgpConditionalAdv.Delete(vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
