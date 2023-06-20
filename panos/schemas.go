package panos

import (
	"log"

	"github.com/fpluchorg/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func templateSchema(ts bool) *schema.Schema {
	ans := &schema.Schema{
		Type:        schema.TypeString,
		Description: "The template.",
		ForceNew:    true,
	}

	if ts {
		ans.Optional = true
		ans.ConflictsWith = []string{"template_stack"}
	} else {
		ans.Required = true
	}

	return ans
}

func templateStackSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Description:   "The template stack.",
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"template"},
	}
}

func deviceGroupSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The device group.",
		ForceNew:    true,
		Default:     "shared",
	}
}

func positionKeywordSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The position keyword for this group of rules",
		ValidateFunc: validateStringIn(movementKeywords()...),
	}
}

func positionReferenceSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The position reference for this group of rules",
	}
}

func rulebaseSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		ForceNew:    true,
		Optional:    true,
		Description: "The rulebase location.",
		Default:     util.PreRulebase,
		ValidateFunc: validateStringIn(
			util.PreRulebase,
			util.Rulebase,
			util.PostRulebase,
		),
	}
}

func templateWithPanoramaSharedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vsys": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
			Default:  "shared",
		},
		"template": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
		},
		"template_stack": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
		},
	}
}

func vsysSchema(v string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "The vsys this object belongs in.",
		Default:     v,
	}
}

func targetSchema(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "NGFW serial numbers and vsys spec.",
		Optional:    true,
		Computed:    computed,
		// TODO(gfreeman): Uncomment once ValidateFunc is supported for TypeSet.
		//ValidateFunc: validateSetKeyIsUnique("serial"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serial": {
					Type:        schema.TypeString,
					Description: "The NGFW serial number.",
					Required:    true,
				},
				"vsys_list": {
					Type:        schema.TypeSet,
					Description: "List of vsys; leave this unspecified if the NGFW is a VM.",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func negateTargetSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}
}

func tagSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "The administrative tags.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func listingSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"total": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of objects present",
		},
		"listing": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Object names",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func saveListing(d *schema.ResourceData, v []string) {
	d.Set("total", len(v))
	if err := d.Set("listing", v); err != nil {
		log.Printf("[WARN] Error setting 'listing' for %q: %s", d.Id(), err)
	}
}

func auditCommentSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The audit comment.",
		Optional:    true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return true
		},
	}
}

func groupTagSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "(PAN-OS 9.0+) The group tag.",
		Optional:    true,
	}
}

func uuidSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "(PAN-OS 9.0+) The PAN-OS UUID.",
		Computed:    true,
	}
}
