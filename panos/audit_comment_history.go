package panos

import (
	"log"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/audit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source.
func dataSourceAuditCommentHistory() *schema.Resource {
	return &schema.Resource{
		Read: readAuditCommentHistory,

		Schema: auditCommentHistorySchema(),
	}
}

func readAuditCommentHistory(d *schema.ResourceData, meta interface{}) error {
	var err error
	var list []audit.Comment

	rType := d.Get("rule_type").(string)
	name := d.Get("name").(string)
	vsys := d.Get("vsys").(string)
	base := d.Get("rulebase").(string)
	dg := d.Get("device_group").(string)
	direction := d.Get("direction").(string)
	nlogs := d.Get("nlogs").(int)
	skip := d.Get("skip").(int)

	id := buildAuditCommentHistoryId(rType, vsys, dg, base, name)
	switch con := meta.(type) {
	case *pango.Firewall:
		switch rType {
		case "security":
			list, err = con.Policies.Security.AuditCommentHistory(vsys, name, direction, nlogs, skip)
		case "pbf":
			list, err = con.Policies.PolicyBasedForwarding.AuditCommentHistory(vsys, name, direction, nlogs, skip)
		case "nat":
			list, err = con.Policies.Nat.AuditCommentHistory(vsys, name, direction, nlogs, skip)
		case "decryption":
			list, err = con.Policies.Decryption.AuditCommentHistory(vsys, name, direction, nlogs, skip)
		}
	case *pango.Panorama:
		switch rType {
		case "security":
			list, err = con.Policies.Security.AuditCommentHistory(dg, base, name, direction, nlogs, skip)
		case "pbf":
			list, err = con.Policies.PolicyBasedForwarding.AuditCommentHistory(dg, base, name, direction, nlogs, skip)
		case "nat":
			list, err = con.Policies.Nat.AuditCommentHistory(dg, base, name, direction, nlogs, skip)
		case "decryption":
			list, err = con.Policies.Decryption.AuditCommentHistory(dg, base, name, direction, nlogs, skip)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("rule_type", rType)
	d.Set("name", name)
	d.Set("vsys", vsys)
	d.Set("rulebase", base)
	d.Set("device_group", dg)
	d.Set("direction", direction)
	d.Set("nlogs", nlogs)
	d.Set("skip", skip)

	info := make([]interface{}, 0, len(list))
	for _, x := range list {
		info = append(info, map[string]interface{}{
			"admin":                  x.Admin,
			"comment":                x.Comment,
			"config_version":         x.ConfigVersion,
			"time_generated":         x.TimeGenerated,
			"time_generated_rfc3339": x.Time.Format(time.RFC3339),
		})
	}

	if err = d.Set("log", info); err != nil {
		log.Printf("[WARN] Error setting 'log' for %q: %s", d.Id(), err)
	}

	return nil
}

// Schema functions.
func auditCommentHistorySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule_type": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "security",
			Description: "The rule type.",
			ValidateFunc: validateStringIn(
				"security",
				"pbf",
				"nat",
				"decryption",
			),
		},
		"vsys":         vsysSchema("vsys1"),
		"device_group": deviceGroupSchema(),
		"rulebase":     rulebaseSchema(),
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The rule name.",
		},
		"direction": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Specify whether logs are shown oldest first (`forward`) or newest first (`backward`).",
			Default:      "backward",
			ValidateFunc: validateStringIn("backward", "forward"),
		},
		"nlogs": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Number of audit comments to return, maximum 5000.",
			Default:      100,
			ValidateFunc: validateIntInRange(1, 5000),
		},
		"skip": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specify the number of logs to skip when doing log retrieval.  This is useful when retrieving logs in batches to skip previously retrieved logs.",
		},

		// Attributes.
		"log": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"admin": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The admin who made the change.",
					},
					"comment": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The audit comment.",
					},
					"config_version": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The config version.",
					},
					"time_generated": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The time generated as reported by PAN-OS.",
					},
					"time_generated_rfc3339": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "An opportunistic representation of time generated in RFC3339. This is created by combining the time_generated with the timezone information of PAN-OS.",
					},
				},
			},
		},
	}
}

// Id functions.
func buildAuditCommentHistoryId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func parseAuditCommentHistoryId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}
