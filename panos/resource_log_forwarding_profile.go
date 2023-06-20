package panos

import (
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/profile/logfwd"
	"github.com/fpluchorg/pango/objs/profile/logfwd/matchlist"
	"github.com/fpluchorg/pango/objs/profile/logfwd/matchlist/action"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLogForwardingProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateLogForwardingProfile,
		Read:   readLogForwardingProfile,
		Update: createUpdateLogForwardingProfile,
		Delete: deleteLogForwardingProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: logForwardingProfileSchema(false),
	}
}

func logForwardingProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"enhanced_logging": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"match_list": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"log_type": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  matchlist.LogTypeTraffic,
						ValidateFunc: validateStringIn(
							matchlist.LogTypeTraffic,
							matchlist.LogTypeThreat,
							matchlist.LogTypeWildfire,
							matchlist.LogTypeUrl,
							matchlist.LogTypeData,
							matchlist.LogTypeGtp,
							matchlist.LogTypeTunnel,
							matchlist.LogTypeAuth,
							matchlist.LogTypeSctp,
							matchlist.LogTypeDecryption,
						),
					},
					"filter": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "All Logs",
					},
					"send_to_panorama": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"snmptrap_server_profiles": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"email_server_profiles": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"syslog_server_profiles": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"http_server_profiles": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"action": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},
								"azure_integration": {
									Type:     schema.TypeList,
									MaxItems: 1,
									Optional: true,
									//ConflictsWith: []string{"match_list.action.tagging_integration"},
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"azure_integration": {
												Type:     schema.TypeBool,
												Optional: true,
												Default:  true,
											},
										},
									},
								},
								"tagging_integration": {
									Type:     schema.TypeList,
									MaxItems: 1,
									Optional: true,
									//ConflictsWith: []string{"match_list.action.azure_integration"},
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:         schema.TypeString,
												Optional:     true,
												Default:      action.ActionAddTag,
												ValidateFunc: validateStringIn(action.ActionAddTag, action.ActionRemoveTag),
											},
											"target": {
												Type:         schema.TypeString,
												Optional:     true,
												Default:      action.TargetSource,
												ValidateFunc: validateStringIn(action.TargetSource, action.TargetDestination),
											},
											"timeout": {
												Type:     schema.TypeInt,
												Optional: true,
											},
											"local_registration": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												/*
													ConflictsWith: []string{
														"match_list.action.tagging_integration.remote_registration",
														"match_list.action.tagging_integration.panorama_registration",
													},
												*/
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"tags": {
															Type:     schema.TypeList,
															Required: true,
															MinItems: 1,
															Elem: &schema.Schema{
																Type: schema.TypeString,
															},
														},
													},
												},
											},
											"remote_registration": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												/*
													ConflictsWith: []string{
														"match_list.action.tagging_integration.local_registration",
														"match_list.action.tagging_integration.panorama_registration",
													},
												*/
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"http_profile": {
															Type:     schema.TypeString,
															Required: true,
														},
														"tags": {
															Type:     schema.TypeList,
															Required: true,
															MinItems: 1,
															Elem: &schema.Schema{
																Type: schema.TypeString,
															},
														},
													},
												},
											},
											"panorama_registration": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												/*
													ConflictsWith: []string{
														"match_list.action.tagging_integration.local_registration",
														"match_list.action.tagging_integration.remote_registration",
													},
												*/
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"tags": {
															Type:     schema.TypeList,
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
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if p {
		ans["device_group"] = deviceGroupSchema()
	} else {
		ans["vsys"] = vsysSchema("vsys1")
	}

	return ans
}

func parseLogForwardingProfile(d *schema.ResourceData) (string, logfwd.Entry, []matchlist.Entry, map[string][]action.Entry) {
	vsys := d.Get("vsys").(string)
	o, ml, mla := loadLogForwardingProfile(d)

	return vsys, o, ml, mla
}

func loadLogForwardingProfile(d *schema.ResourceData) (logfwd.Entry, []matchlist.Entry, map[string][]action.Entry) {
	o := logfwd.Entry{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		EnhancedLogging: d.Get("enhanced_logging").(bool),
	}

	mle := d.Get("match_list").([]interface{})
	if len(mle) == 0 {
		return o, nil, nil
	}

	ml := make([]matchlist.Entry, 0, len(mle))
	mla := make(map[string][]action.Entry)

	for i := range mle {
		mm := mle[i].(map[string]interface{})
		match_entry := matchlist.Entry{
			Name:           mm["name"].(string),
			Description:    mm["description"].(string),
			LogType:        mm["log_type"].(string),
			Filter:         mm["filter"].(string),
			SendToPanorama: mm["send_to_panorama"].(bool),
			SnmpProfiles:   setAsList(mm["snmptrap_server_profiles"].(*schema.Set)),
			EmailProfiles:  setAsList(mm["email_server_profiles"].(*schema.Set)),
			SyslogProfiles: setAsList(mm["syslog_server_profiles"].(*schema.Set)),
			HttpProfiles:   setAsList(mm["http_server_profiles"].(*schema.Set)),
		}
		ml = append(ml, match_entry)
		ael := mm["action"].([]interface{})
		if len(ael) == 0 {
			continue
		}

		action_list := make([]action.Entry, 0, len(ael))
		for j := range ael {
			ae := ael[j].(map[string]interface{})
			action_entry := action.Entry{
				Name: ae["name"].(string),
			}
			if x := asInterfaceMap(ae, "azure_integration"); len(x) != 0 {
				action_entry.ActionType = action.ActionTypeIntegration
				action_entry.Action = action.ActionAzure
			} else if x := asInterfaceMap(ae, "tagging_integration"); len(x) != 0 {
				action_entry.ActionType = action.ActionTypeTagging
				action_entry.Action = x["action"].(string)
				action_entry.Target = x["target"].(string)
				action_entry.Timeout = x["timeout"].(int)
				if y := asInterfaceMap(x, "local_registration"); len(y) != 0 {
					action_entry.Registration = action.RegistrationLocal
					action_entry.Tags = asStringList(y["tags"].([]interface{}))
				} else if y := asInterfaceMap(x, "remote_registration"); len(y) != 0 {
					action_entry.Registration = action.RegistrationRemote
					action_entry.Tags = asStringList(y["tags"].([]interface{}))
					action_entry.HttpProfile = y["http_profile"].(string)
				} else if y := asInterfaceMap(x, "panorama_registration"); len(y) != 0 {
					action_entry.Registration = action.RegistrationPanorama
					action_entry.Tags = asStringList(y["tags"].([]interface{}))
				}
			}
			action_list = append(action_list, action_entry)
		}
		mla[match_entry.Name] = action_list
	}

	return o, ml, mla
}

func saveLogForwardingProfile(d *schema.ResourceData, o logfwd.Entry, ml []matchlist.Entry, mla map[string][]action.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("enhanced_logging", o.EnhancedLogging)

	if len(ml) == 0 {
		d.Set("match_list", nil)
		return
	}

	var action_list []action.Entry

	mle := make([]interface{}, 0, len(ml))
	for _, match_entry := range ml {
		mm := map[string]interface{}{
			"name":                     match_entry.Name,
			"description":              match_entry.Description,
			"log_type":                 match_entry.LogType,
			"filter":                   match_entry.Filter,
			"send_to_panorama":         match_entry.SendToPanorama,
			"snmptrap_server_profiles": listAsSet(match_entry.SnmpProfiles),
			"email_server_profiles":    listAsSet(match_entry.EmailProfiles),
			"syslog_server_profiles":   listAsSet(match_entry.SyslogProfiles),
			"http_server_profiles":     listAsSet(match_entry.HttpProfiles),
		}

		action_list = mla[match_entry.Name]
		if len(action_list) == 0 {
			mm["action"] = nil
		} else {
			ael := make([]interface{}, 0, len(action_list))
			for _, action_entry := range action_list {
				ae := map[string]interface{}{
					"name": action_entry.Name,
				}
				switch action_entry.ActionType {
				case action.ActionTypeIntegration:
					ae["azure_integration"] = []interface{}{
						map[string]interface{}{
							"azure_integration": true,
						},
					}
				case action.ActionTypeTagging:
					ti := map[string]interface{}{
						"action":  action_entry.Action,
						"target":  action_entry.Target,
						"timeout": action_entry.Timeout,
					}
					switch action_entry.Registration {
					case action.RegistrationLocal:
						ti["local_registration"] = []interface{}{
							map[string]interface{}{
								"tags": listAsSet(action_entry.Tags),
							},
						}
					case action.RegistrationPanorama:
						ti["panorama_registration"] = []interface{}{
							map[string]interface{}{
								"tags": listAsSet(action_entry.Tags),
							},
						}
					case action.RegistrationRemote:
						ti["remote_registration"] = []interface{}{
							map[string]interface{}{
								"tags":         listAsSet(action_entry.Tags),
								"http_profile": action_entry.HttpProfile,
							},
						}
					}
					ae["tagging_integration"] = []interface{}{ti}
				}
				ael = append(ael, ae)
			}
			mm["action"] = ael
		}

		mle = append(mle, mm)
	}

	if err := d.Set("match_list", mle); err != nil {
		log.Printf("[WARN] Error setting 'match_list' param for %q: %s", d.Id(), err)
	}
}

func parseLogForwardingProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildLogForwardingProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdateLogForwardingProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o, ml, mla := parseLogForwardingProfile(d)

	if err := fw.Objects.LogForwardingProfile.SetWithoutSubconfig(vsys, o); err != nil {
		return err
	}

	if err := fw.Objects.LogForwardingProfileMatchList.Set(vsys, o.Name, ml...); err != nil {
		return err
	}

	for _, entry := range ml {
		if err := fw.Objects.LogForwardingProfileMatchListAction.Set(vsys, o.Name, entry.Name, mla[entry.Name]...); err != nil {
			return err
		}
	}

	d.SetId(buildLogForwardingProfileId(vsys, o.Name))
	return readLogForwardingProfile(d, meta)
}

func readLogForwardingProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseLogForwardingProfileId(d.Id())

	o, err := fw.Objects.LogForwardingProfile.Get(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	var ml []matchlist.Entry
	var mla map[string][]action.Entry

	mlNames, err := fw.Objects.LogForwardingProfileMatchList.GetList(vsys, name)
	if err != nil {
		return err
	}

	if len(mlNames) > 0 {
		ml = make([]matchlist.Entry, 0, len(mlNames))
		mla = make(map[string][]action.Entry)
		for i := range mlNames {
			mle, err := fw.Objects.LogForwardingProfileMatchList.Get(vsys, name, mlNames[i])
			if err != nil {
				return err
			}
			ml = append(ml, mle)
			aNames, err := fw.Objects.LogForwardingProfileMatchListAction.GetList(vsys, name, mlNames[i])
			if err != nil {
				return err
			}
			if len(aNames) != 0 {
				actionList := make([]action.Entry, 0, len(aNames))
				for j := range aNames {
					ae, err := fw.Objects.LogForwardingProfileMatchListAction.Get(vsys, name, mlNames[i], aNames[j])
					if err != nil {
						return err
					}
					actionList = append(actionList, ae)
				}
				mla[mle.Name] = actionList
			}
		}
	}

	d.Set("vsys", vsys)
	saveLogForwardingProfile(d, o, ml, mla)

	return nil
}

func deleteLogForwardingProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseLogForwardingProfileId(d.Id())

	err := fw.Objects.LogForwardingProfile.Delete(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
