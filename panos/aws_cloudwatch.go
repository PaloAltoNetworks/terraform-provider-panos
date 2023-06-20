package panos

import (
	"github.com/fpluchorg/pango/panosplugin/cloudwatch"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Resource.
func resourceAwsCloudWatch() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateAwsCloudWatch,
		Read:   readAwsCloudWatch,
		Update: createUpdateAwsCloudWatch,
		Delete: deleteAwsCloudWatch,

		Schema: awsCloudWatchSchema(),
	}
}

func createUpdateAwsCloudWatch(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	o := loadAwsCloudWatch(d)

	if err = fw.PanosPlugin.AwsCloudWatch.Edit(o); err != nil {
		return err
	}

	d.SetId(fw.Hostname)
	return readAwsCloudWatch(d, meta)
}

func readAwsCloudWatch(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}

	o, err := fw.PanosPlugin.AwsCloudWatch.Get()
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAwsCloudWatch(d, o)

	return nil
}

func deleteAwsCloudWatch(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}

	err = fw.PanosPlugin.AwsCloudWatch.Delete()
	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func awsCloudWatchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enable AWS CloudWatch setup.",
			Optional:    true,
			Default:     true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "Namespace.",
			Optional:    true,
			Default:     "VMseries",
		},
		"update_interval": {
			Type:        schema.TypeInt,
			Description: "Update time (in min).",
			Optional:    true,
			Default:     5,
		},
	}
}

func loadAwsCloudWatch(d *schema.ResourceData) cloudwatch.Config {
	return cloudwatch.Config{
		Enabled:        d.Get("enabled").(bool),
		Namespace:      d.Get("namespace").(string),
		UpdateInterval: d.Get("update_interval").(int),
	}
}

func saveAwsCloudWatch(d *schema.ResourceData, o cloudwatch.Config) {
	d.Set("enabled", o.Enabled)
	d.Set("namespace", o.Namespace)
	d.Set("update_interval", o.UpdateInterval)
}
