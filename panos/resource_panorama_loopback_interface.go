package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/loopback"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaLoopbackInterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaLoopbackInterface,
		Read:   readPanoramaLoopbackInterface,
		Update: updatePanoramaLoopbackInterface,
		Delete: deletePanoramaLoopbackInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringHasPrefix("loopback."),
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
			"adjust_tcp_mss": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv4_mss_adjust": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv6_mss_adjust": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parsePanoramaLoopbackInterfaceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaLoopbackInterfaceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaLoopbackInterface(d *schema.ResourceData) (string, string, string, loopback.Entry) {
	tmpl := d.Get("template").(string)
	vsys := d.Get("vsys").(string)

	o := loopback.Entry{
		Name:              d.Get("name").(string),
		Comment:           d.Get("comment").(string),
		NetflowProfile:    d.Get("netflow_profile").(string),
		StaticIps:         asStringList(d.Get("static_ips").([]interface{})),
		ManagementProfile: d.Get("management_profile").(string),
		Mtu:               d.Get("mtu").(int),
		AdjustTcpMss:      d.Get("adjust_tcp_mss").(bool),
		Ipv4MssAdjust:     d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:     d.Get("ipv6_mss_adjust").(int),
	}

	return tmpl, "", vsys, o
}

func createPanoramaLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaLoopbackInterface(d)

	if err := pano.Network.LoopbackInterface.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaLoopbackInterfaceId(tmpl, ts, vsys, o.Name))
	return readPanoramaLoopbackInterface(d, meta)
}

func readPanoramaLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaLoopbackInterfaceId(d.Id())

	o, err := pano.Network.LoopbackInterface.Get(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := pano.IsImported(util.InterfaceImport, tmpl, ts, vsys, name)
	if err != nil {
		return err
	}

	d.Set("name", o.Name)
	d.Set("template", tmpl)
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
	d.Set("adjust_tcp_mss", o.AdjustTcpMss)
	d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
	d.Set("ipv6_mss_adjust", o.Ipv6MssAdjust)

	return nil
}

func updatePanoramaLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaLoopbackInterface(d)

	lo, err := pano.Network.LoopbackInterface.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.LoopbackInterface.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaLoopbackInterface(d, meta)
}

func deletePanoramaLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaLoopbackInterfaceId(d.Id())

	err := pano.Network.LoopbackInterface.Delete(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
