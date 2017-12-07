package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/router"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualRouter() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouter,
		Read:   readVirtualRouter,
		Update: updateVirtualRouter,
		Delete: deleteVirtualRouter,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The virtual router's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to import this virtual router into",
			},
			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"static_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"static_ipv6_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospf_int_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospf_ext_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospfv3_int_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospfv3_ext_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ibgp_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ebgp_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"rip_dist": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntInRange(10, 240),
			},
		},
	}
}

func parseVirtualRouterId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildVirtualRouterId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func parseVirtualRouter(d *schema.ResourceData) (string, router.Entry) {
	vsys := d.Get("vsys").(string)
	o := router.Entry{
		Name:           d.Get("name").(string),
		Interfaces:     asStringList(d, "interfaces"),
		StaticDist:     d.Get("static_dist").(int),
		StaticIpv6Dist: d.Get("static_ipv6_dist").(int),
		OspfIntDist:    d.Get("ospf_int_dist").(int),
		OspfExtDist:    d.Get("ospf_ext_dist").(int),
		Ospfv3IntDist:  d.Get("ospfv3_int_dist").(int),
		Ospfv3ExtDist:  d.Get("ospfv3_ext_dist").(int),
		IbgpDist:       d.Get("ibgp_dist").(int),
		EbgpDist:       d.Get("ebgp_dist").(int),
		RipDist:        d.Get("rip_dist").(int),
	}

	return vsys, o
}

func saveDataVirtualRouter(d *schema.ResourceData, vsys string, o router.Entry) {
	d.SetId(buildVirtualRouterId(vsys, o.Name))
	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("interfaces", o.Interfaces)
	d.Set("static_dist", o.StaticDist)
	d.Set("static_ipv6_dist", o.StaticIpv6Dist)
	d.Set("ospf_int_dist", o.OspfIntDist)
	d.Set("ospf_ext_dist", o.OspfExtDist)
	d.Set("ospfv3_int_dist", o.Ospfv3IntDist)
	d.Set("ospfv3_ext_dist", o.Ospfv3ExtDist)
	d.Set("ibgp_dist", o.IbgpDist)
	d.Set("ebgp_dist", o.EbgpDist)
	d.Set("rip_dist", o.RipDist)
}

func createVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseVirtualRouter(d)

	if err := fw.Network.VirtualRouter.Set(vsys, o); err != nil {
		return err
	}

	saveDataVirtualRouter(d, vsys, o)
	return nil
}

func readVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseVirtualRouterId(d.Id())

	o, err := fw.Network.VirtualRouter.Get(name)
	if err != nil {
		d.SetId("")
		return nil
	}

	saveDataVirtualRouter(d, vsys, o)
	return nil
}

func updateVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var err error
	fw := meta.(*pango.Firewall)
	vsys, o := parseVirtualRouter(d)

	lo, err := fw.Network.VirtualRouter.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	err = fw.Network.VirtualRouter.Edit(vsys, lo)

	if err == nil {
		saveDataVirtualRouter(d, vsys, o)
	}
	return err
}

func deleteVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	if name == "default" {
		_ = fw.Network.VirtualRouter.CleanupDefault(vsys)
	} else {
		_ = fw.Network.VirtualRouter.Delete(vsys, name)
	}
	d.SetId("")
	return nil
}
