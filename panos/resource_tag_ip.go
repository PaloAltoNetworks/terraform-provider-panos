package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/userid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTagIp() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateTagIp,
		Read:   readTagIp,
		Update: createUpdateTagIp,
		Delete: deleteTagIp,

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
			"ip": {
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

func buildTagIpId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseTagIpId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parseTagIp(ip string, cur map[string][]string, d *schema.ResourceData) *userid.Message {
	tagList := d.Get("tags").(*schema.Set).List()
	curTags := cur[ip]

	missing := make([]string, 0, len(tagList))
	extras := make([]string, 0, len(curTags))

	// Loop over what the user wants the IP tagged as.
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
			missing = append(missing, tag)
		}
	}

	// Loop over what the IP is actually tagged as right now.
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
		msg.TagIps = []userid.TagIp{
			userid.TagIp{
				Ip:   ip,
				Tags: missing,
			},
		}
	}

	if len(extras) > 0 {
		msg.UntagIps = []userid.UntagIp{
			userid.UntagIp{
				Ip:   ip,
				Tags: extras,
			},
		}
	}

	return msg
}

func createUpdateTagIp(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys := d.Get("vsys").(string)
	ip := d.Get("ip").(string)

	cur, err := fw.UserId.GetIpTags(ip, "", vsys)
	if err != nil {
		return err
	}

	msg := parseTagIp(ip, cur, d)

	if err = fw.UserId.Run(msg, vsys); err != nil {
		return err
	}

	d.SetId(buildTagIpId(vsys, ip))
	return readTagIp(d, meta)
}

func readTagIp(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys, ip := parseTagIpId(d.Id())

	cur, err := fw.UserId.GetIpTags(ip, "", vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("vsys", vsys)
	d.Set("ip", ip)
	if err := d.Set("tags", listAsSet(cur[ip])); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteTagIp(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, ip := parseTagIpId(d.Id())

	cur, err := fw.UserId.GetIpTags(ip, "", vsys)
	if err != nil || len(cur) == 0 {
		d.SetId("")
		return nil
	}

	msg := &userid.Message{
		UntagIps: []userid.UntagIp{
			userid.UntagIp{
				Ip:   ip,
				Tags: cur[ip],
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
