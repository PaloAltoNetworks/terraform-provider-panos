package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/tunnel"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTunnelInterface() *schema.Resource {
	return &schema.Resource{
		Create: createTunnelInterface,
		Read:   readTunnelInterface,
		Update: updateTunnelInterface,
		Delete: deleteTunnelInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringHasPrefix("tunnel."),
			},
			"vsys": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"netflow_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"static_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of static IP addresses",
			},
			"management_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseTunnelInterfaceId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildTunnelInterfaceId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func parseTunnelInterface(d *schema.ResourceData) (string, tunnel.Entry) {
	vsys := d.Get("vsys").(string)
	o := tunnel.Entry{
		Name:              d.Get("name").(string),
		Comment:           d.Get("comment").(string),
		NetflowProfile:    d.Get("netflow_profile").(string),
		StaticIps:         asStringList(d.Get("static_ips").([]interface{})),
		ManagementProfile: d.Get("management_profile").(string),
		Mtu:               d.Get("mtu").(int),
	}

	return vsys, o
}

func createTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseTunnelInterface(d)

	if err := fw.Network.TunnelInterface.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildTunnelInterfaceId(vsys, o.Name))
	return readTunnelInterface(d, meta)
}

func readTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseTunnelInterfaceId(d.Id())

	o, err := fw.Network.TunnelInterface.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := fw.IsImported(util.InterfaceImport, "", "", vsys, name)
	if err != nil {
		return err
	}

	d.Set("name", o.Name)
	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	d.Set("comment", o.Comment)
	d.Set("netflow_profile", o.NetflowProfile)
	if err = d.Set("static_ips", o.StaticIps); err != nil {
		log.Printf("[WARN] Error setting 'static_ips' for %q: %s", d.Id(), err)
	}
	d.Set("management_profile", o.ManagementProfile)
	d.Set("mtu", o.Mtu)

	return nil
}

func updateTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseTunnelInterface(d)

	lo, err := fw.Network.TunnelInterface.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.TunnelInterface.Edit(vsys, lo); err != nil {
		return err
	}

	return readTunnelInterface(d, meta)
}

func deleteTunnelInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	_, name := parseTunnelInterfaceId(d.Id())

	err := fw.Network.TunnelInterface.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
