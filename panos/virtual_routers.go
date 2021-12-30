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

// Data source (listing).
func dataSourceVirtualRouters() *schema.Resource {
	s := listingSchema()
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceVirtualRoutersRead,

		Schema: s,
	}
}

func dataSourceVirtualRoutersRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	id := buildVirtualRouterId(tmpl, ts, "", "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Network.VirtualRouter.GetList()
	case *pango.Panorama:
		listing, err = con.Network.VirtualRouter.GetList(tmpl, ts)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceVirtualRouter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVirtualRouterRead,

		Schema: virtualRouterSchema(false, []string{"vsys"}),
	}
}

func dataSourceVirtualRouterRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o router.Entry

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	id := buildVirtualRouterId(tmpl, ts, "", name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.VirtualRouter.Get(name)
	case *pango.Panorama:
		o, err = con.Network.VirtualRouter.Get(tmpl, ts, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveVirtualRouter(d, o)
	return nil
}

// Resources.
func resourceVirtualRouter() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouter,
		Read:   readVirtualRouter,
		Update: updateVirtualRouter,
		Delete: deleteVirtualRouter,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: virtualRouterSchema(true, []string{"template", "template_stack"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: virtualRouterUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: virtualRouterSchema(true, nil),
	}
}

func resourcePanoramaVirtualRouter() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouter,
		Read:   readVirtualRouter,
		Update: updateVirtualRouter,
		Delete: deleteVirtualRouter,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: virtualRouterSchema(true, []string{"template_stack"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: virtualRouterUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: virtualRouterSchema(true, nil),
	}
}

func virtualRouterUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["template"]; !ok {
		raw["template"] = ""
	}
	if _, ok := raw["template_stack"]; !ok {
		raw["template_stack"] = ""
	}

	return raw, nil
}

func createVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadVirtualRouter(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildVirtualRouterId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.VirtualRouter.Set(vsys, o)
	case *pango.Panorama:
		err = con.Network.VirtualRouter.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readVirtualRouter(d, meta)
}

func readVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var importIsCorrect bool
	var err error
	var o router.Entry

	// Migrate the ID.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 2 {
		// Old NGFW ID.
		d.SetId(buildVirtualRouterId("", "", tok[0], tok[1]))
	}

	tmpl, ts, vsys, name := parseVirtualRouterId(d.Id())
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.VirtualRouter.Get(name)
		if err == nil {
			importIsCorrect, err = con.IsImported(util.VirtualRouterImport, "", "", vsys, name)
		}
	case *pango.Panorama:
		o, err = con.Network.VirtualRouter.Get(tmpl, ts, name)
		if err == nil {
			importIsCorrect, err = con.IsImported(util.VirtualRouterImport, tmpl, ts, vsys, name)
		}
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if importIsCorrect {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	saveVirtualRouter(d, o)

	return nil
}

func updateVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	o := loadVirtualRouter(d)

	vsys := d.Get("vsys").(string)
	tmpl, ts, _, _ := parseVirtualRouterId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Network.VirtualRouter.Get(o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.VirtualRouter.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Network.VirtualRouter.Get(tmpl, ts, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.VirtualRouter.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
	}

	d.SetId(buildVirtualRouterId(tmpl, ts, vsys, o.Name))

	return readVirtualRouter(d, meta)
}

func deleteVirtualRouter(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl, ts, _, name := parseVirtualRouterId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		if name == "default" {
			err = con.Network.VirtualRouter.CleanupDefault()
		} else {
			err = con.Network.VirtualRouter.Delete(name)
		}
	case *pango.Panorama:
		if name == "default" {
			err = con.Network.VirtualRouter.CleanupDefault(tmpl, ts)
		} else {
			err = con.Network.VirtualRouter.Delete(tmpl, ts, name)
		}
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource (entry).
func resourceVirtualRouterEntry() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouterEntry,
		Read:   readVirtualRouterEntry,
		Delete: deleteVirtualRouterEntry,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: virtualRouterEntrySchema([]string{"template", "template_stack"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: virtualRouterUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: virtualRouterEntrySchema(nil),
	}
}

func resourcePanoramaVirtualRouterEntry() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouterEntry,
		Read:   readVirtualRouterEntry,
		Delete: deleteVirtualRouterEntry,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: virtualRouterEntrySchema([]string{"template_stack"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: virtualRouterUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: virtualRouterEntrySchema(nil),
	}
}

func createVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vr := d.Get("virtual_router").(string)
	iface := d.Get("interface").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("virtual_router", vr)
	d.Set("interface", iface)

	id := buildVirtualRouterId(tmpl, ts, vr, iface)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.VirtualRouter.SetInterface(vr, iface)
	case *pango.Panorama:
		err = con.Network.VirtualRouter.SetInterface(tmpl, ts, vr, iface)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readVirtualRouterEntry(d, meta)
}

func readVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o router.Entry

	// Migrate the resource ID.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 2 {
		d.SetId(buildVirtualRouterId("", "", tok[0], tok[1]))
	} else if len(tok) != 4 {
		return fmt.Errorf("Invalid ID (expecting len 2 or len4 ID: %s", d.Id())
	}

	tmpl, ts, vr, iface := parseVirtualRouterId(d.Id())

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("virtual_router", vr)
	d.Set("interface", iface)

	// Two possibilities: either the router isn't present or the interface is
	// not in the virtual router.
	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.VirtualRouter.Get(vr)
	case *pango.Panorama:
		o, err = con.Network.VirtualRouter.Get(tmpl, ts, vr)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	for _, x := range o.Interfaces {
		if x == iface {
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl, ts, vr, iface := parseVirtualRouterId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.VirtualRouter.DeleteInterface(vr, iface)
	case *pango.Panorama:
		err = con.Network.VirtualRouter.DeleteInterface(tmpl, ts, vr, iface)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func virtualRouterSchema(isResource bool, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Description: "The template.",
			Optional:    true,
			ForceNew:    true,
		},
		"template_stack": {
			Type:        schema.TypeString,
			Description: "The template stack.",
			Optional:    true,
			ForceNew:    true,
		},
		"vsys": {
			Type:        schema.TypeString,
			Description: "The vsys.",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name.",
			Required:    true,
			ForceNew:    true,
		},
		"interfaces": {
			Type:        schema.TypeList,
			Description: "List of interfaces in this virtual router.",
			Optional:    true,
			Computed:    true,
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
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "name"})
	}

	return ans
}

func virtualRouterEntrySchema(rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Description: "The template.",
			Optional:    true,
			ForceNew:    true,
		},
		"template_stack": {
			Type:        schema.TypeString,
			Description: "The template stack.",
			Optional:    true,
			ForceNew:    true,
		},
		"virtual_router": {
			Type:        schema.TypeString,
			Description: "The virtual router name.",
			Required:    true,
			ForceNew:    true,
		},
		"interface": {
			Type:        schema.TypeString,
			Description: "The interface name.",
			Required:    true,
			ForceNew:    true,
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	return ans
}

func loadVirtualRouter(d *schema.ResourceData) router.Entry {
	var iList map[string]int

	iMap := d.Get("ecmp_weighted_round_robin_interfaces").(map[string]interface{})
	if len(iMap) > 0 {
		iList = make(map[string]int)
		for key, value := range iMap {
			iList[key] = value.(int)
		}
	}

	return router.Entry{
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
}

func saveVirtualRouter(d *schema.ResourceData, o router.Entry) {
	d.Set("name", o.Name)
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

	if err := d.Set("ecmp_weighted_round_robin_interfaces", bm); err != nil {
		log.Printf("[WARN] Error setting 'ecmp_weighted_round_robin_interfaces' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseVirtualRouterId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildVirtualRouterId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
