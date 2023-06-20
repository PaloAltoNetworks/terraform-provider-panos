package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/interface/loopback"
	"github.com/fpluchorg/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLoopbackInterface() *schema.Resource {
	return &schema.Resource{
		Create: createLoopbackInterface,
		Read:   readLoopbackInterface,
		Update: updateLoopbackInterface,
		Delete: deleteLoopbackInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringHasPrefix("loopback."),
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
			"adjust_tcp_mss": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv4_mss_adjust": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv6_mss_adjust": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseLoopbackInterfaceId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildLoopbackInterfaceId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func parseLoopbackInterface(d *schema.ResourceData) (string, loopback.Entry) {
	vsys := d.Get("vsys").(string)
	o := loopback.Entry{
		Name:              d.Get("name").(string),
		Comment:           d.Get("comment").(string),
		NetflowProfile:    d.Get("netflow_profile").(string),
		StaticIps:         asStringList(d.Get("static_ips").([]interface{})),
		ManagementProfile: d.Get("management_profile").(string),
		Mtu:               d.Get("mtu").(int),
		AdjustTcpMss:      d.Get("adjust_tcp_mss").(bool),
		Ipv4MssAdjust:     d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:     d.Get("ipv6_mss_adjust").(int),
	}

	return vsys, o
}

func createLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseLoopbackInterface(d)

	if err := fw.Network.LoopbackInterface.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildLoopbackInterfaceId(vsys, o.Name))
	return readLoopbackInterface(d, meta)
}

func readLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseLoopbackInterfaceId(d.Id())

	o, err := fw.Network.LoopbackInterface.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
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
	d.Set("adjust_tcp_mss", o.AdjustTcpMss)
	d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
	d.Set("ipv6_mss_adjust", o.Ipv6MssAdjust)

	return nil
}

func updateLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseLoopbackInterface(d)

	lo, err := fw.Network.LoopbackInterface.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.LoopbackInterface.Edit(vsys, lo); err != nil {
		return err
	}

	return readLoopbackInterface(d, meta)
}

func deleteLoopbackInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	_, name := parseLoopbackInterfaceId(d.Id())

	err := fw.Network.LoopbackInterface.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
