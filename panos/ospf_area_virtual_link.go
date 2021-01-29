package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area/vlink"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceOspfAreaVirtualLinks() *schema.Resource {
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
		Read: readDataSourceOspfAreaVirtualLinks,

		Schema: s,
	}
}

func readDataSourceOspfAreaVirtualLinks(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	vr := d.Get("virtual_router").(string)
	area := d.Get("ospf_area").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = base64Encode([]interface{}{
			vr, area,
		})
		listing, err = con.Network.OspfAreaVirtualLink.GetList(vr, area)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = base64Encode([]interface{}{
			tmpl, ts, vr, area,
		})
		listing, err = con.Network.OspfAreaVirtualLink.GetList(tmpl, ts, vr, area)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceOspfAreaVirtualLink() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceOspfAreaVirtualLink,

		Schema: ospfAreaVirtualLinkSchema(false),
	}
}

func readDataSourceOspfAreaVirtualLink(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o vlink.Entry
	var id string
	vr := d.Get("virtual_router").(string)
	area := d.Get("ospf_area").(string)
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAreaVirtualLinkId(vr, area, name)
		o, err = con.Network.OspfAreaVirtualLink.Get(vr, area, name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAreaVirtualLinkId(tmpl, ts, vr, area, name)
		o, err = con.Network.OspfAreaVirtualLink.Get(tmpl, ts, vr, area, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveOspfAreaVirtualLink(d, o)

	return nil
}

// Resource.
func resourceOspfAreaVirtualLink() *schema.Resource {
	return &schema.Resource{
		Create: createOspfAreaVirtualLink,
		Read:   readOspfAreaVirtualLink,
		Update: updateOspfAreaVirtualLink,
		Delete: deleteOspfAreaVirtualLink,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ospfAreaVirtualLinkSchema(true),
	}
}

func createOspfAreaVirtualLink(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	vr := d.Get("virtual_router").(string)
	area := d.Get("ospf_area").(string)
	o := loadOspfAreaVirtualLink(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAreaVirtualLinkId(vr, area, o.Name)
		err = con.Network.OspfAreaVirtualLink.Set(vr, area, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAreaVirtualLinkId(tmpl, ts, vr, area, o.Name)
		err = con.Network.OspfAreaVirtualLink.Set(tmpl, ts, vr, area, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readOspfAreaVirtualLink(d, meta)
}

func readOspfAreaVirtualLink(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o vlink.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, area, name := parseFirewallOspfAreaVirtualLinkId(d.Id())
		o, err = con.Network.OspfAreaVirtualLink.Get(vr, area, name)
	case *pango.Panorama:
		tmpl, ts, vr, area, name := parsePanoramaOspfAreaVirtualLinkId(d.Id())
		o, err = con.Network.OspfAreaVirtualLink.Get(tmpl, ts, vr, area, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveOspfAreaVirtualLink(d, o)
	return nil
}

func updateOspfAreaVirtualLink(d *schema.ResourceData, meta interface{}) error {
	o := loadOspfAreaVirtualLink(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, area, name := parseFirewallOspfAreaVirtualLinkId(d.Id())
		lo, err := con.Network.OspfAreaVirtualLink.Get(vr, area, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfAreaVirtualLink.Edit(vr, area, o); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl, ts, vr, area, name := parsePanoramaOspfAreaVirtualLinkId(d.Id())
		lo, err := con.Network.OspfAreaVirtualLink.Get(tmpl, ts, vr, area, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfAreaVirtualLink.Edit(tmpl, ts, vr, area, o); err != nil {
			return err
		}
	}

	return readOspfAreaVirtualLink(d, meta)
}

func deleteOspfAreaVirtualLink(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, area, name := parseFirewallOspfAreaVirtualLinkId(d.Id())
		err = con.Network.OspfAreaVirtualLink.Delete(vr, area, name)
	case *pango.Panorama:
		tmpl, ts, vr, area, name := parsePanoramaOspfAreaVirtualLinkId(d.Id())
		err = con.Network.OspfAreaVirtualLink.Delete(tmpl, ts, vr, area, name)
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
func ospfAreaVirtualLinkSchema(isResource bool) map[string]*schema.Schema {
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
		"neighbor_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Neighbor ID",
		},
		"transit_area_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Transit area ID",
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
		"auth_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Auth profile",
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

func loadOspfAreaVirtualLink(d *schema.ResourceData) vlink.Entry {
	return vlink.Entry{
		Name:               d.Get("name").(string),
		Enable:             d.Get("enable").(bool),
		NeighborId:         d.Get("neighbor_id").(string),
		TransitAreaId:      d.Get("transit_area_id").(string),
		HelloInterval:      d.Get("hello_interval").(int),
		DeadCounts:         d.Get("dead_counts").(int),
		RetransmitInterval: d.Get("retransmit_interval").(int),
		TransitDelay:       d.Get("transit_delay").(int),
		AuthProfile:        d.Get("auth_profile").(string),
		BfdProfile:         d.Get("bfd_profile").(string),
	}
}

func saveOspfAreaVirtualLink(d *schema.ResourceData, o vlink.Entry) {
	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("neighbor_id", o.NeighborId)
	d.Set("transit_area_id", o.TransitAreaId)
	d.Set("hello_interval", o.HelloInterval)
	d.Set("dead_counts", o.DeadCounts)
	d.Set("retransmit_interval", o.RetransmitInterval)
	d.Set("transit_delay", o.TransitDelay)
	d.Set("auth_profile", o.AuthProfile)
	d.Set("bfd_profile", o.BfdProfile)
}

// Id functions.
func parseFirewallOspfAreaVirtualLinkId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func parsePanoramaOspfAreaVirtualLinkId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildFirewallOspfAreaVirtualLinkId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func buildPanoramaOspfAreaVirtualLinkId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}
