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

func parseVirtualRouterId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildVirtualRouterId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func parseVirtualRouter(d *schema.ResourceData) (string, router.Entry) {
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
