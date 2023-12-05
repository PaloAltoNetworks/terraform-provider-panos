package panos

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/decryption"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceDecryptionRules() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()
	s["rulebase"] = rulebaseSchema()

	return &schema.Resource{
		Read: dataSourceDecryptionRulesRead,

		Schema: s,
	}
}

func dataSourceDecryptionRulesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)

	d.Set("vsys", vsys)
	d.Set("device_group", dg)
	d.Set("rulebase", base)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vsys
		listing, err = con.Policies.Decryption.GetList(vsys)
	case *pango.Panorama:
		id = strings.Join([]string{dg, base}, IdSeparator)
		listing, err = con.Policies.Decryption.GetList(dg, base)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceDecryptionRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDecryptionRuleRead,

		Schema: decryptionRuleGroupSchema(false),
	}
}

func dataSourceDecryptionRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o decryption.Entry

	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("vsys", vsys)
	d.Set("name", name)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = strings.Join([]string{vsys, name}, IdSeparator)
		o, err = con.Policies.Decryption.Get(vsys, name)
	case *pango.Panorama:
		id = strings.Join([]string{dg, base, name}, IdSeparator)
		o, err = con.Policies.Decryption.Get(dg, base, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveDecryptionRules(d, []decryption.Entry{o})

	return nil
}

// Resource.
func resourceDecryptionRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateDecryptionRuleGroup,
		Read:   readDecryptionRuleGroup,
		Update: createUpdateDecryptionRuleGroup,
		Delete: deleteDecryptionRuleGroup,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: decryptionRuleGroupSchema(true),
	}
}

func createUpdateDecryptionRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var prevNames []string

	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	vsys := d.Get("vsys").(string)
	move := movementAtoi(d.Get("position_keyword").(string))
	oRule := d.Get("position_reference").(string)
	rules, auditComments := loadDecryptionRules(d)

	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("vsys", vsys)
	d.Set("position_keyword", movementItoa(move))
	d.Set("position_reference", oRule)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}

	if d.Id() != "" {
		_, _, _, _, _, prevNames = parseDecryptionRuleGroupId(d.Id())
	}

	id := buildDecryptionRuleGroupId(dg, base, vsys, move, oRule, rules)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Decryption.ConfigureRules(vsys, rules, auditComments, false, move, oRule, prevNames)
	case *pango.Panorama:
		err = con.Policies.Decryption.ConfigureRules(dg, base, rules, auditComments, false, move, oRule, prevNames)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readDecryptionRuleGroup(d, meta)
}

func readDecryptionRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []decryption.Entry

	dg, base, vsys, move, oRule, names := parseDecryptionRuleGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Policies.Decryption.GetAll(vsys)
	case *pango.Panorama:
		listing, err = con.Policies.Decryption.GetAll(dg, base)
	}

	if err != nil {
		d.SetId("")
		return nil
	}

	fIdx, oIdx := -1, -1
	for i := range listing {
		if listing[i].Name == names[0] {
			fIdx = i
		} else if listing[i].Name == oRule {
			oIdx = i
		}
		if fIdx != -1 && (oIdx != -1 || oRule == "") {
			break
		}
	}

	if fIdx == -1 {
		// First rule is MIA, but others may be present, so report an
		// empty ruleset to force rules to be recreated.
		d.Set("rule", nil)
		return nil
	} else if oIdx == -1 && movementIsRelative(move) {
		return fmt.Errorf("Can't position group %s %q: rule is not present", movementItoa(move), oRule)
	} else if move == util.MoveTop && fIdx != 0 {
		d.Set("position_keyword", "")
	}

	dlist := make([]decryption.Entry, 0, len(names))
	for i := 0; i+fIdx < len(listing) && i < len(names); i++ {
		if listing[i+fIdx].Name != names[i] {
			break
		}

		dlist = append(dlist, listing[i+fIdx])
	}

	if move == util.MoveBottom && dlist[len(dlist)-1].Name != listing[len(listing)-1].Name {
		d.Set("position_keyword", "")
	}
	saveDecryptionRules(d, dlist)

	return nil
}

func deleteDecryptionRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, base, vsys, _, _, names := parseSecurityRuleGroupId(d.Id())

	ilist := make([]interface{}, 0, len(names))
	for _, x := range names {
		ilist = append(ilist, x)
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Decryption.Delete(vsys, ilist...)
	case *pango.Panorama:
		err = con.Policies.Decryption.Delete(dg, base, ilist...)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func decryptionRuleGroupSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group":       deviceGroupSchema(),
		"rulebase":           rulebaseSchema(),
		"vsys":               vsysSchema("vsys1"),
		"position_keyword":   positionKeywordSchema(),
		"position_reference": positionReferenceSchema(),
		"rule": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The rule name.",
						Required:    true,
					},
					"description": {
						Type:        schema.TypeString,
						Description: "The description.",
						Optional:    true,
					},
					"source_zones": {
						Type:        schema.TypeSet,
						Description: "List of source zones.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"source_addresses": {
						Type:        schema.TypeSet,
						Description: "List of source addresses.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"negate_source": {
						Type:        schema.TypeBool,
						Description: "Negate the source addresses.",
						Optional:    true,
					},
					"source_users": {
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_zones": {
						Type:        schema.TypeSet,
						Description: "List of destination zones.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_addresses": {
						Type:        schema.TypeSet,
						Description: "List of destination addresses.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"negate_destination": {
						Type:        schema.TypeBool,
						Description: "Negate the destination addresses.",
						Optional:    true,
					},
					"tags": tagSchema(),
					"disabled": {
						Type:        schema.TypeBool,
						Description: "Disable this rule.",
						Optional:    true,
					},
					"services": {
						Type:        schema.TypeSet,
						Description: "List of services.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"url_categories": {
						Type:        schema.TypeSet,
						Description: "List of URL categories.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"action": {
						Type:        schema.TypeString,
						Description: "Action to take.",
						Optional:    true,
						Default:     decryption.ActionNoDecrypt,
						ValidateFunc: validateStringIn(
							decryption.ActionNoDecrypt,
							decryption.ActionDecrypt,
							decryption.ActionDecryptAndForward,
							"",
						),
					},
					"decryption_type": {
						Type:        schema.TypeString,
						Description: "The decryption type.",
						Optional:    true,
						ValidateFunc: validateStringIn(
							decryption.DecryptionTypeSslForwardProxy,
							decryption.DecryptionTypeSshProxy,
							decryption.DecryptionTypeSslInboundInspection,
						),
					},
					"ssl_certificate": {
						Type:        schema.TypeString,
						Description: "(PAN-OS 10.1 and below) The SSL certificate.",
						Optional:    true,
					},
					"ssl_certificates": {
						Type:        schema.TypeSet,
						Description: "(PAN-OS 10.2+) List of SSL decryption certs.",
						Optional:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"decryption_profile": {
						Type:        schema.TypeString,
						Description: "The decryption profile.",
						Optional:    true,
					},
					"forwarding_profile": {
						Type:        schema.TypeString,
						Description: "Forwarding profile.",
						Optional:    true,
					},
					"uuid":      uuidSchema(),
					"group_tag": groupTagSchema(),
					"source_hips": {
						Type:        schema.TypeSet,
						Description: "List of source HIP devices.",
						Optional:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_hips": {
						Type:        schema.TypeSet,
						Description: "List of destination HIP devices.",
						Optional:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"log_successful_tls_handshakes": {
						Type:        schema.TypeBool,
						Description: "Log successful TLS handshakes.",
						Optional:    true,
					},
					"log_failed_tls_handshakes": {
						Type:        schema.TypeBool,
						Description: "Log failed TLS handshakes.",
						Optional:    true,
					},
					"log_setting": {
						Type:        schema.TypeString,
						Description: "The log setting.",
						Optional:    true,
					},
					"target":        targetSchema(false),
					"negate_target": negateTargetSchema(),
					"audit_comment": auditCommentSchema(),
				},
			},
		},
	}

	if !isResource {
		delete(ans, "position_keyword")
		delete(ans, "position_reference")
		ans["name"] = &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		}

		computed(ans, "", []string{"device_group", "vsys", "rulebase", "name"})
	}

	return ans
}

func loadDecryptionRules(d *schema.ResourceData) ([]decryption.Entry, map[string]string) {
	auditComments := make(map[string]string)
	rlist := d.Get("rule").([]interface{})
	list := make([]decryption.Entry, 0, len(rlist))
	for i := range rlist {
		x := rlist[i].(map[string]interface{})
		auditComments[x["name"].(string)] = x["audit_comment"].(string)
		list = append(list, decryption.Entry{
			Name:                       x["name"].(string),
			Description:                x["description"].(string),
			SourceZones:                setAsList(x["source_zones"].(*schema.Set)),
			SourceAddresses:            setAsList(x["source_addresses"].(*schema.Set)),
			NegateSource:               x["negate_source"].(bool),
			SourceUsers:                setAsList(x["source_users"].(*schema.Set)),
			DestinationZones:           setAsList(x["destination_zones"].(*schema.Set)),
			DestinationAddresses:       setAsList(x["destination_addresses"].(*schema.Set)),
			NegateDestination:          x["negate_destination"].(bool),
			Tags:                       asStringList(x["tags"].([]interface{})),
			Disabled:                   x["disabled"].(bool),
			Services:                   setAsList(x["services"].(*schema.Set)),
			UrlCategories:              setAsList(x["url_categories"].(*schema.Set)),
			Action:                     x["action"].(string),
			DecryptionType:             x["decryption_type"].(string),
			SslCertificate:             x["ssl_certificate"].(string),
			SslCertificates:            setAsList(x["ssl_certificates"].(*schema.Set)),
			DecryptionProfile:          x["decryption_profile"].(string),
			Targets:                    loadTarget(x["target"]),
			NegateTarget:               x["negate_target"].(bool),
			ForwardingProfile:          x["forwarding_profile"].(string),
			GroupTag:                   x["group_tag"].(string),
			SourceHips:                 setAsList(x["source_hips"].(*schema.Set)),
			DestinationHips:            setAsList(x["destination_hips"].(*schema.Set)),
			LogSuccessfulTlsHandshakes: x["log_successful_tls_handshakes"].(bool),
			LogFailedTlsHandshakes:     x["log_failed_tls_handshakes"].(bool),
			LogSetting:                 x["log_setting"].(string),
		})
	}

	return list, auditComments
}

func dumpDecryptionRule(o decryption.Entry) map[string]interface{} {
	return map[string]interface{}{
		"name":                          o.Name,
		"description":                   o.Description,
		"source_zones":                  listAsSet(o.SourceZones),
		"source_addresses":              listAsSet(o.SourceAddresses),
		"negate_source":                 o.NegateSource,
		"source_users":                  listAsSet(o.SourceUsers),
		"destination_zones":             listAsSet(o.DestinationZones),
		"destination_addresses":         listAsSet(o.DestinationAddresses),
		"negate_destination":            o.NegateDestination,
		"tags":                          o.Tags,
		"disabled":                      o.Disabled,
		"services":                      listAsSet(o.Services),
		"url_categories":                listAsSet(o.UrlCategories),
		"action":                        o.Action,
		"decryption_type":               o.DecryptionType,
		"ssl_certificate":               o.SslCertificate,
		"ssl_certificates":              listAsSet(o.SslCertificates),
		"decryption_profile":            o.DecryptionProfile,
		"target":                        dumpTarget(o.Targets),
		"negate_target":                 o.NegateTarget,
		"forwarding_profile":            o.ForwardingProfile,
		"uuid":                          o.Uuid,
		"group_tag":                     o.GroupTag,
		"source_hips":                   listAsSet(o.SourceHips),
		"destination_hips":              listAsSet(o.DestinationHips),
		"log_successful_tls_handshakes": o.LogSuccessfulTlsHandshakes,
		"log_failed_tls_handshakes":     o.LogFailedTlsHandshakes,
		"log_setting":                   o.LogSetting,
		"audit_comment":                 "",
	}
}

func saveDecryptionRules(d *schema.ResourceData, rules []decryption.Entry) {
	if len(rules) == 0 {
		d.Set("rule", nil)
		return
	}

	list := make([]interface{}, 0, len(rules))
	for _, x := range rules {
		list = append(list, dumpDecryptionRule(x))
	}

	if err := d.Set("rule", list); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildDecryptionRuleGroupId(a, b, c string, d int, e string, f []decryption.Entry) string {
	names := make([]string, 0, len(f))
	for _, x := range f {
		names = append(names, x.Name)
	}
	return strings.Join([]string{a, b, c, strconv.Itoa(d), e, base64Encode(names)}, IdSeparator)
}

func parseDecryptionRuleGroupId(v string) (string, string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[3])
	return t[0], t[1], t[2], move, t[4], base64Decode(t[5])
}
