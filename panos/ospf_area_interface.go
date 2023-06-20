package panos

import (
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/ospf/area/iface"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceOspfAreaInterfaces() *schema.Resource {
	s := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router name",
		},
		"ospf_area": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OSPF area name",
		},
	}

	for key, value := range listingSchema() {
		s[key] = value
	}

	return &schema.Resource{
		Read: readDataSourceOspfAreaInterfaces,

		Schema: s,
	}
}

func readDataSourceOspfAreaInterfaces(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	vr := d.Get("virtual_router").(string)
	area := d.Get("ospf_area").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = base64Encode([]string{
			vr, area,
		})
		listing, err = con.Network.OspfAreaInterface.GetList(vr, area)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = base64Encode([]string{
			tmpl, ts, vr, area,
		})
		listing, err = con.Network.OspfAreaInterface.GetList(tmpl, ts, vr, area)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceOspfAreaInterface() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceOspfAreaInterface,

		Schema: ospfAreaInterfaceSchema(false),
	}
}

func readDataSourceOspfAreaInterface(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o iface.Entry
	var id string
	vr := d.Get("virtual_router").(string)
	area := d.Get("ospf_area").(string)
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAreaInterfaceId(vr, area, name)
		o, err = con.Network.OspfAreaInterface.Get(vr, area, name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAreaInterfaceId(tmpl, ts, vr, area, name)
		o, err = con.Network.OspfAreaInterface.Get(tmpl, ts, vr, area, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveOspfAreaInterface(d, o)

	return nil
}

// Resource.
func resourceOspfAreaInterface() *schema.Resource {
	return &schema.Resource{
		Create: createOspfAreaInterface,
		Read:   readOspfAreaInterface,
		Update: updateOspfAreaInterface,
		Delete: deleteOspfAreaInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ospfAreaInterfaceSchema(true),
	}
}

func createOspfAreaInterface(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	vr := d.Get("virtual_router").(string)
	area := d.Get("ospf_area").(string)
	o := loadOspfAreaInterface(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAreaInterfaceId(vr, area, o.Name)
		err = con.Network.OspfAreaInterface.Set(vr, area, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAreaInterfaceId(tmpl, ts, vr, area, o.Name)
		err = con.Network.OspfAreaInterface.Set(tmpl, ts, vr, area, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readOspfAreaInterface(d, meta)
}

func readOspfAreaInterface(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o iface.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, area, name := parseFirewallOspfAreaInterfaceId(d.Id())
		o, err = con.Network.OspfAreaInterface.Get(vr, area, name)
	case *pango.Panorama:
		tmpl, ts, vr, area, name := parsePanoramaOspfAreaInterfaceId(d.Id())
		o, err = con.Network.OspfAreaInterface.Get(tmpl, ts, vr, area, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveOspfAreaInterface(d, o)
	return nil
}

func updateOspfAreaInterface(d *schema.ResourceData, meta interface{}) error {
	o := loadOspfAreaInterface(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, area, name := parseFirewallOspfAreaInterfaceId(d.Id())
		lo, err := con.Network.OspfAreaInterface.Get(vr, area, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfAreaInterface.Edit(vr, area, o); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl, ts, vr, area, name := parsePanoramaOspfAreaInterfaceId(d.Id())
		lo, err := con.Network.OspfAreaInterface.Get(tmpl, ts, vr, area, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfAreaInterface.Edit(tmpl, ts, vr, area, o); err != nil {
			return err
		}
	}

	return readOspfAreaInterface(d, meta)
}

func deleteOspfAreaInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, area, name := parseFirewallOspfAreaInterfaceId(d.Id())
		err = con.Network.OspfAreaInterface.Delete(vr, area, name)
	case *pango.Panorama:
		tmpl, ts, vr, area, name := parsePanoramaOspfAreaInterfaceId(d.Id())
		err = con.Network.OspfAreaInterface.Delete(tmpl, ts, vr, area, name)
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
func ospfAreaInterfaceSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router",
			ForceNew:    true,
		},
		"ospf_area": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OSPF area name",
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name",
			ForceNew:    true,
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable",
			Default:     true,
		},
		"passive": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Passive",
		},
		"link_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Link type",
			Default:     iface.LinkTypeBroadcast,
			ValidateFunc: validateStringIn(
				iface.LinkTypeBroadcast,
				iface.LinkTypePointToPoint,
				iface.LinkTypePointToMultiPoint,
			),
		},
		"metric": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Metric",
			Default:     10,
		},
		"priority": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Priority",
			Default:     1,
		},
		"hello_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Hello interval in seconds",
			Default:     10,
		},
		"dead_counts": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Dead counts",
			Default:     4,
		},
		"retransmit_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Retransmit interval in seconds",
			Default:     5,
		},
		"transit_delay": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Transit delay in seconds",
			Default:     1,
		},
		"grace_restart_delay": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Graceful restart hello delay in seconds",
			Default:     10,
		},
		"auth_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Auth profile",
		},
		"neighbors": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "(p2mp) List of neighbors",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"bfd_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BFD profile",
		},
	}

	if !isResource {
		computed(ans, "", []string{
			"template", "template_stack",
			"virtual_router", "ospf_area", "name",
		})
	}

	return ans
}

func loadOspfAreaInterface(d *schema.ResourceData) iface.Entry {
	return iface.Entry{
		Name:               d.Get("name").(string),
		Enable:             d.Get("enable").(bool),
		Passive:            d.Get("passive").(bool),
		LinkType:           d.Get("link_type").(string),
		Metric:             d.Get("metric").(int),
		Priority:           d.Get("priority").(int),
		HelloInterval:      d.Get("hello_interval").(int),
		DeadCounts:         d.Get("dead_counts").(int),
		RetransmitInterval: d.Get("retransmit_interval").(int),
		TransitDelay:       d.Get("transit_delay").(int),
		GraceRestartDelay:  d.Get("grace_restart_delay").(int),
		AuthProfile:        d.Get("auth_profile").(string),
		Neighbors:          setAsList(d.Get("neighbors").(*schema.Set)),
		BfdProfile:         d.Get("bfd_profile").(string),
	}
}

func saveOspfAreaInterface(d *schema.ResourceData, o iface.Entry) {
	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("passive", o.Passive)
	d.Set("link_type", o.LinkType)
	d.Set("metric", o.Metric)
	d.Set("priority", o.Priority)
	d.Set("hello_interval", o.HelloInterval)
	d.Set("dead_counts", o.DeadCounts)
	d.Set("retransmit_interval", o.RetransmitInterval)
	d.Set("transit_delay", o.TransitDelay)
	d.Set("grace_restart_delay", o.GraceRestartDelay)
	d.Set("auth_profile", o.AuthProfile)
	if err := d.Set("neighbors", listAsSet(o.Neighbors)); err != nil {
		log.Printf("[WARN] Error setting 'neighbors' for %q: %s", d.Id(), err)
	}
	d.Set("bfd_profile", o.BfdProfile)
}

// Id functions.
func parseFirewallOspfAreaInterfaceId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func parsePanoramaOspfAreaInterfaceId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildFirewallOspfAreaInterfaceId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func buildPanoramaOspfAreaInterfaceId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}
