package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist/action"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaLogForwardingProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaLogForwardingProfile,
		Read:   readPanoramaLogForwardingProfile,
		Update: createUpdatePanoramaLogForwardingProfile,
		Delete: deletePanoramaLogForwardingProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: logForwardingProfileSchema(true),
	}
}

func parsePanoramaLogForwardingProfile(d *schema.ResourceData) (string, logfwd.Entry, []matchlist.Entry, map[string][]action.Entry) {
	dg := d.Get("device_group").(string)
	o, ml, mla := loadLogForwardingProfile(d)

	return dg, o, ml, mla
}

func parsePanoramaLogForwardingProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaLogForwardingProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdatePanoramaLogForwardingProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o, ml, mla := parsePanoramaLogForwardingProfile(d)

	if err := pano.Objects.LogForwardingProfile.SetWithoutSubconfig(dg, o); err != nil {
		return err
	}

	if err := pano.Objects.LogForwardingProfileMatchList.Set(dg, o.Name, ml...); err != nil {
		return err
	}

	for _, entry := range ml {
		if err := pano.Objects.LogForwardingProfileMatchListAction.Set(dg, o.Name, entry.Name, mla[entry.Name]...); err != nil {
			return err
		}
	}

	d.SetId(buildPanoramaLogForwardingProfileId(dg, o.Name))
	return readPanoramaLogForwardingProfile(d, meta)
}

func readPanoramaLogForwardingProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaLogForwardingProfileId(d.Id())

	o, err := pano.Objects.LogForwardingProfile.Get(dg, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	var ml []matchlist.Entry
	var mla map[string][]action.Entry

	mlNames, err := pano.Objects.LogForwardingProfileMatchList.GetList(dg, name)
	if err != nil {
		return err
	}

	if len(mlNames) > 0 {
		ml = make([]matchlist.Entry, 0, len(mlNames))
		mla = make(map[string][]action.Entry)
		for i := range mlNames {
			mle, err := pano.Objects.LogForwardingProfileMatchList.Get(dg, name, mlNames[i])
			if err != nil {
				return err
			}
			ml = append(ml, mle)
			aNames, err := pano.Objects.LogForwardingProfileMatchListAction.GetList(dg, name, mlNames[i])
			if err != nil {
				return err
			}
			if len(aNames) != 0 {
				actionList := make([]action.Entry, 0, len(aNames))
				for j := range aNames {
					ae, err := pano.Objects.LogForwardingProfileMatchListAction.Get(dg, name, mlNames[i], aNames[j])
					if err != nil {
						return err
					}
					actionList = append(actionList, ae)
				}
				mla[mle.Name] = actionList
			}
		}
	}

	d.Set("device_group", dg)
	saveLogForwardingProfile(d, o, ml, mla)

	return nil
}

func deletePanoramaLogForwardingProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaLogForwardingProfileId(d.Id())

	err := pano.Objects.LogForwardingProfile.Delete(dg, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
