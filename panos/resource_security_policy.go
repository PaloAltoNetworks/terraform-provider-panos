package panos

import (
    "github.com/PaloAltoNetworks/pango"
    "github.com/PaloAltoNetworks/pango/poli/security"

    "github.com/hashicorp/terraform/helper/schema"
)


func resourceSecurityPolicy() *schema.Resource {
    return &schema.Resource{
        Create: createSecurityPolicy,
        Read: readSecurityPolicy,
        Update: updateSecurityPolicy,
        Delete: deleteSecurityPolicy,

        Schema: map[string] *schema.Schema{
            "name": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
            },
            "vsys": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
                Description: "The vsys to put this object in",
            },
            "rulebase": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
                Description: "The rulebase (default: rulebase, pre-rulebase, post-rulebase)",
            },
            "type": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
                Description: "Security rule type (default: universal, interzone, intrazone)",
            },
            "description": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
            },
            "tags": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "source_zone": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "source_address": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "negate_source": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
            },
            "source_user": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "hip_profile": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "destination_zone": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "destination_address": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "negate_destination": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
            },
            "application": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "service": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "category": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "action": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Default: "allow",
                Description: "Action (default: allow, deny, drop, reset-client, reset-server, reset-both)",
            },
            "log_setting": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Description: "Log forwarding profile",
            },
            "log_start": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
            },
            "log_end": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
                Default: true,
            },
            "disabled": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
            },
            "schedule": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
            },
            "icmp_unreachable": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
            },
            "disable_server_response_inspection": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
            },
            "group": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "virus": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "spyware": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "vulnerability": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "url_filtering": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "file_blocking": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "wildfire_analysis": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "data_filtering": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func parseSecurityPolicy(d *schema.ResourceData) (string, string, security.Entry) {
    vsys := d.Get("vsys").(string)
    rb := d.Get("rulebase").(string)

    o := security.Entry{
        Name: d.Get("name").(string),
        Type: d.Get("type").(string),
        Description: d.Get("description").(string),
        Tags: asStringList(d, "tags"),
        SourceZone: asStringList(d, "source_zone"),
        SourceAddress: asStringList(d, "source_address"),
        NegateSource: d.Get("negate_source").(bool),
        SourceUser: asStringList(d, "source_user"),
        HipProfile: asStringList(d, "hip_profile"),
        DestinationZone: asStringList(d, "destination_zone"),
        DestinationAddress: asStringList(d, "destination_address"),
        NegateDestination: d.Get("negate_destination").(bool),
        Application: asStringList(d, "application"),
        Service: asStringList(d, "service"),
        Category: asStringList(d, "category"),
        Action: d.Get("action").(string),
        LogSetting: d.Get("log_setting").(string),
        LogStart: d.Get("log_start").(bool),
        LogEnd: d.Get("log_end").(bool),
        Disabled: d.Get("disabled").(bool),
        Schedule: d.Get("schedule").(string),
        IcmpUnreachable: d.Get("icmp_unreachable").(bool),
        DisableServerResponseInspection: d.Get("disable_server_response_inspection").(bool),
        Group: asStringList(d, "group"),
        Virus: asStringList(d, "virus"),
        Spyware: asStringList(d, "spyware"),
        Vulnerability: asStringList(d, "vulnerability"),
        UrlFiltering: asStringList(d, "url_filtering"),
        FileBlocking: asStringList(d, "file_blocking"),
        WildFireAnalysis: asStringList(d, "wildfire_analysis"),
        DataFiltering: asStringList(d, "data_filtering"),
    }

    return vsys, rb, o
}

func saveDataSecurityPolicy(d *schema.ResourceData, o security.Entry) {
    d.SetId(o.Name)
    d.Set("type", o.Type)
    d.Set("description", o.Description)
    d.Set("tags", o.Tags)
    d.Set("source_zone", o.SourceZone)
    d.Set("source_address", o.SourceAddress)
    d.Set("negate_source", o.NegateSource)
    d.Set("source_user", o.SourceUser)
    d.Set("hip_profile", o.HipProfile)
    d.Set("destination_zone", o.DestinationZone)
    d.Set("destination_address", o.DestinationAddress)
    d.Set("negate_destination", o.NegateDestination)
    d.Set("application", o.Application)
    d.Set("service", o.Service)
    d.Set("category", o.Category)
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

func createSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    vsys, rb, o := parseSecurityPolicy(d)
    o.Defaults()

    if err := fw.Policies.Security.VerifiableSet(vsys, rb, o); err != nil {
        return err
    }

    d.SetId(o.Name)
    return nil
}

func readSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    name := d.Get("name").(string)
    vsys := d.Get("vsys").(string)
    rb := d.Get("rulebase").(string)

    o, err := fw.Policies.Security.Get(vsys, rb, name)
    if err != nil {
        d.SetId("")
        return nil
    }

    saveDataSecurityPolicy(d, o)
    return nil
}

func updateSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
    var err error
    fw := meta.(*pango.Firewall)
    vsys, rb, o := parseSecurityPolicy(d)
    o.Defaults()

    lo, err := fw.Policies.Security.Get(vsys, rb, o.Name)
    if err == nil {
        lo.Copy(o)
        err = fw.Policies.Security.Edit(vsys, rb, lo)
    } else {
        err = fw.Policies.Security.Set(vsys, rb, o)
    }

    if err == nil {
        saveDataSecurityPolicy(d, o)
    }
    return err
}

func deleteSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    vsys := d.Get("vsys").(string)
    rb := d.Get("rulebase").(string)
    name := d.Get("name").(string)

    _ = fw.Policies.Security.Delete(vsys, rb, name)
    d.SetId("")
    return nil
}
