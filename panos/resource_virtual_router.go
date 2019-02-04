package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualRouter() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouter,
		Read:   readVirtualRouter,
		Update: updateVirtualRouter,
		Delete: deleteVirtualRouter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"static_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"static_ipv6_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospf_int_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospf_ext_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      110,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospfv3_int_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ospfv3_ext_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      110,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ibgp_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      200,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"ebgp_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: validateIntInRange(10, 240),
			},
			"rip_dist": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      120,
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
		Interfaces:     asStringList(d.Get("interfaces").([]interface{})),
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

func createVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseVirtualRouter(d)

	if err := fw.Network.VirtualRouter.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildVirtualRouterId(vsys, o.Name))
	return readVirtualRouter(d, meta)
}

func readVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseVirtualRouterId(d.Id())

	o, err := fw.Network.VirtualRouter.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := fw.IsImported(util.VirtualRouterImport, "", "", vsys, name)
	if err != nil {
		return err
	}

	d.Set("name", o.Name)
	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	if err := d.Set("interfaces", o.Interfaces); err != nil {
		log.Printf("[WARN] Error setting 'interfaces' for %q: %s", d.Id(), err)
	}
	d.Set("static_dist", o.StaticDist)
	d.Set("static_ipv6_dist", o.StaticIpv6Dist)
	d.Set("ospf_int_dist", o.OspfIntDist)
	d.Set("ospf_ext_dist", o.OspfExtDist)
	d.Set("ospfv3_int_dist", o.Ospfv3IntDist)
	d.Set("ospfv3_ext_dist", o.Ospfv3ExtDist)
	d.Set("ibgp_dist", o.IbgpDist)
	d.Set("ebgp_dist", o.EbgpDist)
	d.Set("rip_dist", o.RipDist)

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
	if err = fw.Network.VirtualRouter.Edit(vsys, lo); err != nil {
		return err
	}

	return readVirtualRouter(d, meta)
}

func deleteVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var err error
	fw := meta.(*pango.Firewall)
	_, name := parseVirtualRouterId(d.Id())

	if name == "default" {
		err = fw.Network.VirtualRouter.CleanupDefault()
	} else {
		err = fw.Network.VirtualRouter.Delete(name)
	}
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
