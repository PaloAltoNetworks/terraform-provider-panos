package panos

import (
	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDhcpInterfaceInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDhcpInterfaceInfoRead,

		Schema: map[string]*schema.Schema{
			"interface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_wins": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_wins": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_nis": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_nis": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ntp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_ntp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pop3_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smtp_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDhcpInterfaceInfoRead(d *schema.ResourceData, meta interface{}) error {
	var err error

	c := meta.(*pango.Firewall)
	i := d.Get("interface").(string)

	info, err := c.GetDhcpInfo(i)
	if err != nil {
		return err
	}

	d.SetId(i)
	d.Set("interface", i)

	keys := []string{
		"state",
		"ip",
		"gateway",
		"server",
		"server_id",
		"primary_dns",
		"secondary_dns",
		"primary_wins",
		"secondary_wins",
		"primary_nis",
		"secondary_nis",
		"primary_ntp",
		"secondary_ntp",
		"pop3_server",
		"smtp_server",
		"dns_suffix",
	}

	for _, k := range keys {
		d.Set(k, info[k])
	}

	return nil
}
