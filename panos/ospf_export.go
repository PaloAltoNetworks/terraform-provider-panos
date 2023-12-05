package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/exp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceOspfExports() *schema.Resource {
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
		Read: readDataSourceOspfExports,

		Schema: s,
	}
}

func readDataSourceOspfExports(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	vr := d.Get("virtual_router").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vr
		listing, err = con.Network.OspfExport.GetList(vr)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = base64Encode([]string{
			tmpl, ts, vr,
		})
		listing, err = con.Network.OspfExport.GetList(tmpl, ts, vr)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceOspfExport() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceOspfExport,

		Schema: ospfExportSchema(false),
	}
}

func readDataSourceOspfExport(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o exp.Entry
	var id string
	vr := d.Get("virtual_router").(string)
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfExportId(vr, name)
		o, err = con.Network.OspfExport.Get(vr, name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfExportId(tmpl, ts, vr, name)
		o, err = con.Network.OspfExport.Get(tmpl, ts, vr, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveOspfExport(d, o)

	return nil
}

// Resource.
func resourceOspfExport() *schema.Resource {
	return &schema.Resource{
		Create: createOspfExport,
		Read:   readOspfExport,
		Update: updateOspfExport,
		Delete: deleteOspfExport,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ospfExportSchema(true),
	}
}

func createOspfExport(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	vr := d.Get("virtual_router").(string)
	o := loadOspfExport(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfExportId(vr, o.Name)
		err = con.Network.OspfExport.Set(vr, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfExportId(tmpl, ts, vr, o.Name)
		err = con.Network.OspfExport.Set(tmpl, ts, vr, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readOspfExport(d, meta)
}

func readOspfExport(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o exp.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfExportId(d.Id())
		o, err = con.Network.OspfExport.Get(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfExportId(d.Id())
		o, err = con.Network.OspfExport.Get(tmpl, ts, vr, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveOspfExport(d, o)
	return nil
}

func updateOspfExport(d *schema.ResourceData, meta interface{}) error {
	o := loadOspfExport(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfExportId(d.Id())
		lo, err := con.Network.OspfExport.Get(vr, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfExport.Edit(vr, o); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfExportId(d.Id())
		lo, err := con.Network.OspfExport.Get(tmpl, ts, vr, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfExport.Edit(tmpl, ts, vr, o); err != nil {
			return err
		}
	}

	return readOspfExport(d, meta)
}

func deleteOspfExport(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfExportId(d.Id())
		err = con.Network.OspfExport.Delete(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfExportId(d.Id())
		err = con.Network.OspfExport.Delete(tmpl, ts, vr, name)
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
func ospfExportSchema(isResource bool) map[string]*schema.Schema {
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
			Description: "The export rule name",
			ForceNew:    true,
		},
		"path_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path type",
			Default:     exp.PathTypeExt2,
			ValidateFunc: validateStringIn(
				exp.PathTypeExt1,
				exp.PathTypeExt2,
			),
		},
		"tag": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Tag",
		},
		"metric": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Metric",
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "virtual_router", "name"})
	}

	return ans
}

func loadOspfExport(d *schema.ResourceData) exp.Entry {
	return exp.Entry{
		Name:     d.Get("name").(string),
		PathType: d.Get("path_type").(string),
		Tag:      d.Get("tag").(string),
		Metric:   d.Get("metric").(int),
	}
}

func saveOspfExport(d *schema.ResourceData, o exp.Entry) {
	d.Set("name", o.Name)
	d.Set("path_type", o.PathType)
	d.Set("tag", o.Tag)
	d.Set("metric", o.Metric)
}

// Id functions.
func parseFirewallOspfExportId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaOspfExportId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildFirewallOspfExportId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaOspfExportId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
