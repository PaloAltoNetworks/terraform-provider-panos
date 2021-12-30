package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	vis "github.com/PaloAltoNetworks/pango/dev/vminfosource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Resource.
func resourceVmInformationSource() *schema.Resource {
	return &schema.Resource{
		Create: createVmInformationSource,
		Read:   readVmInformationSource,
		Update: createVmInformationSource,
		Delete: deleteVmInformationSource,

		Schema: vmInformationSourceSchema(),
	}
}

func createVmInformationSource(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo vis.Entry

	o := loadVmInformationSource(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildVmInformationSourceId(tmpl, vsys, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		if err = con.Device.VmInfoSource.Set(vsys, o); err == nil {
			lo, err = con.Device.VmInfoSource.Get(vsys, o.Name)
		}
	case *pango.Panorama:
		if err = con.Device.VmInfoSource.Set(tmpl, ts, vsys, o); err == nil {
			lo, err = con.Device.VmInfoSource.Get(tmpl, ts, vsys, o.Name)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveVmInformationSourceHashes(d, o, lo)

	return readVmInformationSource(d, meta)
}

func readVmInformationSource(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o vis.Entry

	tmpl, ts, vsys, name := parseVmInformationSourceId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.VmInfoSource.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.VmInfoSource.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveVmInformationSource(d, o)
	return nil
}

func deleteVmInformationSource(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys, name := parseVmInformationSourceId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.VmInfoSource.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.VmInfoSource.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func vmInformationSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("vsys1"),
		"name": {
			Type:        schema.TypeString,
			Description: "The name.",
			Required:    true,
			ForceNew:    true,
		},
		"settings": {
			Type:        schema.TypeMap,
			Description: "Configured and encrypted values.",
			Computed:    true,
			Sensitive:   true,
			Elem: &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
			},
		},

		// Sections.
		"aws_vpc": {
			Type:          schema.TypeList,
			Description:   "AWS VPC information source.",
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"esxi", "vcenter", "google_compute"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:        schema.TypeString,
						Description: "The description.",
						Optional:    true,
					},
					"disabled": {
						Type:        schema.TypeBool,
						Description: "Disabled or not.",
						Optional:    true,
					},
					"source": {
						Type:        schema.TypeString,
						Description: "IP address or name.",
						Required:    true,
					},
					"access_key_id": {
						Type:        schema.TypeString,
						Description: "AWS access key ID.",
						Required:    true,
					},
					"secret_access_key": {
						Type:        schema.TypeString,
						Description: "AWS secret access key.",
						Required:    true,
						Sensitive:   true,
					},
					"update_interval": {
						Type:        schema.TypeInt,
						Description: "Time interval (in sec) for updates.",
						Optional:    true,
						Default:     60,
					},
					"enable_timeout": {
						Type:        schema.TypeBool,
						Description: "Enable vm-info timeout when source is disconnected.",
						Optional:    true,
					},
					"timeout": {
						Type:        schema.TypeInt,
						Description: "The vm-info timeout value (in hours) when source is disconnected.",
						Optional:    true,
						Default:     2,
					},
					"vpc_id": {
						Type:        schema.TypeString,
						Description: "AWS VPC name or ID.",
						Required:    true,
					},
				},
			},
		},

		"esxi": {
			Type:          schema.TypeList,
			Description:   "VMware ESXi information source.",
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"aws_vpc", "vcenter", "google_compute"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:        schema.TypeString,
						Description: "The description.",
						Optional:    true,
					},
					"port": {
						Type:        schema.TypeInt,
						Description: "The port number.",
						Optional:    true,
						Default:     443,
					},
					"disabled": {
						Type:        schema.TypeBool,
						Description: "Disabled or not.",
						Optional:    true,
					},
					"enable_timeout": {
						Type:        schema.TypeBool,
						Description: "Enable vm-info timeout when source is disconnected.",
						Optional:    true,
					},
					"timeout": {
						Type:        schema.TypeInt,
						Description: "The vm-info timeout value (in hours) when source is disconnected.",
						Optional:    true,
						Default:     2,
					},
					"source": {
						Type:        schema.TypeString,
						Description: "IP address or source name for vm-info-source.",
						Required:    true,
					},
					"username": {
						Type:        schema.TypeString,
						Description: "The vm-info-source login username.",
						Required:    true,
					},
					"password": {
						Type:        schema.TypeString,
						Description: "The vm-info-source login password.",
						Required:    true,
					},
					"update_interval": {
						Type:        schema.TypeInt,
						Description: "Time interval (in sec) for updates.",
						Optional:    true,
						Default:     5,
					},
				},
			},
		},

		"vcenter": {
			Type:          schema.TypeList,
			Description:   "VMware vCenter information source.",
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"aws_vpc", "esxi", "google_compute"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:        schema.TypeString,
						Description: "The description.",
						Optional:    true,
					},
					"port": {
						Type:        schema.TypeInt,
						Description: "The port number.",
						Optional:    true,
						Default:     443,
					},
					"disabled": {
						Type:        schema.TypeBool,
						Description: "Disabled or not.",
						Optional:    true,
					},
					"enable_timeout": {
						Type:        schema.TypeBool,
						Description: "Enable vm-info timeout when source is disconnected.",
						Optional:    true,
					},
					"timeout": {
						Type:        schema.TypeInt,
						Description: "The vm-info timeout value (in hours) when source is disconnected.",
						Optional:    true,
						Default:     2,
					},
					"source": {
						Type:        schema.TypeString,
						Description: "IP address or source name for vm-info-source.",
						Required:    true,
					},
					"username": {
						Type:        schema.TypeString,
						Description: "The vm-info-source login username.",
						Required:    true,
					},
					"password": {
						Type:        schema.TypeString,
						Description: "The vm-info-source login password.",
						Required:    true,
					},
					"update_interval": {
						Type:        schema.TypeInt,
						Description: "Time interval (in sec) for updates.",
						Optional:    true,
						Default:     5,
					},
				},
			},
		},

		"google_compute": {
			Type:          schema.TypeList,
			Description:   "Google compute engine information source.",
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"aws_vpc", "esxi", "vcenter"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:        schema.TypeString,
						Description: "The description.",
						Optional:    true,
					},
					"disabled": {
						Type:        schema.TypeBool,
						Description: "Disabled or not.",
						Optional:    true,
					},
					"auth_type": {
						Type:        schema.TypeString,
						Description: "The auth type.",
						Optional:    true,
						Default:     vis.AuthTypeServiceInGce,
						ValidateFunc: validateStringIn(
							vis.AuthTypeServiceInGce,
							vis.AuthTypeServiceAccount,
						),
					},
					"service_account_credential": {
						Type:        schema.TypeString,
						Description: "GCE service account JSON file.",
						Optional:    true,
					},
					"project_id": {
						Type:        schema.TypeString,
						Description: "Google Compute Engine Project-ID.",
						Required:    true,
					},
					"zone_name": {
						Type:        schema.TypeString,
						Description: "Google Compute Engine project zone name.",
						Required:    true,
					},
					"update_interval": {
						Type:        schema.TypeInt,
						Description: "Time interval (in sec) for updates.",
						Optional:    true,
						Default:     5,
					},
					"enable_timeout": {
						Type:        schema.TypeBool,
						Description: "Enable vm-info timeout when source is disconnected.",
						Optional:    true,
					},
					"timeout": {
						Type:        schema.TypeInt,
						Description: "The vm-info timeout value (in hours) when source is disconnected.",
						Optional:    true,
						Default:     2,
					},
				},
			},
		},
	}
}

func loadVmInformationSource(d *schema.ResourceData) vis.Entry {
	ans := vis.Entry{
		Name: d.Get("name").(string),
	}

	if x := configFolder(d, "aws_vpc"); x != nil {
		ans.AwsVpc = &vis.AwsVpc{
			Description:     x["description"].(string),
			Disabled:        x["disabled"].(bool),
			Source:          x["source"].(string),
			AccessKeyId:     x["access_key_id"].(string),
			SecretAccessKey: x["secret_access_key"].(string),
			UpdateInterval:  x["update_interval"].(int),
			EnableTimeout:   x["enable_timeout"].(bool),
			Timeout:         x["timeout"].(int),
			VpcId:           x["vpc_id"].(string),
		}
	}

	if x := configFolder(d, "esxi"); x != nil {
		ans.Esxi = &vis.Esxi{
			Description:    x["description"].(string),
			Port:           x["port"].(int),
			Disabled:       x["disabled"].(bool),
			EnableTimeout:  x["enable_timeout"].(bool),
			Timeout:        x["timeout"].(int),
			Source:         x["source"].(string),
			Username:       x["username"].(string),
			Password:       x["password"].(string),
			UpdateInterval: x["update_interval"].(int),
		}
	}

	if x := configFolder(d, "vcenter"); x != nil {
		ans.Vcenter = &vis.Vcenter{
			Description:    x["description"].(string),
			Port:           x["port"].(int),
			Disabled:       x["disabled"].(bool),
			EnableTimeout:  x["enable_timeout"].(bool),
			Timeout:        x["timeout"].(int),
			Source:         x["source"].(string),
			Username:       x["username"].(string),
			Password:       x["password"].(string),
			UpdateInterval: x["update_interval"].(int),
		}
	}

	if x := configFolder(d, "google_compute"); x != nil {
		ans.GoogleCompute = &vis.GoogleCompute{
			Description:              x["description"].(string),
			Disabled:                 x["disabled"].(bool),
			AuthType:                 x["auth_type"].(string),
			ServiceAccountCredential: x["service_account_credential"].(string),
			ProjectId:                x["project_id"].(string),
			ZoneName:                 x["zone_name"].(string),
			UpdateInterval:           x["update_interval"].(int),
			EnableTimeout:            x["enable_timeout"].(bool),
			Timeout:                  x["timeout"].(int),
		}
	}

	return ans
}

func saveVmInformationSource(d *schema.ResourceData, o vis.Entry) {
	var err error
	settings := d.Get("settings").(map[string]interface{})
	d.Set("name", o.Name)

	if o.AwsVpc == nil {
		d.Set("aws_vpc", nil)
	} else {
		var sk string
		if settings["aws_secret_key"] != nil {
			if o.AwsVpc.SecretAccessKey == settings["aws_secret_key_enc"].(string) {
				sk = settings["aws_secret_key"].(string)
			}
		}
		val := map[string]interface{}{
			"description":     o.AwsVpc.Description,
			"disabled":        o.AwsVpc.Disabled,
			"source":          o.AwsVpc.Source,
			"access_key_id":   o.AwsVpc.AccessKeyId,
			"secret_key_id":   sk,
			"update_interval": o.AwsVpc.UpdateInterval,
			"enable_timeout":  o.AwsVpc.EnableTimeout,
			"timeout":         o.AwsVpc.Timeout,
			"vpc_id":          o.AwsVpc.VpcId,
		}

		if err = d.Set("aws_vpc", []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting 'aws_vpc' for %q: %s", d.Id(), err)
		}
	}

	if o.Esxi == nil {
		d.Set("esxi", nil)
	} else {
		var pwd string
		if settings["esxi_password"] != nil {
			if o.Esxi.Password == settings["esxi_password_enc"].(string) {
				pwd = settings["esxi_password"].(string)
			}
		}
		val := map[string]interface{}{
			"description":     o.Esxi.Description,
			"port":            o.Esxi.Port,
			"disabled":        o.Esxi.Disabled,
			"enable_timeout":  o.Esxi.EnableTimeout,
			"timeout":         o.Esxi.Timeout,
			"source":          o.Esxi.Source,
			"username":        o.Esxi.Username,
			"password":        pwd,
			"update_interval": o.Esxi.UpdateInterval,
		}

		if err = d.Set("esxi", []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting 'esxi' for %q: %s", d.Id(), err)
		}
	}

	if o.Vcenter == nil {
		d.Set("vcenter", nil)
	} else {
		var pwd string
		if settings["vcenter_password"] != nil {
			if o.Vcenter.Password == settings["vcenter_password_enc"].(string) {
				pwd = settings["vcenter_password"].(string)
			}
		}
		val := map[string]interface{}{
			"description":     o.Vcenter.Description,
			"port":            o.Vcenter.Port,
			"disabled":        o.Vcenter.Disabled,
			"enable_timeout":  o.Vcenter.EnableTimeout,
			"timeout":         o.Vcenter.Timeout,
			"source":          o.Vcenter.Source,
			"username":        o.Vcenter.Username,
			"password":        pwd,
			"update_interval": o.Vcenter.UpdateInterval,
		}

		if err = d.Set("vcenter", []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting 'vcenter' for %q: %s", d.Id(), err)
		}
	}

	if o.GoogleCompute == nil {
		d.Set("google_compute", nil)
	} else {
		var creds string
		if settings["gc_creds"] != nil {
			if o.GoogleCompute.ServiceAccountCredential == settings["gc_creds_enc"].(string) {
				creds = settings["gc_creds"].(string)
			}
		}
		val := map[string]interface{}{
			"description":                o.GoogleCompute.Description,
			"disabled":                   o.GoogleCompute.Disabled,
			"auth_type":                  o.GoogleCompute.AuthType,
			"service_account_credential": creds,
			"project_id":                 o.GoogleCompute.ProjectId,
			"zone_name":                  o.GoogleCompute.ZoneName,
			"update_interval":            o.GoogleCompute.UpdateInterval,
			"enable_timeout":             o.GoogleCompute.EnableTimeout,
			"timeout":                    o.GoogleCompute.Timeout,
		}

		if err = d.Set("google_compute", []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting 'google_compute' for %q: %s", d.Id(), err)
		}
	}
}

func saveVmInformationSourceHashes(d *schema.ResourceData, o, enc vis.Entry) {
	ans := make(map[string]interface{})

	if o.AwsVpc != nil && enc.AwsVpc != nil {
		ans["aws_secret_key"] = o.AwsVpc.SecretAccessKey
		ans["aws_secret_key_enc"] = enc.AwsVpc.SecretAccessKey
	}

	if o.Esxi != nil && enc.Esxi != nil {
		ans["esxi_password"] = o.Esxi.Password
		ans["esxi_password_enc"] = enc.Esxi.Password
	}

	if o.Vcenter != nil && enc.Vcenter != nil {
		ans["vcenter_password"] = o.Vcenter.Password
		ans["vcenter_password_enc"] = enc.Vcenter.Password
	}

	if o.GoogleCompute != nil && enc.GoogleCompute != nil {
		ans["gc_creds"] = o.GoogleCompute.ServiceAccountCredential
		ans["gc_creds_enc"] = enc.GoogleCompute.ServiceAccountCredential
	}

	if err := d.Set("settings", ans); err != nil {
		log.Printf("[WARN] Error setting 'settings' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildVmInformationSourceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parseVmInformationSourceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}
