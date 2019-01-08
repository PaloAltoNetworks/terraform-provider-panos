package panos

import (
	"fmt"
	"github.com/PaloAltoNetworks/pango"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDagTags() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateDagTags,
		Read:   readDagTags,
		Update: createUpdateDagTags,
		Delete: deleteDagTags,

		Schema: map[string]*schema.Schema{
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				Description: "The vsys to config DAG tags for",
			},
			"register": {
				Type:     schema.TypeSet,
				Required: true,
				// TODO(gfreeman): Uncomment once ValidateFunc is supported for TypeSet.
				//ValidateFunc: validateSetKeyIsUnique("ip"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tags": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func parseDagTags(cur map[string][]string, d *schema.ResourceData) (*schema.Set, map[string][]string, map[string][]string, *schema.Set, error) {
	dag := d.Get("register").(*schema.Set)
	missingMap := make(map[string][]string)
	overlapMap := make(map[string][]string)
	overlapSet := &schema.Set{F: dag.F}

	osl := dag.List()
	for i := range osl {
		group := osl[i].(map[string]interface{})
		key := group["ip"].(string)
		if _, ok := missingMap[key]; ok {
			return nil, nil, nil, nil, fmt.Errorf("IP %q already defined, please merge these groups", key)
		} else if _, ok := overlapMap[key]; ok {
			return nil, nil, nil, nil, fmt.Errorf("IP %q already defined, please merge these groups", key)
		}
		info := cur[key]
		tl := group["tags"].(*schema.Set).List()
		otags := make([]string, 0, len(tl))
		mtags := make([]string, 0, len(tl))
		for j := range tl {
			tag := tl[j].(string)
			found := false
			for _, v := range info {
				if v == tag {
					found = true
					break
				}
			}
			if found {
				otags = append(otags, tag)
			} else {
				mtags = append(mtags, tag)
			}
		}
		if len(otags) > 0 {
			ogroup := make(map[string]interface{})
			ogroup["ip"] = key
			ogroup["tags"] = listAsSet(otags)
			overlapSet.Add(ogroup)
			overlapMap[key] = otags
		}
		if len(mtags) > 0 {
			missingMap[key] = mtags
		}
	}

	return dag, missingMap, overlapMap, overlapSet, nil
}

func createUpdateDagTags(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Get("vsys").(string)

	cur, err := fw.UserId.Registered("", "", vsys)
	if err != nil {
		return err
	}

	_, missingMap, _, _, err := parseDagTags(cur, d)
	if err != nil {
		return err
	}

	if err = fw.UserId.Run(nil, nil, missingMap, nil, vsys); err != nil {
		return err
	}

	d.SetId(vsys)
	return readDagTags(d, meta)
}

func readDagTags(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Get("vsys").(string)

	cur, err := fw.UserId.Registered("", "", vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}

	_, _, _, overlapSet, err := parseDagTags(cur, d)
	if err != nil {
		return err
	}

	d.Set("vsys", vsys)
	if err := d.Set("register", overlapSet); err != nil {
		log.Printf("[WARN] Error setting 'register' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteDagTags(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Get("vsys").(string)

	cur, err := fw.UserId.Registered("", "", vsys)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, _, overlapMap, _, err := parseDagTags(cur, d)
	if err != nil {
		return err
	}

	// The UserId subsystem doesn't return ObjectNotFound, so we don't need
	// to check for that at this point.
	err = fw.UserId.Run(nil, nil, nil, overlapMap, vsys)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
