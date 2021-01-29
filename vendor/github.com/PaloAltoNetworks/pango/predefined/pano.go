package predefined

import (
	dlpft "github.com/PaloAltoNetworks/pango/predefined/dlp/filetype"
	tdbft "github.com/PaloAltoNetworks/pango/predefined/tdb/filetype"
	"github.com/PaloAltoNetworks/pango/predefined/threat"
	"github.com/PaloAltoNetworks/pango/util"
)

type Panorama struct {
	DlpFileType *dlpft.Panorama
	TdbFileType *tdbft.Panorama
	Threat      *threat.Panorama
}

func PanoramaNamespace(x util.XapiClient) *Panorama {
	return &Panorama{
		DlpFileType: dlpft.PanoramaNamespace(x),
		TdbFileType: tdbft.PanoramaNamespace(x),
		Threat:      threat.PanoramaNamespace(x),
	}
}
