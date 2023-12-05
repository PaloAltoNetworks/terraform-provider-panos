package panos

import (
	"encoding/base64"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/account"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaGcpAccount() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaGcpAccount,
		Read:   readPanoramaGcpAccount,
		Update: updatePanoramaGcpAccount,
		Delete: deletePanoramaGcpAccount,

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
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_account_credential_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      account.Project,
				ValidateFunc: validateStringIn(account.Project, account.Gke),
			},
			"credential_file": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"credential_file_enc": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func parsePanoramaGcpAccount(d *schema.ResourceData) account.Entry {
	o := account.Entry{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		ProjectId:                    d.Get("project_id").(string),
		ServiceAccountCredentialType: d.Get("service_account_credential_type").(string),
		CredentialFile:               base64.StdEncoding.EncodeToString([]byte(d.Get("credential_file").(string))),
	}

	return o
}

func createPanoramaGcpAccount(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	o := parsePanoramaGcpAccount(d)

	if err := pano.Panorama.GcpAccount.Set(o); err != nil {
		return err
	}
	lo, err := pano.Panorama.GcpAccount.Get(o.Name)
	if err != nil {
		return err
	}

	d.SetId(o.Name)
	d.Set("credential_file_enc", lo.CredentialFile)

	return readPanoramaGcpAccount(d, meta)
}

func readPanoramaGcpAccount(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	o, err := pano.Panorama.GcpAccount.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("project_id", o.ProjectId)
	d.Set("service_account_credential_type", o.ServiceAccountCredentialType)
	if d.Get("credential_file_enc").(string) != o.CredentialFile {
		d.Set("credential_file", "(incorrect credentials)")
	}

	return nil
}

func updatePanoramaGcpAccount(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaGcpAccount(d)

	lo, err := pano.Panorama.GcpAccount.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.GcpAccount.Edit(lo); err != nil {
		return err
	}

	eo, err := pano.Panorama.GcpAccount.Get(o.Name)
	if err != nil {
		return err
	}
	d.Set("credential_file_enc", eo.CredentialFile)

	return readPanoramaGcpAccount(d, meta)
}

func deletePanoramaGcpAccount(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	name := d.Id()

	err := pano.Panorama.GcpAccount.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
