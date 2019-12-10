package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer2"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLayer2Subinterface() *schema.Resource {
	return &schema.Resource{
		Create: createLayer2Subinterface,
		Read:   readLayer2Subinterface,
		Update: updateLayer2Subinterface,
		Delete: deleteLayer2Subinterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: layer2SubinterfaceSchema(false),
	}
}

func layer2SubinterfaceSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"interface_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      layer2.EthernetInterface,
			ValidateFunc: validateStringIn(layer2.EthernetInterface, layer2.AggregateInterface),
			ForceNew:     true,
		},
		"parent_interface": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"parent_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      layer2.Layer2,
			ValidateFunc: validateStringIn(layer2.Layer2, layer2.VirtualWire),
			ForceNew:     true,
		},
		"vsys": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "vsys1",
		},
		"tag": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"netflow_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema(false)
	}

	return ans
}

func parseLayer2SubinterfaceId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildLayer2SubinterfaceId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func loadLayer2Subinterface(d *schema.ResourceData) layer2.Entry {
	return layer2.Entry{
		Name:           d.Get("name").(string),
		Tag:            d.Get("tag").(int),
		NetflowProfile: d.Get("netflow_profile").(string),
		Comment:        d.Get("comment").(string),
	}
}

func saveLayer2Subinterface(d *schema.ResourceData, o layer2.Entry) {
	d.Set("name", o.Name)
	d.Set("tag", o.Tag)
	d.Set("netflow_profile", o.NetflowProfile)
	d.Set("comment", o.Comment)
}

func parseLayer2Subinterface(d *schema.ResourceData) (string, string, string, string, layer2.Entry) {
	iType := d.Get("interface_type").(string)
	eth := d.Get("parent_interface").(string)
	mType := d.Get("parent_mode").(string)
	vsys := d.Get("vsys").(string)
	o := loadLayer2Subinterface(d)

	return iType, eth, mType, vsys, o
}

func createLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	iType, eth, mType, vsys, o := parseLayer2Subinterface(d)

	if err := fw.Network.Layer2Subinterface.Set(iType, eth, mType, vsys, o); err != nil {
		return err
	}

	d.SetId(buildLayer2SubinterfaceId(iType, eth, mType, vsys, o.Name))
	return readLayer2Subinterface(d, meta)
}

func readLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	iType, eth, mType, vsys, name := parseLayer2SubinterfaceId(d.Id())

	o, err := fw.Network.Layer2Subinterface.Get(iType, eth, mType, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := fw.IsImported(util.InterfaceImport, "", "", vsys, name)
	if err != nil {
		return err
	}

	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	d.Set("interface_type", iType)
	d.Set("parent_interface", eth)
	d.Set("parent_mode", mType)
	saveLayer2Subinterface(d, o)

	return nil
}

func updateLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	iType, eth, mType, vsys, o := parseLayer2Subinterface(d)

	lo, err := fw.Network.Layer2Subinterface.Get(iType, eth, mType, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.Layer2Subinterface.Edit(iType, eth, mType, vsys, lo); err != nil {
		return err
	}

	d.SetId(buildLayer2SubinterfaceId(iType, eth, mType, vsys, o.Name))
	return readLayer2Subinterface(d, meta)
}

func deleteLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	iType, eth, mType, _, name := parseLayer2SubinterfaceId(d.Id())

	err := fw.Network.Layer2Subinterface.Delete(iType, eth, mType, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
