package predefined

import (
	"github.com/PaloAltoNetworks/pango/objs/app"
	"github.com/PaloAltoNetworks/pango/objs/srvc"
	dlpft "github.com/PaloAltoNetworks/pango/predefined/dlp/filetype"
	tdbft "github.com/PaloAltoNetworks/pango/predefined/tdb/filetype"
	"github.com/PaloAltoNetworks/pango/predefined/threat"
	"github.com/PaloAltoNetworks/pango/util"
)

type Panorama struct {
	Application *app.Predefined
	DlpFileType *dlpft.Panorama
	Services    *srvc.Predefined
	TdbFileType *tdbft.Panorama
	Threat      *threat.Panorama
}

func PanoramaNamespace(x util.XapiClient) *Panorama {
	return &Panorama{
		Application: app.PredefinedNamespace(x),
		DlpFileType: dlpft.PanoramaNamespace(x),
		Services:    srvc.PredefinedNamespace(x),
		TdbFileType: tdbft.PanoramaNamespace(x),
		Threat:      threat.PanoramaNamespace(x),
	}
}
