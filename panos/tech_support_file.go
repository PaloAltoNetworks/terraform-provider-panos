package panos

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source.
func dataSourceTechSupportFile() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceTechSupportFile,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},

			// Local save variables.
			"save_to_file_system": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_system_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// State save variables.
			"save_to_state": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readDataSourceTechSupportFile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var filename string
	var data []byte
	var id string

	timeout := time.Duration(d.Get("timeout").(int)) * time.Second
	localSave := d.Get("save_to_file_system").(bool)
	p := d.Get("file_system_path").(string)
	stateSave := d.Get("save_to_state").(bool)

	d.Set("timeout", timeout)
	d.Set("save_to_file_system", localSave)
	d.Set("file_system_path", p)
	d.Set("save_to_state", stateSave)

	if !localSave && !stateSave {
		return fmt.Errorf("At least one of 'save_to_file_system' or 'save_to_state' must be enabled.")
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		id = con.Hostname
		filename, data, err = con.GetTechSupportFile(timeout)
	case *pango.Panorama:
		id = con.Hostname
		filename, data, err = con.GetTechSupportFile(timeout)
	}

	if err != nil {
		return err
	}

	d.Set("filename", filename)

	if localSave {
		var path string
		if p == "" {
			path = filename
		} else {
			path = filepath.Join(p, filename)
		}

		if err = ioutil.WriteFile(path, data, 0644); err != nil {
			return err
		}
	}

	if stateSave {
		d.Set("content", string(data))
	} else {
		d.Set("content", "")
	}

	d.SetId(id)
	return nil
}
