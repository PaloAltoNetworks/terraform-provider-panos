package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaVirtualRouter() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaVirtualRouter,
		Read:   readPanoramaVirtualRouter,
		Update: updatePanoramaVirtualRouter,
		Delete: deletePanoramaVirtualRouter,

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
				ForceNew: true,
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
			"enable_ecmp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ecmp_symmetric_return": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ecmp_strict_source_path": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ecmp_max_path": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ecmp_load_balance_method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateStringIn(
					router.EcmpLoadBalanceMethodIpModulo,
					router.EcmpLoadBalanceMethodIpHash,
					router.EcmpLoadBalanceMethodWeightedRoundRobin,
					router.EcmpLoadBalanceMethodBalancedRoundRobin,
				),
			},
			"ecmp_hash_source_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ecmp_hash_use_port": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ecmp_hash_seed": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ecmp_weighted_round_robin_interfaces": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func parsePanoramaVirtualRouterId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaVirtualRouterId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaVirtualRouter(d *schema.ResourceData) (string, string, string, router.Entry) {
	tmpl := d.Get("template").(string)
	vsys := d.Get("vsys").(string)
	var iList map[string]int

	iMap := d.Get("ecmp_weighted_round_robin_interfaces").(map[string]interface{})
	if len(iMap) > 0 {
		iList = make(map[string]int)
		for key, value := range iMap {
			iList[key] = value.(int)
		}
	}

	o := router.Entry{
		Name:                             d.Get("name").(string),
		Interfaces:                       asStringList(d.Get("interfaces").([]interface{})),
		StaticDist:                       d.Get("static_dist").(int),
		StaticIpv6Dist:                   d.Get("static_ipv6_dist").(int),
		OspfIntDist:                      d.Get("ospf_int_dist").(int),
		OspfExtDist:                      d.Get("ospf_ext_dist").(int),
		Ospfv3IntDist:                    d.Get("ospfv3_int_dist").(int),
		Ospfv3ExtDist:                    d.Get("ospfv3_ext_dist").(int),
		IbgpDist:                         d.Get("ibgp_dist").(int),
		EbgpDist:                         d.Get("ebgp_dist").(int),
		RipDist:                          d.Get("rip_dist").(int),
		EnableEcmp:                       d.Get("enable_ecmp").(bool),
		EcmpSymmetricReturn:              d.Get("ecmp_symmetric_return").(bool),
		EcmpStrictSourcePath:             d.Get("ecmp_strict_source_path").(bool),
		EcmpMaxPath:                      d.Get("ecmp_max_path").(int),
		EcmpLoadBalanceMethod:            d.Get("ecmp_load_balance_method").(string),
		EcmpHashSourceOnly:               d.Get("ecmp_hash_source_only").(bool),
		EcmpHashUsePort:                  d.Get("ecmp_hash_use_port").(bool),
		EcmpHashSeed:                     d.Get("ecmp_hash_seed").(int),
		EcmpWeightedRoundRobinInterfaces: iList,
	}

	return tmpl, "", vsys, o
}

func createPanoramaVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaVirtualRouter(d)

	if err := pano.Network.VirtualRouter.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaVirtualRouterId(tmpl, ts, vsys, o.Name))
	return readPanoramaVirtualRouter(d, meta)
}

func readPanoramaVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaVirtualRouterId(d.Id())

	o, err := pano.Network.VirtualRouter.Get(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := pano.IsImported(util.VirtualRouterImport, tmpl, ts, vsys, name)
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
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
	d.Set("enable_ecmp", o.EnableEcmp)
	d.Set("ecmp_symmetric_return", o.EcmpSymmetricReturn)
	d.Set("ecmp_strict_source_path", o.EcmpStrictSourcePath)
	d.Set("ecmp_max_path", o.EcmpMaxPath)
	d.Set("ecmp_load_balance_method", o.EcmpLoadBalanceMethod)
	d.Set("ecmp_hash_source_only", o.EcmpHashSourceOnly)
	d.Set("ecmp_hash_use_port", o.EcmpHashUsePort)
	d.Set("ecmp_hash_seed", o.EcmpHashSeed)

	var bm map[string]interface{}
	if len(o.EcmpWeightedRoundRobinInterfaces) > 0 {
		bm = make(map[string]interface{})
		for key, value := range o.EcmpWeightedRoundRobinInterfaces {
			bm[key] = value
		}
	}

	if err = d.Set("ecmp_weighted_round_robin_interfaces", bm); err != nil {
		log.Printf("[WARN] Error setting 'ecmp_weighted_round_robin_interfaces' for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaVirtualRouter(d)

	lo, err := pano.Network.VirtualRouter.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.VirtualRouter.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaVirtualRouter(d, meta)
}

func deletePanoramaVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var err error
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaVirtualRouterId(d.Id())

	if name == "default" {
		err = pano.Network.VirtualRouter.CleanupDefault(tmpl, ts)
	} else {
		err = pano.Network.VirtualRouter.Delete(tmpl, ts, name)
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
