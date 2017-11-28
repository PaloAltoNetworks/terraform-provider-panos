package panos

import (
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
			"permitted_ip": &schema.Schema{
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
		PermittedIp:             asStringList(d, "permitted_ip"),
	}

	return o
}

func saveDataManagementProfile(d *schema.ResourceData, o mngtprof.Entry) {
	d.SetId(o.Name)
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
	d.Set("permitted_ip", o.PermittedIp)
}

func createManagementProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseManagementProfile(d)

	if err := fw.Network.ManagementProfile.Set(o); err != nil {
		return err
	}

	saveDataManagementProfile(d, o)
	return nil
}

func readManagementProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.ManagementProfile.Get(name)
	if err != nil {
		d.SetId("")
		return nil
	}

	saveDataManagementProfile(d, o)
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
	err = fw.Network.ManagementProfile.Edit(lo)
	/*
	   if err == nil {
	       lo.Copy(o)
	       err = fw.Network.ManagementProfile.Edit(lo)
	   } else {
	       err = fw.Network.ManagementProfile.Set(o)
	   }
	*/

	if err == nil {
		saveDataManagementProfile(d, o)
	}

	return err
}

func deleteManagementProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	_ = fw.Network.ManagementProfile.Delete(name)
	d.SetId("")
	return nil
}
