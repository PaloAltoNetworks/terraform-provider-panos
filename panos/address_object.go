package panos

import (
	"log"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceAddressObjects() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceAddressObjectsRead,

		Schema: s,
	}
}

func dataSourceAddressObjectsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.Address.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.Address.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceAddressObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAddressObjectRead,

		Schema: addressObjectSchema(false, true, nil),
	}
}

func dataSourceAddressObjectRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o addr.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAddressObjectId(vsys, name)
		o, err = con.Objects.Address.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAddressObjectId(dg, name)
		o, err = con.Objects.Address.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveAddressObject(d, o)

	return nil
}

// Resource.
func resourceAddressObject() *schema.Resource {
	return &schema.Resource{
		Create: createAddressObject,
		Read:   readAddressObject,
		Update: updateAddressObject,
		Delete: deleteAddressObject,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: addressObjectSchema(true, true, []string{"device_group"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: addressObjectUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: addressObjectSchema(true, true, nil),
	}
}

func resourcePanoramaAddressObject() *schema.Resource {
	return &schema.Resource{
		Create: createAddressObject,
		Read:   readAddressObject,
		Update: updateAddressObject,
		Delete: deleteAddressObject,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: addressObjectSchema(true, true, []string{"vsys"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: addressObjectUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: addressObjectSchema(true, true, nil),
	}
}

func addressObjectUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["vsys"]; ok {
		raw["device_group"] = "shared"
	}
	if _, ok := raw["device_group"]; ok {
		raw["vsys"] = "vsys1"
	}

	return raw, nil
}

func createAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadAddressObject(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAddressObjectId(vsys, o.Name)
		err = con.Objects.Address.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAddressObjectId(dg, o.Name)
		err = con.Objects.Address.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readAddressObject(d, meta)
}

func readAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o addr.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseAddressObjectId(d.Id())
		o, err = con.Objects.Address.Get(vsys, name)
		d.Set("vsys", vsys)
		d.Set("device_group", "shared")
	case *pango.Panorama:
		dg, name := parseAddressObjectId(d.Id())
		o, err = con.Objects.Address.Get(dg, name)
		d.Set("vsys", "vsys1")
		d.Set("device_group", dg)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAddressObject(d, o)
	return nil
}

func updateAddressObject(d *schema.ResourceData, meta interface{}) error {
	o := loadAddressObject(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.Address.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Address.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.Address.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Address.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readAddressObject(d, meta)
}

func deleteAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseAddressObjectId(d.Id())
		err = con.Objects.Address.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseAddressObjectId(d.Id())
		err = con.Objects.Address.Delete(dg, name)
	}

	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Resource (bulk).
func resourceAddressObjects() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateAddressObjects,
		Read:   readAddressObjects,
		Update: createUpdateAddressObjects,
		Delete: deleteAddressObjects,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: addressObjectsSchema(true),
	}
}

func createUpdateAddressObjects(d *schema.ResourceData, meta interface{}) error {
	var err error
	var prevNames []string
	objs := loadAddressObjects(d)

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	if d.Id() != "" {
		_, _, prevNames = parseAddressObjectsId(d.Id())
	}

	id := buildAddressObjectsId(dg, vsys, objs)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.Address.ConfigureGroup(vsys, objs, prevNames)
	case *pango.Panorama:
		err = con.Objects.Address.ConfigureGroup(dg, objs, prevNames)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readAddressObjects(d, meta)
}

func readAddressObjects(d *schema.ResourceData, meta interface{}) error {
	var err error
	var objs, list []addr.Entry

	dg, vsys, names := parseAddressObjectsId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		list, err = con.Objects.Address.GetAll(vsys)
	case *pango.Panorama:
		list, err = con.Objects.Address.GetAll(dg)
	}

	if err != nil {
		return err
	}

	objs = make([]addr.Entry, 0, len(list))
	for _, name := range names {
		for _, x := range list {
			if x.Name == name {
				objs = append(objs, x)
				break
			}
		}
	}

	if len(objs) == 0 {
		d.SetId("")
		return nil
	}

	saveAddressObjects(d, objs)
	return nil
}

func deleteAddressObjects(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, vsys, names := parseAddressObjectsId(d.Id())
	ilist := make([]interface{}, 0, len(names))
	for _, x := range names {
		ilist = append(ilist, x)
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.Address.Delete(vsys, ilist...)
	case *pango.Panorama:
		err = con.Objects.Address.Delete(dg, ilist...)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func addressObjectSchema(isResource, forceNew bool, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      addr.IpNetmask,
			ValidateFunc: validateStringIn(addr.IpNetmask, addr.IpRange, addr.Fqdn, addr.IpWildcard),
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tags": tagSchema(),
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	return ans
}

func addressObjectsSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"object": {
			Type:        schema.TypeSet,
			Description: "An address object spec.",
			MinItems:    1,
			Required:    true,
			Elem: &schema.Resource{
				Schema: addressObjectSchema(true, false, []string{"device_group", "vsys"}),
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group"})
	}

	return ans
}

func loadAddressObject(d *schema.ResourceData) addr.Entry {
	return addr.Entry{
		Name:        d.Get("name").(string),
		Value:       d.Get("value").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        asStringList(d.Get("tags").([]interface{})),
	}
}

func saveAddressObject(d *schema.ResourceData, o addr.Entry) {
	d.Set("name", o.Name)
	d.Set("type", o.Type)
	d.Set("value", o.Value)
	d.Set("description", o.Description)
	if err := d.Set("tags", o.Tags); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
}

func loadAddressObjects(d *schema.ResourceData) []addr.Entry {
	olist := d.Get("object").(*schema.Set).List()
	ans := make([]addr.Entry, 0, len(olist))

	for i := range olist {
		elm := olist[i].(map[string]interface{})
		ans = append(ans, addr.Entry{
			Name:        elm["name"].(string),
			Value:       elm["value"].(string),
			Type:        elm["type"].(string),
			Description: elm["description"].(string),
			Tags:        asStringList(elm["tags"].([]interface{})),
		})
	}

	return ans
}

func saveAddressObjects(d *schema.ResourceData, objs []addr.Entry) {
	items := make([]interface{}, 0, len(objs))

	for _, x := range objs {
		var tlist []interface{}
		if len(x.Tags) > 0 {
			tlist = make([]interface{}, 0, len(x.Tags))
			for _, x := range x.Tags {
				tlist = append(tlist, x)
			}
		}

		items = append(items, map[string]interface{}{
			"name":        x.Name,
			"type":        x.Type,
			"value":       x.Value,
			"description": x.Description,
			"tags":        tlist,
		})
	}

	s := schema.NewSet(
		schema.HashResource(
			addressObjectsSchema(true)["object"].Elem.(*schema.Resource),
		),
		items,
	)

	if err := d.Set("object", s); err != nil {
		log.Printf("[WARN] Error setting 'object' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseAddressObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildAddressObjectId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildAddressObjectsId(a, b string, c []addr.Entry) string {
	list := make([]string, 0, len(c))
	for _, x := range c {
		list = append(list, x.Name)
	}

	return strings.Join([]string{a, b, base64Encode(list)}, IdSeparator)
}

func parseAddressObjectsId(v string) (string, string, []string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], base64Decode(t[2])
}
