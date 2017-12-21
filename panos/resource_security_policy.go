package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: createSecurityPolicy,
		Read:   readSecurityPolicy,
		Update: updateSecurityPolicy,
		Delete: deleteSecurityPolicy,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The vsys to put this object in (default: vsys1)",
			},
			"rulebase": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The rulebase (default: rulebase, pre-rulebase, post-rulebase)",
				ValidateFunc: validateStringIn("rulebase", "pre-rulebase", "post-rulebase"),
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Security rule type (default: universal, interzone, intrazone)",
				ValidateFunc: validateStringIn("universal", "interzone", "intrazone"),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_zone": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_address": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"negate_source": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_user": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hip_profile": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_zone": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_address": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"negate_destination": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"application": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"category": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"action": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "allow",
				Description:  "Action (default: allow, deny, drop, reset-client, reset-server, reset-both)",
				ValidateFunc: validateStringIn("allow", "deny", "drop", "reset-client", "reset-server", "reset-both"),
			},
			"log_setting": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log forwarding profile",
			},
			"log_start": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_end": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"icmp_unreachable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_server_response_inspection": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"virus": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"spyware": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vulnerability": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"url_filtering": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_blocking": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"wildfire_analysis": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_filtering": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func parseSecurityPolicy(d *schema.ResourceData) (string, string, security.Entry) {
	vsys := d.Get("vsys").(string)
	rb := d.Get("rulebase").(string)

	o := security.Entry{
		Name:                            d.Get("name").(string),
		Type:                            d.Get("type").(string),
		Description:                     d.Get("description").(string),
		Tags:                            setAsList(d, "tags"),
		SourceZone:                      asStringList(d, "source_zone"),
		SourceAddress:                   asStringList(d, "source_address"),
		NegateSource:                    d.Get("negate_source").(bool),
		SourceUser:                      asStringList(d, "source_user"),
		HipProfile:                      asStringList(d, "hip_profile"),
		DestinationZone:                 asStringList(d, "destination_zone"),
		DestinationAddress:              asStringList(d, "destination_address"),
		NegateDestination:               d.Get("negate_destination").(bool),
		Application:                     asStringList(d, "application"),
		Service:                         asStringList(d, "service"),
		Category:                        asStringList(d, "category"),
		Action:                          d.Get("action").(string),
		LogSetting:                      d.Get("log_setting").(string),
		LogStart:                        d.Get("log_start").(bool),
		LogEnd:                          d.Get("log_end").(bool),
		Disabled:                        d.Get("disabled").(bool),
		Schedule:                        d.Get("schedule").(string),
		IcmpUnreachable:                 d.Get("icmp_unreachable").(bool),
		DisableServerResponseInspection: d.Get("disable_server_response_inspection").(bool),
		Group:            d.Get("group").(string),
		Virus:            d.Get("virus").(string),
		Spyware:          d.Get("spyware").(string),
		Vulnerability:    d.Get("vulnerability").(string),
		UrlFiltering:     d.Get("url_filtering").(string),
		FileBlocking:     d.Get("file_blocking").(string),
		WildFireAnalysis: d.Get("wildfire_analysis").(string),
		DataFiltering:    d.Get("data_filtering").(string),
	}

	return vsys, rb, o
}

func saveDataSecurityPolicy(d *schema.ResourceData, vsys, rb string, o security.Entry) {
	var err error
	d.SetId(buildSecurityPolicyId(vsys, rb, o.Name))
	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("rulebase", rb)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	if err = d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("source_zone", o.SourceZone); err != nil {
		log.Printf("[WARN] Error setting 'source_zone' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("source_address", o.SourceAddress); err != nil {
		log.Printf("[WARN] Error setting 'source_address' param for %q: %s", d.Id(), err)
	}
	d.Set("negate_source", o.NegateSource)
	if err = d.Set("source_user", o.SourceUser); err != nil {
		log.Printf("[WARN] Error setting 'source_user' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("hip_profile", o.HipProfile); err != nil {
		log.Printf("[WARN] Error setting 'hip_profile' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("destination_zone", o.DestinationZone); err != nil {
		log.Printf("[WARN] Error setting 'destination_zone' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("destination_address", o.DestinationAddress); err != nil {
		log.Printf("[WARN] Error setting 'destination_address' param for %q: %s", d.Id(), err)
	}
	d.Set("negate_destination", o.NegateDestination)
	if err = d.Set("application", o.Application); err != nil {
		log.Printf("[WARN] Error setting 'application' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("service", o.Service); err != nil {
		log.Printf("[WARN] Error setting 'service' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("category", o.Category); err != nil {
		log.Printf("[WARN] Error setting 'category' param for %q: %s", d.Id(), err)
	}
	d.Set("action", o.Action)
	d.Set("log_setting", o.LogSetting)
	d.Set("log_start", o.LogStart)
	d.Set("log_end", o.LogEnd)
	d.Set("disabled", o.Disabled)
	d.Set("schedule", o.Schedule)
	d.Set("icmp_unreachable", o.IcmpUnreachable)
	d.Set("disable_server_response_inspection", o.DisableServerResponseInspection)
	d.Set("group", o.Group)
	d.Set("virus", o.Virus)
	d.Set("spyware", o.Spyware)
	d.Set("vulnerability", o.Vulnerability)
	d.Set("url_filtering", o.UrlFiltering)
	d.Set("file_blocking", o.FileBlocking)
	d.Set("wildfire_analysis", o.WildFireAnalysis)
	d.Set("data_filtering", o.DataFiltering)
}

func parseSecurityPolicyId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildSecurityPolicyId(a, b, c string) string {
	return fmt.Sprintf("%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c)
}

func createSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb, o := parseSecurityPolicy(d)
	o.Defaults()
	/* Defaults() sets LogEnd to true if it is false, but if the user
	   actually wants it as false, we need to reset it to what we were
	   passed from the plan file. */
	o.LogEnd = d.Get("log_end").(bool)

	if err := fw.Policies.Security.VerifiableSet(vsys, rb, o); err != nil {
		return err
	}

	saveDataSecurityPolicy(d, vsys, rb, o)
	return nil
}

func readSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb, name := parseSecurityPolicyId(d.Id())

	o, err := fw.Policies.Security.Get(vsys, rb, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	saveDataSecurityPolicy(d, vsys, rb, o)
	return nil
}

func updateSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error
	fw := meta.(*pango.Firewall)
	vsys, rb, o := parseSecurityPolicy(d)
	o.Defaults()
	/* Defaults() sets LogEnd to true if it is false, but if the user
	   actually wants it as false, we need to reset it to what we were
	   passed from the plan file. */
	o.LogEnd = d.Get("log_end").(bool)

	lo, err := fw.Policies.Security.Get(vsys, rb, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	err = fw.Policies.Security.Edit(vsys, rb, lo)

	if err == nil {
		saveDataSecurityPolicy(d, vsys, rb, o)
	}
	return err
}

func deleteSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb, name := parseSecurityPolicyId(d.Id())

	err := fw.Policies.Security.Delete(vsys, rb, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
