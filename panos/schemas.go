package panos

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func templateSchema(ts bool) *schema.Schema {
	ans := &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,
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
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"template"},
	}
}

func deviceGroupSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		Default:  "shared",
	}
}

func positionKeywordSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "",
		ValidateFunc: validateStringIn(movementKeywords()...),
		ForceNew:     true,
	}
}

func positionReferenceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
}

func rulebaseSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,
		Optional: true,
		Default:  util.PreRulebase,
		ValidateFunc: validateStringIn(
			util.PreRulebase,
			util.Rulebase,
			util.PostRulebase,
		),
	}
}

func vsysSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,
		Optional: true,
		Default:  "vsys1",
	}
}

func targetSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Set:      resourceTargetHash,
		// TODO(gfreeman): Uncomment once ValidateFunc is supported for TypeSet.
		//ValidateFunc: validateSetKeyIsUnique("serial"),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serial": {
					Type:     schema.TypeString,
					Required: true,
				},
				"vsys_list": {
					Type:     schema.TypeSet,
					Optional: true,
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
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}
