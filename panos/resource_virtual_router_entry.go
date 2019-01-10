package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualRouterEntry() *schema.Resource {
	return &schema.Resource{
		Create: createVirtualRouterEntry,
		Read:   readVirtualRouterEntry,
		Delete: deleteVirtualRouterEntry,

		Schema: map[string]*schema.Schema{
			"virtual_router": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func parseVirtualRouterEntryId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildVirtualRouterEntryId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr := d.Get("virtual_router").(string)
	iface := d.Get("interface").(string)

	if err := fw.Network.VirtualRouter.SetInterface(vr, iface); err != nil {
		return err
	}

	d.SetId(buildVirtualRouterEntryId(vr, iface))
	return readVirtualRouterEntry(d, meta)
}

func readVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, iface := parseVirtualRouterEntryId(d.Id())

	// Two possibilities:  either the router itself doesn't exist or the
	// interface isn't present.
	o, err := fw.Network.VirtualRouter.Get(vr)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	for i := range o.Interfaces {
		if o.Interfaces[i] == iface {
			d.Set("virtual_router", vr)
			d.Set("interface", iface)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, iface := parseVirtualRouterEntryId(d.Id())

	if err := fw.Network.VirtualRouter.DeleteInterface(vr, iface); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
