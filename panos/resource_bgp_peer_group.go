package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/peer/group"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpPeerGroup() *schema.Resource {
	return &schema.Resource{
		Create: createBgpPeerGroup,
		Read:   readBgpPeerGroup,
		Update: updateBgpPeerGroup,
		Delete: deleteBgpPeerGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpPeerGroupSchema(false),
	}
}

func bgpPeerGroupSchema(p bool) map[string]*schema.Schema {
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
		"aggregated_confed_as_path": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"soft_reset_with_stored_info": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"type": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      group.TypeEbgp,
			ValidateFunc: validateStringIn(group.TypeEbgp, group.TypeEbgpConfed, group.TypeIbgp, group.TypeIbgpConfed),
		},
		"export_next_hop": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn(group.NextHopOriginal, group.NextHopUseSelf, group.NextHopResolve),
		},
		"import_next_hop": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", group.NextHopOriginal, group.NextHopUsePeer),
		},
		"remove_private_as": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgpPeerGroup(d *schema.ResourceData) (string, group.Entry) {
	vr := d.Get("virtual_router").(string)

	o := group.Entry{
		Name:                    d.Get("name").(string),
		Enable:                  d.Get("enable").(bool),
		AggregatedConfedAsPath:  d.Get("aggregated_confed_as_path").(bool),
		SoftResetWithStoredInfo: d.Get("soft_reset_with_stored_info").(bool),
		Type:                    d.Get("type").(string),
		ExportNextHop:           d.Get("export_next_hop").(string),
		ImportNextHop:           d.Get("import_next_hop").(string),
		RemovePrivateAs:         d.Get("remove_private_as").(bool),
	}

	return vr, o
}

func saveBgpPeerGroup(d *schema.ResourceData, vr string, o group.Entry) {
	d.Set("virtual_router", vr)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("aggregated_confed_as_path", o.AggregatedConfedAsPath)
	d.Set("soft_reset_with_stored_info", o.SoftResetWithStoredInfo)
	d.Set("type", o.Type)
	d.Set("export_next_hop", o.ExportNextHop)
	d.Set("import_next_hop", o.ImportNextHop)
	d.Set("remove_private_as", o.RemovePrivateAs)
}

func parseBgpPeerGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildBgpPeerGroupId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseBgpPeerGroup(d)

	if err := fw.Network.BgpPeerGroup.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildBgpPeerGroupId(vr, o.Name))
	return readBgpPeerGroup(d, meta)
}

func readBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, name := parseBgpPeerGroupId(d.Id())

	o, err := fw.Network.BgpPeerGroup.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpPeerGroup(d, vr, o)

	return nil
}

func updateBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgpPeerGroup(d)

	lo, err := fw.Network.BgpPeerGroup.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpPeerGroup.Edit(vr, lo); err != nil {
		return err
	}

	return readBgpPeerGroup(d, meta)
}

func deleteBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseBgpPeerGroupId(d.Id())

	err := fw.Network.BgpPeerGroup.Delete(vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
