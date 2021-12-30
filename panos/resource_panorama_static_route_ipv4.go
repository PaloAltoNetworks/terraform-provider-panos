package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/route/static/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaStaticRouteIpv4() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaStaticRouteIpv4,
		Read:   readPanoramaStaticRouteIpv4,
		Update: updatePanoramaStaticRouteIpv4,
		Delete: deletePanoramaStaticRouteIpv4,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template_stack"},
			},
			"template_stack": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template"},
			},
			"virtual_router": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipv4.NextHopIpAddress,
				ValidateFunc: validateStringIn(ipv4.NextHopDiscard, ipv4.NextHopIpAddress, ipv4.NextHopNextVr, ""),
			},
			"next_hop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_distance": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"metric": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"route_table": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipv4.RouteTableUnicast,
				ValidateFunc: validateStringIn(ipv4.RouteTableNoInstall, ipv4.RouteTableUnicast, ipv4.RouteTableMulticast, ipv4.RouteTableBoth),
			},
			"bfd_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func parsePanoramaStaticRouteIpv4(d *schema.ResourceData) (string, string, string, ipv4.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vr := d.Get("virtual_router").(string)

	o := ipv4.Entry{
		Name:          d.Get("name").(string),
		Destination:   d.Get("destination").(string),
		Interface:     d.Get("interface").(string),
		Type:          d.Get("type").(string),
		NextHop:       d.Get("next_hop").(string),
		AdminDistance: d.Get("admin_distance").(int),
		Metric:        d.Get("metric").(int),
		RouteTable:    d.Get("route_table").(string),
		BfdProfile:    d.Get("bfd_profile").(string),
	}

	return tmpl, ts, vr, o
}

func parsePanoramaStaticRouteIpv4Id(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaStaticRouteIpv4Id(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaStaticRouteIpv4(d)

	if err := pano.Network.StaticRoute.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaStaticRouteIpv4Id(tmpl, ts, vr, o.Name))
	return readPanoramaStaticRouteIpv4(d, meta)
}

func readPanoramaStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaStaticRouteIpv4Id(d.Id())

	o, err := pano.Network.StaticRoute.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("name", o.Name)
	d.Set("virtual_router", vr)
	d.Set("destination", o.Destination)
	d.Set("interface", o.Interface)
	d.Set("type", o.Type)
	d.Set("next_hop", o.NextHop)
	d.Set("admin_distance", o.AdminDistance)
	d.Set("metric", o.Metric)
	d.Set("route_table", o.RouteTable)
	d.Set("bfd_profile", o.BfdProfile)

	return nil
}

func updatePanoramaStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaStaticRouteIpv4(d)

	lo, err := pano.Network.StaticRoute.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.StaticRoute.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaStaticRouteIpv4(d, meta)
}

func deletePanoramaStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaStaticRouteIpv4Id(d.Id())

	err := pano.Network.StaticRoute.Delete(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
