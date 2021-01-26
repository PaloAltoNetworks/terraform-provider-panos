package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source.
func dataSourceOspf() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceOspf,

		Schema: ospfSchema(false),
	}
}

func readDataSourceOspf(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ospf.Config
	var id string
	vr := d.Get("virtual_router").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vr
		o, err = con.Network.OspfConfig.Get(vr)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfId(tmpl, ts, vr)
		o, err = con.Network.OspfConfig.Get(tmpl, ts, vr)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveOspf(d, o)

	return nil
}

// Resource.
func resourceOspf() *schema.Resource {
	return &schema.Resource{
		Create: createOspf,
		Read:   readOspf,
		Update: updateOspf,
		Delete: deleteOspf,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ospfSchema(true),
	}
}

func createOspf(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	vr := d.Get("virtual_router").(string)
	o := loadOspf(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vr
		err = con.Network.OspfConfig.Set(vr, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfId(tmpl, ts, vr)
		err = con.Network.OspfConfig.Set(tmpl, ts, vr, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readOspf(d, meta)
}

func readOspf(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ospf.Config

	switch con := meta.(type) {
	case *pango.Firewall:
		vr := d.Id()
		o, err = con.Network.OspfConfig.Get(vr)
	case *pango.Panorama:
		tmpl, ts, vr := parsePanoramaOspfId(d.Id())
		o, err = con.Network.OspfConfig.Get(tmpl, ts, vr)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveOspf(d, o)
	return nil
}

func updateOspf(d *schema.ResourceData, meta interface{}) error {
	o := loadOspf(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vr := d.Id()
		lo, err := con.Network.OspfConfig.Get(vr)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfConfig.Edit(vr, o); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl, ts, vr := parsePanoramaOspfId(d.Id())
		lo, err := con.Network.OspfConfig.Get(tmpl, ts, vr)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfConfig.Edit(tmpl, ts, vr, o); err != nil {
			return err
		}
	}

	return readOspf(d, meta)
}

func deleteOspf(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vr := d.Id()
		err = con.Network.OspfConfig.Delete(vr)
	case *pango.Panorama:
		tmpl, ts, vr := parsePanoramaOspfId(d.Id())
		err = con.Network.OspfConfig.Delete(tmpl, ts, vr)
	}

	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Schema handling.
func ospfSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router",
			ForceNew:    true,
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable flag",
		},
		"router_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Router ID",
		},
		"reject_default_route": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Reject default route",
		},
		"allow_redistribute_default_route": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow redistribute default route",
		},
		"rfc_1583": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "RFC 1583",
		},
		"spf_calculation_delay": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "SPF calculation delay",
		},
		"lsa_interval": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "LSA interval",
		},
		"enable_graceful_restart": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable graceful restart",
		},
		"grace_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Grace period",
		},
		"helper_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Helper enable",
		},
		"strict_lsa_checking": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Strict LSA checking",
		},
		"max_neighbor_restart_time": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Max neighbor restart time",
		},
		"bfd_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BFD profile name",
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "virtual_router"})
	}

	return ans
}

func loadOspf(d *schema.ResourceData) ospf.Config {
	return ospf.Config{
		Enable:                        d.Get("enable").(bool),
		RouterId:                      d.Get("router_id").(string),
		RejectDefaultRoute:            d.Get("reject_default_route").(bool),
		AllowRedistributeDefaultRoute: d.Get("allow_redistribute_default_route").(bool),
		Rfc1583:                       d.Get("rfc_1583").(bool),
		SpfCalculationDelay:           d.Get("spf_calculation_delay").(float64),
		LsaInterval:                   d.Get("lsa_interval").(float64),
		EnableGracefulRestart:         d.Get("enable_graceful_restart").(bool),
		GracePeriod:                   d.Get("grace_period").(int),
		HelperEnable:                  d.Get("helper_enable").(bool),
		StrictLsaChecking:             d.Get("strict_lsa_checking").(bool),
		MaxNeighborRestartTime:        d.Get("max_neighbor_restart_time").(int),
		BfdProfile:                    d.Get("bfd_profile").(string),
	}
}

func saveOspf(d *schema.ResourceData, o ospf.Config) {
	d.Set("enable", o.Enable)
	d.Set("router_id", o.RouterId)
	d.Set("reject_default_route", o.RejectDefaultRoute)
	d.Set("allow_redistribute_default_route", o.AllowRedistributeDefaultRoute)
	d.Set("rfc_1583", o.Rfc1583)
	d.Set("spf_calculation_delay", o.SpfCalculationDelay)
	d.Set("lsa_interval", o.LsaInterval)
	d.Set("enable_graceful_restart", o.EnableGracefulRestart)
	d.Set("grace_period", o.GracePeriod)
	d.Set("helper_enable", o.HelperEnable)
	d.Set("strict_lsa_checking", o.StrictLsaChecking)
	d.Set("max_neighbor_restart_time", o.MaxNeighborRestartTime)
	d.Set("bfd_profile", o.BfdProfile)
}

// Id functions.
func parsePanoramaOspfId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaOspfId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}
