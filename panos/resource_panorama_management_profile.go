package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/mngtprof"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaManagementProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaManagementProfile,
		Read:   readPanoramaManagementProfile,
		Update: updatePanoramaManagementProfile,
		Delete: deletePanoramaManagementProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template_stack"},
			},
			"template_stack": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template"},
			},
			"ping": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"telnet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssh": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http_ocsp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"snmp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"response_pages": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userid_service": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userid_syslog_listener_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userid_syslog_listener_udp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permitted_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parsePanoramaManagementProfile(d *schema.ResourceData) (string, string, mngtprof.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	o := mngtprof.Entry{
		Name:                    d.Get("name").(string),
		Ping:                    d.Get("ping").(bool),
		Telnet:                  d.Get("telnet").(bool),
		Ssh:                     d.Get("ssh").(bool),
		Http:                    d.Get("http").(bool),
		HttpOcsp:                d.Get("http_ocsp").(bool),
		Https:                   d.Get("https").(bool),
		Snmp:                    d.Get("snmp").(bool),
		ResponsePages:           d.Get("response_pages").(bool),
		UseridService:           d.Get("userid_service").(bool),
		UseridSyslogListenerSsl: d.Get("userid_syslog_listener_ssl").(bool),
		UseridSyslogListenerUdp: d.Get("userid_syslog_listener_udp").(bool),
		PermittedIps:            asStringList(d.Get("permitted_ips").([]interface{})),
	}

	return tmpl, ts, o
}

func parsePanoramaManagementProfileId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaManagementProfileId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createPanoramaManagementProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaManagementProfile(d)

	if err := pano.Network.ManagementProfile.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaManagementProfileId(tmpl, ts, o.Name))
	return readPanoramaManagementProfile(d, meta)
}

func readPanoramaManagementProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaManagementProfileId(d.Id())

	o, err := pano.Network.ManagementProfile.Get(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("ping", o.Ping)
	d.Set("telnet", o.Telnet)
	d.Set("ssh", o.Ssh)
	d.Set("http", o.Http)
	d.Set("http_ocsp", o.HttpOcsp)
	d.Set("https", o.Https)
	d.Set("snmp", o.Snmp)
	d.Set("response_pages", o.ResponsePages)
	d.Set("userid_service", o.UseridService)
	d.Set("userid_syslog_listener_ssl", o.UseridSyslogListenerSsl)
	d.Set("userid_syslog_listener_udp", o.UseridSyslogListenerUdp)
	if err := d.Set("permitted_ips", o.PermittedIps); err != nil {
		log.Printf("[WARN] Error setting 'permitted_ips' for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaManagementProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaManagementProfile(d)

	lo, err := pano.Network.ManagementProfile.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.ManagementProfile.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaManagementProfile(d, meta)
}

func deletePanoramaManagementProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaManagementProfileId(d.Id())

	err := pano.Network.ManagementProfile.Delete(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
