package panos

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func templateSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"template_stack"},
	}
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
	m := getMovementMap()
	s := make([]string, len(m))
	for _, v := range m {
		s = append(s, v)
	}

	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "",
		ValidateFunc: validateStringIn(s...),
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
