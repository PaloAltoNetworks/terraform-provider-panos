package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/userid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source.
func dataSourceIpTag() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceIpTag,

		Schema: map[string]*schema.Schema{
			"vsys": vsysSchema(),
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optionally filter on just this single IP address",
			},
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optionally filter on just this single tag",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of entries",
			},
			"entries": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of entry specs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The IP's tags",
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

func readDataSourceIpTag(d *schema.ResourceData, meta interface{}) error {
	var ans map[string][]string
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		ans, err = con.UserId.GetIpTags(
			d.Get("ip").(string),
			d.Get("tag").(string),
			d.Get("vsys").(string),
		)
	case *pango.Panorama:
		ans, err = con.UserId.GetIpTags(
			d.Get("ip").(string),
			d.Get("tag").(string),
			d.Get("vsys").(string),
		)
	}

	if err != nil {
		return err
	}

	d.SetId(base64Encode([]interface{}{
		d.Get("ip").(string), d.Get("tag").(string), d.Get("vsys").(string),
	}))

	if len(ans) == 0 {
		d.Set("total", 0)
		d.Set("entries", nil)
	} else {
		d.Set("total", len(ans))
		list := make([]interface{}, 0, len(ans))
		for ip, tags := range ans {
			list = append(list, map[string]interface{}{
				"ip":   ip,
				"tags": listAsSet(tags),
			})
		}

		if err = d.Set("entries", list); err != nil {
			log.Printf("[WARN] Error setting 'entries' for %q: %s", d.Id(), err)
		}
	}

	return nil
}

// Resource.
func resourceIpTag() *schema.Resource {
	return &schema.Resource{
		Create: createIpTag,
		Read:   readIpTag,
		Delete: deleteIpTag,

		Schema: map[string]*schema.Schema{
			"vsys": vsysSchema(),
			"ip": {
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

func createIpTag(d *schema.ResourceData, meta interface{}) error {

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		ip := d.Get("ip").(string)
		tagList := d.Get("tags").(*schema.Set).List()

		cur, err := con.UserId.GetIpTags(ip, "", vsys)
		if err != nil {
			return err
		}
		curTags := cur[ip]

		missing := make([]string, 0, len(tagList))
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
				missing = append(missing, tag)
			}
		}

		if len(missing) > 0 {
			msg := &userid.Message{
				TagIps: []userid.TagIp{
					userid.TagIp{
						Ip:   ip,
						Tags: missing,
					},
				},
			}

			if err = con.UserId.Run(msg, vsys); err != nil {
				return err
			}
		}

		d.SetId(buildIpTagId(vsys, ip, tagList))
	case *pango.Panorama:
		vsys := d.Get("vsys").(string)
		ip := d.Get("ip").(string)
		tagList := d.Get("tags").(*schema.Set).List()

		cur, err := con.UserId.GetIpTags(ip, "", vsys)
		if err != nil {
			return err
		}
		curTags := cur[ip]

		missing := make([]string, 0, len(tagList))
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
				missing = append(missing, tag)
			}
		}

		if len(missing) > 0 {
			msg := &userid.Message{
				TagIps: []userid.TagIp{
					userid.TagIp{
						Ip:   ip,
						Tags: missing,
					},
				},
			}

			if err = con.UserId.Run(msg, vsys); err != nil {
				return err
			}
		}

		d.SetId(buildIpTagId(vsys, ip, tagList))
	}

	return readIpTag(d, meta)
}

func readIpTag(d *schema.ResourceData, meta interface{}) error {
	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, ip, tagList := parseIpTagId(d.Id())

		cur, err := con.UserId.GetIpTags(ip, "", vsys)
		if err != nil || len(cur) == 0 {
			d.SetId("")
			return nil
		}
		curTags := cur[ip]

		list := make([]string, 0, len(tagList))
		for _, tag := range tagList {
			for _, curTag := range curTags {
				if tag == curTag {
					list = append(list, tag)
					break
				}
			}
		}

		d.Set("ip", ip)
		if len(list) == 0 {
			d.Set("tags", nil)
		} else if err = d.Set("tags", listAsSet(list)); err != nil {
			log.Printf("[WARN] Error setting 'tags' for %q: %s", d.Id(), err)
		}
	case *pango.Panorama:
		vsys, ip, tagList := parseIpTagId(d.Id())

		cur, err := con.UserId.GetIpTags(ip, "", vsys)
		if err != nil || len(cur) == 0 {
			d.SetId("")
			return nil
		}
		curTags := cur[ip]

		list := make([]string, 0, len(tagList))
		for _, tag := range tagList {
			for _, curTag := range curTags {
				if tag == curTag {
					list = append(list, tag)
					break
				}
			}
		}

		d.Set("ip", ip)
		if len(list) == 0 {
			d.Set("tags", nil)
		} else if err = d.Set("tags", listAsSet(list)); err != nil {
			log.Printf("[WARN] Error setting 'tags' for %q: %s", d.Id(), err)
		}
	}

	return nil
}

func deleteIpTag(d *schema.ResourceData, meta interface{}) error {
	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, ip, tagList := parseIpTagId(d.Id())

		cur, err := con.UserId.GetIpTags(ip, "", vsys)
		if err != nil {
			return err
		}
		curTags := cur[ip]

		list := make([]string, 0, len(tagList))
		for _, tag := range tagList {
			for _, curTag := range curTags {
				if tag == curTag {
					list = append(list, tag)
					break
				}
			}
		}

		if len(list) > 0 {
			msg := &userid.Message{
				UntagIps: []userid.UntagIp{
					userid.UntagIp{
						Ip:   ip,
						Tags: list,
					},
				},
			}

			if err = con.UserId.Run(msg, vsys); err != nil {
				return err
			}
		}
	case *pango.Panorama:
		vsys, ip, tagList := parseIpTagId(d.Id())

		cur, err := con.UserId.GetIpTags(ip, "", vsys)
		if err != nil {
			return err
		}
		curTags := cur[ip]

		list := make([]string, 0, len(tagList))
		for _, tag := range tagList {
			for _, curTag := range curTags {
				if tag == curTag {
					list = append(list, tag)
					break
				}
			}
		}

		if len(list) > 0 {
			msg := &userid.Message{
				UntagIps: []userid.UntagIp{
					userid.UntagIp{
						Ip:   ip,
						Tags: list,
					},
				},
			}

			if err = con.UserId.Run(msg, vsys); err != nil {
				return err
			}
		}
	}

	d.SetId("")
	return nil
}

// Id functions.
func buildIpTagId(a, b string, c []interface{}) string {
	return strings.Join([]string{a, b, base64Encode(c)}, IdSeparator)
}

func parseIpTagId(v string) (string, string, []string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], base64Decode(t[2])
}
