package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/mngtprof"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceManagementProfile() *schema.Resource {
	return &schema.Resource{
		Create: createManagementProfile,
		Read:   readManagementProfile,
		Update: updateManagementProfile,
		Delete: deleteManagementProfile,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ping": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"telnet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssh": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http_ocsp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"https": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"snmp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"response_pages": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userid_service": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userid_syslog_listener_ssl": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userid_syslog_listener_udp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permitted_ips": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parseManagementProfile(d *schema.ResourceData) mngtprof.Entry {
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

	return o
}

func createManagementProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseManagementProfile(d)

	if err := fw.Network.ManagementProfile.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readManagementProfile(d, meta)
}

func readManagementProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.ManagementProfile.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

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

func updateManagementProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseManagementProfile(d)

	lo, err := fw.Network.ManagementProfile.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.ManagementProfile.Edit(lo); err != nil {
		return err
	}

	return readManagementProfile(d, meta)
}

func deleteManagementProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	err := fw.Network.ManagementProfile.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
