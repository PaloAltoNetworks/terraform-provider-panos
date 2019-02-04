package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaVirtualRouterEntry() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaVirtualRouterEntry,
		Read:   readPanoramaVirtualRouterEntry,
		Delete: deletePanoramaVirtualRouterEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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

func parsePanoramaVirtualRouterEntryId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaVirtualRouterEntryId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl := d.Get("template").(string)
	ts := ""
	vr := d.Get("virtual_router").(string)
	iface := d.Get("interface").(string)

	if err := pano.Network.VirtualRouter.SetInterface(tmpl, ts, vr, iface); err != nil {
		return err
	}

	d.SetId(buildPanoramaVirtualRouterEntryId(tmpl, ts, vr, iface))
	return readPanoramaVirtualRouterEntry(d, meta)
}

func readPanoramaVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, iface := parsePanoramaVirtualRouterEntryId(d.Id())

	// Two possibilities:  either the router itself doesn't exist or the
	// interface isn't present.
	o, err := pano.Network.VirtualRouter.Get(tmpl, ts, vr)
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
			d.Set("template", tmpl)
			//d.Set("template_stack", ts)
			d.Set("virtual_router", vr)
			d.Set("interface", iface)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deletePanoramaVirtualRouterEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, iface := parsePanoramaVirtualRouterEntryId(d.Id())

	if err := pano.Network.VirtualRouter.DeleteInterface(tmpl, ts, vr, iface); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
