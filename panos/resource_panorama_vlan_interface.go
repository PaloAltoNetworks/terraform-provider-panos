package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/vlan"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaVlanInterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaVlanInterface,
		Read:   readPanoramaVlanInterface,
		Update: updatePanoramaVlanInterface,
		Delete: deletePanoramaVlanInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringHasPrefix("vlan."),
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
			"enable_dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_dhcp_default_route": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dhcp_default_route_metric": {
				Type:     schema.TypeInt,
				Optional: true,
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

func parsePanoramaVlanInterfaceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaVlanInterfaceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaVlanInterface(d *schema.ResourceData) (string, string, string, vlan.Entry) {
	tmpl := d.Get("template").(string)
	vsys := d.Get("vsys").(string)

	o := vlan.Entry{
		Name:                   d.Get("name").(string),
		Comment:                d.Get("comment").(string),
		NetflowProfile:         d.Get("netflow_profile").(string),
		StaticIps:              asStringList(d.Get("static_ips").([]interface{})),
		EnableDhcp:             d.Get("enable_dhcp").(bool),
		CreateDhcpDefaultRoute: d.Get("create_dhcp_default_route").(bool),
		DhcpDefaultRouteMetric: d.Get("dhcp_default_route_metric").(int),
		ManagementProfile:      d.Get("management_profile").(string),
		Mtu:                    d.Get("mtu").(int),
		AdjustTcpMss:           d.Get("adjust_tcp_mss").(bool),
		Ipv4MssAdjust:          d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:          d.Get("ipv6_mss_adjust").(int),
	}

	return tmpl, "", vsys, o
}

func createPanoramaVlanInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaVlanInterface(d)

	if err := pano.Network.VlanInterface.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaVlanInterfaceId(tmpl, ts, vsys, o.Name))
	return readPanoramaVlanInterface(d, meta)
}

func readPanoramaVlanInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaVlanInterfaceId(d.Id())

	o, err := pano.Network.VlanInterface.Get(tmpl, ts, name)
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
	d.Set("enable_dhcp", o.EnableDhcp)
	d.Set("create_dhcp_default_route", o.CreateDhcpDefaultRoute)
	d.Set("dhcp_default_route_metric", o.DhcpDefaultRouteMetric)
	d.Set("management_profile", o.ManagementProfile)
	d.Set("mtu", o.Mtu)
	d.Set("adjust_tcp_mss", o.AdjustTcpMss)
	d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
	d.Set("ipv6_mss_adjust", o.Ipv6MssAdjust)

	return nil
}

func updatePanoramaVlanInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaVlanInterface(d)

	lo, err := pano.Network.VlanInterface.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.VlanInterface.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaVlanInterface(d, meta)
}

func deletePanoramaVlanInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaVlanInterfaceId(d.Id())

	err := pano.Network.VlanInterface.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
