package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/group"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceSecurityProfileGroups() *schema.Resource {
	s := listingSchema()
	s["device_group"] = deviceGroupSchema()
	s["vsys"] = vsysSchema("vsys1")

	return &schema.Resource{
		Read: dataSourceSecurityProfileGroupsRead,

		Schema: s,
	}
}

func dataSourceSecurityProfileGroupsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vsys
		listing, err = con.Objects.SecurityProfileGroup.GetList(id)
	case *pango.Panorama:
		id = dg
		listing, err = con.Objects.SecurityProfileGroup.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceSecurityProfileGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityProfileGroupRead,

		Schema: securityProfileGroupSchema(false),
	}
}

func dataSourceSecurityProfileGroupRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o group.Entry

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildSecurityProfileGroupId(dg, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.SecurityProfileGroup.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.SecurityProfileGroup.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveSecurityProfileGroup(d, o)

	return nil
}

// Resource.
func resourceSecurityProfileGroup() *schema.Resource {
	return &schema.Resource{
		Create: createSecurityProfileGroup,
		Read:   readSecurityProfileGroup,
		Update: updateSecurityProfileGroup,
		Delete: deleteSecurityProfileGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: securityProfileGroupSchema(true),
	}
}

func createSecurityProfileGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadSecurityProfileGroup(d)

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildSecurityProfileGroupId(dg, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.SecurityProfileGroup.Set(vsys, o)
	case *pango.Panorama:
		err = con.Objects.SecurityProfileGroup.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readSecurityProfileGroup(d, meta)
}

func readSecurityProfileGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o group.Entry

	dg, vsys, name := parseSecurityProfileGroupId(d.Id())
	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.SecurityProfileGroup.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.SecurityProfileGroup.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSecurityProfileGroup(d, o)
	return nil
}

func updateSecurityProfileGroup(d *schema.ResourceData, meta interface{}) error {
	o := loadSecurityProfileGroup(d)
	dg, vsys, _ := parseSecurityProfileGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Objects.SecurityProfileGroup.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.SecurityProfileGroup.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Objects.SecurityProfileGroup.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.SecurityProfileGroup.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readSecurityProfileGroup(d, meta)
}

func deleteSecurityProfileGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	dg, vsys, name := parseSecurityProfileGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.SecurityProfileGroup.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Objects.SecurityProfileGroup.Delete(dg, name)
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
func securityProfileGroupSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"antivirus_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"anti_spyware_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"vulnerability_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"url_filtering_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"file_blocking_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"data_filtering_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"wildfire_analysis_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"gtp_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"sctp_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	return ans
}

func loadSecurityProfileGroup(d *schema.ResourceData) group.Entry {
	return group.Entry{
		Name:                    d.Get("name").(string),
		AntivirusProfile:        d.Get("antivirus_profile").(string),
		AntiSpywareProfile:      d.Get("anti_spyware_profile").(string),
		VulnerabilityProfile:    d.Get("vulnerability_profile").(string),
		UrlFilteringProfile:     d.Get("url_filtering_profile").(string),
		FileBlockingProfile:     d.Get("file_blocking_profile").(string),
		DataFilteringProfile:    d.Get("data_filtering_profile").(string),
		WildfireAnalysisProfile: d.Get("wildfire_analysis_profile").(string),
		GtpProfile:              d.Get("gtp_profile").(string),
		SctpProfile:             d.Get("sctp_profile").(string),
	}
}

func saveSecurityProfileGroup(d *schema.ResourceData, o group.Entry) {
	d.Set("name", o.Name)
	d.Set("antivirus_profile", o.AntivirusProfile)
	d.Set("anti_spyware_profile", o.AntiSpywareProfile)
	d.Set("vulnerability_profile", o.VulnerabilityProfile)
	d.Set("url_filtering_profile", o.UrlFilteringProfile)
	d.Set("file_blocking_profile", o.FileBlockingProfile)
	d.Set("data_filtering_profile", o.DataFilteringProfile)
	d.Set("wildfire_analysis_profile", o.WildfireAnalysisProfile)
	d.Set("gtp_profile", o.GtpProfile)
	d.Set("sctp_profile", o.SctpProfile)
}

// Id functions.
func parseSecurityProfileGroupId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildSecurityProfileGroupId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}
