package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist/action"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosLogForwardingProfile_basic(t *testing.T) {
	minVersion := version.Number{8, 0, 0, ""}
	minTimeoutVersion := version.Number{9, 0, 0, ""}

	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccPanosVersion.Gte(minVersion) {
		t.Skip("This test is only valid for PAN-OS 8.0+")
	}

	var (
		o    logfwd.Entry
		ml   []matchlist.Entry
		mla  map[string][]action.Entry
		tout int
	)

	if testAccPanosVersion.Gte(minTimeoutVersion) {
		tout = 5
	}

	snmp1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	snmp2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	syslog1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	syslog2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	email1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	email2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	http1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	http2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	tag := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosLogForwardingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogForwardingProfileConfig(snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, name, tout),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLogForwardingProfileExists("panos_log_forwarding_profile.test", &o, &ml, &mla),
					testAccCheckPanosLogForwardingProfileAttributes(&o, &ml, &mla, name, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, tout),
				),
			},
		},
	})
}

func testAccCheckPanosLogForwardingProfileExists(n string, o *logfwd.Entry, ml *[]matchlist.Entry, mla *map[string][]action.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseLogForwardingProfileId(rs.Primary.ID)
		v, err := fw.Objects.LogForwardingProfile.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		var mlLocal []matchlist.Entry
		var mlaLocal map[string][]action.Entry

		mlNames, err := fw.Objects.LogForwardingProfileMatchList.GetList(vsys, name)
		if err != nil {
			return err
		}

		mlLocal = make([]matchlist.Entry, 0, len(mlNames))
		mlaLocal = make(map[string][]action.Entry)

		for i := range mlNames {
			mle, err := fw.Objects.LogForwardingProfileMatchList.Get(vsys, name, mlNames[i])
			if err != nil {
				return err
			}
			mlLocal = append(mlLocal, mle)
			aNames, err := fw.Objects.LogForwardingProfileMatchListAction.GetList(vsys, name, mlNames[i])
			if err != nil {
				return err
			}
			actionList := make([]action.Entry, 0, len(aNames))
			for j := range aNames {
				ae, err := fw.Objects.LogForwardingProfileMatchListAction.Get(vsys, name, mlNames[i], aNames[j])
				if err != nil {
					return err
				}
				actionList = append(actionList, ae)
			}
			mlaLocal[mle.Name] = actionList
		}

		*ml = mlLocal
		*mla = mlaLocal

		return nil
	}
}

func testAccCheckPanosLogForwardingProfileAttributes(o *logfwd.Entry, ml *[]matchlist.Entry, mla *map[string][]action.Entry, name, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag string, tout int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != "lfp acctest" {
			return fmt.Errorf("LFP description is %q", o.Description)
		}

		var m matchlist.Entry
		var actionList []action.Entry
		var a action.Entry

		if len(*ml) != 3 {
			return fmt.Errorf("Match list length is %d", len(*ml))
		}

		m = (*ml)[0]
		if m.Name != "ml-3" {
			return fmt.Errorf("ML1 name is %q", m.Name)
		}
		if m.Description != "ml3 desc" {
			return fmt.Errorf("ML1 description is %q", m.Description)
		}
		if m.LogType != "threat" {
			return fmt.Errorf("ML1 log type is %s", m.LogType)
		}
		if m.SendToPanorama {
			return fmt.Errorf("ML1 send to panorama is %t", m.SendToPanorama)
		}
		if len(m.SnmpProfiles) != 0 {
			return fmt.Errorf("ML1 snmp profiles is %#v", m.SnmpProfiles)
		}
		if len(m.EmailProfiles) != 1 || m.EmailProfiles[0] != email1 {
			return fmt.Errorf("ML1 email profiles is %#v", m.EmailProfiles)
		}
		if len(m.SyslogProfiles) != 2 || (m.SyslogProfiles[0] != syslog1 && m.SyslogProfiles[1] != syslog1) || (m.SyslogProfiles[0] != syslog2 && m.SyslogProfiles[1] != syslog2) {
			return fmt.Errorf("ML1 syslog profiles is %#v", m.SyslogProfiles)
		}
		if len(m.HttpProfiles) != 1 || m.HttpProfiles[0] != http2 {
			return fmt.Errorf("ML1 http profiles is %#v", m.HttpProfiles)
		}

		actionList = (*mla)[m.Name]
		if len(actionList) != 0 {
			return fmt.Errorf("ML1 has action list: %#v", actionList)
		}

		m = (*ml)[1]
		if m.Name != "ml-2" {
			return fmt.Errorf("ML2 name is %q", m.Name)
		}
		if m.Description != "acctest for lfp" {
			return fmt.Errorf("ML2 description is %q", m.Description)
		}
		if m.LogType != "auth" {
			return fmt.Errorf("ML2 log type is %s", m.LogType)
		}
		if m.SendToPanorama {
			return fmt.Errorf("ML2 send to panorama is %t", m.SendToPanorama)
		}
		if len(m.SnmpProfiles) != 2 || (m.SnmpProfiles[0] != snmp1 && m.SnmpProfiles[1] != snmp1) || (m.SnmpProfiles[0] != snmp2 && m.SnmpProfiles[1] != snmp2) {
			return fmt.Errorf("ML2 snmp profiles is %#v", m.SnmpProfiles)
		}
		if len(m.EmailProfiles) != 1 || m.EmailProfiles[0] != email2 {
			return fmt.Errorf("ML2 email profiles is %#v", m.EmailProfiles)
		}
		if len(m.SyslogProfiles) != 0 {
			return fmt.Errorf("ML2 syslog profiles is %#v", m.SyslogProfiles)
		}
		if len(m.HttpProfiles) != 1 || m.HttpProfiles[0] != http1 {
			return fmt.Errorf("ML2 http profiles is %#v", m.HttpProfiles)
		}

		actionList = (*mla)[m.Name]
		if len(actionList) != 0 {
			return fmt.Errorf("ML2 has action list: %#v", actionList)
		}

		m = (*ml)[2]
		if m.Name != "ml-1" {
			return fmt.Errorf("ML3 name is %q", m.Name)
		}
		if m.Description != "rain" {
			return fmt.Errorf("ML3 description is %q", m.Description)
		}
		if m.LogType != "data" {
			return fmt.Errorf("ML3 log type is %s", m.LogType)
		}
		if !m.SendToPanorama {
			return fmt.Errorf("ML3 send to panorama is %t", m.SendToPanorama)
		}
		if len(m.SnmpProfiles) != 0 {
			return fmt.Errorf("ML3 snmp profiles is %#v", m.SnmpProfiles)
		}
		if len(m.EmailProfiles) != 0 {
			return fmt.Errorf("ML3 email profiles is %#v", m.EmailProfiles)
		}
		if len(m.SyslogProfiles) != 0 {
			return fmt.Errorf("ML3 syslog profiles is %#v", m.SyslogProfiles)
		}
		if len(m.HttpProfiles) != 0 {
			return fmt.Errorf("ML3 http profiles is %#v", m.HttpProfiles)
		}

		actionList = (*mla)[m.Name]
		if len(actionList) != 1 {
			return fmt.Errorf("ML3 has action list: %#v", actionList)
		}

		a = actionList[0]
		if a.Name != "act-now" {
			return fmt.Errorf("Action1 name is %q", a.Name)
		}
		if a.ActionType != action.ActionTypeTagging {
			return fmt.Errorf("Action1 action type is %s", a.ActionType)
		}
		if a.Action != action.ActionAddTag {
			return fmt.Errorf("Action1 action is %s", a.Action)
		}
		if a.Target != action.TargetSource {
			return fmt.Errorf("Action1 target is %s", a.Target)
		}
		if a.Registration != action.RegistrationLocal {
			return fmt.Errorf("Action1 reg is %s", a.Registration)
		}
		if len(a.Tags) != 1 || a.Tags[0] != tag {
			return fmt.Errorf("Action1 tags is %#v", a.Tags)
		}
		if a.Timeout != tout {
			return fmt.Errorf("Action1 timeout is %d, not %d", a.Timeout, tout)
		}

		return nil
	}
}

func testAccPanosLogForwardingProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_log_forwarding_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseLogForwardingProfileId(rs.Primary.ID)
			_, err := fw.Objects.LogForwardingProfile.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccLogForwardingProfileConfig(snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, name string, tout int) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

resource "panos_snmptrap_server_profile" "a" {
    name = %q
    v2c_server {
        name = "some server"
        manager = "snmp.example.com"
        community = "public"
    }
}

resource "panos_snmptrap_server_profile" "b" {
    name = %q
    v2c_server {
        name = "server2"
        manager = "snmp.example.com"
        community = "public"
    }
}

resource "panos_syslog_server_profile" "a" {
    name = %q
    syslog_server {
        name = "server3"
        server = "syslog.example.com"
    }
}

resource "panos_syslog_server_profile" "b" {
    name = %q
    syslog_server {
        name = "server4"
        server = "syslog.example.com"
    }
}

resource "panos_email_server_profile" "a" {
    name = %q
    email_server {
        name = "server5"
        display_name = "foobar"
        from_email = "wu@example.com"
        to_email = "tang@example.com"
        email_gateway = "clan.example.com"
    }
}

resource "panos_email_server_profile" "b" {
    name = %q
    email_server {
        name = "server6"
        display_name = "foobar"
        from_email = "black@example.com"
        to_email = "eyed@example.com"
        email_gateway = "peas.example.com"
    }
}

resource "panos_http_server_profile" "a" {
    name = %q
    http_server {
        name = "server7"
        address = "foo.example.com"
        certificate_profile = data.panos_system_info.x.version_major >= 9 ? "None" : ""
        tls_version = data.panos_system_info.x.version_major >= 9 ? "1.2" : ""
    }
}

resource "panos_http_server_profile" "b" {
    name = %q
    http_server {
        name = "server8"
        address = "bar.example.com"
        certificate_profile = data.panos_system_info.x.version_major >= 9 ? "None" : ""
        tls_version = data.panos_system_info.x.version_major >= 9 ? "1.2" : ""
    }
}

resource "panos_administrative_tag" "x" {
    name = %q
    color = "color12"
}

resource "panos_log_forwarding_profile" "test" {
    name = %q
    description = "lfp acctest"
    match_list {
        name = "ml-3"
        description = "ml3 desc"
        log_type = "threat"
        email_server_profiles = [
            panos_email_server_profile.a.name,
        ]
        syslog_server_profiles = [
            panos_syslog_server_profile.a.name,
            panos_syslog_server_profile.b.name,
        ]
        http_server_profiles = [
            panos_http_server_profile.b.name,
        ]
    }
    match_list {
        name = "ml-2"
        description = "acctest for lfp"
        log_type = "auth"
        snmptrap_server_profiles = [
            panos_snmptrap_server_profile.a.name,
            panos_snmptrap_server_profile.b.name,
        ]
        email_server_profiles = [
            panos_email_server_profile.b.name,
        ]
        http_server_profiles = [
            panos_http_server_profile.a.name,
        ]
    }
    match_list {
        name = "ml-1"
        description = "rain"
        log_type = "data"
        send_to_panorama = true
        action {
            name = "act-now"
            tagging_integration {
                timeout = %d
                local_registration {
                    tags = [
                        panos_administrative_tag.x.name,
                    ]
                }
            }
        }
    }
}
`, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, name, tout)
}
