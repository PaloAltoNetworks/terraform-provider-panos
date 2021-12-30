package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceOspfAreas() *schema.Resource {
	s := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router name",
		},
	}

	for key, value := range listingSchema() {
		s[key] = value
	}

	return &schema.Resource{
		Read: readDataSourceOspfAreas,

		Schema: s,
	}
}

func readDataSourceOspfAreas(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	vr := d.Get("virtual_router").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vr
		listing, err = con.Network.OspfArea.GetList(vr)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = base64Encode([]string{
			tmpl, ts, vr,
		})
		listing, err = con.Network.OspfArea.GetList(tmpl, ts, vr)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceOspfArea() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceOspfArea,

		Schema: ospfAreaSchema(false),
	}
}

func readDataSourceOspfArea(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o area.Entry
	var id string
	vr := d.Get("virtual_router").(string)
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAreaId(vr, name)
		o, err = con.Network.OspfArea.Get(vr, name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAreaId(tmpl, ts, vr, name)
		o, err = con.Network.OspfArea.Get(tmpl, ts, vr, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveOspfArea(d, o)

	return nil
}

// Resource.
func resourceOspfArea() *schema.Resource {
	return &schema.Resource{
		Create: createOspfArea,
		Read:   readOspfArea,
		Update: updateOspfArea,
		Delete: deleteOspfArea,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ospfAreaSchema(true),
	}
}

func createOspfArea(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	vr := d.Get("virtual_router").(string)
	o := loadOspfArea(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAreaId(vr, o.Name)
		err = con.Network.OspfArea.Set(vr, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAreaId(tmpl, ts, vr, o.Name)
		err = con.Network.OspfArea.Set(tmpl, ts, vr, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readOspfArea(d, meta)
}

func readOspfArea(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o area.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfAreaId(d.Id())
		o, err = con.Network.OspfArea.Get(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfAreaId(d.Id())
		o, err = con.Network.OspfArea.Get(tmpl, ts, vr, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveOspfArea(d, o)
	return nil
}

func updateOspfArea(d *schema.ResourceData, meta interface{}) error {
	o := loadOspfArea(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfAreaId(d.Id())
		lo, err := con.Network.OspfArea.Get(vr, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfArea.Edit(vr, o); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfAreaId(d.Id())
		lo, err := con.Network.OspfArea.Get(tmpl, ts, vr, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfArea.Edit(tmpl, ts, vr, o); err != nil {
			return err
		}
	}

	return readOspfArea(d, meta)
}

func deleteOspfArea(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfAreaId(d.Id())
		err = con.Network.OspfArea.Delete(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfAreaId(d.Id())
		err = con.Network.OspfArea.Delete(tmpl, ts, vr, name)
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
func ospfAreaSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router",
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name",
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Area type",
			Default:     area.TypeNormal,
			ValidateFunc: validateStringIn(
				area.TypeNormal,
				area.TypeStub,
				area.TypeNssa,
			),
		},
		"accept_summary": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "(stub/nssa) Accept summary",
		},
		"default_route_advertise": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "(stub/nssa) Default route advertise",
		},
		"advertise_metric": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "(stub/nssa) Advertise metric",
		},
		"advertise_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(nssa) Advertise type",
		},
		"ext_range": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(nssa) List of EXT Range specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Network",
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Action",
						Default:     area.ActionAdvertise,
						ValidateFunc: validateStringIn(
							area.ActionAdvertise,
							area.ActionSuppress,
						),
					},
				},
			},
		},
		"range": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of range specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Network",
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Action",
						Default:     area.ActionAdvertise,
						ValidateFunc: validateStringIn(
							area.ActionAdvertise,
							area.ActionSuppress,
						),
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "virtual_router", "name"})
	}

	return ans
}

func loadOspfArea(d *schema.ResourceData) area.Entry {
	var list []interface{}

	var erList []area.Range
	if list = d.Get("ext_range").([]interface{}); len(list) > 0 {
		erList = make([]area.Range, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			erList = append(erList, area.Range{
				Network: elm["network"].(string),
				Action:  elm["action"].(string),
			})
		}
	}

	var rList []area.Range
	if list = d.Get("range").([]interface{}); len(list) > 0 {
		rList = make([]area.Range, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			rList = append(rList, area.Range{
				Network: elm["network"].(string),
				Action:  elm["action"].(string),
			})
		}
	}

	return area.Entry{
		Name:                  d.Get("name").(string),
		Type:                  d.Get("type").(string),
		AcceptSummary:         d.Get("accept_summary").(bool),
		DefaultRouteAdvertise: d.Get("default_route_advertise").(bool),
		AdvertiseMetric:       d.Get("advertise_metric").(int),
		AdvertiseType:         d.Get("advertise_type").(string),
		ExtRanges:             erList,
		Ranges:                rList,
	}
}

func saveOspfArea(d *schema.ResourceData, o area.Entry) {
	d.Set("name", o.Name)
	d.Set("type", o.Type)
	d.Set("accept_summary", o.AcceptSummary)
	d.Set("default_route_advertise", o.DefaultRouteAdvertise)
	d.Set("advertise_metric", o.AdvertiseMetric)
	d.Set("advertise_type", o.AdvertiseType)

	if len(o.ExtRanges) == 0 {
		d.Set("ext_range", nil)
	} else {
		list := make([]interface{}, 0, len(o.ExtRanges))
		for _, x := range o.ExtRanges {
			list = append(list, map[string]interface{}{
				"network": x.Network,
				"action":  x.Action,
			})
		}

		if err := d.Set("ext_range", list); err != nil {
			log.Printf("[WARN] Error setting 'ext_range' for %q: %s", d.Id(), err)
		}
	}

	if len(o.Ranges) == 0 {
		d.Set("range", nil)
	} else {
		list := make([]interface{}, 0, len(o.Ranges))
		for _, x := range o.Ranges {
			list = append(list, map[string]interface{}{
				"network": x.Network,
				"action":  x.Action,
			})
		}

		if err := d.Set("range", list); err != nil {
			log.Printf("[WARN] Error setting 'range' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func parseFirewallOspfAreaId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaOspfAreaId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildFirewallOspfAreaId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaOspfAreaId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
