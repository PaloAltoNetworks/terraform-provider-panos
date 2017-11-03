package main

import (
    "github.com/PaloAltoNetworks/pango"
    "github.com/PaloAltoNetworks/pango/netw/eth"

    "github.com/hashicorp/terraform/helper/schema"
)


func resourceEthernetInterface() *schema.Resource {
    return &schema.Resource{
        Create: createEthernetInterface,
        Read: readEthernetInterface,
        Update: updateEthernetInterface,
        Delete: deleteEthernetInterface,

        Schema: map[string] *schema.Schema{
            "name": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                Description: "The ethernet interface's name",
            },
            "vsys": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                Description: "The vsys to import this ethernet interface into",
            },
            "mode": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                Description: "The interface mode (layer3, layer2, virtual-wire, tap, ha, decrypt-mirror, aggregate-group)",
            },
            "static_ips": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Computed: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
                Description: "Administrative tags for the address object",
            },
            "enable_dhcp": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
                Computed: true,
            },
            "create_dhcp_default_route": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
                Computed: true,
            },
            "dhcp_default_route_metric": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
                Computed: true,
            },
            "ipv6_enabled": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
                Computed: true,
            },
            "management_profile": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "mtu": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
                Computed: true,
            },
            "adjust_tcp_mss": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
                Computed: true,
            },
            "netflow_profile": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "lldp_enabled": &schema.Schema{
                Type: schema.TypeBool,
                Optional: true,
                Computed: true,
            },
            "lldp_profile": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "link_speed": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "link_duplex": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "link_state": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "aggregate_group": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "comment": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "ipv4_mss_adjust": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
                Computed: true,
            },
            "ipv6_mss_adjust": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
                Computed: true,
            },
        },
    }
}

func parseEthernetInterface(d *schema.ResourceData) (string, eth.Entry) {
    vsys := d.Get("vsys").(string)
    o := eth.Entry{
        Name: d.Get("name").(string),
        Mode: d.Get("mode").(string),
        StaticIps: asStringList(d, "static_ips"),
        EnableDhcp: d.Get("enable_dhcp").(bool),
        CreateDhcpDefaultRoute: d.Get("create_dhcp_default_route").(bool),
        DhcpDefaultRouteMetric: d.Get("dhcp_default_route_metric").(int),
        Ipv6Enabled: d.Get("ipv6_enabled").(bool),
        ManagementProfile: d.Get("management_profile").(string),
        Mtu: d.Get("mtu").(int),
        AdjustTcpMss: d.Get("adjust_tcp_mss").(bool),
        NetflowProfile: d.Get("netflow_profile").(string),
        LldpEnabled: d.Get("lldp_enabled").(bool),
        LldpProfile: d.Get("lldp_profile").(string),
        LinkSpeed: d.Get("link_speed").(string),
        LinkDuplex: d.Get("link_duplex").(string),
        LinkState: d.Get("link_state").(string),
        AggregateGroup: d.Get("aggregate_group").(string),
        Comment: d.Get("comment").(string),
        Ipv4MssAdjust: d.Get("ipv4_mss_adjust").(int),
        Ipv6MssAdjust: d.Get("ipv6_mss_adjust").(int),
    }

    return vsys, o
}

func saveDataEthernetInterface(d *schema.ResourceData, o eth.Entry) {
    d.SetId(o.Name)
    d.Set("mode", o.Mode)
    d.Set("static_ips", o.StaticIps)
    d.Set("enable_dhcp", o.EnableDhcp)
    d.Set("create_dhcp_default_route", o.CreateDhcpDefaultRoute)
    d.Set("dhcp_default_route_metric", o.DhcpDefaultRouteMetric)
    d.Set("ipv6_enabled", o.Ipv6Enabled)
    d.Set("management_profile", o.ManagementProfile)
    d.Set("mtu", o.Mtu)
    d.Set("adjust_tcp_mss", o.AdjustTcpMss)
    d.Set("netflow_profile", o.NetflowProfile)
    d.Set("lldp_enabled", o.LldpEnabled)
    d.Set("lldp_profile", o.LldpProfile)
    d.Set("link_speed", o.LinkSpeed)
    d.Set("link_duplex", o.LinkDuplex)
    d.Set("link_state", o.LinkState)
    d.Set("aggregate_group", o.AggregateGroup)
    d.Set("comment", o.Comment)
    d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
    d.Set("ipv6_mss_adjust", o.Ipv6MssAdjust)
}

func createEthernetInterface(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    vsys, o := parseEthernetInterface(d)

    if err := fw.Network.EthernetInterface.Set(vsys, o); err != nil {
        return err
    }

    d.SetId(o.Name)
    return nil
}

func readEthernetInterface(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    name := d.Get("name").(string)

    o, err := fw.Network.EthernetInterface.Get(name)
    if err != nil {
        d.SetId("")
        return nil
    }

    saveDataEthernetInterface(d, o)
    return nil
}

func updateEthernetInterface(d *schema.ResourceData, meta interface{}) error {
    var err error
    fw := meta.(*pango.Firewall)
    vsys, o := parseEthernetInterface(d)

    lo, err := fw.Network.EthernetInterface.Get(o.Name)
    if err == nil {
        lo.Copy(o)
        err = fw.Network.EthernetInterface.Edit(vsys, lo)
    } else {
        err = fw.Network.EthernetInterface.Set(vsys, o)
    }

    if err == nil {
        saveDataEthernetInterface(d, o)
    }
    return err
}

func deleteEthernetInterface(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    vsys := d.Get("vsys").(string)
    name := d.Get("name").(string)

    _ = fw.Network.EthernetInterface.Delete(vsys, name)
    d.SetId("")
    return nil
}
