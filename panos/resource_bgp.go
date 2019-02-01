package panos

import (
	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBgp() *schema.Resource {
	return &schema.Resource{
		Create: createBgp,
		Read:   readBgp,
		Update: updateBgp,
		Delete: deleteBgp,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpSchema(false),
	}
}

func bgpSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"enable": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"router_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"as_number": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"bfd_profile": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"reject_default_route": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"install_route": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"aggregate_med": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"default_local_preference": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "100",
		},
		"as_format": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      bgp.AsFormat2Byte,
			ValidateFunc: validateStringIn(bgp.AsFormat2Byte, bgp.AsFormat4Byte),
		},
		"always_compare_med": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"deterministic_med_comparison": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"ecmp_multi_as": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enforce_first_as": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"enable_graceful_restart": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"stale_route_time": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  120,
		},
		"local_restart_time": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  120,
		},
		"max_peer_restart_time": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  120,
		},
		"reflector_cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"confederation_member_as": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"allow_redistribute_default_route": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema()
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgp(d *schema.ResourceData) (string, bgp.Config) {
	vr := d.Get("virtual_router").(string)

	o := bgp.Config{
		Enable:                        d.Get("enable").(bool),
		RouterId:                      d.Get("router_id").(string),
		AsNumber:                      d.Get("as_number").(string),
		BfdProfile:                    d.Get("bfd_profile").(string),
		RejectDefaultRoute:            d.Get("reject_default_route").(bool),
		InstallRoute:                  d.Get("install_route").(bool),
		AggregateMed:                  d.Get("aggregate_med").(bool),
		DefaultLocalPreference:        d.Get("default_local_preference").(string),
		AsFormat:                      d.Get("as_format").(string),
		AlwaysCompareMed:              d.Get("always_compare_med").(bool),
		DeterministicMedComparison:    d.Get("deterministic_med_comparison").(bool),
		EcmpMultiAs:                   d.Get("ecmp_multi_as").(bool),
		EnforceFirstAs:                d.Get("enforce_first_as").(bool),
		EnableGracefulRestart:         d.Get("enable_graceful_restart").(bool),
		StaleRouteTime:                d.Get("stale_route_time").(int),
		LocalRestartTime:              d.Get("local_restart_time").(int),
		MaxPeerRestartTime:            d.Get("max_peer_restart_time").(int),
		ReflectorClusterId:            d.Get("reflector_cluster_id").(string),
		ConfederationMemberAs:         d.Get("confederation_member_as").(string),
		AllowRedistributeDefaultRoute: d.Get("allow_redistribute_default_route").(bool),
	}

	return vr, o
}

func saveBgp(d *schema.ResourceData, vr string, o bgp.Config) {
	d.Set("virtual_router", vr)

	d.Set("enable", o.Enable)
	d.Set("router_id", o.RouterId)
	d.Set("as_number", o.AsNumber)
	d.Set("bfd_profile", o.BfdProfile)
	d.Set("reject_default_route", o.RejectDefaultRoute)
	d.Set("install_route", o.InstallRoute)
	d.Set("aggregate_med", o.AggregateMed)
	d.Set("default_local_preference", o.DefaultLocalPreference)
	d.Set("as_format", o.AsFormat)
	d.Set("always_compare_med", o.AlwaysCompareMed)
	d.Set("deterministic_med_comparison", o.DeterministicMedComparison)
	d.Set("ecmp_multi_as", o.EcmpMultiAs)
	d.Set("enforce_first_as", o.EnforceFirstAs)
	d.Set("enable_graceful_restart", o.EnableGracefulRestart)
	d.Set("stale_route_time", o.StaleRouteTime)
	d.Set("local_restart_time", o.LocalRestartTime)
	d.Set("max_peer_restart_time", o.MaxPeerRestartTime)
	d.Set("reflector_cluster_id", o.ReflectorClusterId)
	d.Set("confederation_member_as", o.ConfederationMemberAs)
	d.Set("allow_redistribute_default_route", o.AllowRedistributeDefaultRoute)
}

func createBgp(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, o := parseBgp(d)

	if err = fw.Network.BgpConfig.Set(vr, o); err != nil {
		return err
	}

	d.SetId(vr)
	return readBgp(d, meta)
}

func readBgp(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr := d.Id()

	o, err := fw.Network.BgpConfig.Get(vr)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgp(d, vr, o)

	return nil
}

func updateBgp(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, o := parseBgp(d)

	lo, err := fw.Network.BgpConfig.Get(vr)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpConfig.Edit(vr, lo); err != nil {
		return err
	}

	return readBgp(d, meta)
}

func deleteBgp(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr := d.Id()

	err := fw.Network.BgpConfig.Delete(vr)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
