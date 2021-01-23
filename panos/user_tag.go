package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango/userid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source.
func dataSourceUserTag() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceUserTag,

		Schema: map[string]*schema.Schema{
			// Input.
			"vsys": vsysSchema(),
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter on just this user",
			},

			// Output.
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user specs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user",
						},
						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Tags",
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

func readDataSourceUserTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys := d.Get("vsys").(string)
	su := d.Get("user").(string)

	cur, err := fw.UserId.GetUserTags(su, vsys)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(base64Encode([]interface{}{
		vsys, su,
	}))
	if len(cur) == 0 {
		d.Set("users", nil)
		return nil
	}

	data := make([]interface{}, 0, len(cur))
	for user, tags := range cur {
		data = append(data, map[string]interface{}{
			"user": user,
			"tags": listAsSet(tags),
		})
	}

	if err = d.Set("users", data); err != nil {
		log.Printf("[WARN] Error setting 'users' for %q: %s", d.Id(), err)
	}

	return nil
}

// Resource.
func resourceUserTag() *schema.Resource {
	return &schema.Resource{
		Create: createUserTag,
		Read:   readUserTag,
		Delete: deleteUserTag,

		Schema: map[string]*schema.Schema{
			"vsys": vsysSchema(),
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address to tag",
				ForceNew:    true,
			},
			"tags": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Tags",
				MinItems:    1,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func createUserTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys := d.Get("vsys").(string)
	user := d.Get("user").(string)
	tagList := d.Get("tags").(*schema.Set).List()

	cur, err := fw.UserId.GetUserTags(user, vsys)
	if err != nil {
		return err
	}
	curTags := cur[user]

	missing := make([]userid.UserTag, 0, len(tagList))
	for i := range tagList {
		tag := tagList[i].(string)
		var found bool
		for _, x := range curTags {
			if x == tag {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, userid.UserTag{
				Tag: tag,
			})
		}
	}

	if len(missing) > 0 {
		msg := &userid.Message{
			TagUsers: []userid.TagUser{
				userid.TagUser{
					User: user,
					Tags: missing,
				},
			},
		}

		if err = fw.UserId.Run(msg, vsys); err != nil {
			return err
		}
	}

	d.SetId(buildUserTagId(vsys, user, tagList))
	return readUserTag(d, meta)
}

func readUserTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys, user, tags := parseUserTagId(d.Id())

	cur, err := fw.UserId.GetUserTags(user, vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}
	curTags := cur[user]

	overlap := make([]string, 0, len(curTags))
	for _, curTag := range curTags {
		for _, wantTag := range tags {
			if curTag == wantTag {
				overlap = append(overlap, wantTag)
				break
			}
		}
	}

	d.Set("vsys", vsys)
	d.Set("user", user)
	if len(overlap) != 0 {
		if err := d.Set("tags", listAsSet(overlap)); err != nil {
			log.Printf("[WARN] Error setting 'tags' for %q: %s", d.Id(), err)
		}
	} else {
		d.Set("tags", nil)
	}

	return nil
}

func deleteUserTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys, user, tags := parseUserTagId(d.Id())

	cur, err := fw.UserId.GetUserTags(user, vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}
	curTags := cur[user]

	overlap := make([]string, 0, len(curTags))
	for _, curTag := range curTags {
		for _, wantTag := range tags {
			if curTag == wantTag {
				overlap = append(overlap, wantTag)
				break
			}
		}
	}

	if len(overlap) == 0 {
		d.SetId("")
		return nil
	}

	msg := &userid.Message{
		UntagUsers: []userid.UntagUser{
			userid.UntagUser{
				User: user,
				Tags: overlap,
			},
		},
	}

	// The UserId subsystem doesn't return ObjectNotFound, so we don't need
	// to check for that at this point.
	err = fw.UserId.Run(msg, vsys)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// Id functions.
func buildUserTagId(a, b string, c []interface{}) string {
	return strings.Join([]string{a, b, base64Encode(c)}, IdSeparator)
}

func parseUserTagId(v string) (string, string, []string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], base64Decode(t[2])
}
