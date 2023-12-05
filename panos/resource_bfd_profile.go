package panos

import (
	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/bfd"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBfdProfile() *schema.Resource {
	return &schema.Resource{
		Create: createBfdProfile,
		Read:   readBfdProfile,
		Update: updateBfdProfile,
		Delete: deleteBfdProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bfdProfileSchema(false),
	}
}

func bfdProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"mode": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      bfd.ModeActive,
			ValidateFunc: validateStringIn(bfd.ModeActive, bfd.ModePassive),
		},
		"minimum_tx_interval": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1000,
		},
		"minimum_rx_interval": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1000,
		},
		"detection_multiplier": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  3,
		},
		"hold_time": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"minimum_rx_ttl": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBfdProfile(d *schema.ResourceData) bfd.Entry {
	o := bfd.Entry{
		Name:                d.Get("name").(string),
		Mode:                d.Get("mode").(string),
		MinimumTxInterval:   d.Get("minimum_tx_interval").(int),
		MinimumRxInterval:   d.Get("minimum_rx_interval").(int),
		DetectionMultiplier: d.Get("detection_multiplier").(int),
		HoldTime:            d.Get("hold_time").(int),
		MinimumRxTtl:        d.Get("minimum_rx_ttl").(int),
	}

	return o
}

func saveBfdProfile(d *schema.ResourceData, o bfd.Entry) {
	d.Set("name", o.Name)
	d.Set("mode", o.Mode)
	d.Set("minimum_tx_interval", o.MinimumTxInterval)
	d.Set("minimum_rx_interval", o.MinimumRxInterval)
	d.Set("detection_multiplier", o.DetectionMultiplier)
	d.Set("hold_time", o.HoldTime)
	d.Set("minimum_rx_ttl", o.MinimumRxTtl)
}

func createBfdProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseBfdProfile(d)

	if err := fw.Network.BfdProfile.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readBfdProfile(d, meta)
}

func readBfdProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.BfdProfile.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBfdProfile(d, o)

	return nil
}

func updateBfdProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseBfdProfile(d)

	lo, err := fw.Network.BfdProfile.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BfdProfile.Edit(lo); err != nil {
		return err
	}

	return readBfdProfile(d, meta)
}

func deleteBfdProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Id()

	err := fw.Network.BfdProfile.Delete(name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
