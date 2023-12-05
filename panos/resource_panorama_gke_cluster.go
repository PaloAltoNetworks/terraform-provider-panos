package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/gke/cluster"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaGkeCluster() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaGkeCluster,
		Read:   readPanoramaGkeCluster,
		Update: updatePanoramaGkeCluster,
		Delete: deletePanoramaGkeCluster,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gke_cluster_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gcp_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_credential": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func parseGkeClusterId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildGkeClusterId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parsePanoramaGkeCluster(d *schema.ResourceData) (string, cluster.Entry) {
	grp := d.Get("gke_cluster_group").(string)

	o := cluster.Entry{
		Name:              d.Get("name").(string),
		GcpZone:           d.Get("gcp_zone").(string),
		ClusterCredential: d.Get("cluster_credential").(string),
	}

	return grp, o
}

func createPanoramaGkeCluster(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	grp, o := parsePanoramaGkeCluster(d)

	if err := pano.Panorama.GkeCluster.Set(grp, o); err != nil {
		return err
	}

	d.SetId(buildGkeClusterId(grp, o.Name))
	return readPanoramaGkeCluster(d, meta)
}

func readPanoramaGkeCluster(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	grp, name := parseGkeClusterId(d.Id())

	o, err := pano.Panorama.GkeCluster.Get(grp, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("gke_cluster_group", grp)
	d.Set("name", o.Name)
	d.Set("gcp_zone", o.GcpZone)
	d.Set("cluster_credential", o.ClusterCredential)

	return nil
}

func updatePanoramaGkeCluster(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	grp, o := parsePanoramaGkeCluster(d)

	lo, err := pano.Panorama.GkeCluster.Get(grp, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.GkeCluster.Edit(grp, lo); err != nil {
		return err
	}

	return readPanoramaGkeCluster(d, meta)
}

func deletePanoramaGkeCluster(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	grp, name := parseGkeClusterId(d.Id())

	err := pano.Panorama.GkeCluster.Delete(grp, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
