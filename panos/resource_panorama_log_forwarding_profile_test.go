package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist/action"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaLogForwardingProfile_basic(t *testing.T) {
	minVersion := version.Number{8, 0, 0, ""}
	minTimeoutVersion := version.Number{9, 0, 0, ""}

	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
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

	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
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
		CheckDestroy: testAccPanosPanoramaLogForwardingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaLogForwardingProfileConfig(dg, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, name, tout),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLogForwardingProfileExists("panos_panorama_log_forwarding_profile.test", &o, &ml, &mla),
					testAccCheckPanosPanoramaLogForwardingProfileAttributes(&o, &ml, &mla, name, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, tout),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaLogForwardingProfileExists(n string, o *logfwd.Entry, ml *[]matchlist.Entry, mla *map[string][]action.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaLogForwardingProfileId(rs.Primary.ID)
		v, err := pano.Objects.LogForwardingProfile.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		var mlLocal []matchlist.Entry
		var mlaLocal map[string][]action.Entry

		mlNames, err := pano.Objects.LogForwardingProfileMatchList.GetList(dg, name)
		if err != nil {
			return err
		}

		mlLocal = make([]matchlist.Entry, 0, len(mlNames))
		mlaLocal = make(map[string][]action.Entry)

		for i := range mlNames {
			mle, err := pano.Objects.LogForwardingProfileMatchList.Get(dg, name, mlNames[i])
			if err != nil {
				return err
			}
			mlLocal = append(mlLocal, mle)
			aNames, err := pano.Objects.LogForwardingProfileMatchListAction.GetList(dg, name, mlNames[i])
			if err != nil {
				return err
			}
			actionList := make([]action.Entry, 0, len(aNames))
			for j := range aNames {
				ae, err := pano.Objects.LogForwardingProfileMatchListAction.Get(dg, name, mlNames[i], aNames[j])
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

func testAccCheckPanosPanoramaLogForwardingProfileAttributes(o *logfwd.Entry, ml *[]matchlist.Entry, mla *map[string][]action.Entry, name, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag string, tout int) resource.TestCheckFunc {
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

func testAccPanosPanoramaLogForwardingProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_log_forwarding_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaLogForwardingProfileId(rs.Primary.ID)
			_, err := pano.Objects.LogForwardingProfile.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaLogForwardingProfileConfig(dg, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, tag, name string, tout int) string {
	return fmt.Sprintf(`
variable "snmp1" {
    type = string
    default = %q
}

variable "snmp2" {
    type = string
    default = %q
}

variable "syslog1" {
    type = string
    default = %q
}

variable "syslog2" {
    type = string
    default = %q
}

variable "email1" {
    type = string
    default = %q
}

variable "email2" {
    type = string
    default = %q
}

variable "http1" {
    type = string
    default = %q
}

variable "http2" {
    type = string
    default = %q
}

resource "panos_panorama_device_group" "x" {
    name = %q
    description = "lfp acctest"
}

resource "panos_panorama_administrative_tag" "x" {
    device_group = panos_panorama_device_group.x.name
    name = %q
    color = "color12"
}

resource "panos_panorama_log_forwarding_profile" "test" {
    device_group = panos_panorama_device_group.x.name
    name = %q
    description = "lfp acctest"
    match_list {
        name = "ml-3"
        description = "ml3 desc"
        log_type = "threat"
        email_server_profiles = [
            var.email1,
        ]
        syslog_server_profiles = [
            var.syslog1,
            var.syslog2,
        ]
        http_server_profiles = [
            var.http2,
        ]
    }
    match_list {
        name = "ml-2"
        description = "acctest for lfp"
        log_type = "auth"
        snmptrap_server_profiles = [
            var.snmp1,
            var.snmp2,
        ]
        email_server_profiles = [
            var.email2,
        ]
        http_server_profiles = [
            var.http1,
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
                        panos_panorama_administrative_tag.x.name,
                    ]
                }
            }
        }
    }
}
`, snmp1, snmp2, syslog1, syslog2, email1, email2, http1, http2, dg, tag, name, tout)
}
