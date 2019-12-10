package panos

import (
	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/gke/cluster/group"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaGkeClusterGroup() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaGkeClusterGroup,
		Read:   readPanoramaGkeClusterGroup,
		Update: updatePanoramaGkeClusterGroup,
		Delete: deletePanoramaGkeClusterGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gcp_project_credential": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_stack": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func parsePanoramaGkeClusterGroup(d *schema.ResourceData) group.Entry {
	o := group.Entry{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		GcpProjectCredential: d.Get("gcp_project_credential").(string),
		DeviceGroup:          d.Get("device_group").(string),
		TemplateStack:        d.Get("template_stack").(string),
	}

	return o
}

func createPanoramaGkeClusterGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	o := parsePanoramaGkeClusterGroup(d)

	if err := pano.Panorama.GkeClusterGroup.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readPanoramaGkeClusterGroup(d, meta)
}

func readPanoramaGkeClusterGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	o, err := pano.Panorama.GkeClusterGroup.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("gcp_project_credential", o.GcpProjectCredential)
	d.Set("device_group", o.DeviceGroup)
	d.Set("template_stack", o.TemplateStack)

	return nil
}

func updatePanoramaGkeClusterGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaGkeClusterGroup(d)

	lo, err := pano.Panorama.GkeClusterGroup.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.GkeClusterGroup.Edit(lo); err != nil {
		return err
	}

	return readPanoramaGkeClusterGroup(d, meta)
}

func deletePanoramaGkeClusterGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	name := d.Id()

	err := pano.Panorama.GkeClusterGroup.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
