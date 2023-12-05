package panos

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source.
func dataSourceVmAuthKey() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceVmAuthKey,

		Schema: map[string]*schema.Schema{
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of entries",
			},
			"entries": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of vm auth key structs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VM auth key",
						},
						"expiry": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiry time as a string",
						},
						"valid": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If the VM auth key is still valid or if it's expired",
						},
					},
				},
			},
		},
	}
}

func readDataSourceVmAuthKey(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	keys, err := pano.GetVmAuthKeys()
	if err != nil {
		return err
	}

	now, err := pano.Clock()
	if err != nil {
		return err
	}

	d.SetId(pano.Hostname)
	d.Set("total", len(keys))
	if len(keys) == 0 {
		d.Set("entries", nil)
	} else {
		var empty time.Time
		list := make([]interface{}, 0, len(keys))
		for _, key := range keys {
			var valid bool
			if key.Expires != empty && key.Expires.After(now) {
				valid = true
			}
			list = append(list, map[string]interface{}{
				"auth_key": key.AuthKey,
				"expiry":   key.Expiry,
				"valid":    valid,
			})
		}

		if err := d.Set("entries", list); err != nil {
			log.Printf("[WARN] Error setting 'entries' for %q: %s", d.Id(), err)
		}
	}

	return nil
}

// Resource.
func resourceVmAuthKey() *schema.Resource {
	return &schema.Resource{
		Create: createVmAuthKey,
		Read:   readVmAuthKey,
		Delete: deleteVmAuthKey,

		Schema: map[string]*schema.Schema{
			"hours": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     8,
				Description: "The VM auth key lifetime",
			},
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of resource.",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
			},
			"auth_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VM auth key",
			},
			"expiry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiry time as a string",
			},
			"valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the VM auth key is still valid or if it's expired",
			},
		},
	}
}

func createVmAuthKey(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	key, err := pano.CreateVmAuthKey(d.Get("hours").(int))
	if err != nil {
		return err
	}

	d.SetId(key.AuthKey)
	return readVmAuthKey(d, meta)
}

func readVmAuthKey(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	authKey := d.Id()

	now, err := pano.Clock()
	if err != nil {
		return err
	}

	keys, err := pano.GetVmAuthKeys()
	if err != nil {
		return err
	}

	var found bool
	for _, key := range keys {
		if key.AuthKey == authKey {
			var empty time.Time

			d.Set("auth_key", key.AuthKey)
			d.Set("expiry", key.Expiry)
			d.Set("valid", key.Expires != empty && key.Expires.After(now))
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func deleteVmAuthKey(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	authKey := d.Id()
	keys, err := pano.GetVmAuthKeys()
	if err != nil {
		return err
	}

	var found bool
	for _, key := range keys {
		if key.AuthKey == authKey {
			found = true
			break
		}
	}

	if found {
		if err = pano.RevokeVmAuthKey(authKey); err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
