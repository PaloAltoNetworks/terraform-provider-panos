package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango/userid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUseridLogin() *schema.Resource {
	return &schema.Resource{
		Create: createUseridLogin,
		Read:   readUseridLogin,
		Delete: deleteUseridLogin,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				Description: "The vsys to config DAG tags for",
				ForceNew:    true,
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User that should be logged in",
				ForceNew:    true,
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address the user is logging in from",
				ForceNew:    true,
			},
		},
	}
}

func buildUseridLoginId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parseUseridLoginId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func createUseridLogin(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys := d.Get("vsys").(string)
	ip := d.Get("ip").(string)
	user := d.Get("user").(string)

	msg := &userid.Message{
		Logins: []userid.Login{
			userid.Login{
				User: user,
				Ip:   ip,
			},
		},
	}

	if err = fw.UserId.Run(msg, vsys); err != nil {
		return err
	}

	d.SetId(buildUseridLoginId(vsys, ip, user))
	return readUseridLogin(d, meta)
}

func readUseridLogin(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys, ip, _ := parseUseridLoginId(d.Id())

	list, err := fw.UserId.GetLogins(ip, "", vsys)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("vsys", vsys)
	d.Set("ip", ip)
	d.Set("user", list[0].User)

	return nil
}

func deleteUseridLogin(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	vsys, ip, user := parseUseridLoginId(d.Id())

	msg := &userid.Message{
		Logouts: []userid.Logout{
			userid.Logout{
				User: user,
				Ip:   ip,
			},
		},
	}

	err = fw.UserId.Run(msg, vsys)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
