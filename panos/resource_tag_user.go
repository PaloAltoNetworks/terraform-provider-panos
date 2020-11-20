package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/userid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTagUser() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateTagUser,
		Read:   readTagUser,
		Update: createUpdateTagUser,
		Delete: deleteTagUser,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vsys": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "vsys1",
				Description: "The vsys to config DAG tags for",
			},
			"user": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "IP address to tag",
			},
			"tags": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Tags",
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func buildTagUserId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseTagUserId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parseTagUser(user string, cur map[string][]string, d *schema.ResourceData) *userid.Message {
	tagList := d.Get("tags").(*schema.Set).List()
	curTags := cur[user]

	missing := make([]userid.UserTag, 0, len(tagList))
	extras := make([]string, 0, len(curTags))

	// Loop over what the user wants the user tagged as.
	for i := range tagList {
		tag := tagList[i].(string)
		found := false
		for _, v := range curTags {
			if v == tag {
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

	// Loop over what the user is actually tagged as right now.
	for _, curTag := range curTags {
		found := false
		for i := range tagList {
			tag := tagList[i].(string)
			if curTag == tag {
				found = true
				break
			}
		}

		if !found {
			extras = append(extras, curTag)
		}
	}

	msg := &userid.Message{}

	if len(missing) > 0 {
		msg.TagUsers = []userid.TagUser{
			userid.TagUser{
				User: user,
				Tags: missing,
			},
		}
	}

	if len(extras) > 0 {
		msg.UntagUsers = []userid.UntagUser{
			userid.UntagUser{
				User: user,
				Tags: extras,
			},
		}
	}

	return msg
}

func createUpdateTagUser(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys := d.Get("vsys").(string)
	user := d.Get("user").(string)

	cur, err := fw.UserId.GetUserTags(user, vsys)
	if err != nil {
		return err
	}

	msg := parseTagUser(user, cur, d)

	if err = fw.UserId.Run(msg, vsys); err != nil {
		return err
	}

	d.SetId(buildTagUserId(vsys, user))
	return readTagUser(d, meta)
}

func readTagUser(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys, user := parseTagUserId(d.Id())

	cur, err := fw.UserId.GetUserTags(user, vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("vsys", vsys)
	d.Set("user", user)
	if err := d.Set("tags", listAsSet(cur[user])); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteTagUser(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, user := parseTagUserId(d.Id())

	cur, err := fw.UserId.GetUserTags(user, vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}

	msg := &userid.Message{
		UntagUsers: []userid.UntagUser{
			userid.UntagUser{
				User: user,
				Tags: cur[user],
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
