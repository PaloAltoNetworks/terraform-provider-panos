package panos

import (
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/dev/ha"
	halink "github.com/fpluchorg/pango/dev/ha/monitor/link"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schema handling.
func haSchema(isResource bool, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Description: "The template.",
			Optional:    true,
			ForceNew:    true,
		},
		"template_stack": {
			Type:        schema.TypeString,
			Description: "The template stack.",
			Optional:    true,
			ForceNew:    true,
		},
		"vsys": {
			Type:        schema.TypeString,
			Description: "The vsys.",
			Optional:    true,
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable high availability",
		},
		"group_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "HA pair group ID",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HA pair description",
		},
		"mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "High availability mode",
		},
		"peer_ha1_ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HA1 peer ip address",
		},
		"backup_peer_ha1_ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HA1 backup peer ip address",
		},
		"config_sync_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Configuration synchroniation between ha peers",
		},
		"ha1": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA1 local port",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA1 local ip address",
					},
					"netmask": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA1 local netmask",
					},
					"gateway": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA1 local gateway",
					},
					"encryption_enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable HA1 encryption",
					},
					"monitor_hold_time": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "HA1 monitoring hole time (ms)",
					},
				},
			},
		},
		"ha1_backup": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA1 local port",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA1 local ip address",
					},
					"netmask": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA1 local netmask",
					},
					"gateway": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA1 local gateway",
					},
				},
			},
		},
		"ha2": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA2 local port",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA2 local ip address",
					},
					"netmask": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA2 local netmask",
					},
					"gateway": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "HA2 local gateway",
					},
				},
			},
		},
		"ha2_backup": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA2 local port",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA2 local ip address",
					},
					"netmask": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA2 local netmask",
					},
					"gateway": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Backup HA2 local gateway",
					},
				},
			},
		},
		"ha2_state_sync_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable HA2 session synchroniation",
		},
		"ha2_state_sync_transport": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HA2 session synchroniation transport mode",
		},
		"ha2_state_sync_keepalive_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable HA2 session synchroniation keepalive",
		},
		"ha2_state_sync_keepalive_action": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "HA2 session synchroniation keepalive action",
		},
		"ha2_state_sync_keepalive_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "HA2 session synchroniation keepalive threshold (ms)",
		},
		"election_device_priority": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Device priority in HA active/passive election",
		},
		"election_preemptive": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Preemptive HA active/passive election",
		},
		"election_heartbeat_backup": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Preemptive HA active/passive election",
		},
		"link_monitor_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable HA link monitoring",
		},
		"link_monitor_failure_condition": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "any",
			Description: "HA link monitoring failure condition",
		},
		"path_monitor_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable HA path monitoring",
		},
		"path_monitor_failure_condition": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "any",
			Description: "HA path monitoring failure condition",
		},
		"link_group": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Link group name",
					},
					"enable": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Enable link group",
					},
					"failure_condition": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Link group failure condition",
					},
					"interfaces": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Link group member interfaces",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys"})
	}

	return ans
}
func resourceHa() *schema.Resource {
	return &schema.Resource{
		Create: createHa,
		Read:   readHa,
		Update: updateHa,
		Delete: deleteHa,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: virtualRouterSchema(true, []string{"template", "template_stack", "vsys"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: virtualRouterUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: virtualRouterSchema(true, nil),
	}
}

func resourcePanoramaHa() *schema.Resource {
	return &schema.Resource{
		Create: createHa,
		Read:   readHa,
		Update: updateHa,
		Delete: deleteHa,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: virtualRouterSchema(true, []string{}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: virtualRouterUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: virtualRouterSchema(true, nil),
	}
}

func parseHa(d *schema.ResourceData) ha.Config {
	ha1 := ha.Ha1Interface{}
	ha1List := d.Get("ha1").([]interface{})
	for i := range ha1List {
		elm := ha1List[i].(map[string]interface{})
		ha1 = ha.Ha1Interface{
			Port:             elm["port"].(string),
			IpAddress:        elm["ip_address"].(string),
			Netmask:          elm["netmask"].(string),
			Gateway:          elm["gateway"].(string),
			EncryptionEnable: elm["encryption_enable"].(bool),
			MonitorHoldTime:  elm["monitor_hold_time"].(int),
		}
	}

	ha1_backup := ha.Ha1BackupInterface{}
	ha1_backupList := d.Get("ha1_backup").([]interface{})
	for i := range ha1_backupList {
		elm := ha1_backupList[i].(map[string]interface{})
		ha1_backup = ha.Ha1BackupInterface{
			Port:      elm["port"].(string),
			IpAddress: elm["ip_address"].(string),
			Netmask:   elm["netmask"].(string),
			Gateway:   elm["gateway"].(string),
		}
	}

	ha2 := ha.Ha2Interface{}
	ha2List := d.Get("ha2").([]interface{})
	for i := range ha2List {
		elm := ha2List[i].(map[string]interface{})
		ha2 = ha.Ha2Interface{
			Port:      elm["port"].(string),
			IpAddress: elm["ip_address"].(string),
			Netmask:   elm["netmask"].(string),
			Gateway:   elm["gateway"].(string),
		}
	}

	ha2_backup := ha.Ha2Interface{}
	ha2_backupList := d.Get("ha2_backup").([]interface{})
	for i := range ha2_backupList {
		elm := ha2_backupList[i].(map[string]interface{})
		ha2_backup = ha.Ha2Interface{
			Port:      elm["port"].(string),
			IpAddress: elm["ip_address"].(string),
			Netmask:   elm["netmask"].(string),
			Gateway:   elm["gateway"].(string),
		}
	}

	// Link-group list
	link_group := []halink.Entry{}
	link_groupList := d.Get("link_group").([]interface{})
	for i := range link_groupList {
		elm := link_groupList[i].(map[string]interface{})
		o := halink.Entry{
			Name:             elm["name"].(string),
			Enable:           elm["enable"].(bool),
			FailureCondition: elm["failure_condition"].(string),
			Interfaces:       []string{},
		}
		link_group = append(link_group, o)
	}

	o := ha.Config{
		Enable:                 d.Get("enable").(bool),
		GroupId:                d.Get("group_id").(int),
		Description:            d.Get("description").(string),
		Mode:                   d.Get("mode").(string),
		PeerHa1IpAddress:       d.Get("peer_ha1_ip_address").(string),
		BackupPeerHa1IpAddress: d.Get("backup_peer_ha1_ip_address").(string),
		ConfigSyncEnable:       d.Get("config_sync_enable").(bool),

		Ha1:       &ha1,
		Ha1Backup: &ha1_backup,

		Ha2:                            &ha2,
		Ha2Backup:                      &ha2_backup,
		Ha2StateSyncEnable:             d.Get("ha2_state_sync_enable").(bool),
		Ha2StateSyncTransport:          d.Get("ha2_state_sync_transport").(string),
		Ha2StateSyncKeepAliveEnable:    d.Get("ha2_state_sync_keepalive_enable").(bool),
		Ha2StateSyncKeepAliveAction:    d.Get("ha2_state_sync_keepalive_action").(string),
		Ha2StateSyncKeepAliveThreshold: d.Get("ha2_state_sync_keepalive_threshold").(int),

		ElectionDevicePriority:  d.Get("election_device_priority").(string),
		ElectionPreemptive:      d.Get("election_preemptive").(bool),
		ElectionHeartBeatBackup: d.Get("election_heartbeat_backup").(bool),

		LinkMonitorEnable:           d.Get("link_monitor_enable").(bool),
		LinkMonitorFailureCondition: d.Get("link_monitor_failure_condition").(string),

		PathMonitorEnable:           d.Get("path_monitor_enable").(bool),
		PathMonitorFailureCondition: d.Get("path_monitor_failure_condition").(string),
	}

	return o
}

func parseHaId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildHaId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createHa(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := parseHa(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildHaId(tmpl, ts, vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.HaConfig.Edit(o)
	case *pango.Panorama:
		err = con.Device.HaConfig.Edit(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readHa(d, meta)
}

func saveHa(d *schema.ResourceData, o ha.Config) error {
	var err error

	d.Set("enable", o.Enable)
	d.Set("group_id", o.GroupId)
	d.Set("description", o.Description)
	d.Set("mode", o.Mode)
	d.Set("peer_ha1_ip_address", o.PeerHa1IpAddress)
	d.Set("backup_peer_ha1_ip_address", o.BackupPeerHa1IpAddress)
	d.Set("config_sync_enable", o.ConfigSyncEnable)
	d.Set("ha2_state_sync_enable", o.Ha2StateSyncEnable)
	d.Set("ha2_state_sync_transport", o.Ha2StateSyncTransport)
	d.Set("ha2_state_sync_keepalive_enable", o.Ha2StateSyncKeepAliveEnable)
	d.Set("ha2_state_sync_keepalive_action", o.Ha2StateSyncKeepAliveAction)
	d.Set("ha2_state_sync_keepalive_threshold", o.Ha2StateSyncKeepAliveThreshold)
	d.Set("election_device_priority", o.ElectionDevicePriority)
	d.Set("election_preemptive", o.ElectionPreemptive)
	d.Set("election_heartbeat_backup", o.ElectionHeartBeatBackup)
	d.Set("link_monitor_enable", o.LinkMonitorEnable)
	d.Set("link_monitor_failure_condition", o.LinkMonitorFailureCondition)
	d.Set("path_monitor_enable", o.PathMonitorEnable)
	d.Set("path_monitor_failure_condition", o.PathMonitorFailureCondition)

	var ha1List []interface{}
	if o.Ha1 != nil {
		ha1List = make([]interface{}, 0, 1)
		m := map[string]interface{}{
			"port":              o.Ha1.Port,
			"ip_address":        o.Ha1.IpAddress,
			"netmask":           o.Ha1.Netmask,
			"gateway":           o.Ha1.Gateway,
			"encryption_enable": o.Ha1.EncryptionEnable,
			"monitor_hold_time": o.Ha1.MonitorHoldTime,
		}

		ha1List = append(ha1List, m)
	}
	if err = d.Set("ha1", ha1List); err != nil {
		log.Printf("[WARN] Error setting 'ha1' param for %q: %s", d.Id(), err)
	}

	var ha1_backupList []interface{}
	if o.Ha1Backup != nil {
		ha1_backupList = make([]interface{}, 0, 1)
		m := map[string]interface{}{
			"port":       o.Ha1Backup.Port,
			"ip_address": o.Ha1Backup.IpAddress,
			"netmask":    o.Ha1Backup.Netmask,
			"gateway":    o.Ha1Backup.Gateway,
		}

		ha1_backupList = append(ha1_backupList, m)
	}
	if err = d.Set("ha1_backup", ha1_backupList); err != nil {
		log.Printf("[WARN] Error setting 'ha1_backup' param for %q: %s", d.Id(), err)
	}

	var ha2List []interface{}
	if o.Ha2 != nil {
		ha2List = make([]interface{}, 0, 1)
		m := map[string]interface{}{
			"port":       o.Ha2.Port,
			"ip_address": o.Ha2.IpAddress,
			"netmask":    o.Ha2.Netmask,
			"gateway":    o.Ha2.Gateway,
		}

		ha2List = append(ha1List, m)
	}
	if err = d.Set("ha2", ha2List); err != nil {
		log.Printf("[WARN] Error setting 'ha2' param for %q: %s", d.Id(), err)
	}

	var ha2_backupList []interface{}
	if o.Ha2Backup != nil {
		ha2_backupList = make([]interface{}, 0, 1)
		m := map[string]interface{}{
			"port":       o.Ha1Backup.Port,
			"ip_address": o.Ha1Backup.IpAddress,
			"netmask":    o.Ha1Backup.Netmask,
			"gateway":    o.Ha1Backup.Gateway,
		}

		ha2_backupList = append(ha2_backupList, m)
	}
	if err = d.Set("ha2_backup", ha2_backupList); err != nil {
		log.Printf("[WARN] Error setting 'ha2_backup' param for %q: %s", d.Id(), err)
	}

	return nil
}

func readHa(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ha.Config

	tmpl, ts, vsys := parseHaId(d.Id())
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.HaConfig.Get()
	case *pango.Panorama:
		o, err = con.Device.HaConfig.Get(tmpl, ts, vsys)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveHa(d, o)

	return nil
}

func updateHa(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := parseHa(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Device.HaConfig.Get()
		if err != nil {
			return err
		}
		// Todo: Should implement a merge
		lo.Copy(o)
		if err = con.Device.HaConfig.Edit(lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Device.HaConfig.Get(tmpl, ts, vsys)
		if err != nil {
			return err
		}
		// Todo: Should implement a merge
		lo.Copy(o)
		if err = con.Device.HaConfig.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	return readHa(d, meta)
}

func deleteHa(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := parseHa(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Device.HaConfig.Get()
		if err != nil {
			return err
		}

		// Todo: Should implement a merge
		lo.Copy(o)
		err = con.Device.HaConfig.Delete()
	case *pango.Panorama:
		lo, err := con.Device.HaConfig.Get(tmpl, ts, vsys)
		if err != nil {
			return err
		}
		// Todo: Should implement a merge
		lo.Copy(o)
		err = con.Device.HaConfig.Delete(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
