package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/route/static/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceStaticRouteIpv4() *schema.Resource {
	return &schema.Resource{
		Create: createStaticRouteIpv4,
		Read:   readStaticRouteIpv4,
		Update: updateStaticRouteIpv4,
		Delete: deleteStaticRouteIpv4,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func parseStaticRouteIpv4(d *schema.ResourceData) (string, ipv4.Entry) {
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

	return vr, o
}

func parseStaticRouteIpv4Id(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildStaticRouteIpv4Id(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseStaticRouteIpv4(d)

	if err := fw.Network.StaticRoute.Set(vr, o); err != nil {
		return err
	}

	d.SetId(buildStaticRouteIpv4Id(vr, o.Name))
	return readStaticRouteIpv4(d, meta)
}

func readStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, name := parseStaticRouteIpv4Id(d.Id())

	o, err := fw.Network.StaticRoute.Get(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

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

func updateStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseStaticRouteIpv4(d)

	lo, err := fw.Network.StaticRoute.Get(vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.StaticRoute.Edit(vr, lo); err != nil {
		return err
	}

	return readStaticRouteIpv4(d, meta)
}

func deleteStaticRouteIpv4(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, name := parseStaticRouteIpv4Id(d.Id())

	err := fw.Network.StaticRoute.Delete(vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
