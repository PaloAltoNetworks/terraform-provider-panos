package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/dhcp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schema handling.
func dhcpSchema(isResource bool, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Description: "The template.",
			Optional:    true,
			ForceNew:    true,
		},
		"template_stack": {
			Type:        schema.TypeString,
			Description: "The template stack.",
			Optional:    true,
			ForceNew:    true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"relay": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ipv4_enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "DHCP relay enabled for ipv4",
					},
					"ipv4_servers": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "DHCP servers ipv4 address list",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"ipv6_enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "DHCP relay enabled for ipv4",
					},
					"ipv6_servers": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "DHCP servers ipv6 address list",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"server": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "DHCP server ipv6 server name",
								},
								"interface": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "DHCP server ipv6 interface name",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys"})
	}

	return ans
}

func dhcpUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["template"]; !ok {
		raw["template"] = ""
	}
	if _, ok := raw["template_stack"]; !ok {
		raw["template_stack"] = ""
	}

	return raw, nil
}

func resourceDHCP() *schema.Resource {
	return &schema.Resource{
		Create: createDHCP,
		Read:   readDHCP,
		Update: updateDHCP,
		Delete: deleteDHCP,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: haSchema(true, []string{"template", "template_stack", "vsys"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: dhcpUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: dhcpSchema(true, nil),
	}
}

func resourcePanoramaDHCP() *schema.Resource {
	return &schema.Resource{
		Create: createDHCP,
		Read:   readDHCP,
		Update: updateDHCP,
		Delete: deleteDHCP,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: haSchema(true, []string{}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: dhcpUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: dhcpSchema(true, nil),
	}
}

func parseDHCPId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildDHCPId(a, b, c string) string {
	return fmt.Sprintf("%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c)
}

func parseDHCP(d *schema.ResourceData) dhcp.Entry {
	relay := dhcp.Relay{}
	relayList := d.Get("relay").([]interface{})
	for i := range relayList {
		log.Printf("%v", relayList[i])
		elm := relayList[i].(map[string]interface{})
		relay = dhcp.Relay{
			Ipv4Enabled: elm["ipv4_enable"].(bool),
			Ipv4Servers: setAsList(elm["ipv4_servers"].(*schema.Set)),
			Ipv6Enabled: elm["ipv6_enable"].(bool),
		}
		/* TODO : find a way to get ipv6 servers here
		ipv6Servers := []dhcp.Ipv6Server{relayList[i]["ipv6_servers"]}
		for j := range relayList[i]["ipv6_servers"] {}
		*/
	}

	o := dhcp.Entry{
		Name:  d.Get("name").(string),
		Relay: &relay,
	}

	return o
}

func createDHCP(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := parseDHCP(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("name", name)

	id := buildDHCPId(tmpl, ts, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.Dhcp.Edit(o)
	case *pango.Panorama:
		err = con.Network.Dhcp.Edit(tmpl, ts, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readDHCP(d, meta)
}

func saveDHCP(d *schema.ResourceData, o dhcp.Entry) error {
	d.Set("name", o.Name)

	if o.Relay != nil {
		m := map[string]interface{}{
			"ipv4_enable":  o.Relay.Ipv4Enabled,
			"ipv4_servers": o.Relay.Ipv4Servers,
			"ipv6_enable":  o.Relay.Ipv6Enabled,
			// TODO "ipv6_servers":  o.Relay.Ipv6Servers,
		}

		d.Set("relay", m)
	}

	return nil
}

func readDHCP(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o dhcp.Entry

	tmpl, ts, name := parseDHCPId(d.Id())
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.Dhcp.Get(name)
	case *pango.Panorama:
		o, err = con.Network.Dhcp.Get(tmpl, ts, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	err = saveDHCP(d, o)
	if err != nil {
		return err
	}

	return nil
}

func updateDHCP(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := parseDHCP(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("name", name)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Network.Dhcp.Get(name)
		if err != nil {
			return err
		}
		// Todo: Should implement a merge
		lo.Copy(o)
		if err = con.Network.Dhcp.Edit(lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Network.Dhcp.Get(tmpl, ts, name)
		if err != nil {
			return err
		}
		// Todo: Should implement a merge
		lo.Copy(o)
		if err = con.Network.Dhcp.Edit(tmpl, ts, lo); err != nil {
			return err
		}
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	return readDHCP(d, meta)
}

func deleteDHCP(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := parseDHCP(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("name", name)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Network.Dhcp.Get(name)
		if err != nil {
			return err
		}

		// Todo: Should implement a merge
		lo.Copy(o)
		err = con.Network.Dhcp.Delete()
	case *pango.Panorama:
		lo, err := con.Network.Dhcp.Get(tmpl, ts, name)
		if err != nil {
			return err
		}
		// Todo: Should implement a merge
		lo.Copy(o)
		err = con.Network.Dhcp.Delete(tmpl, ts, name)
	}

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
