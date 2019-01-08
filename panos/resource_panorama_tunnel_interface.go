package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/tunnel"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaTunnelInterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaTunnelInterface,
		Read:   readPanoramaTunnelInterface,
		Update: updatePanoramaTunnelInterface,
		Delete: deletePanoramaTunnelInterface,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringHasPrefix("tunnel."),
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"netflow_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"static_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of static IP addresses",
			},
			"management_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parsePanoramaTunnelInterfaceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaTunnelInterfaceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaTunnelInterface(d *schema.ResourceData) (string, string, string, tunnel.Entry) {
	tmpl := d.Get("template").(string)
	vsys := d.Get("vsys").(string)

	o := tunnel.Entry{
		Name:              d.Get("name").(string),
		Comment:           d.Get("comment").(string),
		NetflowProfile:    d.Get("netflow_profile").(string),
		StaticIps:         asStringList(d.Get("static_ips").([]interface{})),
		ManagementProfile: d.Get("management_profile").(string),
		Mtu:               d.Get("mtu").(int),
	}

	return tmpl, "", vsys, o
}

func createPanoramaTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaTunnelInterface(d)

	if err := pano.Network.TunnelInterface.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaTunnelInterfaceId(tmpl, ts, vsys, o.Name))
	return readPanoramaTunnelInterface(d, meta)
}

func readPanoramaTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaTunnelInterfaceId(d.Id())

	o, err := pano.Network.TunnelInterface.Get(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := pano.IsImported(util.InterfaceImport, tmpl, ts, vsys, name)
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("name", o.Name)
	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	d.Set("comment", o.Comment)
	d.Set("netflow_profile", o.NetflowProfile)
	if err = d.Set("static_ips", o.StaticIps); err != nil {
		log.Printf("[WARN] Error setting 'static_ips' for %q: %s", d.Id(), err)
	}
	d.Set("management_profile", o.ManagementProfile)
	d.Set("mtu", o.Mtu)

	return nil
}

func updatePanoramaTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaTunnelInterface(d)

	lo, err := pano.Network.TunnelInterface.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.TunnelInterface.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaTunnelInterface(d, meta)
}

func deletePanoramaTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaTunnelInterfaceId(d.Id())

	err := pano.Network.TunnelInterface.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
