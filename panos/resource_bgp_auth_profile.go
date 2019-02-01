package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/profile/auth"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBgpAuthProfile() *schema.Resource {
	return &schema.Resource{
		Create: createBgpAuthProfile,
		Read:   readBgpAuthProfile,
		Update: updateBgpAuthProfile,
		Delete: deleteBgpAuthProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAuthProfileSchema(false),
	}
}

func bgpAuthProfileSchema(p bool) map[string]*schema.Schema {
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
		"secret": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"secret_enc": &schema.Schema{
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}

	if p {
		ans["template"] = templateSchema()
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgpAuthProfile(d *schema.ResourceData) (string, auth.Entry) {
	vr := d.Get("virtual_router").(string)

	o := auth.Entry{
		Name:   d.Get("name").(string),
		Secret: d.Get("secret").(string),
	}

	return vr, o
}

func saveBgpAuthProfile(d *schema.ResourceData, vr string, o auth.Entry) {
	d.Set("virtual_router", vr)
	d.Set("name", o.Name)

	if d.Get("secret_enc").(string) != o.Secret {
		d.Set("secret", "(incorrect)")
	}
}

func parseBgpAuthProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildBgpAuthProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseBgpAuthProfile(d)

	if err := fw.Network.BgpAuthProfile.Set(vr, o); err != nil {
		return err
	}

	lo, err := fw.Network.BgpAuthProfile.Get(vr, o.Name)
	if err != nil {
		return err
	}

	d.SetId(buildBgpAuthProfileId(vr, o.Name))
	d.Set("secret_enc", lo.Secret)

	return readBgpAuthProfile(d, meta)
}

func readBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, name := parseBgpAuthProfileId(d.Id())

	o, err := fw.Network.BgpAuthProfile.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpAuthProfile(d, vr, o)

	return nil
}

func updateBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpAuthProfile(d)

	lo, err := fw.Network.BgpAuthProfile.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpAuthProfile.Edit(vr, lo); err != nil {
		return err
	}

	eo, err := fw.Network.BgpAuthProfile.Get(vr, o.Name)
	if err != nil {
		return err
	}

	d.Set("secret_enc", eo.Secret)
	return readBgpAuthProfile(d, meta)
}

func deleteBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpAuthProfileId(d.Id())

	err := fw.Network.BgpAuthProfile.Delete(vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
